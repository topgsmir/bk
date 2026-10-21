package manage

import (
	"testing"
	"time"
)

func TestPackage3WaveHandshakeRejectsStaleReports(t *testing.T) {
	c := &externalCoordinator{ctCoordinator: &ctCoordinator{tok: "fixture-token"}, issues: map[string]string{}}
	c.enableWaves()
	ready := c.setWave(0, 0, 4)
	if got := c.answer("wave wrong-token", "127.0.0.1"); got != "" {
		t.Fatal("unauthenticated wave read")
	}
	if got := c.answer("wave fixture-token", "127.0.0.1"); got != "wave 0 0 4" {
		t.Fatal(got)
	}
	issue := externalIssueReport(map[string]string{"d3-tcp-tcp": "fixture failure"})
	if got := c.answer("wave-ready fixture-token 0 "+issue, "127.0.0.1"); got != "ok" {
		t.Fatal(got)
	}
	select {
	case <-ready:
	default:
		t.Fatal("peer readiness not delivered")
	}
	if got := c.answer("wave-final fixture-token 0 "+issue, "127.0.0.1"); got != "wait" {
		t.Fatal("early final accepted", got)
	}
	c.finishMeasurements()
	if got := c.answer("wave-final fixture-token 0 "+issue, "127.0.0.1"); got != "ok" {
		t.Fatal(got)
	}
	next := c.setWave(1, 4, 8)
	if got := c.answer("wave-ready fixture-token 0 "+issue, "127.0.0.1"); got != "wait" {
		t.Fatal("stale readiness accepted")
	}
	if got := c.answer("wave-final fixture-token 0 "+issue, "127.0.0.1"); got != "wait" {
		t.Fatal("stale final accepted")
	}
	select {
	case <-next:
		t.Fatal("stale report advanced next wave")
	default:
	}
	if c.peerIssues()["d3-tcp-tcp"] != "fixture failure" {
		t.Fatal("lost earlier diagnostic")
	}
}

func TestPackage3PlanDeadlineFitsPeerAcceptanceWindow(t *testing.T) {
	s, err := startExternalTestReadyIPv6("127.0.0.1", "127.0.0.2", []string{"d3-tcp"}, nil, nil, nil, [2]string{})
	if err != nil {
		t.Fatal(err)
	}
	defer s.close()
	until := time.Unix(s.plan.Until, 0)
	if until.After(time.Now().Add(2 * time.Hour)) {
		t.Fatal("package 3 plan would be rejected as too far in the future")
	}
	if until.Before(time.Now().Add(externalJoinWait + 70*time.Minute)) {
		t.Fatal("full package 3 matrix expires too early")
	}
}
