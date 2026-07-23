# Vible C++ SDK

A header-only C++17 client for the Vible OpenAI-compatible `/v1/*` API, plus
runnable samples. Vible is a desktop app that exposes an OpenAI-style HTTP API
either locally (`http://127.0.0.1:7770/v1`) or in the cloud
(`https://<workspace>.api.vible.space/v1`). This SDK is a thin client over that
API.

## Requirements

- A C++17 compiler
- **libcurl** (system) — HTTP, SSE streaming, audio upload, calls
  - macOS: `brew install curl` (or use the system libcurl)
  - Debian/Ubuntu: `apt-get install libcurl4-openssl-dev`
- **nlohmann/json** — pulled automatically via CMake `FetchContent` (needs
  network access on the first configure)
- CMake ≥ 3.14

## Build

```sh
cd sdk/cpp
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build -j
```

This builds all samples into `build/`: `chat`, `stream`, `models`, `speech`,
`transcribe`, `converse`, `call`.

## Configuration

The client resolves config in its constructor, in this order:

| Source | Keys |
| --- | --- |
| Environment | `VIBLE_API_KEY`, `VIBLE_BASE_URL` (default `http://127.0.0.1:7770/v1`), `VIBLE_HOST` (default `127.0.0.1:7770`), `VIBLE_MODEL` (default `vible-agent`), `VIBLE_CALL_KEY` |
| Local app | If `VIBLE_API_KEY` is empty, `GET http://{HOST}/api/apikey` returns `{key, base, ws, model}` (operator-only loopback endpoint — handy for local dev). The call key is similarly auto-resolved from `POST http://{HOST}/api/calls/key`. |

So on the same machine as the Vible app, no env vars are needed. For the cloud,
set `VIBLE_API_KEY` and `VIBLE_BASE_URL`.

```cpp
#include "vible.hpp"

vible::Client c;                       // resolve from env / local app
// or:
vible::Client c("https://acme.api.vible.space/v1", "vib-…");
```

## Usage

```cpp
vible::Global g;  // curl_global_init / cleanup (make one in main)
vible::Client c;

// Chat (non-stream)
std::string reply = c.chat({{"system", "Be terse."}, {"user", "Hi"}});

// Chat (stream)
c.chat_stream({{"user", "Count to 5"}},
              [](const std::string& tok) { std::cout << tok << std::flush; });

// Models
for (auto& id : c.models()) std::cout << id << "\n";

// Speech (TTS) → wav bytes (or play on device)
auto wav = c.speech("Hello", {.voice = "clone:Felix"});
// vible::SpeechOptions{ voice, lang, speed, speak } — speak=true plays on device.

// Transcribe (STT)
std::string text = c.transcribe(wav_bytes, /*language=*/"en");

// Converse (audio in → agent → spoken reply on device)
auto cv = c.converse(wav_bytes);   // cv.transcript, cv.spoken

// Call (cloud phone — uses the CALL key, not the vib- key)
auto call = c.call("Your deploy is done", /*title=*/"CI");
```

All methods throw `vible::Error` (carrying `.status`, the HTTP status, or `0`
for a transport failure) on error.

### Voice mapping (`speech`)

The `voice` field is interpreted server-side:

- `""` or an OpenAI built-in name (`alloy`, `echo`, `fable`, …) → default voice
- `"clone:NAME"` → a local voice clone named `NAME`
- any other string → a Vible speaker name (e.g. `Vivian`)

## Realtime WebSocket — out of scope

The C++ SDK does **not** implement the realtime WebSocket endpoint
(`/v1/realtime`), since that needs a WebSocket library. Bring your own
(e.g. [Boost.Beast](https://github.com/boostorg/beast) or
[libwebsockets](https://libwebsockets.org/)). The protocol is:

- Connect to `ws://{HOST}/v1/realtime` (or the `wss://…` cloud URL), with the
  `Authorization: Bearer {vib-key}` header.
- Send `{"type":"response.create","messages":[{"role":"user","content":"…"}]}`.
- Receive `response.output_text.delta` events (token deltas),
  `response.completed` when done, and `error` on failure.

## API surface

```cpp
namespace vible {
  struct Message { std::string role, content; };
  struct ChatOptions   { std::string model; double temperature; int max_tokens; };
  struct SpeechOptions { std::string voice, lang; double speed; bool speak; };
  struct ConverseResult { std::string transcript; bool spoken; };
  struct CallResult     { std::string call_id; bool ringing, delivered; };
  struct Error : std::runtime_error { long status; };

  class Client {
    Client();                                  // resolve config
    Client(std::string base, std::string key);

    std::string chat(const std::vector<Message>&, const ChatOptions& = {});
    void chat_stream(const std::vector<Message>&,
                     std::function<void(const std::string&)> on_delta,
                     const ChatOptions& = {});
    std::vector<std::string> models();
    std::vector<uint8_t> speech(const std::string& input, const SpeechOptions& = {});
    std::string transcribe(const std::vector<uint8_t>& audio,
                           const std::string& language = "");
    ConverseResult converse(const std::vector<uint8_t>& audio,
                            const std::string& language = "");
    CallResult call(const std::string& prompt, const std::string& title = "",
                    const std::string& caller = "", const std::string& context = "");
  };
}
```

## Try it

```sh
./build/chat "Tell me a fact about steel"
./build/stream "Write a haiku"
./build/models
./build/speech "Hello from Vible"      # writes out.wav
./build/transcribe                     # reads out.wav
./build/converse                       # reads out.wav, agent replies aloud
./build/call "Build finished"          # needs a call key
```
