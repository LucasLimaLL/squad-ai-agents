package orchestrator

import "testing"

func TestLoopRunner_NoProgressEscalates(t *testing.T) {
	budget := NewCycleBudget("pagamentos", "pagamentos-123", 10)
	runner := &AgentLoopRunner{
		AgentSubPartKey: "dev:pagamentos-123-front",
		Budget:          budget,
		Progress:        &NoProgressDetector{},
		MaxLocalCycles:  5,
	}

	result, err := runner.Run(func(cycle int) (string, *ExitResult, error) {
		// simula o mesmo erro se repetindo -> deve escalar por sem-progresso
		return "same-build-error", nil, nil
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if result.Status != ExitEscalated {
		t.Fatalf("esperado ExitEscalated, obtido %v", result.Status)
	}
	if result.Cycle != 2 {
		t.Fatalf("esperado escalar no ciclo 2 (2ª repetição), obtido %d", result.Cycle)
	}
}

func TestLoopRunner_SuccessOnFirstMatch(t *testing.T) {
	budget := NewCycleBudget("pagamentos", "pagamentos-123", 10)
	runner := &AgentLoopRunner{
		AgentSubPartKey: "dev:pagamentos-123-back",
		Budget:          budget,
		Progress:        &NoProgressDetector{},
		MaxLocalCycles:  5,
	}

	result, err := runner.Run(func(cycle int) (string, *ExitResult, error) {
		return "build-ok", &ExitResult{Status: ExitSuccess}, nil
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if result.Status != ExitSuccess {
		t.Fatalf("esperado ExitSuccess, obtido %v", result.Status)
	}
}

func TestCycleBudget_AggregateOverflowAcrossSubParts(t *testing.T) {
	budget := NewCycleBudget("pagamentos", "pagamentos-123", 3)
	if err := budget.Consume("dev:front"); err != nil {
		t.Fatalf("não deveria estourar ainda: %v", err)
	}
	if err := budget.Consume("dev:back"); err != nil {
		t.Fatalf("não deveria estourar ainda: %v", err)
	}
	if err := budget.Consume("security:front"); err != nil {
		t.Fatalf("não deveria estourar ainda: %v", err)
	}
	if err := budget.Consume("security:back"); err == nil {
		t.Fatalf("esperava erro de budget agregado estourado")
	}
}

func TestFanInGate_HoldsOnPartialFailure(t *testing.T) {
	gate := NewFanInGate("pagamentos-123", []string{"pagamentos-123-front", "pagamentos-123-back"})
	gate.SubParts["pagamentos-123-front"].DevSuccess = true
	gate.SubParts["pagamentos-123-front"].SecurityOK = true
	gate.SubParts["pagamentos-123-front"].FinopsOK = true
	// back ainda não terminou
	if gate.ReadyForIntegratedQA() {
		t.Fatalf("fan-in não deveria liberar com sub-parte incompleta")
	}

	gate.SubParts["pagamentos-123-back"].DevSuccess = true
	gate.SubParts["pagamentos-123-back"].SecurityOK = true
	gate.SubParts["pagamentos-123-back"].FinopsOK = true
	if !gate.ReadyForIntegratedQA() {
		t.Fatalf("fan-in deveria liberar com todas as sub-partes prontas")
	}
}
