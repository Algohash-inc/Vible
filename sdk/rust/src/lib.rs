//! Vible Rust SDK —— OpenAI 相容 /v1 API 嘅 blocking client。
//!
//! ```no_run
//! use vible::{Vible, Msg};
//! let v = Vible::new().unwrap();
//! println!("{}", v.chat(&[Msg::user("用一句話自我介紹。")], None).unwrap());
//! ```

use std::io::{BufRead, BufReader, Write};
use serde_json::{json, Value};

const OAI_VOICES: &[&str] = &[
    "alloy", "echo", "fable", "onyx", "nova", "shimmer",
    "ash", "ballad", "coral", "sage", "verse", "marin", "cedar",
];

/// API / 連線錯誤。
#[derive(Debug)]
pub struct VibleError(pub String);

impl std::fmt::Display for VibleError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result { write!(f, "{}", self.0) }
}
impl std::error::Error for VibleError {}
impl From<reqwest::Error> for VibleError {
    fn from(e: reqwest::Error) -> Self { VibleError(e.to_string()) }
}

type Result<T> = std::result::Result<T, VibleError>;

/// 一條對話訊息。
#[derive(Clone)]
pub struct Msg {
    pub role: String,
    pub content: String,
}
impl Msg {
    pub fn user(content: &str) -> Self { Msg { role: "user".into(), content: content.into() } }
    pub fn system(content: &str) -> Self { Msg { role: "system".into(), content: content.into() } }
    pub fn assistant(content: &str) -> Self { Msg { role: "assistant".into(), content: content.into() } }
    fn to_json(list: &[Msg]) -> Value {
        Value::Array(list.iter().map(|m| json!({ "role": m.role, "content": m.content })).collect())
    }
}

/// 語音回合結果。
#[derive(Debug)]
pub struct ConverseResult {
    pub transcript: String,
    pub spoken: bool,
}

fn env(k: &str) -> String { std::env::var(k).unwrap_or_default().trim().to_string() }

/// 指向一個 Vible 工作區嘅 client。
pub struct Vible {
    pub base: String,
    pub key: String,
    pub ws: String,
    pub model: String,
    pub host: String,
    pub call_key: String,
    http: reqwest::blocking::Client,
}

impl Vible {
    /// 解析設定：env 優先；冇 `VIBLE_API_KEY` 就向本機 app `GET /api/apikey` 攞。
    pub fn new() -> Result<Vible> {
        let host = { let h = env("VIBLE_HOST"); if h.is_empty() { "127.0.0.1:7770".to_string() } else { h } };
        let mut key = env("VIBLE_API_KEY");
        let mut base = env("VIBLE_BASE_URL");
        let mut ws = env("VIBLE_WS_URL");
        let mut model = env("VIBLE_MODEL");
        let call_key = env("VIBLE_CALL_KEY");
        let http = reqwest::blocking::Client::new();

        if key.is_empty() {
            let url = format!("http://{host}/api/apikey");
            let d: Value = http.get(&url).send().and_then(|r| r.json()).map_err(|e| {
                VibleError(format!(
                    "攞唔到金鑰：{e}。開住 Vible app（{host}）或 set VIBLE_API_KEY / VIBLE_BASE_URL。"
                ))
            })?;
            if key.is_empty() { key = d["key"].as_str().unwrap_or("").into(); }
            if base.is_empty() { base = d["base"].as_str().unwrap_or("").into(); }
            if ws.is_empty() { ws = d["ws"].as_str().unwrap_or("").into(); }
            if model.is_empty() { model = d["model"].as_str().unwrap_or("").into(); }
        }
        if base.is_empty() { base = format!("http://{host}/v1"); }
        if ws.is_empty() {
            ws = format!("{}/v1/realtime", base.replacen("http", "ws", 1).trim_end_matches("/v1"));
        }
        if model.is_empty() { model = "vible-agent".into(); }

        Ok(Vible { base, key, ws, model, host, call_key, http })
    }

