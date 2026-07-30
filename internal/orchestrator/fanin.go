package orchestrator

// SubPartStatus acompanha o progresso de uma sub-parte do fan-out de Dev
// (ex: "pagamentos-123-front") através dos gates que rodam por sub-parte
// antes do fan-in liberar o QA integrado.
type SubPartStatus struct {
	SubIssueID string
	DevSuccess bool
	SecurityOK bool
	FinopsOK   bool
	Escalated  bool // se true, fan-in nunca libera até resolver
}

func (s SubPartStatus) Ready() bool {
	return s.DevSuccess && s.SecurityOK && s.FinopsOK && !s.Escalated
}

// FanInGate decide quando todas as sub-partes de uma feature estão prontas
// para o QA rodar o teste integrado. Falha parcial (uma sub-parte escalada
// ou não aprovada) segura o fan-in inteiro — não avança com sucesso parcial.
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

// ReadyForIntegratedQA retorna true somente quando TODAS as sub-partes
// estão prontas (Dev + Security + Finops aprovados, sem escalação pendente).
func (g *FanInGate) ReadyForIntegratedQA() bool {
	for _, sp := range g.SubParts {
		if !sp.Ready() {
			return false
		}
	}
	return true
}

// OnArchitectReconsult reseta o progresso das sub-partes afetadas quando o
// Architect dispara uma atualização centralizada — a reconsulta é sempre
// broadcast para TODOS os Devs do fan-out, e Security/Finops precisam
// reaprovar (aprovação anterior não é considerada válida).
func (g *FanInGate) OnArchitectReconsult(budget *CycleBudget) {
	for subIssueID, sp := range g.SubParts {
		sp.DevSuccess = false
		sp.SecurityOK = false
		sp.FinopsOK = false
		// o retrabalho conta contra o budget daquela sub-parte,
		// não é um ciclo "gratuito".
		_ = budget.Consume(subIssueID + ":architect_reconsult")
	}
}
