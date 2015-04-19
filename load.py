from concurrent.futures import ThreadPoolExecutor
from requests_futures.sessions import FuturesSession

session = FuturesSession(executor=ThreadPoolExecutor(max_workers=1000))
# first request is started in background

for x in xrange(1,1000):
    session.get('http://192.168.1.69:9000/api/1/hello')
