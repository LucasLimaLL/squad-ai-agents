// Package auditoria implementa o Harness do agente Auditoria: satélite
// que registra o trilho de auditoria (decisões, escalações) de uma
// sub-parte. Não bloqueia o pipeline — roda uma única passada por ciclo.
package auditoria

import (
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents"
	"github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

var _ agents.Harness = (*Harness)(nil)

// Harness é o esqueleto do agente Auditoria: Plan/Act/Observe
// determinísticos, sem chamar um motor de IA ainda. Leitura cross-módulo
// completa (exceção documentada em internal/agents.ReadScope, junto com
// Architect).
type Harness struct{}

func New() *Harness { return &Harness{} }

func (h *Harness) Config() agents.AgentConfig {
	return agents.AgentConfig{
		ID: "auditoria",
		ReadScope: agents.ReadScope{
			OwnModuleFull:   true,
			CrossModuleMode: "full",
		},
		MaxCyclesBeforeEscalate: 1,
		HardCapCycles:           1,
	}
}

func (h *Harness) Plan(ctx agents.AgentContext) (agents.Action, error) {
	return agents.Action{
		Kind:    "log_audit_trail",
		Payload: map[string]any{"sub_issue_id": ctx.SubIssueID},
	}, nil
}

func (h *Harness) Act(ctx agents.AgentContext, action agents.Action) (agents.Observation, error) {
	return agents.Observation{Signature: "auditoria:logged"}, nil
}

func (h *Harness) Observe(ctx agents.AgentContext, obs agents.Observation) (orchestrator.ExitResult, error) {
	return orchestrator.ExitResult{
		Status:      orchestrator.ExitSuccess,
		AgentID:     "auditoria",
		SubIssueID:  ctx.SubIssueID,
		ParentIssue: ctx.ParentIssueID,
	}, nil
}
