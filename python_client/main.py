import requests

from os import getenv
api_key = getenv("API_KEY")
host = getenv("HOST").rstrip("/")

headers = {
    "API-Key": api_key
}

data = {"link": "https://www.worldometers.info/languages/how-many-letters-alphabet/"}

r = requests.post(
    host + "/update",
    headers=headers,
    json=data
)


print(r.content)
