package orchestrator

import "fmt"

type CycleBudget struct {
	ModuleID     string
	IssueID      string
	MaxAggregate int
	perAgent     map[string]int
	total        int
}

func NewCycleBudget(moduleID, issueID string, maxAggregate int) *CycleBudget {
	return &CycleBudget{
		ModuleID:     moduleID,
		IssueID:      issueID,
		MaxAggregate: maxAggregate,
		perAgent:     make(map[string]int),
	}
}

func (b *CycleBudget) Consume(agentSubPartKey string) error {
	b.perAgent[agentSubPartKey]++
	b.total++
	if b.total > b.MaxAggregate {
		return fmt.Errorf(
			"budget agregado estourado para issue %s no módulo %s: %d/%d ciclos (última consumida por %s)",
			b.IssueID, b.ModuleID, b.total, b.MaxAggregate, agentSubPartKey,
		)
	}
	return nil
}

func (b *CycleBudget) UsedBy(agentSubPartKey string) int {
	return b.perAgent[agentSubPartKey]
}

func (b *CycleBudget) Total() int {
	return b.total
}

type NoProgressDetector struct {
	lastSignature string
	repeatCount   int
}

func (d *NoProgressDetector) Observe(signature string) (noProgress bool) {
	if signature == d.lastSignature {
		d.repeatCount++
	} else {
		d.repeatCount = 0
		d.lastSignature = signature
	}
	return d.repeatCount >= 1
}
