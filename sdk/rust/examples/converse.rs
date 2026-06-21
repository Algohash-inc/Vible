// 語音回合：你講 → agent 喺裝置應你。先 cargo run --example speech 整 out.wav。
//   cargo run --example converse
use vible::Vible;

fn main() {
    let path = std::env::args().nth(1).unwrap_or_else(|| "out.wav".to_string());
    let audio = std::fs::read(&path).unwrap_or_else(|e| {
        eprintln!("讀唔到 {path}：{e}（先 cargo run --example speech）");
        std::process::exit(1);
    });
    let v = Vible::new().expect("config");
    let r = v.converse(&audio, None).expect("converse");
    println!("你講咗：{}", r.transcript);
    println!("{}", if r.spoken { "agent 回覆緊（喺裝置出聲）…" } else { "（靜音 / 冇嘢講）" });
}
