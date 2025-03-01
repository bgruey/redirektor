"""
Local testing script, usage:
    source ../.env && python main.py <root api from server startup>
"""

import base64
import time
from os import getenv
from sys import argv

import requests

host = "http://" + getenv("HOST").rstrip("/")

root_headers = {"Api-Key": argv[1]}


def create_key() -> str:
    r = requests.post(host + "/key", headers=root_headers)
    print(r.content)

    return r.json()["api_key"]


headers = {"Api-Key": create_key()}


data = {
    "link": "https://www.worldometers.info/languages/how-many-letters-alphabet/?a="
    + f"{time.time()}"
}


def print_link_response():
    r = requests.post(host + "/link", headers=headers, json=data)
    if r.status_code != 201:
        print(r.content)
        return

    r_data = r.json()
    short_url = r_data["short_url"]
    image_filename = short_url.replace("/", "_") + ".png"
    print(f"Url {short_url} qr code saved to {image_filename}")

    with open(image_filename, "wb") as f:
        f.write(base64.b64decode(r_data["qrcode_b64"]))


def delete_key():
    r = requests.delete(
        host + "/key",
        headers=root_headers,
        json={"api_key": headers["Api-Key"], "deleted_at": int(time.time())},
    )
    print(r.content)


print_link_response()
print_link_response()

delete_key()

time.sleep(1)
print_link_response()
