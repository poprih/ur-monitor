# UR Monitor

A monitoring system built with Go that integrates with LINE Messaging API.

## Overview

UR Monitor is a server application that provides monitoring capabilities and LINE bot integration. It's built using Go and uses PostgreSQL (Neon) for data storage, deployed on Vercel.

## Features

- LINE Bot integration for notifications and interactions
- PostgreSQL database powered by Neon
- RESTful API endpoints
- Vercel deployment

## Prerequisites

- Go 1.22.3 or later
- PostgreSQL (Neon) database
- LINE Messaging API credentials
- Vercel account (for deployment)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/poprih/ur-monitor.git
cd ur-monitor
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
export LINE_CHANNEL_ACCESS_TOKEN=your_channel_token
export DATABASE_URL=your_neon_postgres_url
```

## Development

This project is designed to be deployed on Vercel. For local development, you can use the Vercel CLI to run the application locally:

1. Install Vercel CLI:

```bash
npm i -g vercel
```

2. Start local development server:

```bash
vercel dev
```

## Deployment

The API is deployed on Vercel. To deploy:

1. Install Vercel CLI:

```bash
npm i -g vercel
```

2. Deploy:

```bash
vercel
```

## Project Structure

```
.
├── api/         # API handlers and routes
├── cmd/         # Command line tools and scripts
├── db/          # Database related code
├── lib/         # Library and utility functions
├── pkg/         # Internal packages
└── README.md    # This file
```

## License

MIT License

---

# UR Monitor (日本語)

LINE Messaging API と統合された Go で構築された監視システム。

## 概要

UR Monitor は、Go で構築され、PostgreSQL（Neon）をデータストレージとして使用するサーバーアプリケーションです。LINE ボットとの統合機能を提供し、Vercel にデプロイされています。

## 機能

- 通知やインタラクションのための LINE ボット統合
- Neon による PostgreSQL データベース
- RESTful API エンドポイント
- Vercel デプロイメント

## 必要条件

- Go 1.22.3 以上
- PostgreSQL（Neon）データベース
- LINE Messaging API の認証情報
- Vercel アカウント（デプロイ用）

## インストール

1. リポジトリをクローン:

```bash
git clone https://github.com/poprih/ur-monitor.git
cd ur-monitor
```

2. 依存関係をインストール:

```bash
go mod download
```

3. 環境変数を設定:

```bash
export LINE_CHANNEL_ACCESS_TOKEN=your_channel_token
export DATABASE_URL=your_neon_postgres_url
```

## 開発

このプロジェクトは Vercel にデプロイするように設計されています。ローカル開発には Vercel CLI を使用してアプリケーションを実行できます：

1. Vercel CLI をインストール:

```bash
npm i -g vercel
```

2. ローカル開発サーバーを起動:

```bash
vercel dev
```

## デプロイメント

API は Vercel にデプロイされています。デプロイするには:

1. Vercel CLI をインストール:

```bash
npm i -g vercel
```

2. デプロイ:

```bash
vercel
```

## プロジェクト構造

```
.
├── api/         # APIハンドラーとルート
├── cmd/         # コマンドラインツールとスクリプト
├── db/          # データベース関連のコード
├── lib/         # ライブラリとユーティリティ関数
├── pkg/         # 内部パッケージ
└── README.md    # このファイル
```

## ライセンス

MIT ライセンス
