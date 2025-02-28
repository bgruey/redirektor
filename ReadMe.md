# Redirektor

A Go server to redirect from a shortened url to another url. It has an API to add links programatically, using an API token. See `python_client/main.py` for an example.

## Generating API Tokens

No functionality!

If the server does not have one when it starts, it creates one, and logs a single API token from the database.

Without any other auth, it didn't make sense to make an endpoint to create new tokens.