    fn check(resp: reqwest::blocking::Response) -> Result<reqwest::blocking::Response> {
        if resp.status().is_success() { return Ok(resp); }
        let status = resp.status();
        let body = resp.text().unwrap_or_default();
        let msg = serde_json::from_str::<Value>(&body)
            .ok()
            .and_then(|v| v["error"]["message"].as_str().map(String::from))
            .unwrap_or(body);
        Err(VibleError(format!("HTTP {status}: {msg}")))
    }

    // ── Chat ────────────────────────────────────────────────────────
    /// 非串流：一次過攞齊回覆文字。
    pub fn chat(&self, messages: &[Msg], model: Option<&str>) -> Result<String> {
        let body = json!({ "model": model.unwrap_or(&self.model), "messages": Msg::to_json(messages), "stream": false });
        let resp = Self::check(self.http.post(format!("{}/chat/completions", self.base))
            .bearer_auth(&self.key).json(&body).send()?)?;
        let v: Value = resp.json()?;
        Ok(v["choices"][0]["message"]["content"].as_str().unwrap_or("").to_string())
    }

    /// SSE 串流：逐 token 叫 `on_delta`，最後回完整文字。
    pub fn chat_stream<F: FnMut(&str)>(&self, messages: &[Msg], mut on_delta: F) -> Result<String> {
        let body = json!({ "model": self.model, "messages": Msg::to_json(messages), "stream": true });
        let resp = Self::check(self.http.post(format!("{}/chat/completions", self.base))
            .bearer_auth(&self.key).json(&body).send()?)?;
        let mut full = String::new();
        for line in BufReader::new(resp).lines() {
            let line = line.unwrap_or_default();
            let Some(p) = line.strip_prefix("data:") else { continue };
            let p = p.trim();
            if p == "[DONE]" { break; }
            if let Ok(v) = serde_json::from_str::<Value>(p) {
                if let Some(d) = v["choices"][0]["delta"]["content"].as_str() {
                    full.push_str(d);
                    on_delta(d);
                }
            }
        }
        Ok(full)
    }

    pub fn models(&self) -> Result<Vec<String>> {
        let resp = Self::check(self.http.get(format!("{}/models", self.base)).bearer_auth(&self.key).send()?)?;
        let v: Value = resp.json()?;
        Ok(v["data"].as_array().map(|a| {
            a.iter().filter_map(|m| m["id"].as_str().map(String::from)).collect()
        }).unwrap_or_default())
    }

    // ── 語音 I/O ────────────────────────────────────────────────────
    // OpenAI 內建聲名/留空 → 預設聲；clone:NAME → 本機 clone；其餘 → Vible speaker。
    fn map_voice(voice: &str, body: &mut Value) {
        let v = voice.trim();
        if v.to_lowercase().starts_with("clone:") {
            body["voice"] = json!(v);
        } else if !v.is_empty() && !OAI_VOICES.contains(&v.to_lowercase().as_str()) {
            body["qwen_voice"] = json!(v);
        } else {
            body["no_clone"] = json!(true);
        }
    }

    /// TTS：text → audio/wav bytes。
    pub fn speech(&self, input: &str, voice: &str) -> Result<Vec<u8>> {
        let mut body = json!({ "model": "vible-tts", "input": input });
        Self::map_voice(voice, &mut body);
        let resp = Self::check(self.http.post(format!("{}/audio/speech", self.base))
            .bearer_auth(&self.key).json(&body).send()?)?;
        Ok(resp.bytes()?.to_vec())
    }

    /// `speak:true` —— 唔回 bytes，改為喺呢部機直接講出。
    pub fn speak(&self, input: &str, voice: &str) -> Result<()> {
        let mut body = json!({ "model": "vible-tts", "input": input, "speak": true });
        Self::map_voice(voice, &mut body);
        Self::check(self.http.post(format!("{}/audio/speech", self.base))
            .bearer_auth(&self.key).json(&body).send()?)?;
        Ok(())
    }

