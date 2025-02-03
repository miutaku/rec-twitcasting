# rec-twitcasting

This application is for recording live streams being broadcasted on twitcasting.tv.

# Supporting env

- OS: Linux
- arch: arm64,x64

# How to use (Quick Start)

Please install Docker (it might work with podman etc., but this is unverified).

## Preparation

- Please obtain a client ID and client secret from [API Key](#api-key) section. You'll need to specify them as `<YOUR_TWITCASTING_CLIENT_ID>` and `<TWITCASTING_CLIENT_SECRET>`.
- Specify the IP address or domain name of the server machine where this application will run as `<BACKEND_SERVER>`.
  - If you want to run both the frontend (web server application) and backend (API server application) on the same server machine (all-in-one).

## Do it !

```shell
$ cp .env_sample .env
$ sed -i 's/<YOUR_TWITCASTING_CLIENT_ID>/__YOUR_TWITCASTING_CLIENT_ID__/g'
$ sed -i 's/<TWITCASTING_CLIENT_SECRET>/__TWITCASTING_CLIENT_SECRET__/g'
$ sed -i 's/<BACKEND_SERVER>/__YOUR_SERVER_IP_OR_FQDN__/g'
$ docker compose up -d
```

access to `http://__YOUR_SERVER_IP_OR_FQDN__:3000`

## Hands-on Demo

![Hands-on Demo](hands_on.gif)

# API key

You can get the API key from the official page [here](https://twitcasting.tv/developerapp.php).

# Diagram

![Docker Compose Diagram](docker-compose.drawio.png)

`k8s.drawio` is a conceptual diagram for Kubernetes support that we plan to add in the future.

# License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.