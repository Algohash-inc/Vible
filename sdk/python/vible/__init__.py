"""
Vible Python SDK —— OpenAI 相容 /v1 API 嘅輕量 client。

核心（chat / stream / models / speech / transcribe / converse / call）只用 Python
標準庫（urllib），零相依、即裝即用。realtime WebSocket 用到第三方 `websockets`
（選擇性，requirements.txt 有）。

    from vible import Vible
    v = Vible()                       # 自動解析金鑰（env 或本機 app）
    print(v.chat("用一句話自我介紹。"))
"""
from __future__ import annotations

import json
import os
import urllib.request
import urllib.error
from typing import Callable, Iterable

__all__ = ["Vible", "VibleError"]
__version__ = "1.0.0"

_OAI_VOICES = {
    "alloy", "echo", "fable", "onyx", "nova", "shimmer",
    "ash", "ballad", "coral", "sage", "verse", "marin", "cedar",
}


class VibleError(Exception):
    """API / 連線錯誤。`status` 帶 HTTP 狀態碼（若有）。"""

    def __init__(self, message: str, status: int | None = None):
        super().__init__(message)
        self.status = status


def _norm_messages(messages) -> list[dict]:
    """方便用法：傳 str 當單條 user message；傳 list 就照用。"""
    if isinstance(messages, str):
        return [{"role": "user", "content": messages}]
    return list(messages)


