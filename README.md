# rec-twitcasting-wine

This application is for recording live streams being broadcasted on twitcasting.tv.

# Supporting env

- OS: Linux
- arch: arm64,x64

# How to use

```shell
docker compose up -d
```

# API key

You can get the API key from the official page [here](https://twitcasting.tv/developerapp.php).

# Setting env

If you do not want to run it with Docker, set the environment variables.

```shell
export TWITCASTING_CLIENT_ID=<YOUR_ID>
export TWITCASTING_CLIENT_SECRET=<YOUR_SECRET>
export OUTPUT_DIR=<YOUR_RECORDING_DIR_PATH> # if not set, default parameter (./recorded)
```
