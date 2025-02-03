# rec-twitcasting

このアプリケーションは、twitcasting.tvで配信されているライブストリームを録画するためのものです。

[English version here](README.md)

# 対応環境

- OS: Linux
- アーキテクチャ: arm64,x64

# 使用方法（クイックスタート）

Dockerをインストールしておいてください（podmanなども動作する可能性がありますが、未検証です）。
[公式インストールガイド](https://docs.docker.jp/linux/step_one.html)
[公式インストールスクリプト(英語)](https://github.com/docker/docker-install)

## 準備

- [APIキー](#apiキー)セクションからクライアントIDとクライアントシークレットを取得してください。それぞれ`<YOUR_TWITCASTING_CLIENT_ID>`と`<TWITCASTING_CLIENT_SECRET>`として指定する必要があります。
- このアプリケーションを実行するサーバーマシンのIPアドレスまたはドメイン名を`<BACKEND_SERVER>`として指定してください。
    - フロントエンド（Webサーバーアプリケーション）とバックエンド（APIサーバーアプリケーション）を同じサーバーマシンで実行する場合（オールインワン）。

## 実行手順

```shell
$ cp .env_sample .env
$ sed -i 's/<YOUR_TWITCASTING_CLIENT_ID>/__YOUR_TWITCASTING_CLIENT_ID__/g' .env
$ sed -i 's/<TWITCASTING_CLIENT_SECRET>/__TWITCASTING_CLIENT_SECRET__/g' .env
$ sed -i 's/<BACKEND_SERVER>/__YOUR_SERVER_IP_OR_FQDN__/g' .env
$ docker compose up -d
```

`http://__YOUR_SERVER_IP_OR_FQDN__:3000`にアクセスしてください。

## ハンズオン

![ハンズオン](hands_on.gif)

# APIキー

APIキーは[twitcasting.tv 公式ページ](https://twitcasting.tv/developerapp.php)から取得できます。

# 構成図

![Docker Compose 構成図](docker-compose.drawio.png)

`k8s.drawio`は将来追加予定のKubernetesサポートの概念図です。

# ライセンス

このプロジェクトはMITライセンスの下で提供されています。

詳細は[LICENSE](LICENSE)ファイルを参照してください。