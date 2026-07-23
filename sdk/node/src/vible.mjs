// Vible Node.js SDK —— OpenAI 相容 /v1 API 嘅輕量 client。
// 零相依：Node 18+ 有原生 fetch，Node 21+/22 有原生 WebSocket（realtime() 要 21+）。
//
//   import { Vible } from './src/vible.mjs';
//   const v = await Vible.create();
//   console.log(await v.chat('用一句話自我介紹。'));

const OAI_VOICES = new Set([
  'alloy', 'echo', 'fable', 'onyx', 'nova', 'shimmer',
  'ash', 'ballad', 'coral', 'sage', 'verse', 'marin', 'cedar',
]);

export class VibleError extends Error {
  constructor(message, status) { super(message); this.name = 'VibleError'; this.status = status; }
}

const normMessages = (m) => (typeof m === 'string' ? [{ role: 'user', content: m }] : m);

export class Vible {
  constructor(cfg = {}) {
    this.host = cfg.host || process.env.VIBLE_HOST || '127.0.0.1:7770';
    this.key = (cfg.apiKey || process.env.VIBLE_API_KEY || '').trim();
    this.base = (cfg.baseUrl || process.env.VIBLE_BASE_URL || '').trim();
    this.ws = (cfg.wsUrl || process.env.VIBLE_WS_URL || '').trim();
    this.model = (cfg.model || process.env.VIBLE_MODEL || '').trim();
    this.callKey = (cfg.callKey || process.env.VIBLE_CALL_KEY || '').trim();
  }

  // async 解析金鑰，所以用 factory（constructor 唔可以 await）。
  static async create(cfg = {}) {
    const v = new Vible(cfg);
    await v._resolve();
    return v;
  }

  async _resolve() {
    // 冇明文金鑰就向本機 app 攞（GET /api/apikey，只限本機 operator）。
    if (!this.key) {
      let d;
      try {
        const r = await fetch(`http://${this.host}/api/apikey`);
        if (!r.ok) throw new Error('HTTP ' + r.status);
        d = await r.json();
      } catch (e) {
        throw new VibleError(`攞唔到金鑰：${e.message}。開住 Vible app（${this.host}）或 set VIBLE_API_KEY / VIBLE_BASE_URL。`);
      }
      this.key ||= d.key; this.base ||= d.base; this.ws ||= d.ws; this.model ||= d.model;
    }
    if (!this.base) this.base = `http://${this.host}/v1`;
    if (!this.ws) this.ws = this.base.replace('http', 'ws').replace(/\/v1$/, '') + '/v1/realtime';
    if (!this.model) this.model = 'vible-agent';
    return this;
  }

  _headers(extra, key) {
    return { Authorization: `Bearer ${key || this.key}`, ...(extra || {}) };
  }

  async _err(r) {
    const body = await r.text().catch(() => '');
    let msg = body;
    try { msg = JSON.parse(body).error?.message || body; } catch {}
    return new VibleError(`HTTP ${r.status}: ${msg}`, r.status);
  }

  // ── Chat ──────────────────────────────────────────────────────────
  async chat(messages, { model, ...params } = {}) {
    const r = await fetch(`${this.base}/chat/completions`, {
      method: 'POST',
      headers: this._headers({ 'Content-Type': 'application/json' }),
      body: JSON.stringify({ model: model || this.model, messages: normMessages(messages), stream: false, ...params }),
    });
    if (!r.ok) throw await this._err(r);
    const d = await r.json();
    return d.choices?.[0]?.message?.content ?? '';
  }

  // SSE 串流：逐 token 叫 onDelta，最後回完整文字。
  async chatStream(messages, onDelta, { model, ...params } = {}) {
    const r = await fetch(`${this.base}/chat/completions`, {
      method: 'POST',
      headers: this._headers({ 'Content-Type': 'application/json' }),
      body: JSON.stringify({ model: model || this.model, messages: normMessages(messages), stream: true, ...params }),
    });
    if (!r.ok) throw await this._err(r);
    const reader = r.body.getReader();
    const dec = new TextDecoder();
    let buf = '', full = '';
    for (;;) {
      const { done, value } = await reader.read();
      if (done) break;
      buf += dec.decode(value, { stream: true });
      let i;
      while ((i = buf.indexOf('\n\n')) >= 0) {
        const block = buf.slice(0, i); buf = buf.slice(i + 2);
        for (const line of block.split('\n')) {
          if (!line.startsWith('data:')) continue;
          const p = line.slice(5).trim();
          if (p === '[DONE]') return full;
          const delta = JSON.parse(p).choices?.[0]?.delta?.content;
          if (delta) { full += delta; onDelta?.(delta); }
        }
      }
    }
    return full;
  }

