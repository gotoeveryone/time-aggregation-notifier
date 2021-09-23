from datetime import datetime
import json

import requests
from requests.models import Response
from chalicelib.ssm import SSMClient


class SlackClient:
    def __init__(self, ssm_client: SSMClient, debug=False):
        self._webhook_url = ssm_client.get_parameter(
            'slack_time_aggregation_notifier_webhook_url' if not debug else 'slack_test_webhook_url')

    def send(self, start: datetime, end: datetime, message: str) -> Response:
        """
        Send message to Slack
        """
        return requests.post(self._webhook_url, data=json.dumps({
            'text': '{subject}\n```{message}```'.format(
                subject='*Summary ({start}-{end})*'.format(
                    start=start.strftime('%Y-%m-%d'),
                    end=end.strftime('%Y-%m-%d'),
                ),
                message=message,
            ),
            'username': 'time-aggregation-notifier',
        }))
