// Package finops implementa o Harness do agente Finops: verifica custo
// estimado da sub-parte contra o threshold configurado. Nunca escala
// sozinho — estourar o threshold é sempre decisão humana explícita.
package finops

import (
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents"
	"github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

var _ agents.Harness = (*Harness)(nil)

// Harness é o esqueleto do agente Finops: Plan/Act/Observe
// determinísticos, sem chamar um motor de IA ainda. Sem categoria de
// escalação própria — ver NeverBlocksAlone em Config().
type Harness struct{}

func New() *Harness { return &Harness{} }

func (h *Harness) Config() agents.AgentConfig {
	return agents.AgentConfig{
		ID: "finops",
		ReadScope: agents.ReadScope{
			OwnModuleFull:   true,
			CrossModuleMode: "policies_only",
		},
		MaxCyclesBeforeEscalate: 3,
		HardCapCycles:           3,
		NeverBlocksAlone:        true,
	}
}

func (h *Harness) Plan(ctx agents.AgentContext) (agents.Action, error) {
	return agents.Action{
		Kind:    "check_cost_threshold",
		Payload: map[string]any{"sub_issue_id": ctx.SubIssueID},
	}, nil
}

func (h *Harness) Act(ctx agents.AgentContext, action agents.Action) (agents.Observation, error) {
	return agents.Observation{Signature: "finops:within_budget"}, nil
}

func (h *Harness) Observe(ctx agents.AgentContext, obs agents.Observation) (orchestrator.ExitResult, error) {
	return orchestrator.ExitResult{
		Status:      orchestrator.ExitSuccess,
		AgentID:     "finops",
		SubIssueID:  ctx.SubIssueID,
		ParentIssue: ctx.ParentIssueID,
	}, nil
}
