import os, time, random
from flask import Flask, request
from database import get_population_from_db
from api_client import send_request_to_backend
from logger import setup_logger

app = Flask(__name__)
logger = setup_logger()

@app.route('/', methods=['GET'])
def main():
    # Get Data.
    pref = request.args.get('pref')
    logger.info(f"リクエスト受信: {pref}")  
    headers = dict(request.headers)
    logger.info(f"ヘッダー: {headers}")

    # Query DB ( MySQL ).
    population = get_population_from_db(pref)
    
    # Request
    fruit = send_request_to_backend()
    if population is None: return "Connection Error"

    return population + " 人ですがみんな " + fruit + " が好きな"

if __name__ == '__main__':
    host, port = os.getenv('CNDT_EASTERN_API_HOST', '0.0.0.0'), os.getenv('CNDT_EASTERN_API_PORT', 8089)
    app.run(host=host, port=port)