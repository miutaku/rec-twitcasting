# rec-twitcasting

このアプリケーションは、twitcasting.tvで配信されているライブストリームを録画するためのものです。

[English version here](README.md)

# 重要事項

**必ず読んでください。**

## あくまで個人利用にとどめてください

当たり前ですが配信者が録画を公開しないのに**勝手に転載したりしてはいけません。**

悪用が見られた場合は**本レポジトリを非公開(private)レポジトリ**にする可能性があります。

## 自己責任でお願いします

そもそもこのアプリケーションはMITライセンスの下で開発されているため、自己責任で使ってください。
本アプリケーションの利用によって起こるいかなる問題においても、開発者は一切の責任を負わないものとします。

## issueを書く場合

バグ報告は大歓迎です！

ただし、

- 原則としてissueテンプレートを使う。
  - バグ報告(Bug report)
  - 機能追加リクエスト(Feature Request)
- リスペクトのない報告をしない。
- すでにcloseされた他のissueを含め、既知の問題か確認する。

これらは守ってください。
守られていない場合はcloseする可能性があります。

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
$ sed -i 's/__YOUR_TWITCASTING_CLIENT_ID__/<YOUR_TWITCASTING_CLIENT_ID>/g' .env
$ sed -i 's/__YOUR_TWITCASTING_CLIENT_SECRET__/<TWITCASTING_CLIENT_SECRET>/g' .env
$ sed -i 's/__YOUR_SERVER_IP_OR_FQDN__/<BACKEND_SERVER>/g' .env
$ docker compose up -d
```

`http://<BACKEND_SERVER>:3000`にアクセスしてください。

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
