// Chat（非串流）。  cargo run --example chat
use vible::{Msg, Vible};

fn main() {
    let v = Vible::new().expect("config");
    println!("{}", v.chat(&[Msg::user("用一句話自我介紹。")], None).expect("chat"));
}
