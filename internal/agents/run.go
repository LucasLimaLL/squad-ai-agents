package agents

import "github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"

// Run liga um Harness concreto ao AgentLoopRunner: cada ciclo chama
// Plan -> Act -> Observe e traduz o resultado para o formato que o
// runner espera (exit == nil enquanto Observe não retornar um Status
// terminal).
func Run(h Harness, budget *orchestrator.CycleBudget, subPartKey, parentIssue, subIssue string) (orchestrator.ExitResult, error) {
	runner := &orchestrator.AgentLoopRunner{
		AgentSubPartKey: subPartKey,
		Budget:          budget,
		Progress:        &orchestrator.NoProgressDetector{},
		MaxLocalCycles:  h.Config().HardCapCycles,
	}

	ctx := AgentContext{SubIssueID: subIssue, ParentIssueID: parentIssue}

	return runner.Run(func(cycle int) (string, *orchestrator.ExitResult, error) {
		ctx.Cycle = cycle

		action, err := h.Plan(ctx)
		if err != nil {
			return "", nil, err
		}

		obs, err := h.Act(ctx, action)
		if err != nil {
			return "", nil, err
		}
		ctx.PreviousResult = &obs

		result, err := h.Observe(ctx, obs)
		if err != nil {
			return "", nil, err
		}
		if result.Status != "" {
			return obs.Signature, &result, nil
		}
		return obs.Signature, nil, nil
	})
}
