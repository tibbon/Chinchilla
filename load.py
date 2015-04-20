from concurrent.futures import ThreadPoolExecutor
from requests_futures.sessions import FuturesSession
import time

session = FuturesSession(executor=ThreadPoolExecutor(max_workers=1000))
# first request is started in background

for x in xrange(1,100):
    session.get('http://localhost:8080/api/1/hello')
    session.get('http://localhost:8080/api/2/hello')
    session.get('http://localhost:8080/api/3/hello')
    time.sleep(3)


