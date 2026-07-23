// Package vible is a thin, zero-dependency Go client for the Vible
// OpenAI-compatible API (chat, streaming, models, TTS, STT, converse,
// cloud phone calls, and a hand-rolled realtime WebSocket).
//
// The desktop app exposes /v1/* either locally (http://127.0.0.1:7770/v1)
// or in the cloud (https://<workspace>.api.vible.space/v1). Auth is a
// "vib-" bearer key, except cloud phone calls which use a separate call key.
package vible

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client talks to one Vible endpoint. Construct it with New(); the fields are
// exported so callers can override anything after resolution if they need to.
type Client struct {
	Base    string // {scheme}://host/v1
	Key     string // vib- bearer key
	WS      string // realtime websocket URL ({base} with http->ws, /v1->/v1/realtime)
	Model   string // default model for Chat/Stream
	Host    string // host:port used by the operator-only local endpoints
	CallKey string // separate bearer key for cloud phone /calls
	HTTP    *http.Client
}

// Message is one chat turn.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ConverseResult is what /audio/converse returns: the recognized transcript
// plus whether the agent spoke its reply on the device.
type ConverseResult struct {
	Transcript string `json:"transcript"`
	Spoken     bool   `json:"spoken"`
}

// CallResult is the 202 body from /calls.
type CallResult struct {
	CallID    string `json:"call_id"`
	Ringing   bool   `json:"ringing"`
	Delivered bool   `json:"delivered"`
}

// options holds the per-call knobs. All fields are optional; zero values mean
// "let the server decide". Callers set them through the Option helpers below.
type options struct {
	Model       string
	Temperature *float64
	TopP        *float64
	MaxTokens   *int
	// TTS-specific
	Voice string
	Speed *float64
	Speak bool
	// STT / converse
	Language string
	// Call-specific
	Title   string
	Caller  string
	Context string
}

// Option mutates a per-call options struct (functional-options pattern).
type Option func(*options)

// WithModel overrides the model for this call only.
func WithModel(m string) Option { return func(o *options) { o.Model = m } }

// WithTemperature sets sampling temperature.
func WithTemperature(t float64) Option { return func(o *options) { o.Temperature = &t } }

// WithTopP sets nucleus sampling.
func WithTopP(p float64) Option { return func(o *options) { o.TopP = &p } }

// WithMaxTokens caps the response length.
func WithMaxTokens(n int) Option { return func(o *options) { o.MaxTokens = &n } }

// WithVoice picks the TTS voice. See Speech for the mapping rules.
func WithVoice(v string) Option { return func(o *options) { o.Voice = v } }

// WithSpeed sets TTS playback speed (1.0 = normal).
func WithSpeed(s float64) Option { return func(o *options) { o.Speed = &s } }

// WithSpeak makes TTS play on the device and return JSON instead of WAV bytes.
func WithSpeak() Option { return func(o *options) { o.Speak = true } }

// WithLanguage hints the STT/converse language (e.g. "en", "zh").
func WithLanguage(l string) Option { return func(o *options) { o.Language = l } }

// WithTitle sets a call title.
func WithTitle(t string) Option { return func(o *options) { o.Title = t } }

// WithCaller sets the displayed caller for a call.
func WithCaller(c string) Option { return func(o *options) { o.Caller = c } }

// WithContext attaches extra context to a call.
func WithContext(c string) Option { return func(o *options) { o.Context = c } }

func applyOpts(opts []Option) options {
	var o options
	for _, fn := range opts {
		fn(&o)
	}
	return o
}

// New builds a Client by resolving configuration from the environment, falling
// back to the operator-only local /api/apikey endpoint when no key is set.
func New() (*Client, error) {
	c := &Client{
		Base:    envOr("VIBLE_BASE_URL", "http://127.0.0.1:7770/v1"),
		Key:     os.Getenv("VIBLE_API_KEY"),
		WS:      os.Getenv("VIBLE_WS_URL"),
		Model:   envOr("VIBLE_MODEL", "vible-agent"),
		Host:    envOr("VIBLE_HOST", "127.0.0.1:7770"),
		CallKey: os.Getenv("VIBLE_CALL_KEY"),
		HTTP:    &http.Client{Timeout: 120 * time.Second},
	}

	// 無 vib- key 時，向本機 operator-only 端點攞 key 同其它預設值。
	// 呢個端點只喺本機可用，所以只適合本地開發。
	if c.Key == "" {
		if err := c.fetchLocalKey(); err != nil {
			return nil, err
		}
	}

	// WS 預設值：由 base 推導 —— http->ws，並把尾段 /v1 換成 /v1/realtime。
	// 必須喺 fetchLocalKey 之後做，因為嗰度可能已經填咗 c.WS。
	if c.WS == "" {
		c.WS = deriveWS(c.Base)
	}
	return c, nil
}

