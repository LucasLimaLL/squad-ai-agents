package agents

import "github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"

type ReadScope struct {
	OwnModuleFull   bool
	CrossModuleMode string
}

type AgentConfig struct {
	ID                             string
	ReadScope                      ReadScope
	MaxCyclesBeforeEscalate        int
	HardCapCycles                  int
	EscalatesOnCriticalImmediately bool
	NeverBlocksAlone               bool
}

type Harness interface {
	Config() AgentConfig
	Plan(ctx AgentContext) (Action, error)
	Act(ctx AgentContext, action Action) (Observation, error)
	Observe(ctx AgentContext, obs Observation) (orchestrator.ExitResult, error)
}

type AgentContext struct {
	SubIssueID       string
	ParentIssueID    string
	Cycle            int
	GraphVersionHash string
	PreviousResult   *Observation
}

type Action struct {
	Kind    string
	Payload map[string]any
}

type Observation struct {
	Signature string
	Data      map[string]any
}
