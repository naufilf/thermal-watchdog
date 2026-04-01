package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/naufilf/thermal-watchdog/proto/telemetry"
)

// Implements gRPC interface that is auto-generated
type WatchdogServer struct {
	pb.UnimplementedThermalWatchdogServer
}

// The method gRPC auto calls when a node connects
func (s *WatchdogServer) StreamTelemetry(stream pb.ThermalWatchdog_StreamTelemetryServer) error {
	var messagesProcessed int32 = 0;

	// Loop that infinitely catches data
	for {
		// Blocks server until new message arrives from hardware
		report, err := stream.Recv()

		// Case 1: Hardware cleanly closed connection (End of File)
		if err == io.EOF {
			log.Printf("Stream cleanly closed. Messages processed: %d", messagesProcessed)
			// Sending final stream response back
			return stream.SendAndClose(&pb.StreamResponse{MessagesProcessed: messagesProcessed})
		}

		// Case 2: Network dropped, or the node crashed ungracefully
		if err != nil {
			log.Printf("Stream abruptly terminated with error: %v", err)
			return err
		}

		// Case 3: Succesfuly recieved a telemetry report
		messagesProcessed++
		log.Printf("[Node: %s] Temp: %.2f°C", report.NodeId, report.GpuTemperatureCelsius)

	} 
}

func main() {
	// Open standard TCP port for server to listen on
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen in on port: %v", err)
	}

	// gRPC server engine creation
	grpcServer := grpc.NewServer()

	// Register the custom Watchdog implementation with the engine
	pb.RegisterThermalWatchdogServer(grpcServer, &WatchdogServer{})

	log.Println("Thermal Watchdog listening on port 50051...")

	// Start serving (This blocks forever)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