func (c *Client) fetchLocalKey() error {
	u := fmt.Sprintf("http://%s/api/apikey", c.Host)
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return fmt.Errorf("vible: no VIBLE_API_KEY and local key fetch failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vible: local key endpoint %s returned %s", u, resp.Status)
	}
	var k struct {
		Key   string `json:"key"`
		Base  string `json:"base"`
		WS    string `json:"ws"`
		Model string `json:"model"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&k); err != nil {
		return fmt.Errorf("vible: decode local key: %w", err)
	}
	if k.Key == "" {
		return errors.New("vible: local key endpoint returned empty key")
	}
	c.Key = k.Key
	// 只喺 env 未提供時，先採用端點回傳嘅值，唔好覆蓋使用者明示嘅設定。
	if k.Base != "" && os.Getenv("VIBLE_BASE_URL") == "" {
		c.Base = k.Base
	}
	if k.WS != "" && os.Getenv("VIBLE_WS_URL") == "" {
		c.WS = k.WS
	}
	if k.Model != "" && os.Getenv("VIBLE_MODEL") == "" {
		c.Model = k.Model
	}
	return nil
}

// deriveWS turns an https?://.../v1 base into a ws(s)://.../v1/realtime URL.
func deriveWS(base string) string {
	ws := base
	if strings.HasPrefix(ws, "https://") {
		ws = "wss://" + strings.TrimPrefix(ws, "https://")
	} else if strings.HasPrefix(ws, "http://") {
		ws = "ws://" + strings.TrimPrefix(ws, "http://")
	}
	ws = strings.TrimRight(ws, "/")
	if strings.HasSuffix(ws, "/v1") {
		ws = strings.TrimSuffix(ws, "/v1") + "/v1/realtime"
	}
	return ws
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (c *Client) modelFor(o options) string {
	if o.Model != "" {
		return o.Model
	}
	return c.Model
}

// chatBody assembles the JSON request shared by Chat and ChatStream.
func (c *Client) chatBody(messages []Message, o options, stream bool) map[string]any {
	body := map[string]any{
		"model":    c.modelFor(o),
		"messages": messages,
		"stream":   stream,
	}
	if o.Temperature != nil {
		body["temperature"] = *o.Temperature
	}
	if o.TopP != nil {
		body["top_p"] = *o.TopP
	}
	if o.MaxTokens != nil {
		body["max_tokens"] = *o.MaxTokens
	}
	return body
}

func (c *Client) newJSONReq(method, path string, body any) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.Base+path, r)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	return req, nil
}

// readErr drains a non-2xx response into a useful error.
func readErr(resp *http.Response) error {
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("vible: %s -> %s: %s", resp.Request.URL.Path, resp.Status, strings.TrimSpace(string(b)))
}

// Chat sends a non-streaming chat completion and returns the assistant text.
func (c *Client) Chat(messages []Message, opts ...Option) (string, error) {
	o := applyOpts(opts)
	req, err := c.newJSONReq(http.MethodPost, "/chat/completions", c.chatBody(messages, o, false))
	if err != nil {
		return "", err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", readErr(resp)
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("vible: chat returned no choices")
	}
	return out.Choices[0].Message.Content, nil
}

// ChatStream streams a chat completion, invoking onDelta for each token.
func (c *Client) ChatStream(messages []Message, onDelta func(string), opts ...Option) error {
	o := applyOpts(opts)
	req, err := c.newJSONReq(http.MethodPost, "/chat/completions", c.chatBody(messages, o, true))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return readErr(resp)
	}

	// SSE：逐行讀，只理 "data: " 開頭嘅行；"data: [DONE]" 為結束標記。
	// 用 bufio.Scanner 並調大 buffer，因為單一 token 行偶爾會比預設 64KB 嘅
	// 行上限細好多，但個別 chunk 仍可能較長，加大係穩陣做法。
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "data:") {
			continue // 略過空白行同 SSE 嘅 event:/id: 等欄位
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			return nil
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // 容忍非預期/keep-alive 行，唔好因為一行壞咗就中斷成串
		}
		if len(chunk.Choices) > 0 {
			if d := chunk.Choices[0].Delta.Content; d != "" {
				onDelta(d)
			}
		}
	}
	return sc.Err()
}

// Models lists available model ids.
func (c *Client) Models() ([]string, error) {
	req, err := c.newJSONReq(http.MethodGet, "/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, readErr(resp)
	}
	var out struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(out.Data))
	for _, m := range out.Data {
		ids = append(ids, m.ID)
	}
	return ids, nil
}

// Speech synthesizes input to a WAV byte slice.
//
// Voice mapping (handled server-side, passed through verbatim here):
//   - OpenAI built-in names (alloy/echo/fable/onyx/nova/shimmer/...) or "" -> default voice
//   - "clone:NAME" -> local clone named NAME
//   - any other string -> a Vible speaker name
//
// If WithSpeak() is set the device plays the audio and the server returns JSON
// {spoken:true}; in that case the returned []byte is the raw JSON, not WAV.
func (c *Client) Speech(input string, opts ...Option) ([]byte, error) {
	o := applyOpts(opts)
	body := map[string]any{
		"model": "vible-tts",
		"input": input,
		"voice": o.Voice, // 空字串交畀伺服器當「預設聲」處理，唔喺客戶端硬塞名
	}
	if o.Speed != nil {
		body["speed"] = *o.Speed
	}
	if o.Language != "" {
		body["lang"] = o.Language
	}
	if o.Speak {
		body["speak"] = true
	}
	req, err := c.newJSONReq(http.MethodPost, "/audio/speech", body)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, readErr(resp)
	}
	// speak:true 會回 JSON 而唔係 WAV bytes；兩種情況都直接回傳原始 body，
	// 由呼叫者按 WithSpeak 決定點解讀。
	return io.ReadAll(resp.Body)
}

// Transcribe sends raw audio bytes (declared as audio/wav) and returns text.
func (c *Client) Transcribe(audio []byte, opts ...Option) (string, error) {
	o := applyOpts(opts)
	path := "/audio/transcriptions"
	if o.Language != "" {
		path += "?language=" + url.QueryEscape(o.Language)
	}
	req, err := http.NewRequest(http.MethodPost, c.Base+path, bytes.NewReader(audio))
	if err != nil {
		return "", err
	}
	// 端點收原始 audio body（唔係 multipart），所以 Content-Type 標 audio/wav。
	req.Header.Set("Content-Type", "audio/wav")
	req.Header.Set("Authorization", "Bearer "+c.Key)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", readErr(resp)
	}
	var out struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Text, nil
}

// Converse sends raw audio and the agent speaks its reply on the device. It
// returns the recognized transcript plus whether a reply was spoken.
func (c *Client) Converse(audio []byte, opts ...Option) (ConverseResult, error) {
	o := applyOpts(opts)
	path := "/audio/converse"
	if o.Language != "" {
		path += "?language=" + url.QueryEscape(o.Language)
	}
	req, err := http.NewRequest(http.MethodPost, c.Base+path, bytes.NewReader(audio))
	if err != nil {
		return ConverseResult{}, err
	}
	req.Header.Set("Content-Type", "audio/wav")
	req.Header.Set("Authorization", "Bearer "+c.Key)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return ConverseResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ConverseResult{}, readErr(resp)
	}
	var out ConverseResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ConverseResult{}, err
	}
	return out, nil
}

// Call places a cloud phone call. NOTE: this uses the call key, not the vib-
// key. If CallKey is empty it is fetched from the local operator endpoint.
func (c *Client) Call(prompt string, opts ...Option) (CallResult, error) {
	o := applyOpts(opts)
	if c.CallKey == "" {
		if err := c.fetchCallKey(); err != nil {
			return CallResult{}, err
		}
	}
	body := map[string]any{"prompt": prompt}
	if o.Title != "" {
		body["title"] = o.Title
	}
	if o.Caller != "" {
		body["caller"] = o.Caller
	}
	if o.Context != "" {
		body["context"] = o.Context
	}
	b, err := json.Marshal(body)
	if err != nil {
		return CallResult{}, err
	}
	req, err := http.NewRequest(http.MethodPost, c.Base+"/calls", bytes.NewReader(b))
	if err != nil {
		return CallResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// 重點：/calls 用 call key 做 Bearer，唔係 vib- key。
	req.Header.Set("Authorization", "Bearer "+c.CallKey)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return CallResult{}, err
	}
	defer resp.Body.Close()
	// 成功係 202 Accepted（電話響緊／已派發），唔係 200。
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return CallResult{}, readErr(resp)
	}
	var out CallResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return CallResult{}, err
	}
	return out, nil
}

func (c *Client) fetchCallKey() error {
	u := fmt.Sprintf("http://%s/api/calls/key", c.Host)
	resp, err := c.HTTP.Post(u, "application/json", nil)
	if err != nil {
		return fmt.Errorf("vible: no VIBLE_CALL_KEY and call key fetch failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vible: call key endpoint %s returned %s", u, resp.Status)
	}
	var k struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&k); err != nil {
		return fmt.Errorf("vible: decode call key: %w", err)
	}
	if k.Key == "" {
		return errors.New("vible: call key endpoint returned empty key")
	}
	c.CallKey = k.Key
	return nil
}

// ---------------------------------------------------------------------------
// Realtime: a minimal, hand-rolled RFC 6455 WebSocket client (zero deps).
// We only need: client-side handshake, send one masked text frame, read
// server text frames until response.completed / error.
// ---------------------------------------------------------------------------

// Realtime opens a websocket, sends the messages as one response.create, and
// invokes onDelta for each streamed token until the response completes.
func (c *Client) Realtime(messages []Message, onDelta func(string)) error {
	wsURL := c.WS + "?api_key=" + url.QueryEscape(c.Key)
	u, err := url.Parse(wsURL)
	if err != nil {
		return err
	}

	conn, err := wsDial(u, c.HTTP)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 送出 response.create：包住所有 messages，由伺服器產生回應。
	reqMsg := map[string]any{"type": "response.create", "messages": messages}
	payload, err := json.Marshal(reqMsg)
	if err != nil {
		return err
	}
	if err := wsWriteText(conn, payload); err != nil {
		return err
	}

	for {
		data, err := wsReadText(conn)
		if err != nil {
			return err
		}
		var ev struct {
			Type  string `json:"type"`
			Delta string `json:"delta"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(data, &ev); err != nil {
			continue // 略過唔識解嘅控制訊息
		}
		switch ev.Type {
		case "response.output_text.delta":
			if ev.Delta != "" {
				onDelta(ev.Delta)
			}
		case "response.completed":
			return nil
		case "error":
			if ev.Error != nil {
				return fmt.Errorf("vible realtime: %s", ev.Error.Message)
			}
			return errors.New("vible realtime: error")
		}
		// session.created 等其它事件直接略過，繼續讀。
	}
}

