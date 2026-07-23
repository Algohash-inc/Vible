# Vible Go SDK

A thin, **zero-dependency** Go client for the Vible OpenAI-compatible API:
chat (sync + streaming), models, text-to-speech, speech-to-text, voice
converse, cloud phone calls, and a realtime WebSocket. Standard library only.

## Install

```bash
go get github.com/Algohash-inc/Vible/sdk/go
```

```go
import vible "github.com/Algohash-inc/Vible/sdk/go"
```

## Configuration

`vible.New()` resolves configuration in this order:

| Setting   | Env var          | Default                       |
|-----------|------------------|-------------------------------|
| API key   | `VIBLE_API_KEY`  | (fetched from local endpoint) |
| Base URL  | `VIBLE_BASE_URL` | `http://127.0.0.1:7770/v1`    |
| WS URL    | `VIBLE_WS_URL`   | base with `http→ws`, `/v1→/v1/realtime` |
| Host      | `VIBLE_HOST`     | `127.0.0.1:7770`              |
| Model     | `VIBLE_MODEL`    | `vible-agent`                 |
| Call key  | `VIBLE_CALL_KEY` | (fetched from local endpoint) |

If `VIBLE_API_KEY` is empty, the client calls the operator-only local endpoint
`GET http://{HOST}/api/apikey`, which returns `{key, base, ws, model}`. This is
only reachable on the local machine, so it is fine for local development but not
for the cloud (`https://<workspace>.api.vible.space/v1`), where you must set
`VIBLE_API_KEY` (and `VIBLE_BASE_URL`) yourself.

Auth is a `vib-` bearer key for everything **except** cloud phone calls, which
use a separate **call key** (`VIBLE_CALL_KEY`, or fetched from
`POST http://{HOST}/api/calls/key`).

## Quick start

```go
c, err := vible.New()
if err != nil { log.Fatal(err) }

reply, err := c.Chat([]vible.Message{
    {Role: "user", Content: "Hello!"},
})
fmt.Println(reply)
```

## API

```go
func New() (*Client, error)

func (c *Client) Chat(messages []Message, opts ...Option) (string, error)
func (c *Client) ChatStream(messages []Message, onDelta func(string), opts ...Option) error
func (c *Client) Models() ([]string, error)
func (c *Client) Speech(input string, opts ...Option) ([]byte, error)
func (c *Client) Transcribe(audio []byte, opts ...Option) (string, error)
func (c *Client) Converse(audio []byte, opts ...Option) (ConverseResult, error)
func (c *Client) Call(prompt string, opts ...Option) (CallResult, error)
func (c *Client) Realtime(messages []Message, onDelta func(string)) error
```

Options (functional): `WithModel`, `WithTemperature`, `WithTopP`,
`WithMaxTokens`, `WithVoice`, `WithSpeed`, `WithSpeak`, `WithLanguage`,
`WithTitle`, `WithCaller`, `WithContext`.

### Voice mapping (TTS)

`Speech` passes `voice` through to the server, which interprets it as:

- OpenAI built-in name (`alloy`, `echo`, `fable`, `onyx`, `nova`, `shimmer`, …)
  or empty → **default voice**
- `clone:NAME` → a **local clone** named `NAME`
- any other string → a **Vible speaker** name

With `WithSpeak()` the device plays the audio and `Speech` returns the JSON
`{"spoken":true}` as bytes instead of WAV.

## Examples

Each example is `go run .`-able:

```bash
go run ./examples/chat
go run ./examples/stream
go run ./examples/models
go run ./examples/speech        # writes out.wav
go run ./examples/transcribe    # reads out.wav (or: go run ./examples/transcribe path.wav)
go run ./examples/converse      # reads out.wav
go run ./examples/call
go run ./examples/realtime
```

Run `speech` first so `transcribe` and `converse` have an `out.wav` to read.

## Realtime WebSocket

The realtime client is a small **hand-rolled RFC 6455** implementation (no
third-party WS library). It performs the opening handshake, sends one masked
client text frame containing `response.create`, and reads server frames,
forwarding `response.output_text.delta` tokens to your callback until
`response.completed` (or surfacing `error`).
