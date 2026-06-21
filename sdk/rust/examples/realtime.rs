// Realtime WebSocket（全雙工）。  cargo run --example realtime
use vible::{Msg, Vible};

fn main() {
    let v = Vible::new().expect("config");
    v.realtime(&[Msg::user("講三樣關於海洋嘅有趣事。")], |d| {
        print!("{d}");
        vible::flush();
    }).expect("realtime");
    println!("\n[完成]");
}