// wsDial performs the RFC 6455 opening handshake over a raw TCP/TLS conn.
func wsDial(u *url.URL, hc *http.Client) (net.Conn, error) {
	host := u.Host
	var conn net.Conn
	var err error
	switch u.Scheme {
	case "ws":
		if !strings.Contains(host, ":") {
			host += ":80"
		}
		conn, err = net.DialTimeout("tcp", host, 15*time.Second)
	case "wss":
		if !strings.Contains(host, ":") {
			host += ":443"
		}
		// 用 tls.Dial，ServerName 設為純主機名（去掉 port）以正確驗證憑證。
		conn, err = tlsDial(host)
	default:
		return nil, fmt.Errorf("vible realtime: unsupported scheme %q", u.Scheme)
	}
	if err != nil {
		return nil, err
	}

	// Sec-WebSocket-Key 係 16 bytes 隨機值嘅 base64；伺服器會用佢計回應雜湊。
	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		conn.Close()
		return nil, err
	}
	secKey := base64.StdEncoding.EncodeToString(keyBytes)

	reqPath := u.RequestURI()
	if reqPath == "" {
		reqPath = "/"
	}
	handshake := "GET " + reqPath + " HTTP/1.1\r\n" +
		"Host: " + u.Host + "\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + secKey + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n\r\n"
	if _, err := conn.Write([]byte(handshake)); err != nil {
		conn.Close()
		return nil, err
	}

	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, &http.Request{Method: "GET"})
	if err != nil {
		conn.Close()
		return nil, err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		conn.Close()
		return nil, fmt.Errorf("vible realtime: handshake failed: %s", resp.Status)
	}
	// 驗證 Sec-WebSocket-Accept = base64(sha1(key + GUID))，確認對端真係 WS 伺服器。
	if got := resp.Header.Get("Sec-WebSocket-Accept"); got != wsAccept(secKey) {
		conn.Close()
		return nil, errors.New("vible realtime: bad Sec-WebSocket-Accept")
	}
	// 若握手回應已 buffer 咗多餘 bytes，包成一個會先讀 buffer 嘅 conn。
	if br.Buffered() > 0 {
		peek, _ := br.Peek(br.Buffered())
		return &bufConn{Conn: conn, pre: append([]byte(nil), peek...)}, nil
	}
	return conn, nil
}