class Vible:
    """指向一個 Vible 工作區嘅 client。"""

    def __init__(
        self,
        api_key: str | None = None,
        base_url: str | None = None,
        ws_url: str | None = None,
        model: str | None = None,
        host: str | None = None,
        call_key: str | None = None,
    ):
        self.host = host or os.environ.get("VIBLE_HOST", "127.0.0.1:7770")
        self.key = (api_key or os.environ.get("VIBLE_API_KEY", "")).strip()
        self.base = (base_url or os.environ.get("VIBLE_BASE_URL", "")).strip()
        self.ws = (ws_url or os.environ.get("VIBLE_WS_URL", "")).strip()
        self.model = (model or os.environ.get("VIBLE_MODEL", "")).strip()
        self.call_key = (call_key or os.environ.get("VIBLE_CALL_KEY", "")).strip()

        # 冇明文金鑰就向本機 app 攞（GET /api/apikey 只限本機 operator，會連
        # base / ws / model 一齊回）。換機自動重生嘅金鑰亦由呢度攞到最新一條。
        if not self.key:
            try:
                d = self._get_json(f"http://{self.host}/api/apikey")
            except Exception as e:  # noqa: BLE001 — 攞唔到就畀清楚指示
                raise VibleError(
                    f"攞唔到金鑰：{e}。開住 Vible 桌面 app（{self.host}），"
                    f"或 set VIBLE_API_KEY / VIBLE_BASE_URL。"
                )
            self.key = self.key or d.get("key", "")
            self.base = self.base or d.get("base", "")
            self.ws = self.ws or d.get("ws", "")
            self.model = self.model or d.get("model", "")

        if not self.base:
            self.base = f"http://{self.host}/v1"
        if not self.ws:
            # http(s)://…/v1 → ws(s)://…/v1/realtime
            self.ws = self.base.replace("http", "ws", 1).rstrip("/")
            self.ws = self.ws[:-3] + "/v1/realtime" if self.ws.endswith("/v1") else self.ws + "/v1/realtime"
        if not self.model:
            self.model = "vible-agent"

    # ── 內部 HTTP 工具 ──────────────────────────────────────────────────
    def _headers(self, extra: dict | None = None, key: str | None = None) -> dict:
        h = {"Authorization": f"Bearer {key or self.key}"}
        if extra:
            h.update(extra)
        return h

    def _request(self, method, url, *, headers=None, data=None):
        req = urllib.request.Request(url, method=method, headers=headers or {}, data=data)
        try:
            return urllib.request.urlopen(req)
        except urllib.error.HTTPError as e:
            body = e.read().decode("utf-8", "replace")
            msg = body
            try:
                msg = json.loads(body).get("error", {}).get("message", body)
            except Exception:  # noqa: BLE001
                pass
            raise VibleError(f"HTTP {e.code}: {msg}", status=e.code)

    def _get_json(self, url, headers=None):
        with self._request("GET", url, headers=headers) as r:
            return json.load(r)

    # ── Chat ────────────────────────────────────────────────────────────
    def chat(self, messages, *, model: str | None = None, **params) -> str:
        """非串流：一次過攞齊回覆文字。"""
        body = {"model": model or self.model, "messages": _norm_messages(messages), "stream": False}
        body.update(params)
        with self._request(
            "POST", f"{self.base}/chat/completions",
            headers=self._headers({"Content-Type": "application/json"}),
            data=json.dumps(body).encode(),
        ) as r:
            d = json.load(r)
        return d["choices"][0]["message"]["content"]

    def chat_stream(self, messages, on_delta: Callable[[str], None] | None = None, *,
                    model: str | None = None, **params) -> str:
        """SSE 串流：每粒 token 即時 callback on_delta，最後 return 完整文字。"""
        body = {"model": model or self.model, "messages": _norm_messages(messages), "stream": True}
        body.update(params)
        full = []
        with self._request(
            "POST", f"{self.base}/chat/completions",
            headers=self._headers({"Content-Type": "application/json"}),
            data=json.dumps(body).encode(),
        ) as r:
            for raw in r:  # urllib 逐行 stream
                line = raw.decode("utf-8").rstrip("\n")
                if not line.startswith("data:"):
                    continue
                payload = line[5:].strip()
                if payload == "[DONE]":
                    break
                delta = json.loads(payload)["choices"][0]["delta"].get("content")
                if delta:
                    full.append(delta)
                    if on_delta:
                        on_delta(delta)
        return "".join(full)

    def models(self) -> list[str]:
        d = self._get_json(f"{self.base}/models", headers=self._headers())
        return [m["id"] for m in d.get("data", [])]

    # ── 語音 I/O ────────────────────────────────────────────────────────
    @staticmethod
    def _map_voice(voice: str, body: dict) -> dict:
        """OpenAI 內建聲名/留空 → 預設聲；clone:NAME → 本機 clone；其餘 → Vible speaker。"""
        v = (voice or "").strip()
        if v.lower().startswith("clone:"):
            body["voice"] = v
        elif v and v.lower() not in _OAI_VOICES:
            body["qwen_voice"] = v
        else:
            body["no_clone"] = True
        return body

    def speech(self, input: str, *, voice: str = "alloy", speed: float | None = None,
               lang: str | None = None, speak: bool = False):
        """TTS：text → audio/wav bytes。speak=True 改為喺裝置直接講出，回 dict。"""
        body = self._map_voice(voice, {"model": "vible-tts", "input": input})
        if speed is not None:
            body["speed"] = speed
        if lang:
            body["lang"] = lang
        if speak:
            body["speak"] = True
        with self._request(
            "POST", f"{self.base}/audio/speech",
            headers=self._headers({"Content-Type": "application/json"}),
            data=json.dumps(body).encode(),
        ) as r:
            ct = r.headers.get("Content-Type", "")
            data = r.read()
        if "application/json" in ct:   # speak:true → {spoken:true}
            return json.loads(data)
        return data                    # audio/wav bytes

    def transcribe(self, audio: bytes, *, language: str | None = None,
                   content_type: str = "audio/wav") -> str:
        """STT：raw 音訊 bytes → 文字。"""
        url = f"{self.base}/audio/transcriptions"
        if language:
            url += f"?language={language}"
        with self._request(
            "POST", url,
            headers=self._headers({"Content-Type": content_type}),
            data=audio,
        ) as r:
            return json.load(r).get("text", "")

    def converse(self, audio: bytes, *, language: str | None = None,
                 content_type: str = "audio/wav") -> dict:
        """音訊 → STT → agent → 喺裝置講出回覆。即時回 {transcript, spoken}。"""
        url = f"{self.base}/audio/converse"
        if language:
            url += f"?language={language}"
        with self._request(
            "POST", url,
            headers=self._headers({"Content-Type": content_type}),
            data=audio,
        ) as r:
            return json.load(r)

    # ── 雲端來電 ────────────────────────────────────────────────────────
    def mint_call_key(self) -> str:
        """向本機 operator 端點鑄一條 call key（綁住呢部機 device_id）。"""
        with self._request("POST", f"http://{self.host}/api/calls/key") as r:
            return json.load(r)["key"]

    def call(self, prompt: str, *, title: str | None = None, caller: str | None = None,
             context=None, call_key: str | None = None) -> dict:
        """外部 agent 主動打電話畀用戶。認證用 call key（唔係 vib- 金鑰）。"""
        ck = call_key or self.call_key or self.mint_call_key()
        body = {"prompt": prompt}
        if title:
            body["title"] = title
        if caller:
            body["caller"] = caller
        if context is not None:
            body["context"] = context
        with self._request(
            "POST", f"{self.base}/calls",
            headers=self._headers({"Content-Type": "application/json"}, key=ck),
            data=json.dumps(body).encode(),
        ) as r:
            return json.load(r)

    # ── Realtime WebSocket（需 `websockets`）────────────────────────────
    def realtime(self, messages, on_delta: Callable[[str], None]) -> str:
        """全雙工：送一條 turn，逐 token callback，回完整文字。需要 `pip install websockets`。"""
        try:
            import asyncio
            import websockets
        except ImportError:
            raise VibleError("realtime 需要 `pip install websockets`。")

        url = f"{self.ws}?api_key={self.key}"
        full = []

        async def _run():
            async with websockets.connect(url) as ws:
                await ws.send(json.dumps({"type": "response.create", "messages": _norm_messages(messages)}))
                async for raw in ws:
                    m = json.loads(raw)
                    t = m.get("type")
                    if t == "response.output_text.delta":
                        full.append(m["delta"])
                        on_delta(m["delta"])
                    elif t == "response.completed":
                        break
                    elif t == "error":
                        raise VibleError(m["error"]["message"])

        asyncio.run(_run())
        return "".join(full)
