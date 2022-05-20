# Fichi's URL shortener

Too much disposable [hitechs](https://www.webos-hightechs.co.il) income and this is what you get.

The system is is comprised of 2 parts

* Web server for redirecting and API operations
* Firestore NoSQL database for storing the mappings

## Terraform deployment

Terraform is used to deploy everything from scratch

Here's a table of the Terraform variables:
| Variable               | Type    | Default               | Description                                                                                                         |
|------------------------|---------|-----------------------|---------------------------------------------------------------------------------------------------------------------|
| `container_group_name` | String  | `"fichis-cont-group"` | Name for container group in which the containers will run.                                                          |
| `storage_account_name` | String  | `"fichisfiles"`       | Name of the storage account for Redis persistence dumps.                                                            |
| `file_share_name`      | String  | `"redisfs"`           | Name of file share in which Redis files will reside.                                                                |
| `resource_group_name`  | String  | `"fichis-app-rg"`     | Name of Resource Group in which to create everything.                                                               |
| `azure_region`         | String  | `"westeurope"`        | Azure region in which to create all of the resources.                                                               |
| `tls_enabled`          | Boolean | `false`               | Tells the web server whether to listen on HTTP or HTTPS.                                                            |
| `http_port`            | Number  | `80`                  | HTTP port to listen on.                                                                                             |
| `https_port`           | Number  | `443`                 | HTTPS port to listen on.                                                                                            |
| `redis_host`           | String  | `"localhost"`         | Hostname of Redis server, if both containers run in the same container group,   they can communicate via localhost. |
| `redis_port`           | Number  | `6379`                | Port on which the Redis server listens.                                                                             |
| `certificate_file`     | String  | N/A                   | Path to certificate file to pass to the web server if listening on HTTPS.                                           |
| `key_file`             | String  | N/A                   | Path to key file to pass to the web server if listening on HTTPS.                                                   |

## App documentation

The app itself looks for a few enviroment variables to determine how and what it should do.

| Variable          | Description                                                                                                    |
|-------------------|----------------------------------------------------------------------------------------------------------------|
| `FICHIS_DOMAIN_NAME` | The domain name for the URL shortener (e.g. fich.is) |
| `FICHIS_PROBE_PATH` | Path in which the health probe will listen |
| `FICHIS_HTTP_PORT`  | HTTP port to listen on.                                                                                        |
| `FICHIS_HTTPS_PORT` | HTTPS port to listen on.                                                                                       |
| `FICHIS_TLS_ON`     | If this value is equal to "yes",   the server will listen on HTTPS. Otherwise, the server will listen on HTTP. |
| `FICHIS_GOOGLE_APPLICATION_CREDENTIALS_FILE_PATH` | Where the Google Service Account credentials file resides |
| `FICHIS_GOOGLE_PROJECT_ID` | The Project ID of the Google Project in which Firestore was created |
| `FICHIS_API_VALIDATION_ON` | If set to "yes", an API_KEY will need to be passed to all API operations |
| `FICHIS_API_KEY` | The API key if API validation is enabled |

## API Endpoints

All API calls require you to pass the `api_key` parameter as well

| Endpoint | Method | Description | Parameters |
|----------|--------|-------------|------------|
`/api/create/ShortenedLink` | `POST` | Used to create shortened links | `url` - The URL to shorten<br />`id` - The id to assign the shortened URL |
`/api/create/AutoShortenedLink` | `POST` | Used to create shortened links without specifying an `id` | `url` - the URL to shorten |
`/api/get/ShortenedLink` | `GET` | Used to retrieve the URL of a shortened link by it's `id` | `id` - The ID of the link to retrieve |
`/api/delete/ShortenedLink` | `DELETE` | Used to delete a shortened URL by its `id` |
