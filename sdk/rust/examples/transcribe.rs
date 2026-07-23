// STT：音訊 → 文字。先 cargo run --example speech 整 out.wav。
//   cargo run --example transcribe   （或 cargo run --example transcribe -- 你.wav）
use vible::Vible;

fn main() {
    let path = std::env::args().nth(1).unwrap_or_else(|| "out.wav".to_string());
    let audio = std::fs::read(&path).unwrap_or_else(|e| {
        eprintln!("讀唔到 {path}：{e}（先 cargo run --example speech）");
        std::process::exit(1);
    });
    let v = Vible::new().expect("config");
    println!("{}", v.transcribe(&audio, Some("yue")).expect("transcribe"));
}
