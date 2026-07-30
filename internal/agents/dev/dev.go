// Package dev implementa o Harness do agente Dev: implementa e builda
// uma sub-parte do fan-out (ex: "pagamentos-123-front").
package dev

import (
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents"
	"github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

var _ agents.Harness = (*Harness)(nil)

// Harness é o esqueleto do agente Dev: Plan/Act/Observe determinísticos,
// sem chamar um motor de IA ainda. Categorias de escalação:
// "no_progress_same_error", "budget_exceeded". O fan-out por sub-parte
// é responsabilidade de quem chama (uma instância de loop por
// sub-issue), não do Harness em si.
type Harness struct{}

func New() *Harness { return &Harness{} }

func (h *Harness) Config() agents.AgentConfig {
	return agents.AgentConfig{
		ID: "dev",
		ReadScope: agents.ReadScope{
			OwnModuleFull:   true,
			CrossModuleMode: "interface_summary",
		},
		MaxCyclesBeforeEscalate: 3,
		HardCapCycles:           5,
	}
}

func (h *Harness) Plan(ctx agents.AgentContext) (agents.Action, error) {
	return agents.Action{
		Kind:    "implement_and_build",
		Payload: map[string]any{"sub_issue_id": ctx.SubIssueID},
	}, nil
}

func (h *Harness) Act(ctx agents.AgentContext, action agents.Action) (agents.Observation, error) {
	return agents.Observation{Signature: "dev:build_ok"}, nil
}

func (h *Harness) Observe(ctx agents.AgentContext, obs agents.Observation) (orchestrator.ExitResult, error) {
	return orchestrator.ExitResult{
		Status:      orchestrator.ExitSuccess,
		AgentID:     "dev",
		SubIssueID:  ctx.SubIssueID,
		ParentIssue: ctx.ParentIssueID,
	}, nil
}
