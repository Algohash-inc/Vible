# Vible Node.js SDK

OpenAI 相容 Vible `/v1` API 嘅輕量 client。**零相依**：Node 18+ 有原生 `fetch`，
Node 21+/22 有原生 `WebSocket`（`realtime()` 要 21+）。

## 用法

```js
import { Vible } from './src/vible.mjs';

const v = await Vible.create();                       // 自動解析金鑰（env 或本機 app）
console.log(await v.chat('用一句話自我介紹。'));        // 非串流

await v.chatStream('數 1 到 5。', d => process.stdout.write(d));   // SSE 串流

const wav = await v.speech('你好。', { voice: 'alloy' });          // TTS → Uint8Array(wav)
console.log(await v.transcribe(wav));                              // STT
console.log(await v.converse(wav));                               // 語音回合
console.log(await v.call('提你開會。', { title: '提醒' }));         // 雲端來電
await v.realtime('講個故事。', d => process.stdout.write(d));       // WebSocket
```

`Vible.create()` 自動解析設定（見 repo 根 [`sdk/README.md`](../README.md) 嘅 env 表）。
要明文傳：`new Vible({ apiKey, baseUrl })` 然後 `await v._resolve()`，或直接 `Vible.create({ apiKey, baseUrl })`。

## Samples

```bash
cd node
node samples/chat.mjs        # 或 npm run chat
node samples/stream.mjs      # npm run stream
node samples/models.mjs
node samples/speech.mjs      # → out.wav
node samples/transcribe.mjs  # 讀 out.wav
node samples/converse.mjs
node samples/call.mjs
node samples/realtime.mjs
```
