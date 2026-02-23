package http

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type realtimeInboundMessage struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type realtimeEvent struct {
	Type          string `json:"type"`
	Author        string `json:"author,omitempty"`
	Text          string `json:"text,omitempty"`
	TurnComplete  bool   `json:"turnComplete,omitempty"`
	Interrupted   bool   `json:"interrupted,omitempty"`
	InputMimeType string `json:"inputMimeType,omitempty"`
}

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func (h *Handler) handleRealtimeWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/ws/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.NotFound(w, r)
		return
	}
	h.handleRealtimeWS(w, r, parts[0], parts[1])
}

func (h *Handler) handleRealtimeWebSocketSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	h.handleRealtimeWS(w, r, "session-client", sessionID)
}

func (h *Handler) handleRealtimeWS(w http.ResponseWriter, r *http.Request, _ string, sessionID string) {
	rw, conn, err := upgradeToWebSocket(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	if _, err := h.app.GetSession(r.Context(), sessionID); err != nil {
		_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "session not found"})
		_ = rw.Flush()
		return
	}

	_ = writeWSTextJSON(rw, realtimeEvent{Type: "ready"})
	_ = rw.Flush()

	for {
		opcode, payload, err := readWSFrame(rw.Reader)
		if err != nil {
			return
		}

		switch opcode {
		case 0x2: // binary
			_ = payload
			_ = writeWSTextJSON(rw, realtimeEvent{Type: "audio_ack", InputMimeType: "audio/pcm"})
			_ = rw.Flush()
			continue
		case 0x8: // close
			_ = writeWSClose(rw, "session closed")
			_ = rw.Flush()
			return
		case 0x1: // text
		default:
			_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "unsupported websocket opcode"})
			_ = rw.Flush()
			continue
		}

		var message realtimeInboundMessage
		if err := json.Unmarshal(payload, &message); err != nil {
			_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "invalid JSON message"})
			_ = rw.Flush()
			continue
		}

		switch message.Type {
		case "text":
			reply, err := h.app.SendMessage(context.Background(), sessionID, message.Text)
			if err != nil {
				_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": err.Error()})
				_ = rw.Flush()
				continue
			}
			_ = writeWSTextJSON(rw, realtimeEvent{Type: "event", Author: "assistant", Text: reply, TurnComplete: true})
		case "audio":
			if message.Data == "" {
				_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "audio data is required"})
				_ = rw.Flush()
				continue
			}
			if _, err := base64.StdEncoding.DecodeString(message.Data); err != nil {
				_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "invalid base64 audio payload"})
				_ = rw.Flush()
				continue
			}
			_ = writeWSTextJSON(rw, realtimeEvent{Type: "audio_ack", InputMimeType: "audio/pcm"})
		case "image":
			if message.Data == "" {
				_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "image data is required"})
				_ = rw.Flush()
				continue
			}
			if _, err := base64.StdEncoding.DecodeString(message.Data); err != nil {
				_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "invalid base64 image payload"})
				_ = rw.Flush()
				continue
			}
			_ = writeWSTextJSON(rw, realtimeEvent{Type: "image_ack"})
		case "activity_start":
			if !getenvBool("ENABLE_MANUAL_ACTIVITY_SIGNALS", false) {
				_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "activity_start ignored: manual activity signals disabled"})
				_ = rw.Flush()
				continue
			}
			_ = writeWSTextJSON(rw, realtimeEvent{Type: "activity_start_ack"})
		case "activity_end":
			if !getenvBool("ENABLE_MANUAL_ACTIVITY_SIGNALS", false) {
				_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "activity_end ignored: manual activity signals disabled"})
				_ = rw.Flush()
				continue
			}
			_ = writeWSTextJSON(rw, realtimeEvent{Type: "activity_end_ack", TurnComplete: true})
		case "close":
			_ = writeWSClose(rw, "session closed")
			_ = rw.Flush()
			return
		default:
			_ = writeWSTextJSON(rw, map[string]any{"type": "error", "errorMessage": "unsupported websocket message type"})
		}
		_ = rw.Flush()
	}
}

func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*bufio.ReadWriter, net.Conn, error) {
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return nil, nil, errors.New("missing websocket upgrade header")
	}
	if !strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") {
		return nil, nil, errors.New("missing connection upgrade header")
	}
	key := strings.TrimSpace(r.Header.Get("Sec-WebSocket-Key"))
	if key == "" {
		return nil, nil, errors.New("missing Sec-WebSocket-Key")
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("websocket unsupported")
	}
	conn, rw, err := hijacker.Hijack()
	if err != nil {
		return nil, nil, err
	}

	h := sha1.Sum([]byte(key + wsGUID))
	accept := base64.StdEncoding.EncodeToString(h[:])
	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", accept)
	if _, err := rw.WriteString(response); err != nil {
		conn.Close()
		return nil, nil, err
	}
	if err := rw.Flush(); err != nil {
		conn.Close()
		return nil, nil, err
	}
	return rw, conn, nil
}

func readWSFrame(r *bufio.Reader) (byte, []byte, error) {
	head := make([]byte, 2)
	if _, err := io.ReadFull(r, head); err != nil {
		return 0, nil, err
	}
	opcode := head[0] & 0x0F
	masked := (head[1] & 0x80) != 0
	if !masked {
		return 0, nil, errors.New("client frames must be masked")
	}

	payloadLen := int(head[1] & 0x7F)
	switch payloadLen {
	case 126:
		ext := make([]byte, 2)
		if _, err := io.ReadFull(r, ext); err != nil {
			return 0, nil, err
		}
		payloadLen = int(binary.BigEndian.Uint16(ext))
	case 127:
		ext := make([]byte, 8)
		if _, err := io.ReadFull(r, ext); err != nil {
			return 0, nil, err
		}
		payloadLen = int(binary.BigEndian.Uint64(ext))
	}

	mask := make([]byte, 4)
	if _, err := io.ReadFull(r, mask); err != nil {
		return 0, nil, err
	}
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, nil, err
	}
	for i := range payload {
		payload[i] ^= mask[i%4]
	}
	return opcode, payload, nil
}

func writeWSFrame(w *bufio.Writer, opcode byte, payload []byte) error {
	if len(payload) > 125 {
		return errors.New("payload too large for basic websocket frame")
	}
	header := []byte{0x80 | opcode, byte(len(payload))}
	if _, err := w.Write(header); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

func writeWSTextJSON(w *bufio.ReadWriter, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return writeWSFrame(w.Writer, 0x1, payload)
}

func writeWSClose(w *bufio.ReadWriter, reason string) error {
	payload := make([]byte, 2+len(reason))
	binary.BigEndian.PutUint16(payload[:2], 1000)
	copy(payload[2:], reason)
	if len(payload) > 125 {
		payload = payload[:125]
	}
	if err := writeWSFrame(w.Writer, 0x8, payload); err != nil {
		return err
	}
	return w.Flush()
}
