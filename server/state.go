package main

import(
	"sync"
	"time"
	"log"
)

// Respresents the current snapshot of a single TPU server in memory
type TempReading struct {
	LatestTemp float32
	LastSeen time.Time
}

// Maintains a sliding window of snapshots of a node
type NodeState struct {
	History []TempReading
}

// Holds global map of all nodes connected to ingestion service
type StateManager struct {
	// mutex that locks out all threads when a thread is actively writing
	mu sync.RWMutex
	// node = nodeId, NodeState (identifier for TPU and its state at a certain time)
	nodes map[string]*NodeState
}

// Constructer fxn to initalize map
func NewStateManager() *StateManager {
	return &StateManager{
		nodes: make(map[string]*NodeState),
	}
}

// Update method called by the gRPC stream handler whenever a new temp arrives
func (sm *StateManager) Update(nodeId string, temp float32) {
	// Lock nodes map
	sm.mu.Lock()
	// Ensures mutex will not stay locked forever
	defer sm.mu.Unlock()

	// Case 1: If node doesn't exist in the map yet, create it
	if _, exists:= sm.nodes[nodeId]; !exists {
		sm.nodes[nodeId] = &NodeState{}
	}

	// Case 2: Update node status
	history := sm.nodes[nodeId].History
	history = append(history, TempReading{LatestTemp: temp, LastSeen: time.Now()})

	// Remove snapshots older than 10 seconds from History
	cutoffIndex := 0
	for index, reading := range history {
		if (time.Since(reading.LastSeen) > 10 * time.Second) {
			// The reading is too old, the cutoff point moves up
			cutoffIndex = index + 1
		} else {
			// Hit the first time valid reading
			break
		}
	}

	// Update slice header stored in map to be same as our time adjusted slice header
	sm.nodes[nodeId].History = history[cutoffIndex:]
	
	// Logging
	log.Printf("Map stored node [%s], there are currently %v snapshots of the node", nodeId, len(sm.nodes[nodeId].History))

	// Critical overheatting error log
	if (calcAvg(sm.nodes[nodeId].History) > 90.0) {
		log.Printf("[CRITICAL] NODE [%s] IS OVERHEATING", nodeId)
	}
}

// Safely reads a node's current state
func (sm *StateManager) Get(nodeId string) (*NodeState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	state, exist := sm.nodes[nodeId]
	return state, exist
}

func calcAvg(s []TempReading) float32 {
	var totalTemp float32 = 0.0
	for _, temp := range s {
		totalTemp += temp.LatestTemp 
	}

	return totalTemp / float32(len(s))
}