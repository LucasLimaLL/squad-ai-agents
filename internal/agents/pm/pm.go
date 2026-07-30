// Package pm implementa o Harness do agente PM: traduz o escopo do
// Analista em critérios de aceite testáveis.
package pm

import (
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents"
	"github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

var _ agents.Harness = (*Harness)(nil)

// Harness é o esqueleto do agente PM: Plan/Act/Observe determinísticos,
// sem chamar um motor de IA ainda. Categoria de escalação:
// "criteria_not_testable".
type Harness struct{}

func New() *Harness { return &Harness{} }

func (h *Harness) Config() agents.AgentConfig {
	return agents.AgentConfig{
		ID: "pm",
		ReadScope: agents.ReadScope{
			OwnModuleFull:   true,
			CrossModuleMode: "none",
		},
		MaxCyclesBeforeEscalate: 3,
		HardCapCycles:           5,
	}
}

func (h *Harness) Plan(ctx agents.AgentContext) (agents.Action, error) {
	return agents.Action{
		Kind:    "define_acceptance_criteria",
		Payload: map[string]any{"sub_issue_id": ctx.SubIssueID},
	}, nil
}

func (h *Harness) Act(ctx agents.AgentContext, action agents.Action) (agents.Observation, error) {
	return agents.Observation{Signature: "pm:criteria_defined"}, nil
}

func (h *Harness) Observe(ctx agents.AgentContext, obs agents.Observation) (orchestrator.ExitResult, error) {
	return orchestrator.ExitResult{
		Status:      orchestrator.ExitSuccess,
		AgentID:     "pm",
		SubIssueID:  ctx.SubIssueID,
		ParentIssue: ctx.ParentIssueID,
	}, nil
}
