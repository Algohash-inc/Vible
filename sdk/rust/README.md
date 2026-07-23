# Vible Rust SDK

OpenAI 相容 Vible `/v1` API 嘅 blocking client crate。唔需要 async runtime
（`reqwest` blocking 行 HTTP/SSE/audio/calls，`tungstenite` 自帶 blocking client 行 WebSocket）。

## 用法

```rust
use vible::{Vible, Msg};

let v = Vible::new()?;                                   // 自動解析金鑰（env 或本機 app）
println!("{}", v.chat(&[Msg::user("用一句話自我介紹。")], None)?);

v.chat_stream(&[Msg::user("數 1 到 5。")], |d| print!("{d}"))?;   // SSE 串流

let wav = v.speech("你好。", "alloy")?;                   // TTS → Vec<u8>(wav)
println!("{}", v.transcribe(&wav, Some("yue"))?);        // STT
let r = v.converse(&wav, None)?;                         // 語音回合
println!("{:?}", v.call("提你開會。", Some("提醒"))?);    // 雲端來電
v.realtime(&[Msg::user("講個故事。")], |d| print!("{d}"))?;       // WebSocket
```

加落你個專案：

```toml
[dependencies]
vible = { path = "path/to/sdk/rust" }   # 或日後 crates.io
```

設定解析見 repo 根 [`sdk/README.md`](../README.md) 嘅 env 表。

## Examples

```bash
cd rust
cargo run --example chat
cargo run --example stream
cargo run --example models
cargo run --example speech        # → out.wav
cargo run --example transcribe    # 讀 out.wav（或 -- 你.wav）
cargo run --example converse
cargo run --example call
cargo run --example realtime
```
