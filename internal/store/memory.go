package store

import (
	"sync"
	"time"
)

type AgentStatus struct {
	AgentId      string
	CpuPercent   float64
	MemPercent   float64
	UploadKbps   float64
	DownloadKbps float64
	LastUpdated  int64
}

var (
	mu     sync.RWMutex
	agents = make(map[string]*AgentStatus)
)

func UpdateStatus(status *AgentStatus) {
	mu.Lock()
	defer mu.Unlock()
	status.LastUpdated = time.Now().Unix()
	agents[status.AgentId] = status
}

func GetAllAgents() []*AgentStatus {
	mu.RLock()
	defer mu.RUnlock()
	list := []*AgentStatus{}
	for _, a := range agents {
		list = append(list, a)
	}
	return list
}

func GetAgentCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(agents)
}
