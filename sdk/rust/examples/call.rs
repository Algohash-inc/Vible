// 雲端來電：外部 agent 打電話畀用戶。call key 自動鑄；要裝置連住 hub 至會真響。
//   cargo run --example call
use vible::Vible;

fn main() {
    let v = Vible::new().expect("config");
    match v.call("提你 3 點開個會。", Some("會議提醒")) {
        Ok(r) => println!("{r}"),   // {call_id, ringing, delivered}
        Err(e) => println!("call 失敗：{e}"),
    }
}
