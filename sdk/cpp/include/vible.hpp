// vible.hpp — header-only C++17 client for the Vible OpenAI-compatible API.
//
// A thin client over a Vible desktop workspace's /v1/* HTTP API
// (local 127.0.0.1:7770, or cloud https://<workspace>.api.vible.space).
// Auth is a `vib-` bearer key; cloud phone calls use a separate call key.
//
// Dependencies: libcurl (HTTP/SSE/audio/calls) + nlohmann/json (JSON).
//
// Usage:
//   #include "vible.hpp"
//   vible::Client c;                                  // resolves config from env / local app
//   std::string reply = c.chat({{"user", "Hello"}});
//
// All methods throw vible::Error on transport or HTTP-status failure.

#ifndef VIBLE_HPP
#define VIBLE_HPP

#include <cstdint>
#include <cstdlib>
#include <functional>
#include <stdexcept>
#include <string>
#include <vector>

#include <curl/curl.h>
#include <nlohmann/json.hpp>

namespace vible {

// 所有錯誤（傳輸層失敗或 HTTP 非 2xx）都用呢個 type 拋出。status==0 代表
// libcurl/連線層出錯（未收到 HTTP 回應）；否則就係伺服器回嘅 HTTP 狀態碼。
struct Error : std::runtime_error {
  long status;  // HTTP status, or 0 for a transport-level failure
  explicit Error(const std::string& msg, long st = 0)
      : std::runtime_error(msg), status(st) {}
};

struct Message {
  std::string role;
  std::string content;
};

// chat / chat_stream 嘅可選參數。temperature/max_tokens 留 <0 / <=0 即「唔傳」，
// 等伺服器用自己預設（唔好硬塞 0 落上游覆蓋咗預設值）。
struct ChatOptions {
  std::string model;          // override the resolved default model
  double temperature = -1.0;  // <0  → omit
  int max_tokens = 0;         // <=0 → omit
};

struct SpeechOptions {
  std::string voice;     // ""/alloy/echo/... → default; "clone:NAME" → local clone;
                         // any other string → Vible speaker name
  std::string lang;      // optional language hint
  double speed = 0.0;    // <=0 → omit (server default speed)
  bool speak = false;    // true → play on the device, return no audio bytes
};

struct ConverseResult {
  std::string transcript;
  bool spoken = false;
};

struct CallResult {
  std::string call_id;
  bool ringing = false;
  bool delivered = false;
};

class Client {
 public:
  // 預設建構：由環境變數 + 本機 operator 端點解析設定（見 resolve()）。
  Client() { resolve(); }

  // 明確指定 base_url + vib- key（其餘設定仍可由環境補上）。
  Client(std::string base, std::string key)
      : base_(std::move(base)), key_(std::move(key)) {
    host_ = env_("VIBLE_HOST", "127.0.0.1:7770");
    model_ = env_("VIBLE_MODEL", "vible-agent");
    call_key_ = env_("VIBLE_CALL_KEY", "");
  }

  const std::string& base() const { return base_; }
  const std::string& model() const { return model_; }

  // ── 1. Chat (non-stream) ─────────────────────────────────────────────
  // POST {base}/chat/completions, returns choices[0].message.content.
  std::string chat(const std::vector<Message>& messages,
                   const ChatOptions& opt = {}) {
    nlohmann::json body = build_chat_body(messages, opt, /*stream=*/false);
    std::string resp;
    long code = http_post_json(base_ + "/chat/completions", body.dump(), key_,
                               &resp, nullptr);
    if (code < 200 || code >= 300) throw http_error("chat", code, resp);
    auto j = parse(resp, "chat");
    return j.at("choices").at(0).at("message").at("content").get<std::string>();
  }

  // ── 2. Chat (stream / SSE) ───────────────────────────────────────────
  // POST {base}/chat/completions with stream:true. Invokes on_delta once per
  // token as it arrives.
  void chat_stream(const std::vector<Message>& messages,
                   std::function<void(const std::string&)> on_delta,
                   const ChatOptions& opt = {}) {
    nlohmann::json body = build_chat_body(messages, opt, /*stream=*/true);

    // SSE 解析狀態：libcurl 嘅 write callback 一次只交到一截任意大細嘅 byte，
    // 邊界可以切喺一行/一個 event 中間。所以累積落 sse.buffer，逐個 "\n\n"
    // 分隔嘅完整 event block 先處理；"data: [DONE]" 為終止符。
    SseState sse;
    sse.on_delta = std::move(on_delta);

    std::string err;
    long code = http_post_json(base_ + "/chat/completions", body.dump(), key_,
                               nullptr, &sse, &err);
    if (!err.empty()) throw Error("chat_stream: " + err, 0);
    if (code < 200 || code >= 300)
      throw http_error("chat_stream", code, sse.buffer);
  }

