import boto3


class SSMClient:
    def __init__(self):
        self._ssm = boto3.client('ssm')

    def get_parameter(self, key: str) -> str:
        """
        Get value from Parameter Store
        """
        response = self._ssm.get_parameter(
            Name=key,
            WithDecryption=True
        )
        return response['Parameter']['Value']
