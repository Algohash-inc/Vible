package vible

import (
	"crypto/tls"
	"net"
	"strings"
	"time"
)

// tlsDial opens a TLS connection for the wss:// realtime transport. ServerName
// 設為去 port 後嘅主機名，確保憑證 SAN 驗證正確。
func tlsDial(hostPort string) (net.Conn, error) {
	serverName := hostPort
	if i := strings.LastIndex(serverName, ":"); i >= 0 {
		serverName = serverName[:i]
	}
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	return tls.DialWithDialer(dialer, "tcp", hostPort, &tls.Config{ServerName: serverName})
}
