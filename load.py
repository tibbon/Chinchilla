import requests


def main():
    Url = "local"

    for x in xrange(1, 25):
        obj = requests.get("http://localhost:8080/api/1/hello")
        print obj.text


if __name__ == '__main__':
        main()
