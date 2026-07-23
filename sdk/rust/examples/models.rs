// 列出工作區可用 model。  cargo run --example models
use vible::Vible;

fn main() {
    let v = Vible::new().expect("config");
    for m in v.models().expect("models") {
        println!("{m}");
    }
}
