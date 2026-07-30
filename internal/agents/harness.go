package agents

import "github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"

// ReadScope define o que um agente pode ler do grafo de conhecimento:
// sempre leitura completa do próprio módulo, e leitura cross-módulo
// sempre parcial/funcional (nunca "acesso total a tudo"), exceto para
// Auditoria e Architect.
type ReadScope struct {
	OwnModuleFull   bool
	CrossModuleMode string // "none" | "interface_summary" | "policies_only" | "infra_deps" | "full" | "broad_search"
}

// AgentConfig é a configuração estática de harness de um agente —
// o que ele pode ler, e os limites do loop dele.
type AgentConfig struct {
	ID                              string
	ReadScope                       ReadScope
	MaxCyclesBeforeEscalate         int
	HardCapCycles                   int  // 0 = sem hard cap próprio (usa só o budget agregado)
	EscalatesOnCriticalImmediately  bool // ex: Security
	NeverBlocksAlone                bool // ex: Finops
}

// Harness é a interface que todo agente do squad-ai implementa.
// O motor de IA por trás (Claude, Devin, etc.) é um detalhe de
// implementação de Act — o orquestrador só depende desta interface.
type Harness interface {
	Config() AgentConfig

	// Plan decide a próxima ação com base no contexto atual.
	Plan(ctx AgentContext) (Action, error)

	// Act executa a ação decidida (chama sandbox, GitHub API, etc.)
	Act(ctx AgentContext, action Action) (Observation, error)

	// Observe interpreta o resultado e decide se o critério de saída
	// foi atendido, se deve refinar, ou se deve escalar.
	Observe(ctx AgentContext, obs Observation) (orchestrator.ExitResult, error)
}

// AgentContext é o payload de contexto passado a cada ciclo — o
// conteúdo exato varia por agente (ver harness de cada um), mas todo
// agente recebe ao menos isto:
type AgentContext struct {
	SubIssueID       string
	ParentIssueID    string
	Cycle            int
	GraphVersionHash string // melhoria #5: versão do grafo lida nesse ciclo
	PreviousResult   *Observation
}

type Action struct {
	Kind    string // ex: "run_build", "comment_pr", "read_module_summary"
	Payload map[string]any
}

type Observation struct {
	Signature string // usada pelo NoProgressDetector (hash do erro/resultado)
	Data      map[string]any
}
