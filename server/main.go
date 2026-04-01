package main

import (
	"fmt"
	pb "github.com/naufilf/thermal-watchdog/proto/telemetry"
)

func main() {
	report := &pb.TelemetryReport{NodeId: "test-node-01"}
	fmt.Printf(("Bazel is working Node: %s\n"), report.NodeId)
}