  // ── 3. Models ────────────────────────────────────────────────────────
  // GET {base}/models → list of model ids.
  std::vector<std::string> models() {
    std::string resp;
    long code = http_get(base_ + "/models", key_, &resp);
    if (code < 200 || code >= 300) throw http_error("models", code, resp);
    auto j = parse(resp, "models");
    std::vector<std::string> out;
    for (const auto& m : j.at("data")) out.push_back(m.at("id").get<std::string>());
    return out;
  }

  // ── 4. Speech (TTS) ──────────────────────────────────────────────────
  // POST {base}/audio/speech. Returns audio/wav bytes — UNLESS opt.speak is
  // true, in which case the audio plays on the device and an empty vector is
  // returned (the server replies with JSON {spoken:true}, no bytes).
  std::vector<uint8_t> speech(const std::string& input,
                              const SpeechOptions& opt = {}) {
    nlohmann::json body{{"model", "vible-tts"}, {"input", input}};
    body["voice"] = opt.voice;  // voice mapping is done server-side (see header doc)
    if (!opt.lang.empty()) body["lang"] = opt.lang;
    if (opt.speed > 0.0) body["speed"] = opt.speed;
    if (opt.speak) body["speak"] = true;

    std::vector<uint8_t> bytes;
    std::string ctype;
    long code = http_post_json_bytes(base_ + "/audio/speech", body.dump(), key_,
                                     &bytes, &ctype);
    if (code < 200 || code >= 300) {
      throw http_error("speech", code,
                       std::string(bytes.begin(), bytes.end()));
    }
    // speak:true 時伺服器回 JSON {spoken:true} 而唔係 wav；冇 bytes 可回。
    if (ctype.find("application/json") != std::string::npos) return {};
    return bytes;
  }

  // ── 5. Transcribe (STT) ──────────────────────────────────────────────
  // POST {base}/audio/transcriptions with the RAW audio bytes as the body and
  // Content-Type: audio/wav. Returns the recognised text.
  std::string transcribe(const std::vector<uint8_t>& audio,
                         const std::string& language = "") {
    std::string url = base_ + "/audio/transcriptions";
    if (!language.empty()) url += "?language=" + url_encode(language);
    std::string resp;
    long code = http_post_raw(url, audio, "audio/wav", key_, &resp);
    if (code < 200 || code >= 300) throw http_error("transcribe", code, resp);
    return parse(resp, "transcribe").at("text").get<std::string>();
  }

  // ── 6. Converse ──────────────────────────────────────────────────────
  // POST {base}/audio/converse with raw audio. STT → Vible agent → device TTS.
  // Returns the transcript immediately; reply is spoken on the device.
  ConverseResult converse(const std::vector<uint8_t>& audio,
                          const std::string& language = "") {
    std::string url = base_ + "/audio/converse";
    if (!language.empty()) url += "?language=" + url_encode(language);
    std::string resp;
    long code = http_post_raw(url, audio, "audio/wav", key_, &resp);
    if (code < 200 || code >= 300) throw http_error("converse", code, resp);
    auto j = parse(resp, "converse");
    ConverseResult r;
    r.transcript = j.value("transcript", std::string());
    r.spoken = j.value("spoken", false);
    return r;
  }

  // ── 7. Call (cloud phone) ────────────────────────────────────────────
  // POST {base}/calls. Authenticated with the CALL key (NOT the vib- key):
  // calls are a separate hub capability that maps a call key → a device, so
  // the bearer here is call_key_, not key_. Returns 202 {call_id,...}.
  CallResult call(const std::string& prompt, const std::string& title = "",
                  const std::string& caller = "", const std::string& context = "") {
    if (call_key_.empty())
      throw Error("call: no call key (set VIBLE_CALL_KEY or run locally)", 0);
    nlohmann::json body{{"prompt", prompt}};
    if (!title.empty()) body["title"] = title;
    if (!caller.empty()) body["caller"] = caller;
    if (!context.empty()) body["context"] = context;
    std::string resp;
    long code = http_post_json(base_ + "/calls", body.dump(), call_key_, &resp,
                               nullptr);
    if (code < 200 || code >= 300) throw http_error("call", code, resp);
    auto j = parse(resp, "call");
    CallResult r;
    r.call_id = j.value("call_id", std::string());
    r.ringing = j.value("ringing", false);
    r.delivered = j.value("delivered", false);
    return r;
  }

