import datetime
import os
from urllib.parse import urljoin
import boto3
from chalice import Chalice, Cron
import requests

app = Chalice(app_name='time-aggregation-notifier')
app.debug = True if os.getenv('DEBUG', False) else False


def get_redmine_data(url: str, params: dict):
    """ Redmine の API から JSON 形式でデータを取得 """
    res = requests.get(url, params)
    return res.json()


def get_client_time_entries(issues: dict):
    """ Issue 一覧から該当するカスタムフィールド単位の作業時間を取得 """
    params = {
        'key': os.getenv('REDMINE_KEY'),
        'issue_id': ','.join(map(str, issues.keys())),
        'status_id': '*',
    }
    url = urljoin(os.getenv('REDMINE_URL'), 'issues.json')
    field_name = os.getenv('REDMINE_CUSTOM_FIELD_NAME', 'client')
    data = get_redmine_data(url, params)
    page = data['total_count'] // data['limit'] + 1
    clients = {}
    for p in range(1, page + 1):
        if p > 1:
            params['offset'] = data['limit'] * p
            data = get_redmine_data(url, params)
        for r in data['issues']:
            client = list(filter(lambda x: x['name'] == field_name, r['custom_fields']))[0] or None
            if not client:
                continue
            hours = sum(issues[r['id']])
            clients[client['value']] = hours if client['value'] not in clients else sum([clients[client['value']], hours])

    # 合計作業時間の降順でソート
    return sorted(clients.items(), key=lambda x: -x[1])


def get_issue_time_entries(start: datetime.datetime, end: datetime.datetime):
    """ Issue 単位の作業時間を取得 """
    params = {
        'key': os.getenv('REDMINE_KEY'),
        'project': os.getenv('REDMINE_PROJECT'),
        'from': start.strftime('%Y-%m-%d'),
        'to': end.strftime('%Y-%m-%d')
    }
    url = urljoin(os.getenv('REDMINE_URL'), 'time_entries.json')
    data = get_redmine_data(url, params)
    page = data['total_count'] // data['limit'] + 1
    issues = {}
    for p in range(1, page + 1):
        if p > 1:
            params['offset'] = data['limit'] * (p - 1)
            data = get_redmine_data(url, params)
        for r in data['time_entries']:
            issue_id = r['issue']['id'] or None
            if not issue_id:
                continue
            if issue_id not in issues:
                issues[issue_id] = [r['hours']]
            else:
                issues[issue_id].append(r['hours'])

    return issues


def send_notification(start, end, today, message):
    """ 通知を送る """
    dest = os.getenv('SEND_NOTIFYCATION')
    if dest == 'sns':
        period = '集計期間: {start}～{end}'.format(
            start=start.strftime('%Y-%m-%d'),
            end=end.strftime('%Y-%m-%d'),
        )
        sns_client = boto3.client('sns')
        sns_client.publish(
            TopicArn=os.getenv('SNS_TOPIC_ARN'),
            Subject='【自動通知】{date}_稼働時間集計'.format(date=today.strftime('%Y%m%d')),
            Message='{period}\n{message}'.format(period=period, message=message),
        )
    elif dest == 'chatwork':
        requests.post('https://api.chatwork.com/v2/rooms/{room_id}/messages'.format(
            room_id=os.getenv('CHATWORK_ROOM_ID'),
        ), headers={
            'X-ChatworkToken': os.getenv('CHATWORK_API_TOKEN'),
        }, data={
            'body': '[info][title]{subject}[/title]{message}[/info]'.format(
                subject='稼働時間集計 ({start}-{end})'.format(start=start.strftime('%Y-%m-%d'), end=end.strftime('%Y-%m-%d')),
                message=message,
            ),
        })
    else:
        app.log.info(message)


@app.schedule(Cron(0, 0, '?', '*', 'MON', '*'))
def lambda_handler(event, context={}):
    # 集計開始日・終了日を決定
    today = datetime.datetime.today()
    base_date = today
    if os.getenv('BASE_DATE'):
        base_date = datetime.datetime.strptime(os.getenv('BASE_DATE'), '%Y-%m-%d')

    # 基準日前日から BACK_DATE に設定した日数戻った日付までの期間を集計対象とする
    end = base_date - datetime.timedelta(days=1)
    start = end - datetime.timedelta(days=int(os.getenv('BACK_DATE', 6)))

    issues = get_issue_time_entries(start, end)
    clients = get_client_time_entries(issues)

    messages = []
    total = 0
    for name, hours in clients:
        messages.append('{name}: {hours}h'.format(name=name, hours=hours))
        total += hours

    messages.append('合計: {hours}h'.format(hours=total))

    message = '\n'.join(messages)

    send_notification(start, end, today, message)

    return clients
