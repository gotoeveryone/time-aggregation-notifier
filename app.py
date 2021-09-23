import datetime
import os
from chalice import Chalice, Cron
from chalicelib.slack import SlackClient
from chalicelib.ssm import SSMClient
from chalicelib.redmine import RedmineClient

app = Chalice(app_name='time-aggregation-notifier')
app.debug = True if os.getenv('DEBUG', False) else False


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

    ssm_client = SSMClient()

    redmine_client = RedmineClient(ssm_client=ssm_client)
    issues = redmine_client.get_issue_time_entries(start, end)
    clients = redmine_client.get_client_time_entries(issues)

    messages = []
    total = 0
    for name, hours in clients:
        messages.append('{name}: {hours}h'.format(name=name, hours=hours))
        total += hours

    messages.append('合計: {hours}h'.format(hours=total))

    message = '\n'.join(messages)

    client = SlackClient(ssm_client=ssm_client, debug=app.debug)
    res = client.send(start, end, message)

    return {
        'start': start.strftime('%Y-%m-%d'),
        'end': end.strftime('%Y-%m-%d'),
        'send_status': res.status_code,
    }
