import datetime
import os
from urllib.parse import urljoin
import boto3
from chalice import Chalice, Cron
import requests

app = Chalice(app_name='time-aggregation-notifier')
app.debug = True if os.environ.get('DEBUG', False) else False


def get_redmine_data(url: str, params: dict):
    """ Redmine の API から JSON 形式でデータを取得 """
    res = requests.get(url, params)
    return res.json()


def get_client_time_entries(issues: dict):
    """ Issue 一覧から該当するクライアント単位の作業時間を取得 """
    params = {
        'key': os.environ.get('REDMINE_KEY'),
        'issue_id': ','.join(map(str, issues.keys())),
        'status_id': '*',
    }
    url = urljoin(os.environ.get('REDMINE_URL'), 'issues.json')
    data = get_redmine_data(url, params)
    page = data['total_count'] // data['limit'] + 1
    clients = {}
    for p in range(1, page + 1):
        if p > 1:
            params['offset'] = data['limit'] * p
            data = get_redmine_data(url, params)
        for r in data['issues']:
            client = list(filter(lambda x: x['name'] == 'クライアント', r['custom_fields']))[0] or None
            if not client:
                continue
            hours = sum(issues[r['id']])
            clients[client['value']] = hours if client['value'] not in clients else sum([clients[client['value']], hours])

    # 合計作業時間の降順でソート
    return sorted(clients.items(), key=lambda x: -x[1])


def get_issue_time_entries(start: datetime.datetime, end: datetime.datetime):
    """ Issue 単位の作業時間を取得 """
    params = {
        'key': os.environ.get('REDMINE_KEY'),
        'project': os.environ.get('REDMINE_PROJECT'),
        'from': start.strftime('%Y-%m-%d'),
        'to': end.strftime('%Y-%m-%d')
    }
    url = urljoin(os.environ.get('REDMINE_URL'), 'time_entries.json')
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


@app.schedule(Cron(0, 20, '?', '*', 'SUN', '*'))
def lambda_handler(event, context={}):
    # 集計開始日・終了日を決定
    end = datetime.datetime.today().astimezone(datetime.timezone(datetime.timedelta(hours=9), name='JST'))
    if os.environ.get('START_DATE'):
        end = datetime.datetime.strptime(os.environ.get('START_DATE'), '%Y-%m-%d')
    start = end - datetime.timedelta(days=int(os.environ.get('BACK_DATE', 7)))

    issues = get_issue_time_entries(start, end)
    clients = get_client_time_entries(issues)

    messages = []
    for name, hours in clients:
        messages.append('{name}: {hours}h'.format(name=name, hours=hours))

    message = '\n'.join(messages)
    if not app.debug:
        # SNS へ連携
        sns_client = boto3.client('sns')
        sns_client.publish(
            TopicArn=os.environ.get('SNS_TOPIC_ARN'),
            Subject='集計結果',
            Message=message,
        )

    return clients
