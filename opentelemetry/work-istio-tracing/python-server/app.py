import os, time, random
from flask import Flask, request
from database import get_population_from_db
from api_client import send_request_to_backend
from logger import setup_logger

from opentelemetry import trace
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from opentelemetry.propagators.b3 import B3MultiFormat

app = Flask(__name__)
logger = setup_logger()

@app.route('/', methods=['GET'])
def main():
    # Get Data.
    pref = request.args.get('pref')
    logger.info(f"リクエスト受信: {pref}")  
    
    # ヘッダーからキャリアを作成
    headers = dict(request.headers)
    logger.info(f"ヘッダー: {headers}")
    carrier = {key.lower(): value for key, value in headers.items()}
    ctx = trace.get_current_span()

    # W3C Trace Contextが存在する場合
    if 'traceparent' in carrier:
        logger.info("プロパゲーターは W3C Trace Context やで！")
        ctx = TraceContextTextMapPropagator().extract(carrier=carrier)
        logger.info("W3C Trace Context 時の context は、 {}".format(ctx))
    # B3ヘッダーが存在する場合
    elif 'x-b3-traceid' in carrier:
        logger.info("プロパゲーターは B3 Header やで！")
        b3_format = B3MultiFormat()
        ctx = b3_format.extract(carrier=carrier)
        logger.info("B3 時の context は、 {}".format(ctx))

    # Query DB ( MySQL ).
    population = get_population_from_db(pref)
    
    # Request
    fruit = send_request_to_backend(ctx)
    if population is None: return "Connection Error"

    return population + " 人ですがみんな " + fruit + " が好きな"

if __name__ == '__main__':
    host, port = os.getenv('CNDT_EASTERN_API_HOST', '0.0.0.0'), os.getenv('CNDT_EASTERN_API_PORT', 8089)
    app.run(host=host, port=port)