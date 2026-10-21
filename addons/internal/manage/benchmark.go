package manage

import (
	"fmt"

	"time"
)

// PathQuality is the measured state of the link to a peer.
type PathQuality struct {
	Target string

	Sent     int
	Received int

	Min    time.Duration
	Avg    time.Duration
	Max    time.Duration
	Jitter time.Duration // mean deviation from the average

	// Err is set when the probe could not run at all (an unusable address).
	Err error
}

// LossPercent returns the share of probes that never completed.
func (p PathQuality) LossPercent() float64 {
	if p.Sent == 0 {
		return 0
	}
	return float64(p.Sent-p.Received) / float64(p.Sent) * 100
}

// Usable reports whether enough probes came back to trust the numbers.
func (p PathQuality) Usable() bool { return p.Err == nil && p.Received >= 2 }

// FECPlan is a forward-error-correction ratio: for every Data packets the link
// carries Parity extra ones, so up to Parity losses in each group of Data+Parity
// are repaired without waiting a round trip for a retransmit. Both ends must run
// the same ratio.
type FECPlan struct {
	Data   int
	Parity int
	Why    string
}

// Set reports whether the plan carries a real ratio.
func (f FECPlan) Set() bool { return f.Data > 0 && f.Parity > 0 }

// Ratio renders the plan as "data:parity".
func (f FECPlan) Ratio() string { return fmt.Sprintf("%d:%d", f.Data, f.Parity) }

// RecommendFEC sizes the FEC ratio to the measured loss. This is the single
// setting that most decides how a KCP link feels on a lossy route: too little
// parity and lost packets still stall waiting for a retransmit, too much and the
// parity traffic itself congests the link.
//
// The rule is that parity/(data+parity) has to exceed the loss rate with
// headroom, because loss arrives in bursts rather than evenly, so each tier
// tolerates comfortably more than the loss that selects it. The ratios are the
// ones proven on this kind of route; the data shard count drops as parity climbs
// so a very lossy link is not also carrying huge FEC groups.
func RecommendFEC(q PathQuality) FECPlan {
	loss := q.LossPercent()
	switch {
	case !q.Usable():
		// Nothing trustworthy to size against: the gaming default, which every
		// KCP preset already carries.
		return FECPlan{10, 3, "link not measurable — the default gaming ratio"}
	case loss < 1:
		return FECPlan{20, 5, fmt.Sprintf("%.0f%% loss — a light 25%% parity is plenty on a clean route", loss)}
	case loss < 5:
		return FECPlan{10, 5, fmt.Sprintf("%.0f%% loss — 10:5 repairs bursts up to ~33%%", loss)}
	case loss < 15:
		return FECPlan{8, 8, fmt.Sprintf("%.0f%% loss — 8:8 repairs bursts up to ~50%%", loss)}
	default:
		return FECPlan{4, 8, fmt.Sprintf("%.0f%% loss — 4:8 trades 200%% overhead for repairing up to ~67%%", loss)}
	}
}

// KeepAlivePlan is a keepalive/heartbeat pair derived from a real measurement
// of the link rather than from a fixed default.
type KeepAlivePlan struct {
	KeepAlive int // seconds
	Heartbeat int // seconds
	Why       string
}

// RecommendKeepAlive derives the liveness timers from the measured link.
//
// The fixed defaults (75s / 40s) have to be safe for the worst link anyone
// might have, which makes them slow to notice a dead peer on a good one. Once
// the round trip and its variability are known, the timers can be tightened or
// loosened to match:
//
//   - The heartbeat must comfortably exceed the worst round trip seen, or a
//     slow-but-alive peer gets declared dead and the tunnel restarts for nothing.
//   - Beyond that, shorter is better: it is how quickly a genuinely dropped
//     tunnel is noticed and rebuilt.
func RecommendKeepAlive(q PathQuality) KeepAlivePlan {
	// Without a trustworthy measurement, keep the conservative defaults.
	if !q.Usable() {
		return KeepAlivePlan{KeepAlive: 75, Heartbeat: 40,
			Why: "the link could not be measured, so the safe defaults are kept"}
	}

	// Budget on the worst observed round trip plus the jitter, not the average:
	// the timer has to survive the bad samples, not the typical ones.
	worst := q.Max + q.Jitter
	if lossy := q.LossPercent() >= 2; lossy {
		// On a lossy link a heartbeat can itself be the packet that goes
		// missing, so leave room for a retry before declaring the peer gone.
		worst *= 2
	}

	heartbeat := int(worst.Seconds()) + 5
	switch {
	case heartbeat < 10:
		heartbeat = 10 // never hammer a healthy link
	case heartbeat > 60:
		heartbeat = 60 // past this, a drop takes too long to notice
	}
	keepalive := heartbeat * 2

	return KeepAlivePlan{
		KeepAlive: keepalive,
		Heartbeat: heartbeat,
		Why: fmt.Sprintf("worst round trip seen was %s with ±%s jitter and %.0f%% loss",
			shortDur(q.Max), shortDur(q.Jitter), q.LossPercent()),
	}
}

// shortDur formats a duration the way a person reads latency.
func shortDur(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}
