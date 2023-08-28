package server

import (
	"context"

	"k8s.io/klog/v2"
	pkgagent "sigs.k8s.io/apiserver-network-proxy/pkg/agent"
)

type AgentIdBackendManager struct {
	*DefaultBackendStorage
}

var _ BackendManager = &AgentIdBackendManager{}

func NewAgentIdBackendManager() *AgentIdBackendManager {
	return &AgentIdBackendManager{
		DefaultBackendStorage: NewDefaultBackendStorage(
			[]pkgagent.IdentifierType{pkgagent.UID},
		),
	}
}

func (dibm *AgentIdBackendManager) Backend(ctx context.Context) (Backend, error) {
	dibm.mu.RLock()
	defer dibm.mu.RUnlock()

	if len(dibm.backends) == 0 {
		return nil, &ErrNotFound{}
	}

	agentID := ctx.Value(agentId).(string)
	if agentID != "" {
		bes, exist := dibm.backends[agentID]
		if exist && len(bes) > 0 {
			klog.V(4).InfoS("Get the backend through the AgentIdBackendManager", "agentID", agentID)
			return dibm.backends[agentID][0], nil
		}
	}

	klog.V(4).Infof("Get the backend through the AgentIdBackendManager failed", "agentID", agentID)
	return nil, &ErrNotFound{}
}
