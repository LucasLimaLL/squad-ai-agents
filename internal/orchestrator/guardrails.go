package orchestrator

import "fmt"

// CycleBudget controla o teto de ciclos por issue, agregado por módulo/pipeline
// configurada — não é um número fixo para todos os módulos. Escalação conta
// como ciclo (decisão explícita do projeto: escalação não é neutra).
type CycleBudget struct {
	ModuleID     string
	IssueID      string
	MaxAggregate int            // teto somado entre todos os agentes/sub-partes da issue
	perAgent     map[string]int // ciclos consumidos por agente/sub-parte
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

// Consume registra um ciclo gasto por um agente/sub-parte (incluindo
// escalações, que contam como ciclo). Retorna erro se estourar o budget
// agregado da issue.
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

// NoProgressDetector identifica quando um agente está repetindo o mesmo
// erro/resultado sem avançar — sinal de loop morto, não de progresso lento.
// Regra do projeto: mesmo erro/resultado repetido 2x seguidas.
type NoProgressDetector struct {
	lastSignature string
	repeatCount   int
}

// Observe recebe uma assinatura do resultado do ciclo atual (ex: hash da
// mensagem de erro, ou do motivo de reprovação) e retorna true se isso já
// caracteriza "sem progresso" (2 repetições seguidas da mesma assinatura).
func (d *NoProgressDetector) Observe(signature string) (noProgress bool) {
	if signature == d.lastSignature {
		d.repeatCount++
	} else {
		d.repeatCount = 0
		d.lastSignature = signature
	}
	return d.repeatCount >= 1 // 2ª ocorrência da mesma assinatura (repeatCount chega a 1)
}
