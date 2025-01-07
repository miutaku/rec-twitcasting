# rec-twitcasting-wine

This application is for recording live streams being broadcasted on twitcasting.tv.

# Supporting env

- OS: Linux
- arch: arm64,x64

# how to use

```shell
docker compose up -d
```

# api key

You can get the API key from the official page [here](https://twitcasting.tv/developerapp.php).

# setting env

If you do not want to run it with Docker, set the environment variables.

```shell
export TWITCASTING_CLIENT_ID=<YOUR_ID>
export TWITCASTING_CLIENT_SECRET=<YOUR_SECRET>
export OUTPUT_DIR=<YOUR_RECORDING_DIR_PATH> # if not set, default parameter (./recorded)
```
