package manage

import (
	"bufio"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

// The original coordinator and wire protocol remain untouched. This adapter
// adds an authenticated readiness/diagnostic report on the same TCP/UDP port.
type externalCoordinator struct {
	*ctCoordinator
	reportMu                      sync.Mutex
	issues                        map[string]string
	ready                         bool
	measured                      bool
	finalReport                   chan struct{}
	finalOnce                     sync.Once
	waves                         bool
	waveIndex, waveStart, waveEnd int
	waveReady                     chan struct{}
	waveJoined                    bool
}

func startExternalCoordinator(port int, token string) (*externalCoordinator, error) {
	base := &ctCoordinator{tok: token, joined: make(chan struct{}), fetched: make(chan struct{}), spoofArrived: make(chan int, 1), pmtu: make(chan int, 1)}
	var e error
	base.tcp, e = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		return nil, e
	}
	base.udp, e = net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if e != nil {
		base.tcp.Close()
		return nil, e
	}
	c := &externalCoordinator{ctCoordinator: base, issues: map[string]string{}, finalReport: make(chan struct{})}
	go func() {
		for {
			conn, e := base.tcp.Accept()
			if e != nil {
				return
			}
			go func() {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
				line, e := bufio.NewReader(io.LimitReader(conn, 64<<10)).ReadString('\n')
				if e != nil {
					return
				}
				host, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
				if reply := c.answer(line, host); reply != "" {
					_, _ = io.WriteString(conn, reply+"\n")
				}
			}()
		}
	}()
	go func() {
		buffer := make([]byte, 64<<10)
		for {
			n, from, e := base.udp.ReadFrom(buffer)
			if e != nil {
				return
			}
			host, _, _ := net.SplitHostPort(from.String())
			if reply := c.answer(string(buffer[:n]), host); reply != "" {
				_, _ = base.udp.WriteTo([]byte(reply+"\n"), from)
			}
		}
	}()
	return c, nil
}
func (c *externalCoordinator) answer(line, from string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 || subtle.ConstantTimeCompare([]byte(fields[1]), []byte(c.tok)) != 1 {
		return ""
	}
	if reply, handled := c.waveAnswer(fields); handled {
		return reply
	}
	if fields[0] == "ready" || fields[0] == "diagnostic" || fields[0] == "final" {
		if len(fields) != 3 {
			return ""
		}
		body, e := unGzipB64(fields[2])
		if e != nil || len(body) > 256<<10 {
			return ""
		}
		var issues map[string]string
		if json.Unmarshal(body, &issues) != nil || len(issues) > 256 {
			return ""
		}
		for name, why := range issues {
			if len(name) > 50 || len(why) > 1024 {
				return ""
			}
		}
		c.reportMu.Lock()
		if fields[0] == "final" && !c.measured {
			c.reportMu.Unlock()
			return "wait"
		}
		for name, why := range issues {
			c.issues[name] = why
		}
		if fields[0] == "ready" {
			c.ready = true
		}
		c.reportMu.Unlock()
		if fields[0] == "final" {
			c.finalOnce.Do(func() { close(c.finalReport) })
		}
		if fields[0] == "ready" {
			return c.ctCoordinator.answer("hello "+c.tok, from)
		}
		return "ok"
	}
	// A v1 peer cannot start probing before the v2 readiness report arrives.
	if fields[0] == "hello" {
		c.reportMu.Lock()
		ready := c.ready
		c.reportMu.Unlock()
		if !ready {
			return "update"
		}
	}
	if fields[0] == "result" {
		c.reportMu.Lock()
		measured := c.measured
		c.reportMu.Unlock()
		c.ctCoordinator.mu.Lock()
		pending := c.verdict == ""
		live := c.live
		c.ctCoordinator.mu.Unlock()
		if measured && pending && live != nil {
			b, _ := json.Marshal(live())
			return "measured " + gzipB64(b)
		}
	}
	return c.ctCoordinator.answer(line, from)
}
func (c *externalCoordinator) peerIssues() map[string]string {
	c.reportMu.Lock()
	defer c.reportMu.Unlock()
	copy := map[string]string{}
	for name, why := range c.issues {
		copy[name] = why
	}
	return copy
}
func externalIssueReport(issues map[string]string) string {
	bounded := map[string]string{}
	for name, why := range issues {
		if len(why) > 1024 {
			why = why[:1024]
		}
		bounded[name] = why
	}
	b, _ := json.Marshal(bounded)
	return gzipB64(b)
}

func (c *externalCoordinator) finishMeasurements() {
	c.reportMu.Lock()
	c.measured = true
	c.reportMu.Unlock()
}
