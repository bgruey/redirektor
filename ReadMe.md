# Redirektor

A Go server to redirect from a shortened url to another url. It has an API to add links programatically, using an API Key. See `python_client/main.py` for an example.

## API Keys

At startup the server creates a root API Token than can be used to create additional API Keys. It prints this token to the logs.

To authenticate, add the following header: `"Api-Key": "<api key: string>"`. The field can be configured in `server/api/redirect/consts.go`.

## API Endpoints

### Root `/<link hash: string>`

This endpoint is the shortened url that will redirect to a longer url designated by the `<link hash>`.

### Key `/key`

This endpoint requires an API Key with root privileges.

#### POST

This method creates a new key (without root privileges). No post data is required.

##### Response:
```json
{
    "api_key": <api key: string>
}
```

#### DELETE

This method deletes an API key at the specified time. If `deleted_at` is not supplied, it is deleted at `0` seconds from the epoch.

##### Body
```json
{
    "api_key": <api key: string>,
    "deleted_at": <unix timestamp: int64>
}
```

##### Response
```json
{
    "api_key": <api key: string>,
    "deleted_at": <unix timestamp: int64>
}
```


### Link `/link`

This endpoint requires an API key created by the `/key` endpoint, with a root API key being dissallowed.

#### POST

This method creates a link redirect. The action is idempotent, with subsequent calls for the same link returning the same url. 

The hash is defined as the Base64 representation via [RFC 4648.5](https://en.wikipedia.org/wiki/Base64#Variants_summary_table) of the SHA256 of the link supplied, using the minimum number of characters to make the hash unique in the database. This means the first hash will always be a single character.

##### Request
```json
{
    "link": <link text: string>
}
```

##### Response
```json
{
    "short_url":"<root url>/<hash: string>"
}
```

## Environment Variables

For local testing, the following `.env` file works with the supplied `docker-compose.yaml` file.
```bash
export POSTGRES_HOST=redirektor-db
export POSTGRES_PASSWORD=testingpw
export POSTGRES_USER=postgres
export POSTGRES_DB=postgres
export HOST="localhost:8080"
```
