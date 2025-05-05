package store

import (
	"sync"
	"time"

	"nova-panel/pb"
)

type AgentStatus struct {
	Id         int32
	Host       *pb.HostInfo
	State      *pb.StateInfo
	LastActive time.Time
}

var (
	mu     sync.RWMutex
	agents = make(map[int32]*AgentStatus)
)

func UpdateStatus(status *pb.StatusRequest) {
	mu.Lock()
	defer mu.Unlock()
	agents[status.Id] = &AgentStatus{
		Id:         status.Id,
		Host:       status.Host,
		State:      status.State,
		LastActive: status.LastActive.AsTime(),
	}
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
