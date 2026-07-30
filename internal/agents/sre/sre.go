// Package sre implementa o Harness do agente SRE: revisa prontidão
// operacional (observabilidade, rollback, dependências de infra) da
// sub-parte.
package sre

import (
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents"
	"github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

var _ agents.Harness = (*Harness)(nil)

// Harness é o esqueleto do agente SRE: Plan/Act/Observe determinísticos,
// sem chamar um motor de IA ainda. Categoria de escalação:
// "operational_item_unresolved".
type Harness struct{}

func New() *Harness { return &Harness{} }

func (h *Harness) Config() agents.AgentConfig {
	return agents.AgentConfig{
		ID: "sre",
		ReadScope: agents.ReadScope{
			OwnModuleFull:   true,
			CrossModuleMode: "infra_deps",
		},
		MaxCyclesBeforeEscalate: 3,
		HardCapCycles:           5,
	}
}

func (h *Harness) Plan(ctx agents.AgentContext) (agents.Action, error) {
	return agents.Action{
		Kind:    "review_operational_readiness",
		Payload: map[string]any{"sub_issue_id": ctx.SubIssueID},
	}, nil
}

func (h *Harness) Act(ctx agents.AgentContext, action agents.Action) (agents.Observation, error) {
	return agents.Observation{Signature: "sre:operational_ready"}, nil
}

func (h *Harness) Observe(ctx agents.AgentContext, obs agents.Observation) (orchestrator.ExitResult, error) {
	return orchestrator.ExitResult{
		Status:      orchestrator.ExitSuccess,
		AgentID:     "sre",
		SubIssueID:  ctx.SubIssueID,
		ParentIssue: ctx.ParentIssueID,
	}, nil
}
