"""
    Local testing script, usage:
        source ../.env && python main.py <root api from server startup>
"""
import requests

from os import getenv
from sys import argv
import time

host = "http://" + getenv("HOST").rstrip("/")

root_headers = {"Api-Key": argv[1]}

def create_key() -> str:
    r = requests.post(
        host + "/key",
        headers=root_headers
    )
    print(r.content)

    return r.json()["api_key"]


headers = {
    "Api-Key": create_key()
}


data = {"link": "https://www.worldometers.info/languages/how-many-letters-alphabet/?a=" + f"{time.time()}"}


def print_link_response():
    r = requests.post(
        host + "/link",
        headers=headers,
        json=data
    )

    print(r.content)

def delete_key():
    r = requests.delete(
        host + "/key",
        headers=root_headers,
        json={
            "api_key": headers["Api-Key"],
            "deleted_at": int(time.time())
        }
    )

print_link_response()
print_link_response()

delete_key()

time.sleep(2)
print_link_response()

