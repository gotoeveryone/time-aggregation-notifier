import json
from app import lambda_handler

dummy_event = {
    'account': 'admin',
    'detail': {},
    'detail-type': 'Scheduled Event',
    'id': 'dummy',
    'region': 'ap-northeast-1',
    'resources': [],
    'source': 'aws.events',
    'time': '2019-06-24T01:23:45Z',
    'version': '1.0'
}

# Run at local
if __name__ == '__main__':
    print(json.dumps(lambda_handler(dummy_event, context={}), ensure_ascii=False, indent=2))
