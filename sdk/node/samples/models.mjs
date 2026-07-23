// 列出工作區可用 model。  node samples/models.mjs
import { Vible } from '../src/vible.mjs';

const v = await Vible.create();
for (const m of await v.models()) console.log(m);
