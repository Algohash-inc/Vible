// 雲端來電：外部 agent 打電話畀用戶。call key 自動鑄；要裝置連住 hub 至會真響。
//   node samples/call.mjs
import { Vible } from '../src/vible.mjs';

const v = await Vible.create();
console.log(await v.call('提你 3 點開個會。', { title: '會議提醒' }));   // {call_id, ringing, delivered}
