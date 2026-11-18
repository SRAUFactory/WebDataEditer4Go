# WebDataEditer4Go

WebAssembly (Go) と素の HTML / CSS / JavaScript で構築された、ブラウザ完結型のデータ編集・変換 Web アプリケーションです。
[SRAUFactory/web-data-editor](https://github.com/SRAUFactory/web-data-editor) をベースに、npm/フロントエンドフレームワーク依存を排除し、Go 標準ライブラリによる堅牢なパースと GitHub Pages への静的配信を実現しています。

---

## 目的

- **依存関係ゼロの持続可能な運用**:
  - npm などの外部パッケージマネージャーや頻繁なアップデート対応を不要にし、長期的なメンテナンスコストを最小化します。
- **堅牢なデータ処理**:
  - セル内改行や複雑なエスケープシーケンスを含む CSV/TSV を、Go の `encoding/csv` 等の標準ライブラリを用いて正確に処理します。
- **GitHub Pages での完全動作**:
  - バックエンドサーバーを必要とせず、すべての処理がクライアントブラウザ上の WebAssembly (WASM) で完結します。

---

## 主な機能

- **新規作成**: 指定した行数・列数 (1〜50) の空テーブル作成
- **ファイル読み込み**: CSV / TSV / JSON のパースおよびテーブル展開
- **テーブル編集**:
  - 各セルの直接編集
  - 行・列の挿入 (前後)、削除、連番採番、複製コピー
- **ファイル保存・変換**:
  - CSV / TSV / JSON / Markdown 形式でのエクスポート・ダウンロード
  - 改行コード (LF / CRLF / CR) の選択
  - Markdown のクリップボードコピー
- **組み込みマニュアル**: 操作方法の確認

---

## ディレクトリ構成

```
.
├── cmd/
│   └── wasm/           # Go WebAssembly エントリポイント (syscall/js)
├── internal/
│   ├── format/         # CSV / TSV / JSON / Markdown のフォーマッタ・パーサー
│   └── table/          # 二次元配列テーブルデータ構造・操作ロジック
├── web/                # 静的 Web アセット (HTML, CSS, JS, wasm_exec.js)
├── tools/              # ビルド・ユーティリティスクリプト (マニュアルHTML変換等)
├── testdata/           # 各種形式のテスト用サンプルデータ
└── docs/
    └── spec/           # 仕様駆動開発 (SDD) 仕様書
        ├── features.md          # 機能一覧・UI挙動仕様
        └── wasm_js_boundary.md  # Go WASM - JS 境界インターフェース仕様
```

---

## 開発手順

### 必要な環境
- Go 1.26 以上

### WASM バイナリのビルド
```bash
GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/wasm
```

### wasm_exec.js の準備
Go 公式ランタイムスクリプトを `web/` 配下に配置します。
```bash
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/
```

### ローカルサーバーの起動
ローカルで静的ファイルサーバーを起動して確認します (例: Python または Go)。
```bash
# Python の場合
python3 -m http.server 8080 --directory web

# ブラウザでアクセス
open http://localhost:8080
```

### テスト実行
```bash
go test ./...
```

---

## デプロイ手順 (GitHub Pages)

GitHub Actions を通じて `master` ブランチへのプッシュ時に自動デプロイされます。

1. **Build**:
   - Go 環境のセットアップ
   - `GOOS=js GOARCH=wasm` による `web/main.wasm` の生成
   - `$(go env GOROOT)/lib/wasm/wasm_exec.js` の `web/` 配下へのコピー
   - マニュアル Markdown から静的 HTML への事前ビルド
2. **Deploy**:
   - `actions/upload-pages-artifact` により `web/` をアップロード
   - `actions/deploy-pages` により GitHub Pages へ配信