 private:
  std::string base_, key_, host_, model_, call_key_;

  // ── config resolution ────────────────────────────────────────────────
  // 順序：環境變數行先；VIBLE_API_KEY 空 → 打本機 operator-only 端點
  // GET http://{HOST}/api/apikey 攞 {key,base,ws,model}（本機開發方便用，
  // 雲端唔會有呢條，所以攞唔到只當無聲失敗，留低環境變數預設值）。
  void resolve() {
    base_ = env_("VIBLE_BASE_URL", "http://127.0.0.1:7770/v1");
    key_ = env_("VIBLE_API_KEY", "");
    host_ = env_("VIBLE_HOST", "127.0.0.1:7770");
    model_ = env_("VIBLE_MODEL", "vible-agent");
    call_key_ = env_("VIBLE_CALL_KEY", "");

    if (key_.empty()) {
      std::string resp;
      long code = http_get("http://" + host_ + "/api/apikey", "", &resp);
      if (code >= 200 && code < 300) {
        try {
          auto j = nlohmann::json::parse(resp);
          key_ = j.value("key", key_);
          base_ = j.value("base", base_);
          model_ = j.value("model", model_);
        } catch (...) { /* 唔係 JSON / 唔通 → 沿用環境預設 */ }
      }
    }
    // call key 同樣可由本機 operator 端點補上（POST /api/calls/key → {key}）。
    if (call_key_.empty()) {
      std::string resp;
      long code = http_post_json("http://" + host_ + "/api/calls/key", "{}", "",
                                 &resp, nullptr);
      if (code >= 200 && code < 300) {
        try { call_key_ = nlohmann::json::parse(resp).value("key", std::string()); }
        catch (...) {}
      }
    }
  }

  static std::string env_(const char* name, const std::string& def) {
    const char* v = std::getenv(name);
    return (v && *v) ? std::string(v) : def;
  }

  nlohmann::json build_chat_body(const std::vector<Message>& messages,
                                 const ChatOptions& opt, bool stream) const {
    nlohmann::json msgs = nlohmann::json::array();
    for (const auto& m : messages)
      msgs.push_back({{"role", m.role}, {"content", m.content}});
    nlohmann::json body{
        {"model", opt.model.empty() ? model_ : opt.model},
        {"messages", msgs},
        {"stream", stream},
    };
    if (opt.temperature >= 0.0) body["temperature"] = opt.temperature;
    if (opt.max_tokens > 0) body["max_tokens"] = opt.max_tokens;
    return body;
  }

  static nlohmann::json parse(const std::string& s, const char* what) {
    try {
      return nlohmann::json::parse(s);
    } catch (const std::exception& e) {
      throw Error(std::string(what) + ": bad JSON response: " + e.what(), 0);
    }
  }

  static Error http_error(const char* what, long code, const std::string& body) {
    std::string msg = std::string(what) + ": HTTP " + std::to_string(code);
    // 盡量抽出 OpenAI 風格 {error:{message}} 畀人睇得明；抽唔到就附原文截斷。
    try {
      auto j = nlohmann::json::parse(body);
      if (j.contains("error") && j["error"].is_object() &&
          j["error"].contains("message"))
        msg += " — " + j["error"]["message"].get<std::string>();
      else if (j.contains("error") && j["error"].is_string())
        msg += " — " + j["error"].get<std::string>();
    } catch (...) {
      if (!body.empty()) msg += " — " + body.substr(0, 300);
    }
    return Error(msg, code);
  }

  // ── libcurl plumbing ─────────────────────────────────────────────────

  // 普通 string body 嘅 write callback（accumulate 落 std::string）。
  static size_t write_str(char* ptr, size_t size, size_t nmemb, void* ud) {
    auto* s = static_cast<std::string*>(ud);
    s->append(ptr, size * nmemb);
    return size * nmemb;
  }

