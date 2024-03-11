import logging
import json
import os
import requests
import time

import boto3

logger = logging.getLogger()
layer_name = os.environ.get("LAYER_NAME_FOR_SDK")

# seems to improve response time if here
# also generates 2nd trace
client = boto3.client('lambda')
if layer_name:
    layer_vers = client.get_layer_version(
        LayerName=layer_name,
        VersionNumber=1
    )


def lambda_handler(event, context):
    if event.get("exc"):
        raise Exception("Test exc was raised")

    if event.get("timeout"):
        time.sleep(10)

    # here for predictable response_time metrics
    time.sleep(1)

    logger.warning("First log inside handler")

    # api gateway to instrumented java
    java_api = os.environ.get("JAVA_API_RAISES")
    req_headers = {}
    resp_headers = {}
    if java_api:
        resp = requests.get(java_api)
        req_headers = dict(resp.request.headers)
        resp_headers = dict(resp.headers)

    if layer_name:
        layer_vers = client.get_layer_version(
            LayerName=layer_name,
            VersionNumber=1
        )
    else:
        layer_vers = "123"

    logger.warning("Second log inside handler")

    path = os.environ.get("SETTINGS_FILEPATH")
    if path and os.path.exists(path):
        f = open(path, "r")
        return {
            'statusCode': 200,
            'body': {
                'SettingsFile': f.read(),
                'ApmLayerVersion': layer_vers,
                'RequestHeaders': req_headers,
                'ResponseHeaders': resp_headers,
            }
        }
    else:
        return {
            'statusCode': 200,
            'body': {
                'SettingsFile': json.dumps(f"{path} does not exist!"),
                'ApmLayerVersion': layer_vers,
                'RequestHeaders': req_headers,
                'ResponseHeaders': resp_headers,
            }
        }
