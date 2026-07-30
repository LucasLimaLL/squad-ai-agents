package agents_test

import (
	"testing"

	"github.com/LucasLimaLL/squad-ai-agents/internal/agents"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/analista"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/architect"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/auditoria"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/dev"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/documentador"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/finops"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/pm"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/qa"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/security"
	"github.com/LucasLimaLL/squad-ai-agents/internal/agents/sre"
	"github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

// TestFanOutFanIn_DevSecurityFinops roda o loop completo (Plan->Act->Observe
// via agents.Run) para os agentes Dev, Security e Finops, nas duas
// sub-partes de uma issue fictícia, e confirma que o FanInGate só libera
// o QA integrado depois que TODAS as sub-partes passaram pelos três.
func TestFanOutFanIn_DevSecurityFinops(t *testing.T) {
	const parentIssue = "pagamentos-123"
	subParts := []string{"pagamentos-123-front", "pagamentos-123-back"}

	budget := orchestrator.NewCycleBudget("pagamentos", parentIssue, 100)
	gate := orchestrator.NewFanInGate(parentIssue, subParts)

	devHarness := dev.New()
	secHarness := security.New()
	finHarness := finops.New()

	for _, sub := range subParts {
		if gate.ReadyForIntegratedQA() {
			t.Fatalf("fan-in não deveria estar pronto antes de processar %s", sub)
		}

		result, err := agents.Run(devHarness, budget, "dev:"+sub, parentIssue, sub)
		if err != nil {
			t.Fatalf("dev retornou erro em %s: %v", sub, err)
		}
		if result.Status != orchestrator.ExitSuccess {
			t.Fatalf("dev não teve sucesso em %s: status=%v", sub, result.Status)
		}
		gate.SubParts[sub].DevSuccess = true

		result, err = agents.Run(secHarness, budget, "security:"+sub, parentIssue, sub)
		if err != nil {
			t.Fatalf("security retornou erro em %s: %v", sub, err)
		}
		if result.Status != orchestrator.ExitSuccess {
			t.Fatalf("security não teve sucesso em %s: status=%v", sub, result.Status)
		}
		gate.SubParts[sub].SecurityOK = true

		result, err = agents.Run(finHarness, budget, "finops:"+sub, parentIssue, sub)
		if err != nil {
			t.Fatalf("finops retornou erro em %s: %v", sub, err)
		}
		if result.Status != orchestrator.ExitSuccess {
			t.Fatalf("finops não teve sucesso em %s: status=%v", sub, result.Status)
		}
		gate.SubParts[sub].FinopsOK = true
	}

	if !gate.ReadyForIntegratedQA() {
		t.Fatalf("fan-in deveria estar pronto: dev+security+finops OK em todas as sub-partes")
	}

	if got := budget.UsedBy("dev:" + subParts[0]); got != 1 {
		t.Fatalf("esperado 1 ciclo consumido por dev:%s, obtido %d", subParts[0], got)
	}
}

// TestAllAgents_Config confirma que os 10 agentes existem, cada um com
// ID único, e que os dois exemplos documentados no core (Security escala
// imediato em achado crítico; Finops nunca bloqueia sozinho) estão
// refletidos na config.
func TestAllAgents_Config(t *testing.T) {
	harnesses := []agents.Harness{
		analista.New(),
		pm.New(),
		architect.New(),
		dev.New(),
		qa.New(),
		security.New(),
		sre.New(),
		finops.New(),
		auditoria.New(),
		documentador.New(),
	}

	seen := make(map[string]bool, len(harnesses))
	for _, h := range harnesses {
		cfg := h.Config()
		if cfg.ID == "" {
			t.Fatalf("agente sem ID configurado: %#v", cfg)
		}
		if seen[cfg.ID] {
			t.Fatalf("ID de agente duplicado: %s", cfg.ID)
		}
		seen[cfg.ID] = true

		if !cfg.ReadScope.OwnModuleFull {
			t.Fatalf("%s: OwnModuleFull deveria ser sempre true", cfg.ID)
		}
	}
	if len(seen) != 10 {
		t.Fatalf("esperado 10 agentes distintos, obtido %d", len(seen))
	}

	if !security.New().Config().EscalatesOnCriticalImmediately {
		t.Fatalf("security deveria ter EscalatesOnCriticalImmediately=true")
	}
	if !finops.New().Config().NeverBlocksAlone {
		t.Fatalf("finops deveria ter NeverBlocksAlone=true")
	}
}
