import json

import boto3

s3_client = boto3.client('s3')


def lambda_handler(event, context):
    # Permissions test: write object to S3
    bucket_name = "tammy-test-log-bucket"
    file_name = "my_log_from_lambda.log"
    file_content = f"ReqId {context.aws_request_id} | Logged from Lambda!"
    s3_client.put_object(Body=file_content, Bucket=bucket_name, Key=file_name)

    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }
