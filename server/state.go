package main

import(
	"sync"
	"time"
	"log"
)

// Respresents the current snapshot of a single TPU server in memory
type NodeState struct {
	LatestTemp float32
	LastSeen time.Time
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
	sm.nodes[nodeId].LatestTemp = temp
	sm.nodes[nodeId].LastSeen = time.Now()
	
	// Logging
	log.Printf("Map stored node [%s] with temp: %.2f and last seen: %v", nodeId, sm.nodes[nodeId].LatestTemp, sm.nodes[nodeId].LastSeen)
}

// Safely reads a node's current state
func (sm *StateManager) Get(nodeId string) (*NodeState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	state, exist := sm.nodes[nodeId]
	return state, exist
}