package orchestrator

type ExitStatus string

const (
	ExitSuccess   ExitStatus = "success"
	ExitEscalated ExitStatus = "escalated"
	ExitAborted   ExitStatus = "aborted"
)

type ExitResult struct {
	Status      ExitStatus
	AgentID     string
	SubIssueID  string
	ParentIssue string
	Cycle       int
	Escalation  *EscalationReason
}

func (r ExitResult) IsTerminal() bool {
	return true
}
