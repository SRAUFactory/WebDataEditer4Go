# WebDataEditer4Go 仕様書: Go (WASM) と JavaScript の境界設計

本仕様書は、Go WebAssembly (WASM) 側とブラウザ (JavaScript / DOM) 側の責務分担、状態保持モデル、および公開関数・データ送受信インターフェースを定義するものです。

---

## 1. 責務分担と状態管理方針

### 1.1. 状態管理の基本原則
- **状態 (State) の主所有権**: **JavaScript (ブラウザメモリ)** が保持する。
  - 理由:
    1. DOM レンダリング、ユーザー入力 (input イベント)、仮想スクロールや UI イベントとの結合度が高い。
    2. Go WASM ヒープ内に状態を持つと、ブラウザ再読込や WASM 初期化タイミングによる同期ズレやメモリリークの管理が複雑になる。
    3. Go WASM 側を「ステートレスなデータ変換・テーブル操作エンジン」として設計することで、純粋関数としてテストが容易になり、並行処理や保守性が格段に向上する。
- **Go (WASM) の役割**:
  - CSV / TSV / JSON / Markdown のパースおよびシリアライズ
  - セル内改行やエスケープ処理 (RFC 4180、Markdown仕様準拠)
  - テーブル操作 (行・列の挿入、削除、採番、複製) のロジック演算
- **JavaScript の役割**:
  - UI イベントのハンドリング (ボタン押下、入力変更、ファイル選択、ダイアログ表示)
  - DOM 構築と更新
  - File API (読み込み) および Blob / File System Access API (ダウンロード・クリップボード書き込み)

---

## 2. 公開関数インターフェース (WASM API)

Go WASM はグローバルオブジェクト `window.WebDataEditor` 配下に以下の関数を公開する。
データの受け渡しは、Go の `syscall/js` を介して、プリミティブ型および JSON 文字列、または JS 配列/オブジェクトとして行う。

```typescript
interface WebDataEditorAPI {
  /**
   * ファイル内容 (テキスト) をパースし、二次元文字配列を返す
   * @param content ファイル文字列
   * @param format "CSV" | "TSV" | "JSON"
   * @param lineEnding "LF" | "CRLF" | "CR"
   * @returns JSON 文字列: { success: true, data: string[][] } または { success: false, error: string }
   */
  parseData(content: string, format: string, lineEnding: string): string;

  /**
   * 二次元文字配列を指定形式のテキスト文字列にシリアライズする
   * @param rowsJson 二次元配列の JSON 文字列 (string[][])
   * @param format "CSV" | "TSV" | "JSON" | "Markdown"
   * @param lineEnding "LF" | "CRLF" | "CR"
   * @returns JSON 文字列: { success: true, content: string } または { success: false, error: string }
   */
  serializeData(rowsJson: string, format: string, lineEnding: string): string;

  /**
   * 指定したテーブル操作を適用した新しい二次元配列を返す
   * @param rowsJson 現在の二次元配列の JSON 文字列 (string[][])
   * @param action "insertRow" | "deleteRow" | "insertCol" | "deleteCol" | "assignRow" | "assignCol" | "copyRow" | "copyCol"
   * @param paramsJson 操作パラメータ JSON { index: number, position?: "before" | "after" }
   * @returns JSON 文字列: { success: true, data: string[][] } または { success: false, error: string }
   */
  applyTableOperation(rowsJson: string, action: string, paramsJson: string): string;
}
```

---

## 3. データフォーマットとエラーハンドリング

### 3.1. 戻り値フォーマット
Go から JS への戻り値は、一貫した JSON エンベロープを使用する。

- **成功時 (Success)**:
  ```json
  {
    "success": true,
    "data": [
      ["ヘッダー1", "ヘッダー2"],
      ["値1", "値2"]
    ]
  }
  ```
  または (シリアライズ時):
  ```json
  {
    "success": true,
    "content": "ヘッダー1,ヘッダー2\n値1,値2"
  }
  ```

- **失敗時 (Error)**:
  ```json
  {
    "success": false,
    "error": "CSV parse error at line 3: quote not closed"
  }
  ```

### 3.2. 改行コードの取り扱い
- 改行コードパラメータ:
  - `"LF"`: `\n` (0x0A)
  - `"CRLF"`: `\r\n` (0x0D, 0x0A)
  - `"CR"`: `\r` (0x0D)
- パース時:
  - CSV: Go 標準の `encoding/csv` は `\r\n` と `\n` を透過的に処理可能。CR のみで区切られたデータについては、WASM 側で入力時に統一正規化を行う。
- シリアライズ時:
  - 指定された改行コード文字を用いて行を結合する。

---

## 4. WASM ライフサイクルと初期化シーケンス

1. ブラウザが `index.html` をロード。
2. Go 公式のランタイムヘルパー `wasm_exec.js` を読み込む。
3. `WebAssembly.instantiateStreaming` (または fallback の `fetch` + `WebAssembly.instantiate`) により `main.wasm` をコンパイル・インスタンス化する。
4. `go.run(instance)` を非同期実行し、Go 側の `main()` が `js.Global().Set("WebDataEditor", ...)` を実行。
5. JS 側で `window.WebDataEditorReady` イベントまたはコールバックを発火し、UI コントロールを活性化する。
