// Package documentador implementa o Harness do agente Documentador:
// satélite que gera/atualiza documentação a partir do que foi
// implementado. Não bloqueia o pipeline — roda uma única passada por
// ciclo.
package documentador

import (
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents"
	"github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

var _ agents.Harness = (*Harness)(nil)

// Harness é o esqueleto do agente Documentador: Plan/Act/Observe
// determinísticos, sem chamar um motor de IA ainda. Único agente com
// CrossModuleMode "broad_search" — precisa buscar contexto amplo pra
// documentar integrações entre módulos.
type Harness struct{}

func New() *Harness { return &Harness{} }

func (h *Harness) Config() agents.AgentConfig {
	return agents.AgentConfig{
		ID: "documentador",
		ReadScope: agents.ReadScope{
			OwnModuleFull:   true,
			CrossModuleMode: "broad_search",
		},
		MaxCyclesBeforeEscalate: 1,
		HardCapCycles:           1,
	}
}

func (h *Harness) Plan(ctx agents.AgentContext) (agents.Action, error) {
	return agents.Action{
		Kind:    "generate_docs",
		Payload: map[string]any{"sub_issue_id": ctx.SubIssueID},
	}, nil
}

func (h *Harness) Act(ctx agents.AgentContext, action agents.Action) (agents.Observation, error) {
	return agents.Observation{Signature: "documentador:docs_generated"}, nil
}

func (h *Harness) Observe(ctx agents.AgentContext, obs agents.Observation) (orchestrator.ExitResult, error) {
	return orchestrator.ExitResult{
		Status:      orchestrator.ExitSuccess,
		AgentID:     "documentador",
		SubIssueID:  ctx.SubIssueID,
		ParentIssue: ctx.ParentIssueID,
	}, nil
}
