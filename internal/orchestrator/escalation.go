package orchestrator

type EscalationReason struct {
	Agent    string
	SubPart  string
	Category string
	Evidence string
	Tags     []string
	FreeText string
}

var KnownEscalationCategories = map[string][]string{
	"analista":  {"scope_unresolved"},
	"pm":        {"criteria_not_testable"},
	"architect": {"cross_module_conflict_unresolved"},
	"dev":       {"no_progress_same_error", "budget_exceeded"},
	"qa":        {"same_criterion_failing_repeatedly"},
	"security":  {"critical_finding", "recurring_finding"},
	"sre":       {"operational_item_unresolved"},
}
