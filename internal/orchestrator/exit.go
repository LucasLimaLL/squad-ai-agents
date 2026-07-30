package orchestrator

// ExitStatus representa como um ciclo de agente terminou.
// Todo agente do squad-ai encerra sua execução em um destes 3 estados —
// nunca em um estado indefinido/silencioso.
type ExitStatus string

const (
	// ExitSuccess: critério de saída do agente foi atendido.
	ExitSuccess ExitStatus = "success"

	// ExitEscalated: escalado para decisão humana em tela.
	// Ao decidir, o fluxo retoma o MESMO agente, no MESMO ponto do loop.
	// Conta como 1 ciclo consumido do budget — não é neutro.
	ExitEscalated ExitStatus = "escalated"

	// ExitAborted: cancelamento humano da demanda inteira.
	// Propaga para todos os agentes/sub-partes envolvidos na feature.
	ExitAborted ExitStatus = "aborted"
)

// ExitResult é o resultado padronizado que todo agente retorna ao final
// de cada ciclo do loop plan->act->observe->refine.
type ExitResult struct {
	Status      ExitStatus
	AgentID     string
	SubIssueID  string // ex: "pagamentos-123-front"
	ParentIssue string // ex: "pagamentos-123"
	Cycle       int    // número do ciclo em que esse resultado ocorreu
	Escalation  *EscalationReason
}

// IsTerminal indica se esse resultado encerra o loop daquele agente
// (todos os três tipos são terminais para o ciclo atual; ExitEscalated
// é terminal só até a decisão humana chegar, quando o loop é retomado
// com um novo ExitResult).
func (r ExitResult) IsTerminal() bool {
	return true
}
