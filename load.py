import requests


def main():
    Url = "local"

    for x in xrange(1, 200):
        requests.get("http://localhost:8080/api/1/hello")


if __name__ == '__main__':
        main()
