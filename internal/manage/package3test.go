package manage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/topgsmir/bk/internal/externaltunnel"
	"io"
	"strings"
	"sync"
	"time"
)

// The full matrix runs four forwarding cases at a time. Both ends retire each
// previous wave before advancing; a small VPS does not start 88 cores at once.
func (c *externalCoordinator) enableWaves() {
	c.reportMu.Lock()
	defer c.reportMu.Unlock()
	c.waves = true
	c.waveIndex = -1
}
func (c *externalCoordinator) setWave(index, start, end int) <-chan struct{} {
	c.reportMu.Lock()
	defer c.reportMu.Unlock()
	c.waveIndex, c.waveStart, c.waveEnd = index, start, end
	c.waveReady, c.waveJoined = make(chan struct{}), false
	c.measured = false
	c.finalReport, c.finalOnce = make(chan struct{}), sync.Once{}
	return c.waveReady
}
func (c *externalCoordinator) waveAnswer(f []string) (string, bool) {
	if f[0] != "wave" && f[0] != "wave-ready" && f[0] != "wave-final" {
		return "", false
	}
	c.reportMu.Lock()
	defer c.reportMu.Unlock()
	if !c.waves {
		return "", true
	}
	if f[0] == "wave" {
		if len(f) != 2 {
			return "", true
		}
		return fmt.Sprintf("wave %d %d %d", c.waveIndex, c.waveStart, c.waveEnd), true
	}
	if len(f) != 4 {
		return "", true
	}
	var index int
	if _, e := fmt.Sscan(f[2], &index); e != nil || index != c.waveIndex {
		return "wait", true
	}
	if f[0] == "wave-final" && !c.measured {
		return "wait", true
	}
	b, e := unGzipB64(f[3])
	if e != nil || len(b) > 256<<10 {
		return "", true
	}
	var issues map[string]string
	if json.Unmarshal(b, &issues) != nil || len(issues) > 256 {
		return "", true
	}
	for name, why := range issues {
		if len(name) > 50 || len(why) > 1024 {
			return "", true
		}
	}
	for name, why := range issues {
		c.issues[name] = why
	}
	if f[0] == "wave-ready" && !c.waveJoined {
		c.waveJoined = true
		close(c.waveReady)
	}
	if f[0] == "wave-final" {
		c.finalOnce.Do(func() { close(c.finalReport) })
	}
	return "ok", true
}
func (s *externalTest) runPackage3Waves(ctx context.Context, progress func(int, ConnTestResult)) []ConnTestResult {
	rows := make([]ConnTestResult, len(s.cases))
	for i, c := range s.cases {
		rows[i] = ConnTestResult{Kind: c.kind, Transport: c.tr, Status: ctTesting, Total: connTestSoak}
		if progress != nil {
			progress(i, rows[i])
		}
	}
	for start, wave := 0, 0; start < len(s.cases); wave = wave + 1 {
		end := start + s.plan.BatchSize
		if end > len(s.cases) {
			end = len(s.cases)
		}
		for i := start; i < end; i++ {
			if externaltunnel.L2TPIPsec(s.plan.Cases[i].Kind) {
				if i == start {
					end = start + 1
				} else {
					end = i
				}
				break
			}
		}
		part := &externalTest{dir: s.dir, plan: s.plan, coord: s.coord, cases: s.cases[start:end]}
		ready := s.coord.setWave(wave, start, end)
		for i, c := range part.cases {
			if c.skip != "" {
				continue
			}
			spec := s.plan.Cases[start+i]
			if e := externaltunnel.Check(spec); e != nil {
				c.skip = e.Error()
				continue
			}
			engine, e := startExternalEngine(s.dir, spec)
			if e != nil {
				c.skip = e.Error()
				continue
			}
			part.engines = append(part.engines, engine)
			s.engines = append(s.engines, engine)
		}
		timer := time.NewTimer(2 * time.Minute)
		select {
		case <-ready:
		case <-ctx.Done():
			for _, c := range part.cases {
				if c.skip == "" {
					c.skip = ctx.Err().Error()
				}
			}
		case <-timer.C:
			for _, c := range part.cases {
				if c.skip == "" {
					c.skip = "peer did not prepare this test wave"
				}
			}
		}
		timer.Stop()
		update := func(i int, r ConnTestResult) {
			if progress != nil {
				progress(start+i, r)
			}
		}
		measured := part.runMeasurements(ctx, update, false)
		copy(rows[start:end], measured)
		stopAdditionalEngines(part.engines)
		start = end
	}
	s.best = ctComputeBest(rows, s.cases, 0, s.dir)
	s.coord.publish(rows, s.best)
	return rows
}
func runPackage3Peer(ctx context.Context, a connTestAddr, p externalPlan, dir string, preparation map[string]string, out io.Writer, live func([]ConnTestResult)) ([]ConnTestResult, ConnTestBest, error) {
	if p.BatchSize != 4 {
		return nil, ConnTestBest{}, fmt.Errorf("incompatible package 3 batching")
	}
	local := map[string]string{}
	var engines []*ctEngine
	defer func() { stopAdditionalEngines(engines) }()
	if _, e := ctAsk(a.Host, a.Coord, "ready "+a.Tok+" "+externalIssueReport(local)); e != nil {
		return nil, ConnTestBest{}, e
	}
	wave := -1
	for ctx.Err() == nil {
		r, e := ctAsk(a.Host, a.Coord, "wave "+a.Tok)
		if e == nil && strings.HasPrefix(r, "wave ") {
			var index, start, end int
			if _, e = fmt.Sscanf(r, "wave %d %d %d", &index, &start, &end); e != nil {
				return nil, ConnTestBest{}, e
			}
			if index >= 0 && index != wave {
				if index != wave+1 || start < 0 || end > len(p.Cases) || end <= start || end-start > p.BatchSize {
					return nil, ConnTestBest{}, fmt.Errorf("invalid package 3 wave")
				}
				stopAdditionalEngines(engines)
				engines = nil
				for _, original := range p.Cases[start:end] {
					tr := externalTransport(original)
					spec := original.Mirror()
					spec.Secret = ctCaseToken(p.Token, "extra", tr) + ctCaseToken(p.Token, "extra-key", tr)
					if why := preparation[spec.Kind]; why != "" {
						local[tr] = why
						continue
					}
					if e := externaltunnel.Check(spec); e != nil {
						local[tr] = e.Error()
						continue
					}
					engine, e := startExternalEngine(dir, spec)
					if e != nil {
						local[tr] = e.Error()
						continue
					}
					engines = append(engines, engine)
				}
				fmt.Fprintf(out, "Package 3: testing cases %d–%d of %d.\n", start+1, end, len(p.Cases))
				wave = index
			}
			if wave >= 0 {
				_, _ = ctAsk(a.Host, a.Coord, fmt.Sprintf("wave-ready %s %d %s", a.Tok, wave, externalIssueReport(local)))
			}
		}
		r, e = ctAsk(a.Host, a.Coord, "result "+a.Tok)
		if e == nil && (strings.HasPrefix(r, "live ") || strings.HasPrefix(r, "measured ")) {
			prefix := "live "
			measured := strings.HasPrefix(r, "measured ")
			if measured {
				prefix = "measured "
			}
			b, err := unGzipB64(strings.TrimPrefix(r, prefix))
			var rows []ConnTestResult
			if err == nil && json.Unmarshal(b, &rows) == nil {
				if live != nil {
					live(rows)
				}
				if measured {
					for _, row := range rows {
						if row.Status != ctOK {
							for _, engine := range engines {
								if strings.HasSuffix(engine.name, "-"+row.Transport) {
									if why := externalEngineDiagnostic(engine); why != "" {
										local[row.Transport] = why
									}
								}
							}
						}
					}
					_, _ = ctAsk(a.Host, a.Coord, fmt.Sprintf("wave-final %s %d %s", a.Tok, wave, externalIssueReport(local)))
				}
			}
		}
		if e == nil && strings.HasPrefix(r, "done ") {
			b, e := unGzipB64(strings.TrimPrefix(r, "done "))
			if e != nil {
				return nil, ConnTestBest{}, e
			}
			var verdict ctVerdict
			if e := json.Unmarshal(b, &verdict); e != nil {
				return nil, ConnTestBest{}, e
			}
			return verdict.Results, verdict.Best, nil
		}
		ctSleep(ctx, time.Second)
	}
	return nil, ConnTestBest{}, ctx.Err()
}
