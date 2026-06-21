// STT：音訊 → 文字。先跑 speech.mjs 整 out.wav。
//   node samples/transcribe.mjs   （或 node samples/transcribe.mjs 你自己.wav）
import { readFileSync } from 'node:fs';
import { Vible } from '../src/vible.mjs';

const path = process.argv[2] || 'out.wav';
let audio;
try { audio = readFileSync(path); }
catch { console.error(`讀唔到 ${path}（先跑 samples/speech.mjs）`); process.exit(1); }

const v = await Vible.create();
console.log(await v.transcribe(audio, { language: 'yue' }));
