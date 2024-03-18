import os
import requests
from logger import setup_logger

logger = setup_logger()

def send_request_to_backend():
    base_url = os.getenv('FRUIT_SERVER_ADDRESS', 'http://query-fruit.demo.svc.cluster.local:8080')
    url = f"{base_url}"
    logger.info(f'URL: {url}')
    
    try:  
        response = requests.get(url, verify=False)

        return response.text

    except requests.exceptions.ConnectionError as e:
        logger.warning(f"Connection error: {e}")
        
        return None