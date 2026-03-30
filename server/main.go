package main

import (
	pb "github.com/naufilf/thermal-watchdog/proto/telemetry"
)

func main() {
	report := &pb.TelemetryReport{NodeId: "test-node-01"}
	fmf.Printf(("Bazel is working Node: %s\n"), report.NodeId)
}
