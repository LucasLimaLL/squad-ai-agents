package orchestrator

type StepFunc func(cycle int) (signature string, exit *ExitResult, err error)

type AgentLoopRunner struct {
	AgentSubPartKey string
	Budget          *CycleBudget
	Progress        *NoProgressDetector
	MaxLocalCycles  int
}

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