// wsAccept computes the RFC 6455 server accept value for a client key.
func wsAccept(key string) string {
	const guid = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.Sum([]byte(key + guid))
	return base64.StdEncoding.EncodeToString(h[:])
}

// bufConn lets us prepend bytes that bufio already read past during handshake.
type bufConn struct {
	net.Conn
	pre []byte
}

func (b *bufConn) Read(p []byte) (int, error) {
	if len(b.pre) > 0 {
		n := copy(p, b.pre)
		b.pre = b.pre[n:]
		return n, nil
	}
	return b.Conn.Read(p)
}

// wsWriteText writes a single masked text frame (opcode 0x1, FIN set).
// 客戶端送出嘅 frame 一定要 masked（RFC 6455 §5.3），否則伺服器會斷線。
func wsWriteText(conn net.Conn, payload []byte) error {
	var buf bytes.Buffer
	buf.WriteByte(0x81) // FIN(0x80) | text opcode(0x1)

	n := len(payload)
	switch {
	case n <= 125:
		buf.WriteByte(byte(n) | 0x80) // 設 MASK bit(0x80)
	case n <= 0xFFFF:
		buf.WriteByte(126 | 0x80)
		var ext [2]byte
		binary.BigEndian.PutUint16(ext[:], uint16(n))
		buf.Write(ext[:])
	default:
		buf.WriteByte(127 | 0x80)
		var ext [8]byte
		binary.BigEndian.PutUint64(ext[:], uint64(n))
		buf.Write(ext[:])
	}

	// 4-byte 隨機 mask key，逐 byte 同 payload XOR。
	var mask [4]byte
	if _, err := rand.Read(mask[:]); err != nil {
		return err
	}
	buf.Write(mask[:])
	masked := make([]byte, n)
	for i := 0; i < n; i++ {
		masked[i] = payload[i] ^ mask[i%4]
	}
	buf.Write(masked)

	_, err := conn.Write(buf.Bytes())
	return err
}

