import os
import mysql.connector
from logger import setup_logger

logger = setup_logger()

# Get Population from MySQL.
def get_population_from_db(pref):
    config = {
        'user': os.getenv('DB_USER', 'root'),
        'password': os.getenv('DB_PASSWOR', 'password'),
        'host': os.getenv('DB_HOST', '127.0.0.1'),
        'database': os.getenv('DB_NAME', 'population'),
        'port': os.getenv('DB_PORT', 3306),
        'raise_on_warnings': True,
    }
    
    try: 
        cnx = mysql.connector.connect(**config)
        cursor = cnx.cursor()

        query = "SELECT population FROM population WHERE prefecture = %s"
        cursor.execute(query, (pref,))    
        
        result = cursor.fetchone()

        population = str(result[0]) if result else "Not Found"
        logger.info(f"DB からデータ取得: {population}")
    finally:
        cursor.close()
        cnx.close()
        
    return population