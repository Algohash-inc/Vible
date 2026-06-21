// TTS：text → out.wav。voice："alloy"/留空→預設；"clone:名"→clone；speaker 名→qwen_voice。
//   node samples/speech.mjs
import { writeFileSync } from 'node:fs';
import { Vible } from '../src/vible.mjs';

const v = await Vible.create();
const wav = await v.speech('你好，我係你個 agent 嘅聲音。', { voice: 'alloy' });
writeFileSync('out.wav', wav);
console.log(`已儲存 out.wav（${wav.length} bytes）。播放：afplay out.wav`);