    /// STT：raw 音訊 bytes → 文字。
    pub fn transcribe(&self, audio: &[u8], language: Option<&str>) -> Result<String> {
        let mut url = format!("{}/audio/transcriptions", self.base);
        if let Some(l) = language { url.push_str(&format!("?language={l}")); }
        let resp = Self::check(self.http.post(url).bearer_auth(&self.key)
            .header("Content-Type", "audio/wav").body(audio.to_vec()).send()?)?;
        let v: Value = resp.json()?;
        Ok(v["text"].as_str().unwrap_or("").to_string())
    }

    /// 音訊 → STT → agent → 喺裝置講出回覆。回 {transcript, spoken}。
    pub fn converse(&self, audio: &[u8], language: Option<&str>) -> Result<ConverseResult> {
        let mut url = format!("{}/audio/converse", self.base);
        if let Some(l) = language { url.push_str(&format!("?language={l}")); }
        let resp = Self::check(self.http.post(url).bearer_auth(&self.key)
            .header("Content-Type", "audio/wav").body(audio.to_vec()).send()?)?;
        let v: Value = resp.json()?;
        Ok(ConverseResult {
            transcript: v["transcript"].as_str().unwrap_or("").to_string(),
            spoken: v["spoken"].as_bool().unwrap_or(false),
        })
    }

    // ── 雲端來電 ────────────────────────────────────────────────────
    /// 向本機 operator 端點鑄一條 call key。
    pub fn mint_call_key(&self) -> Result<String> {
        let resp = Self::check(self.http.post(format!("http://{}/api/calls/key", self.host)).send()?)?;
        let v: Value = resp.json()?;
        Ok(v["key"].as_str().unwrap_or("").to_string())
    }

    /// 外部 agent 打電話畀用戶。認證用 call key（唔係 vib- 金鑰）；冇就自動鑄一條。
    pub fn call(&self, prompt: &str, title: Option<&str>) -> Result<Value> {
        let ck = if self.call_key.is_empty() { self.mint_call_key()? } else { self.call_key.clone() };
        let mut body = json!({ "prompt": prompt });
        if let Some(t) = title { body["title"] = json!(t); }
        let resp = Self::check(self.http.post(format!("{}/calls", self.base))
            .bearer_auth(&ck).json(&body).send()?)?;
        Ok(resp.json()?)
    }

    // ── Realtime WebSocket ──────────────────────────────────────────
    /// 全雙工：送一條 turn，逐 token 叫 `on_delta`，回完整文字。
    pub fn realtime<F: FnMut(&str)>(&self, messages: &[Msg], mut on_delta: F) -> Result<String> {
        use tungstenite::{connect, Message};
        let url = format!("{}?api_key={}", self.ws, self.key);
        let (mut sock, _) = connect(url).map_err(|e| VibleError(format!("WS 連接失敗：{e}")))?;
        let create = json!({ "type": "response.create", "messages": Msg::to_json(messages) });
        sock.send(Message::Text(create.to_string().into())).map_err(|e| VibleError(e.to_string()))?;

        let mut full = String::new();
        loop {
            match sock.read() {
                Ok(Message::Text(t)) => {
                    let m: Value = serde_json::from_str(t.as_str()).unwrap_or(Value::Null);
                    match m["type"].as_str() {
                        Some("response.output_text.delta") => {
                            let d = m["delta"].as_str().unwrap_or("");
                            full.push_str(d);
                            on_delta(d);
                        }
                        Some("response.completed") => break,
                        Some("error") => return Err(VibleError(m["error"]["message"].as_str().unwrap_or("").to_string())),
                        _ => {}
                    }
                }
                Ok(Message::Close(_)) | Err(_) => break,
                _ => {}
            }
        }
        Ok(full)
    }
}

// flush helper 畀 examples 用（避免每個 example 都 import std::io::Write）。
#[doc(hidden)]
pub fn flush() { std::io::stdout().flush().ok(); }
