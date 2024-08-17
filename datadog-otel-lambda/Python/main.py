import os
import json
import time
import random
import boto3
from flask import Flask, request, jsonify
from ddtrace import tracer, patch_all, patch
from ddtrace.propagation.http import HTTPPropagator

# Datadogのトレースを有効化
patch_all()
patch(flask=True)

app = Flask(__name__)

lambda_client = boto3.client('lambda', region_name='ap-northeast-1')

def create_data(length):
    span = tracer.trace("create_data")
    data = ''.join(random.choice('01') for _ in range(length))
    span.finish()
    return data

def invoke_lambda_function(function_name, data):
    span = tracer.trace(f"invoke_lambda_function: {function_name}")
    
    # OpenTelemetry のコンテキストを Datadog の形式に変換
    headers = {}
    propagator = HTTPPropagator()
    propagator.inject(span.context, headers)
    
    payload = {
        "bitstring": data,
        "carrier": json.dumps(headers)
    }

    response = lambda_client.invoke(
        FunctionName=function_name,
        InvocationType='RequestResponse',
        Payload=json.dumps(payload)
    )
    
    response_payload = json.loads(response['Payload'].read())
    span.finish()
    return response_payload

@app.route("/health", methods=["GET"])
def health_check():
    return "OK", 200

@app.route("/", methods=["GET"])
def handle_request():
    length_str = request.args.get("length")
    if not length_str:
        return "Invalid length parameter", 400
    
    try:
        length = int(length_str)
    except ValueError:
        return "Invalid length parameter", 400

    data = create_data(length)
    print("bitstring: ", data)

    function_name1 = os.getenv("FUNCTION_NAME_1")
    function_name2 = os.getenv("FUNCTION_NAME_2")

    response1 = invoke_lambda_function(function_name1, data)
    response2 = invoke_lambda_function(function_name2, data)

    response = {
        "response1": response1,
        "response2": response2,
    }

    return jsonify(response)

if __name__ == "__main__":
    port = int(os.getenv("SERVER_PORT", 8080))
    print("Server Port:", port)
    app.run(host="0.0.0.0", port=port)