  // binary body 嘅 write callback（音訊 bytes，絕不可經 std::string utf8 來回）。
  static size_t write_bytes(char* ptr, size_t size, size_t nmemb, void* ud) {
    auto* v = static_cast<std::vector<uint8_t>*>(ud);
    v->insert(v->end(), reinterpret_cast<uint8_t*>(ptr),
              reinterpret_cast<uint8_t*>(ptr) + size * nmemb);
    return size * nmemb;
  }

  // SSE 串流狀態 + write callback：逐 chunk 餵入、切完整 event、抽
  // choices[0].delta.content 回呼。
  struct SseState {
    std::string buffer;
    std::function<void(const std::string&)> on_delta;
    bool done = false;
  };
  static size_t write_sse(char* ptr, size_t size, size_t nmemb, void* ud) {
    auto* st = static_cast<SseState*>(ud);
    size_t n = size * nmemb;
    st->buffer.append(ptr, n);
    // 一個 SSE event 以空行（\n\n）作結；逐個完整 event 食出嚟。
    size_t pos;
    while ((pos = st->buffer.find("\n\n")) != std::string::npos) {
      std::string block = st->buffer.substr(0, pos);
      st->buffer.erase(0, pos + 2);
      size_t lstart = 0;
      while (lstart <= block.size()) {
        size_t lend = block.find('\n', lstart);
        std::string line = block.substr(
            lstart, lend == std::string::npos ? std::string::npos : lend - lstart);
        lstart = (lend == std::string::npos) ? block.size() + 1 : lend + 1;
        if (line.rfind("data:", 0) != 0) continue;
        std::string data = line.substr(5);
        // 去前導空白（OpenAI 慣例 "data: {json}"）。
        size_t b = data.find_first_not_of(" \t");
        if (b == std::string::npos) continue;
        data = data.substr(b);
        if (data == "[DONE]") { st->done = true; return n; }
        try {
          auto ev = nlohmann::json::parse(data);
          const auto& ch = ev.at("choices").at(0);
          if (ch.contains("delta") && ch["delta"].contains("content") &&
              ch["delta"]["content"].is_string())
            st->on_delta(ch["delta"]["content"].get<std::string>());
        } catch (...) { /* keep-alive / 非完整 / 非 token chunk → 略過 */ }
      }
    }
    return n;
  }

  struct CurlHandle {
    CURL* h;
    curl_slist* hdr = nullptr;
    CurlHandle() : h(curl_easy_init()) {
      if (!h) throw Error("curl_easy_init failed", 0);
    }
    ~CurlHandle() {
      if (hdr) curl_slist_free_all(hdr);
      if (h) curl_easy_cleanup(h);
    }
  };

  static void add_auth(CurlHandle& c, const std::string& key) {
    if (!key.empty()) c.hdr = curl_slist_append(c.hdr, ("Authorization: Bearer " + key).c_str());
  }

  long http_get(const std::string& url, const std::string& key,
                std::string* out) const {
    CurlHandle c;
    add_auth(c, key);
    curl_easy_setopt(c.h, CURLOPT_URL, url.c_str());
    curl_easy_setopt(c.h, CURLOPT_HTTPHEADER, c.hdr);
    curl_easy_setopt(c.h, CURLOPT_WRITEFUNCTION, write_str);
    curl_easy_setopt(c.h, CURLOPT_WRITEDATA, out);
    return perform(c, "GET " + url, nullptr);
  }

  long http_post_json(const std::string& url, const std::string& body,
                      const std::string& key, std::string* out, SseState* sse,
                      std::string* err = nullptr) const {
    CurlHandle c;
    c.hdr = curl_slist_append(c.hdr, "Content-Type: application/json");
    add_auth(c, key);
    curl_easy_setopt(c.h, CURLOPT_URL, url.c_str());
    curl_easy_setopt(c.h, CURLOPT_POST, 1L);
    curl_easy_setopt(c.h, CURLOPT_POSTFIELDS, body.c_str());
    curl_easy_setopt(c.h, CURLOPT_POSTFIELDSIZE, static_cast<long>(body.size()));
    curl_easy_setopt(c.h, CURLOPT_HTTPHEADER, c.hdr);
    if (sse) {
      curl_easy_setopt(c.h, CURLOPT_WRITEFUNCTION, write_sse);
      curl_easy_setopt(c.h, CURLOPT_WRITEDATA, sse);
    } else {
      curl_easy_setopt(c.h, CURLOPT_WRITEFUNCTION, write_str);
      curl_easy_setopt(c.h, CURLOPT_WRITEDATA, out);
    }
    return perform(c, "POST " + url, err);
  }

