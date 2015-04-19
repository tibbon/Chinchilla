from concurrent.futures import ThreadPoolExecutor
from requests_futures.sessions import FuturesSession

session = FuturesSession(executor=ThreadPoolExecutor(max_workers=1000))
# first request is started in background

for x in xrange(1,100):
    session.get('http://localhost:8080/api/1/hello')
# second requests is started immediately

# wait for the first request to complete, if it hasn't already
# response_one = future_one.result()
# print('response one status: {0}'.format(response_one.status_code))
# print(response_one.content)
# # wait for the second request to complete, if it hasn't already
# response_two = future_two.result()
# print('response two status: {0}'.format(response_two.status_code))
# print(response_two.content)