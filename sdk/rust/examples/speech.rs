// TTS：text → out.wav。voice："alloy"/留空→預設；"clone:名"→clone；speaker 名→qwen_voice。
//   cargo run --example speech
use std::io::Write;
use vible::Vible;

fn main() {
    let v = Vible::new().expect("config");
    let wav = v.speech("你好，我係你個 agent 嘅聲音。", "alloy").expect("speech");
    std::fs::File::create("out.wav").and_then(|mut f| f.write_all(&wav)).expect("write");
    println!("已儲存 out.wav（{} bytes）。播放：afplay out.wav", wav.len());
}
