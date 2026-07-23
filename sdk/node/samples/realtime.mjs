// Realtime WebSocket（全雙工）—— Node 21+ 原生 WebSocket。  node samples/realtime.mjs
import { Vible } from '../src/vible.mjs';

const v = await Vible.create();
await v.realtime('講三樣關於海洋嘅有趣事。', (d) => process.stdout.write(d));
process.stdout.write('\n[完成]\n');
