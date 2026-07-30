package orchestrator

type SubPartStatus struct {
	SubIssueID string
	DevSuccess bool
	SecurityOK bool
	FinopsOK   bool
	Escalated  bool
}

func (s SubPartStatus) Ready() bool {
	return s.DevSuccess && s.SecurityOK && s.FinopsOK && !s.Escalated
}

type FanInGate struct {
	ParentIssueID string
	SubParts      map[string]*SubPartStatus
}

func NewFanInGate(parentIssueID string, subIssueIDs []string) *FanInGate {
	subParts := make(map[string]*SubPartStatus, len(subIssueIDs))
	for _, id := range subIssueIDs {
		subParts[id] = &SubPartStatus{SubIssueID: id}
	}
	return &FanInGate{ParentIssueID: parentIssueID, SubParts: subParts}
}

func (g *FanInGate) ReadyForIntegratedQA() bool {
	for _, sp := range g.SubParts {
		if !sp.Ready() {
			return false
		}
	}
	return true
}

func (g *FanInGate) OnArchitectReconsult(budget *CycleBudget) {
	for subIssueID, sp := range g.SubParts {
		sp.DevSuccess = false
		sp.SecurityOK = false
		sp.FinopsOK = false
		_ = budget.Consume(subIssueID + ":architect_reconsult")
	}
}
