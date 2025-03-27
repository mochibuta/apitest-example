# apitest example

## 実行に必要なツール

- Go 1.24 or later
- Docker & Docker Compose
- mise

## 実行方法

1. リポジトリのクローン
```bash
git clone https://github.com/yourusername/apitest-example.git
cd apitest-example
```

2. 環境のセットアップ
```bash
mise install
```

3. APIサーバー開発時
```bash
# サーバー起動（with air）
mise run dev

# 稼働中サーバへのrunnによるAPIテスト実行
mise run runn
```

4. APIテスト実行
```bash
mise run apitest
```
