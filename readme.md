# Fichi's URL shortener

Too much disposable [hitechs](https://www.webos-hightechs.co.il) income and this is what you get.


The system is built from 2 containers:
* HTTPS server
* Redis


# Terraform deployment
Terraform is used to deploy both of them in Azure Container Instances.

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



# App documentation

The app itself looks for a few enviroment variables to determine how and what it should do.

| Variable          | Description                                                                                                    |
|-------------------|----------------------------------------------------------------------------------------------------------------|
| `FICHIS_HTTP_PORT`  | HTTP port to listen on.                                                                                        |
| `FICHIS_HTTPS_PORT` | HTTPS port to listen on.                                                                                       |
| `FICHIS_TLS_ON`     | If this value is equal to "yes",   the server will listen on HTTPS. Otherwise, the server will listen on HTTP. |
| `FICHIS_REDIS_HOST` | Hostname on which the Redis server resides.                                                                    |
| `FICHIS_REDIS_PORT` | Port on which the Redis server listens.                                                                        |

## API Endpoints

All API calls require you to pass the `api_key` parameter as well

| Endpoint | Method | Description | Parameters |
------------------------------------------------
`/api/create/ShortenedLink` | `POST` | Used to create shortened links | `url` -
The URL to shorten\n`id` - The id to assign the shortened URL ("fich.is/{id}")

