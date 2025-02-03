# rec-twitcasting

This application is for recording live streams being broadcasted on twitcasting.tv.

[日本語版はこちら](README_ja.md)

# Important Information

**Please read this carefully.**

## Please Limit to Personal Use Only

Obviously, you must not repost recordings without the streamer's permission.

If any misuse is observed, this repository may be made private.

## Use at Your Own Risk

Since this application is developed under the MIT License, please use it at your own risk.
The developer assumes no responsibility for any problems that may occur due to the use of this application.

## When Creating Issues

Bug reports are welcome!

However, please follow these guidelines:

- Use the issue template
  - Bug report
  - Feature request
- Make respectful reports
- Check if the issue is already known (including previously closed issues)

These rules must be followed.
Issues may be closed if they do not comply.

# Supporting env

- OS: Linux
- arch: arm64,x64

# How to use (Quick Start)

Please install Docker (it might work with podman etc., but this is unverified).
[Official install script](https://github.com/docker/docker-install)

## Preparation

- Please obtain a client ID and client secret from [API Key](#api-key) section. You'll need to specify them as `<YOUR_TWITCASTING_CLIENT_ID>` and `<TWITCASTING_CLIENT_SECRET>`.
- Specify the IP address or domain name of the server machine where this application will run as `<BACKEND_SERVER>`.
  - If you want to run both the frontend (web server application) and backend (API server application) on the same server machine (all-in-one).

## Do it !

```shell
$ cp .env_sample .env
$ sed -i 's/<YOUR_TWITCASTING_CLIENT_ID>/__YOUR_TWITCASTING_CLIENT_ID__/g' .env
$ sed -i 's/<TWITCASTING_CLIENT_SECRET>/__TWITCASTING_CLIENT_SECRET__/g' .env
$ sed -i 's/<BACKEND_SERVER>/__YOUR_SERVER_IP_OR_FQDN__/g' .env
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