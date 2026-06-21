// Chat（SSE 串流）—— 逐 token 印。  cargo run --example stream
use vible::{Msg, Vible};

fn main() {
    let v = Vible::new().expect("config");
    v.chat_stream(&[Msg::user("數 1 到 5，每個數字之間停一下。")], |d| {
        print!("{d}");
        vible::flush();
    }).expect("stream");
    println!();
}
