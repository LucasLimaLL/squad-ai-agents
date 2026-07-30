package orchestrator

// EscalationReason é o formato padronizado de motivo de escalação,
// usado tanto na tela de decisão humana quanto no log da Auditoria.
type EscalationReason struct {
	Agent    string   // quem escalou (ex: "dev", "security")
	SubPart  string   // ex: "pagamentos-123-front"
	Category string   // categoria fixa por agente, para agregação (ex: "no_progress", "critical_finding")
	Evidence string   // evidência objetiva (ex: "mesmo erro de build 2x seguidas")
	Tags     []string // flexível, para classificação adicional
	FreeText string   // contexto legível para o humano na tela
}

// Categorias fixas conhecidas por agente. Mantidas como referência —
// cada harness de agente deve escolher a partir de um conjunto restrito
// para permitir agregação confiável na Auditoria.
var KnownEscalationCategories = map[string][]string{
	"analista":  {"scope_unresolved"},
	"pm":        {"criteria_not_testable"},
	"architect": {"cross_module_conflict_unresolved"},
	"dev":       {"no_progress_same_error", "budget_exceeded"},
	"qa":        {"same_criterion_failing_repeatedly"},
	"security":  {"critical_finding", "recurring_finding"},
	"sre":       {"operational_item_unresolved"},
	// finops nunca escala sozinho — bloqueio fora do threshold é sempre
	// decisão humana explícita, não uma escalação por falta de progresso.
}