  // POST JSON → binary response (audio/wav bytes). 同時抽出回應 Content-Type，
  // 等 caller 分得到 wav vs speak:true 嗰陣回嘅 JSON。
  long http_post_json_bytes(const std::string& url, const std::string& body,
                            const std::string& key, std::vector<uint8_t>* out,
                            std::string* ctype) const {
    CurlHandle c;
    c.hdr = curl_slist_append(c.hdr, "Content-Type: application/json");
    add_auth(c, key);
    curl_easy_setopt(c.h, CURLOPT_URL, url.c_str());
    curl_easy_setopt(c.h, CURLOPT_POST, 1L);
    curl_easy_setopt(c.h, CURLOPT_POSTFIELDS, body.c_str());
    curl_easy_setopt(c.h, CURLOPT_POSTFIELDSIZE, static_cast<long>(body.size()));
    curl_easy_setopt(c.h, CURLOPT_HTTPHEADER, c.hdr);
    curl_easy_setopt(c.h, CURLOPT_WRITEFUNCTION, write_bytes);
    curl_easy_setopt(c.h, CURLOPT_WRITEDATA, out);
    long code = perform(c, "POST " + url, nullptr);
    char* ct = nullptr;
    if (curl_easy_getinfo(c.h, CURLINFO_CONTENT_TYPE, &ct) == CURLE_OK && ct)
      *ctype = ct;
    return code;
  }

  // POST 原始 audio bytes 做 body（STT / converse）：唔包 multipart，直接攞
  // raw buffer 上傳 + 指定 Content-Type（伺服器照單收 raw audio）。
  long http_post_raw(const std::string& url, const std::vector<uint8_t>& body,
                     const std::string& content_type, const std::string& key,
                     std::string* out) const {
    CurlHandle c;
    c.hdr = curl_slist_append(c.hdr, ("Content-Type: " + content_type).c_str());
    add_auth(c, key);
    curl_easy_setopt(c.h, CURLOPT_URL, url.c_str());
    curl_easy_setopt(c.h, CURLOPT_POST, 1L);
    curl_easy_setopt(c.h, CURLOPT_POSTFIELDS,
                     body.empty() ? "" : reinterpret_cast<const char*>(body.data()));
    curl_easy_setopt(c.h, CURLOPT_POSTFIELDSIZE, static_cast<long>(body.size()));
    curl_easy_setopt(c.h, CURLOPT_HTTPHEADER, c.hdr);
    curl_easy_setopt(c.h, CURLOPT_WRITEFUNCTION, write_str);
    curl_easy_setopt(c.h, CURLOPT_WRITEDATA, out);
    return perform(c, "POST " + url, nullptr);
  }

  // 統一 perform：執行請求、回 HTTP 狀態碼。傳輸層出錯時，若 caller 提供 err
  // 指標就寫低錯誤（畀串流路徑分辨「連唔上」vs「HTTP 錯」），否則直接拋。
  long perform(CurlHandle& c, const std::string& what, std::string* err) const {
    curl_easy_setopt(c.h, CURLOPT_USERAGENT, "vible-cpp/1.0");
    curl_easy_setopt(c.h, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(c.h, CURLOPT_CONNECTTIMEOUT, 15L);
    CURLcode rc = curl_easy_perform(c.h);
    if (rc != CURLE_OK) {
      std::string m = what + ": " + curl_easy_strerror(rc);
      if (err) { *err = m; return 0; }
      throw Error(m, 0);
    }
    long code = 0;
    curl_easy_getinfo(c.h, CURLINFO_RESPONSE_CODE, &code);
    return code;
  }

  static std::string url_encode(const std::string& s) {
    CURL* h = curl_easy_init();
    char* e = curl_easy_escape(h, s.c_str(), static_cast<int>(s.size()));
    std::string out = e ? e : s;
    if (e) curl_free(e);
    if (h) curl_easy_cleanup(h);
    return out;
  }
};

// libcurl 全域初始化 RAII helper（可選）。多執行緒程式應喺 main 開頭整一個。
struct Global {
  Global() { curl_global_init(CURL_GLOBAL_DEFAULT); }
  ~Global() { curl_global_cleanup(); }
};

}  // namespace vible

#endif  // VIBLE_HPP