  async models() {
    const r = await fetch(`${this.base}/models`, { headers: this._headers() });
    if (!r.ok) throw await this._err(r);
    return (await r.json()).data?.map((m) => m.id) ?? [];
  }

  // ── 語音 I/O ──────────────────────────────────────────────────────
  // OpenAI 內建聲名/留空 → 預設聲；clone:NAME → 本機 clone；其餘 → Vible speaker。
  static _mapVoice(voice, body) {
    const v = (voice || '').trim();
    if (/^clone:/i.test(v)) body.voice = v;
    else if (v && !OAI_VOICES.has(v.toLowerCase())) body.qwen_voice = v;
    else body.no_clone = true;
    return body;
  }

  // TTS：回 Uint8Array(wav)。speak:true → 喺裝置講出，回 {spoken:true}。
  async speech(input, { voice = 'alloy', speed, lang, speak = false } = {}) {
    const body = Vible._mapVoice(voice, { model: 'vible-tts', input });
    if (speed != null) body.speed = speed;
    if (lang) body.lang = lang;
    if (speak) body.speak = true;
    const r = await fetch(`${this.base}/audio/speech`, {
      method: 'POST', headers: this._headers({ 'Content-Type': 'application/json' }), body: JSON.stringify(body),
    });
    if (!r.ok) throw await this._err(r);
    if ((r.headers.get('content-type') || '').includes('application/json')) return r.json();
    return new Uint8Array(await r.arrayBuffer());
  }

  // STT：raw 音訊 bytes → 文字。
  async transcribe(audio, { language, contentType = 'audio/wav' } = {}) {
    const url = `${this.base}/audio/transcriptions` + (language ? `?language=${language}` : '');
    const r = await fetch(url, { method: 'POST', headers: this._headers({ 'Content-Type': contentType }), body: audio });
    if (!r.ok) throw await this._err(r);
    return (await r.json()).text ?? '';
  }

  // 音訊 → STT → agent → 喺裝置講出回覆。回 {transcript, spoken}。
  async converse(audio, { language, contentType = 'audio/wav' } = {}) {
    const url = `${this.base}/audio/converse` + (language ? `?language=${language}` : '');
    const r = await fetch(url, { method: 'POST', headers: this._headers({ 'Content-Type': contentType }), body: audio });
    if (!r.ok) throw await this._err(r);
    return r.json();
  }

  // ── 雲端來電 ──────────────────────────────────────────────────────
  async mintCallKey() {
    const r = await fetch(`http://${this.host}/api/calls/key`, { method: 'POST' });
    if (!r.ok) throw await this._err(r);
    return (await r.json()).key;
  }

  // 認證用 call key（唔係 vib- 金鑰）。冇就自動鑄一條。
  async call(prompt, { title, caller, context, callKey } = {}) {
    const ck = callKey || this.callKey || (await this.mintCallKey());
    const body = { prompt };
    if (title) body.title = title;
    if (caller) body.caller = caller;
    if (context !== undefined) body.context = context;
    const r = await fetch(`${this.base}/calls`, {
      method: 'POST', headers: this._headers({ 'Content-Type': 'application/json' }, ck), body: JSON.stringify(body),
    });
    if (!r.ok) throw await this._err(r);
    return r.json();
  }

  // ── Realtime WebSocket（Node 21+ 原生 WebSocket）──────────────────
  realtime(messages, onDelta) {
    return new Promise((resolve, reject) => {
      const ws = new WebSocket(`${this.ws}?api_key=${this.key}`);
      let full = '';
      ws.addEventListener('open', () =>
        ws.send(JSON.stringify({ type: 'response.create', messages: normMessages(messages) })));
      ws.addEventListener('message', (e) => {
        const m = JSON.parse(e.data);
        if (m.type === 'response.output_text.delta') { full += m.delta; onDelta?.(m.delta); }
        else if (m.type === 'response.completed') { ws.close(); resolve(full); }
        else if (m.type === 'error') { ws.close(); reject(new VibleError(m.error.message)); }
      });
      ws.addEventListener('error', (e) => reject(new VibleError(e.message || 'WS error')));
    });
  }
}

export default Vible;
