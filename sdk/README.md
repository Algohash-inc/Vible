<div align="center">

# Vible SDK

**把任何 agent 指向你嘅 Vible 工作區 —— 佢就能講、能聽、能對話，甚至能打電話畀你。**

OpenAI 相容 `/v1` API · 五語言 client library + samples

</div>

Vible 桌面 app（[下載 DMG](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg)）
本身開咗一個 **OpenAI 相容嘅 `/v1/*` API 埠**（本機 `127.0.0.1:7770`，公網工作區
`https://<你>.api.vible.space`）。呢個 `sdk/` 資料夾就係佢嘅官方 SDK —— 每種語言都有一個
**可重用嘅 client library** 加一套 **可即跑嘅 samples**。

| 語言 | Client | 資料夾 | 相依 |
|------|--------|--------|------|
| **Python** | `from vible import Vible` | [`python/`](python/) | 核心零相依（標準庫） |
| **Node.js** | `import { Vible }` | [`node/`](node/) | 零相依（Node 21+） |
| **Rust** | `vible::Vible` | [`rust/`](rust/) | `reqwest` + `tungstenite` |
| **Go** | `vible.New()` | [`go/`](go/) | 零相依（標準庫） |
| **C++** | `vible::Client` | [`cpp/`](cpp/) | libcurl + nlohmann/json |

全部 **OpenAI 相容**：既有嘅 SDK、工具、提示詞照用，淨係換 `base_url` + `api_key`。

---

## 端點一覽

| 能力 | 端點 | Client method |
|------|------|---------------|
| **Chat · HTTP** | `POST /v1/chat/completions`（`stream:false`） | `chat()` |
| **Chat · SSE** | `POST /v1/chat/completions`（`stream:true`） | `chat_stream()` |
| **Chat · WebSocket** | `GET /v1/realtime` | `realtime()` |
| **Models** | `GET /v1/models` | `models()` |
| **語音輸出（TTS）** | `POST /v1/audio/speech` | `speech()` |
| **語音輸入（STT）** | `POST /v1/audio/transcriptions` | `transcribe()` |
| **語音回合** | `POST /v1/audio/converse` | `converse()` |
| **雲端來電** | `POST /v1/calls` | `call()` |

> Chat / models / audio 認證用 **`vib-` 金鑰**；只有 `/v1/calls` 用另一種 **call key**（見下）。
> C++ 未內建 WebSocket（`realtime()`），其餘四種語言都有。

---

## 攞金鑰（零設定）

開 **Vible 桌面 app → 左邊「API」分頁** 撳「複製」攞 `vib-…` 金鑰 + Base URL。

> **本機跑 demo 乜都唔使設定** —— 每個 client 會自動由 `http://127.0.0.1:7770/api/apikey`
> 攞返 `key / base / ws / model`（呢個端點只限本機 operator）。
> 連雲端工作區 / 部署 / CI 時，先 `cp .env.example .env` 並 export `VIBLE_API_KEY` + `VIBLE_BASE_URL`。

金鑰係**每部機系統自動產生 + 加密儲存**（DEK 綁機，config 畀人偷走都解唔到），
換機自動重生；app「API」分頁可隨時「重新產生」令舊鑰即時失效。

---

## 各語言 30 秒上手

**Python**
```bash
cd python && python3 samples/chat.py
```
**Node.js**（Node 21+）
```bash
cd node && node samples/chat.mjs
```
**Rust**
```bash
cd rust && cargo run --example chat
```
**Go**
```bash
cd go && go run ./examples/chat
```
**C++**（需 libcurl）
```bash
cd cpp && cmake -S . -B build && cmake --build build && ./build/chat
```

每個資料夾嘅 `README.md` 有完整 method 簽名同所有 samples 嘅跑法。

---

## 語音 I/O 細節

**TTS `/v1/audio/speech`** body `{ model:"vible-tts", input, voice, speed?, lang?, speak? }`，回 `audio/wav`：

- `voice`：OpenAI 內建名（`alloy`…）或留空 → **Vible 預設聲**（本機免額度）；
  `"clone:你個clone名"` → 本機 voice clone；其餘字串 → **Vible speaker 名**（`qwen_voice`）。
- `speak:true` → 唔回 bytes，直接喺**呢部機講出嚟**（畀你個 agent 把聲）。

**STT `/v1/audio/transcriptions`** 收 `multipart/form-data`（`file`）或**直接 raw audio bytes**
（SDK 用後者，最簡單）+ `Content-Type: audio/wav`，可選 `?language=`，回 `{ text }`。

**`/v1/audio/converse`** 收音訊 → STT → 餵 Vible agent → 經 speakqueue **喺裝置出聲回覆**，
即時回 `{ transcript, spoken }`。

---

## 雲端來電 `/v1/calls`

畀外部 agent 主動 ring 用戶（原生響鈴 + web SSE）。認證用 **call key**（唔係 `vib-` 金鑰）：

1. 攞 call key：app loopback `POST /api/calls/key`（SDK 會自動鑄），或 set `VIBLE_CALL_KEY`。
2. `POST /v1/calls`，`Authorization: Bearer <call key>`，body `{ prompt, title?, caller?, context? }`。
3. 回 `202 { call_id, ringing, delivered }`。

> 要裝置已連住 hub（device-link）至會真係響；單機 demo 會回 `delivered:false`。

---

## WebSocket 協定（`/v1/realtime`）

```
送  {"type":"response.create","messages":[{"role":"user","content":"…"}]}
收  {"type":"session.created"}                                連線就緒
    {"type":"response.output_text.delta","delta":"…"}         逐 token
    {"type":"response.completed"}                              完成
    {"type":"error","error":{"message":"…"}}                   失敗
```
金鑰三種帶法：`?api_key=vib-…`、`Authorization: Bearer vib-…`、或瀏覽器 `Sec-WebSocket-Protocol: vib-…`。

---

## 環境變數

| 變數 | 預設 | 說明 |
|------|------|------|
| `VIBLE_API_KEY` | （本機自動攞） | 工作區 `vib-` 金鑰 |
| `VIBLE_BASE_URL` | `http://127.0.0.1:7770/v1` | API base（雲端 `https://<你>.api.vible.space/v1`） |
| `VIBLE_WS_URL` | 由 base 推導 | WebSocket URL |
| `VIBLE_HOST` | `127.0.0.1:7770` | 只換 host、保留 `/v1`、`/api/apikey` 路徑 |
| `VIBLE_MODEL` | 由 app 報 / `vible-agent` | 指定 model |
| `VIBLE_CALL_KEY` | （自動鑄） | `/v1/calls` 用嘅 call key |

## 認證 & 錯誤

- 金鑰錯／無 → HTTP `401`；WebSocket 直接握手失敗。
- 上游錯誤（例如 credit 不足 `402`）原樣以 OpenAI 錯誤格式 + 對應 HTTP status 回傳；
  各 client 會 throw / return 一個帶 `status` 嘅錯誤。
- `/v1/*` 全部 CORS 開放（`OPTIONS` preflight 已處理）。
