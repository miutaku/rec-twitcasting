# TwitCasting Recording Application

This application is a job that periodically executes an API based on the userdata information in the DB for the TwitCasting recording application.

This application was created because Docker Compose does not provide a feature equivalent to cronjob.

# Flow

1. Retrieve the TwitCasting usernames from the DB where `recording_state` is `false`.
2. Request api-rec-twitcasting to start recording the broadcasts of the retrieved users.
