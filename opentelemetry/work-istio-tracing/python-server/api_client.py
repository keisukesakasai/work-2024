import os
import requests
from logger import setup_logger
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from opentelemetry.baggage.propagation import W3CBaggagePropagator
from opentelemetry.propagators.b3 import B3MultiFormat

logger = setup_logger()

def send_request_to_backend(ctx):
    base_url = os.getenv('FRUIT_SERVER_ADDRESS', 'http://query-fruit.demo.svc.cluster.local:8080')
    url = f"{base_url}"
    logger.info(f'URL: {url}')
    
    try:
        headers = {}
        W3CBaggagePropagator().inject(headers, ctx)
        TraceContextTextMapPropagator().inject(headers, ctx)        
        b3_format = B3MultiFormat()
        b3_format.inject(headers, ctx)
        
        logger.info(f"送信時のヘッダー: {headers}")
        
        response = requests.get(url, headers=headers, verify=False)

        return response.text

    except requests.exceptions.ConnectionError as e:
        logger.warning(f"Connection error: {e}")
        
        return None