# Vible Python SDK

OpenAI 相容 Vible `/v1` API 嘅輕量 client。核心**零相依**（只用標準庫）。

## 安裝

```bash
pip install .                 # 由呢個資料夾安裝做 `vible` 套件
# 或開發直接用：唔裝都得，samples 由同層 import
pip install websockets        # 只有 realtime() 要
```

## 用法

```python
from vible import Vible

v = Vible()                                  # 自動解析金鑰（env 或本機 app）
print(v.chat("用一句話自我介紹。"))           # 非串流

v.chat_stream("數 1 到 5。", on_delta=print)   # SSE 串流，逐 token callback

wav = v.speech("你好。", voice="alloy")        # TTS → wav bytes
open("out.wav", "wb").write(wav)

print(v.transcribe(open("out.wav","rb").read()))            # STT
print(v.converse(open("out.wav","rb").read()))             # 語音回合（agent 喺裝置應你）
print(v.call("提你開會。", title="提醒"))                    # 雲端來電
v.realtime("講個故事。", on_delta=print)        # WebSocket（需 websockets）
```

`Vible()` 會自動解析設定（見 repo 根 [`sdk/README.md`](../README.md) 嘅 env 表）：
有 `VIBLE_API_KEY` 就用 env；冇就向本機 app `GET /api/apikey` 攞 `key/base/ws/model`。

## Samples

```bash
cd python
python3 samples/chat.py
python3 samples/stream.py
python3 samples/models.py
python3 samples/speech.py        # → out.wav
python3 samples/transcribe.py    # 讀 out.wav
python3 samples/converse.py
python3 samples/call.py
pip install websockets && python3 samples/realtime.py
```
