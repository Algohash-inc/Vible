// Chat（非串流）。  node samples/chat.mjs
import { Vible } from '../src/vible.mjs';

const v = await Vible.create();
console.log(await v.chat('用一句話自我介紹。'));
