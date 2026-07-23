// Chat（SSE 串流）—— 逐 token 印。  node samples/stream.mjs
import { Vible } from '../src/vible.mjs';

const v = await Vible.create();
await v.chatStream('數 1 到 5，每個數字之間停一下。', (d) => process.stdout.write(d));
process.stdout.write('\n');
