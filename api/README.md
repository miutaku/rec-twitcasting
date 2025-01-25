# TwitCasting Recording API

This API provides a way to check if a TwitCasting user is live streaming and record their stream using `ffmpeg`. The stream is saved in a specified output directory.

## Table of Contents

- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [Endpoints](#endpoints)
- [Recording Output](#recording-output)
- [How to Run](#how-to-run)
- [License](#license)

## Setup

1. Clone this repository and go to this directory.

2. Install dependencies:

Ensure `ffmpeg` is installed on your system. If not, install it:

```bash
# For macOS
brew install ffmpeg
# For Debian based
sudo apt update
sudo apt install ffmpeg
# For RHEL based
sudo dnf install ffmpeg
```

3. Set up environment variables in `.env` or export them:

```bash
export TWITCASTING_CLIENT_ID=<your_client_id>
export TWITCASTING_CLIENT_SECRET=<your_client_secret>
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=user
export DB_PASSWORD=password
export DB_NAME=dbname
export DB_TABLE_NAME=tablename
export OUTPUT_DIR=./recorded
export LOG_LEVEL=debug
```

4. Run the server:

```bash
go run main.go
```

## Environment Variables

| Variable                | Description                                                                 | Default Value       |
|-------------------------|-----------------------------------------------------------------------------|---------------------|
| `TWITCASTING_CLIENT_ID`| Your TwitCasting Client ID.                                                 | None (required)     |
| `TWITCASTING_CLIENT_SECRET`| Your TwitCasting Client Secret.                                         | None (required)     |
| `DB_HOST`               | Database host.                                                            | `localhost`         |
| `DB_PORT`               | Database port.                                                            | `5432`              |
| `DB_USER`               | Database user.                                                            | `user`              |
| `DB_PASSWORD`           | Database password.                                                        | `password`          |
| `DB_NAME`               | Database name.                                                            | `dbname`            |
| `DB_TABLE_NAME`         | Database table name.                                                      | `tablename`         |
| `OUTPUT_DIR`            | Directory to save recorded videos.                                        | `./recorded`        |
| `LOG_LEVEL`             | Set to `debug` to see detailed logs.                                      | None                |

## Endpoints

### `GET /check-live`

Checks if a TwitCasting user is live streaming and records their stream if live.

#### Request

- **URL Query Parameters**:
  - `username` (string): The TwitCasting username to check and record.

#### Example

```bash
curl "http://localhost:8080/check-live?username=<twitcasting_username>"
```

#### Responses

- **200 OK**:
  - User is not live streaming:
    ```
    User is not live streaming.
    ```
  - Recording finished:
    ```
    Recording finished. Saved as: <output_file_path>
    ```
- **400 Bad Request**:
  - Missing `username` parameter:
    ```
    username parameter is required
    ```
- **500 Internal Server Error**:
  - Errors while checking stream or updating recording state:
    ```
    Failed to get current live information: <error_details>
    ```

## Recording Output

Recorded streams are saved in the directory specified by the `OUTPUT_DIR` environment variable. The directory structure is as follows:

```
recorded/
  └── <username>/
      └── <YYYY-MM-DD>/
          └── <HH-MM>_<title>.mp4
```

- `<username>`: The TwitCasting username.
- `<YYYY-MM-DD>`: The date of the recording.
- `<HH-MM>`: The time of the recording.
- `<title>`: The sanitized title of the live stream.

## How to Run

1. Ensure all environment variables are correctly set.
2. Start the server:

```bash
go run main.go
```

3. Access the endpoint to check and record a live stream:

```bash
curl "http://localhost:8080/check-live?username=<twitcasting_username>"
```
