from datetime import datetime
import os
from urllib.parse import urljoin

import requests
from chalicelib.ssm import SSMClient


class RedmineClient:
    def __init__(self, ssm_client: SSMClient):
        self._url = ssm_client.get_parameter('redmine_url')
        self._key = ssm_client.get_parameter('redmine_key')
        self._project = os.getenv('REDMINE_PROJECT', 'job')
        self._custom_field_name = os.getenv('REDMINE_CUSTOM_FIELD_NAME', 'client')

    def get_client_time_entries(self, issues: dict):
        """ Issue 一覧から該当するカスタムフィールド単位の作業時間を取得 """
        params = {
            'issue_id': ','.join(map(str, issues.keys())),
            'status_id': '*',
        }
        url = urljoin(self._url, 'issues.json')
        data = self.get_data(url, params)
        page = data['total_count'] // data['limit'] + 1
        clients = {}
        for p in range(1, page + 1):
            if p > 1:
                params['offset'] = data['limit'] * p
                data = self.get_data(url, params)
            for r in data['issues']:
                client = list(filter(lambda x: x['name'] == self._custom_field_name, r['custom_fields']))[0] or None
                if not client:
                    continue
                hours = sum(issues[r['id']])
                clients[client['value']] = hours if client['value'] not in clients else sum([clients[client['value']], hours])

        # 合計作業時間の降順でソート
        return sorted(clients.items(), key=lambda x: -x[1])

    def get_issue_time_entries(self, start: datetime, end: datetime):
        """ Issue 単位の作業時間を取得 """
        params = {
            'project': self._project,
            'from': start.strftime('%Y-%m-%d'),
            'to': end.strftime('%Y-%m-%d')
        }
        url = urljoin(self._url, 'time_entries.json')
        data = self.get_data(url, params)
        page = data['total_count'] // data['limit'] + 1
        issues = {}
        for p in range(1, page + 1):
            if p > 1:
                params['offset'] = data['limit'] * (p - 1)
                data = self.get_data(url, params)
            for r in data['time_entries']:
                issue_id = r['issue']['id'] or None
                if not issue_id:
                    continue
                if issue_id not in issues:
                    issues[issue_id] = [r['hours']]
                else:
                    issues[issue_id].append(r['hours'])

        return issues

    def get_data(self, url: str, params: dict):
        """ Redmine の API から JSON 形式でデータを取得 """
        params.update({'key': self._key})
        res = requests.get(url, params)
        return res.json()
