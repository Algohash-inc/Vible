// 語音回合：你講 → agent 喺裝置應你。先跑 speech.mjs 整 out.wav。
//   node samples/converse.mjs
import { readFileSync } from 'node:fs';
import { Vible } from '../src/vible.mjs';

const path = process.argv[2] || 'out.wav';
let audio;
try { audio = readFileSync(path); }
catch { console.error(`讀唔到 ${path}（先跑 samples/speech.mjs）`); process.exit(1); }

const v = await Vible.create();
const d = await v.converse(audio);
console.log('你講咗：', d.transcript ?? '');
console.log(d.spoken ? 'agent 回覆緊（喺裝置出聲）…' : '（靜音 / 冇嘢講）');