// wsReadText reads frames and returns the next text/binary payload, replying to
// pings/closes as needed. Reassembles fragmented frames.
func wsReadText(conn net.Conn) ([]byte, error) {
	var assembled []byte
	for {
		hdr := make([]byte, 2)
		if _, err := io.ReadFull(conn, hdr); err != nil {
			return nil, err
		}
		fin := hdr[0]&0x80 != 0
		opcode := hdr[0] & 0x0F
		masked := hdr[1]&0x80 != 0
		length := int(hdr[1] & 0x7F)

		switch length {
		case 126:
			ext := make([]byte, 2)
			if _, err := io.ReadFull(conn, ext); err != nil {
				return nil, err
			}
			length = int(binary.BigEndian.Uint16(ext))
		case 127:
			ext := make([]byte, 8)
			if _, err := io.ReadFull(conn, ext); err != nil {
				return nil, err
			}
			length = int(binary.BigEndian.Uint64(ext))
		}

		var maskKey [4]byte
		if masked { // 伺服器 frame 正常唔 mask，但為穩陣仍處理
			if _, err := io.ReadFull(conn, maskKey[:]); err != nil {
				return nil, err
			}
		}
		payload := make([]byte, length)
		if length > 0 {
			if _, err := io.ReadFull(conn, payload); err != nil {
				return nil, err
			}
			if masked {
				for i := range payload {
					payload[i] ^= maskKey[i%4]
				}
			}
		}

		switch opcode {
		case 0x0, 0x1, 0x2: // continuation / text / binary
			assembled = append(assembled, payload...)
			if fin {
				return assembled, nil
			}
		case 0x8: // close
			return nil, io.EOF
		case 0x9: // ping -> 回 pong（masked），維持連線
			_ = wsWriteControl(conn, 0xA, payload)
		case 0xA: // pong -> 忽略
		}
	}
}

// wsWriteControl writes a masked control frame (used for pong).
func wsWriteControl(conn net.Conn, opcode byte, payload []byte) error {
	var buf bytes.Buffer
	buf.WriteByte(0x80 | opcode)
	n := len(payload)
	buf.WriteByte(byte(n) | 0x80) // control payload <=125 by spec
	var mask [4]byte
	if _, err := rand.Read(mask[:]); err != nil {
		return err
	}
	buf.Write(mask[:])
	for i := 0; i < n; i++ {
		buf.WriteByte(payload[i] ^ mask[i%4])
	}
	_, err := conn.Write(buf.Bytes())
	return err
}
