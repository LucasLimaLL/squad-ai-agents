package orchestrator

// AgentLoopRunner executa o loop plan->act->observe->refine de UM agente
// (ou de UMA sub-parte do fan-out, que tem seu próprio budget/detector
// isolado) até atingir um ExitResult terminal.
//
// A interface Harness fica em internal/agents para não criar dependência
// circular; aqui recebemos funções já resolvidas para manter este pacote
// livre da concretização de cada agente.
type StepFunc func(cycle int) (signature string, exit *ExitResult, err error)

type AgentLoopRunner struct {
	AgentSubPartKey string // ex: "dev:pagamentos-123-front"
	Budget          *CycleBudget
	Progress        *NoProgressDetector
	MaxLocalCycles  int // limite próprio do agente, além do budget agregado
}

// Run roda o loop até um ExitResult terminal (success, escalated, aborted)
// ou até estourar o budget agregado / limite local — o que vier primeiro
// vira uma escalação automática.
func (r *AgentLoopRunner) Run(step StepFunc) (ExitResult, error) {
	cycle := 0
	for {
		cycle++
		if err := r.Budget.Consume(r.AgentSubPartKey); err != nil {
			return ExitResult{Status: ExitEscalated, Cycle: cycle}, err
		}
		if r.MaxLocalCycles > 0 && cycle > r.MaxLocalCycles {
			return ExitResult{Status: ExitEscalated, Cycle: cycle}, nil
		}

		signature, exit, err := step(cycle)
		if err != nil {
			return ExitResult{}, err
		}
		if exit != nil {
			exit.Cycle = cycle
			return *exit, nil
		}

		// sem exit terminal ainda: checa "sem progresso" antes do próximo ciclo
		if r.Progress.Observe(signature) {
			return ExitResult{
				Status: ExitEscalated,
				Cycle:  cycle,
				Escalation: &EscalationReason{
					Category: "no_progress_same_error",
					Evidence: "mesma assinatura de resultado repetida",
				},
			}, nil
		}
	}
}
