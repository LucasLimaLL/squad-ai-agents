# squad-ai-agents

Núcleo de orquestração da squad-ai: agentes de IA (Analista, PM, Architect,
Dev, QA, Security, SRE, Finops como gates; Auditoria e Documentador como
satélites) que executam um fluxo autônomo plan → act → observe → refine,
com critérios de saída explícitos e guardrails de budget para nunca ficar
em loop infinito.

## Requisitos

- Go 1.22+

## Build, vet e testes

```bash
go build ./...
go vet ./...
go test ./... -v
```

## Estrutura

- `internal/orchestrator` — núcleo independente de agente específico:
  tipos de saída de ciclo, budget de ciclos, detecção de sem-progresso,
  gate de fan-in e o loop runner.
- `internal/agents` — interface `Harness` comum a todos os agentes; cada
  agente concreto (dev, qa, security, ...) implementa essa interface no
  seu próprio pacote/harness.

## Conceitos principais

### ExitStatus / ExitResult

Todo ciclo de agente termina em um destes três estados — nunca em um
estado indefinido (`internal/orchestrator/exit.go`):

- `ExitSuccess` — critério de saída do agente foi atendido.
- `ExitEscalated` — escalado para decisão humana em tela; ao decidir, o
  fluxo retoma o **mesmo** agente no **mesmo** ponto do loop. Conta como
  1 ciclo consumido do budget — não é neutro.
- `ExitAborted` — cancelamento humano da demanda inteira; propaga para
  todos os agentes/sub-partes envolvidos na feature.

### CycleBudget e NoProgressDetector

`internal/orchestrator/guardrails.go`:

- `CycleBudget` limita o total de ciclos **agregados** por issue/módulo
  (configurável por `NewCycleBudget(moduleID, issueID, maxAggregate)`).
  Escalação também consome ciclo.
- `NoProgressDetector` sinaliza "sem progresso" quando a mesma
  assinatura de resultado (ex.: hash do erro de build) se repete 2x
  seguidas — sinal de loop morto, não de progresso lento.

### FanInGate

`internal/orchestrator/fanin.go`: segura o fan-in de um fan-out de Dev
até que **todas** as sub-partes estejam com Dev + Security + Finops
aprovados e sem escalação pendente (`ReadyForIntegratedQA`). Uma
reconsulta do Architect (`OnArchitectReconsult`) reseta as aprovações de
todas as sub-partes e consome budget de cada uma — o retrabalho não é
gratuito.

### AgentLoopRunner

`internal/orchestrator/loop.go`: executa o loop plan→act→observe→refine
de um agente (ou de uma sub-parte isolada do fan-out) até um
`ExitResult` terminal, checando budget e sem-progresso a cada ciclo.
Recebe a lógica de cada ciclo via `StepFunc`:

```go
type StepFunc func(cycle int) (signature string, exit *ExitResult, err error)
```

`exit == nil` significa "ainda sem resultado terminal, continue o
loop"; retornar um `*ExitResult` encerra o loop naquele ciclo.

### Harness

`internal/agents/harness.go`: interface que todo agente concreto
implementa — `Config()` (read scope, limites de ciclo), `Plan`, `Act` e
`Observe`. O motor de IA por trás de cada agente (Claude, Devin, etc.) é
um detalhe de implementação de `Act`; o orquestrador só depende desta
interface.

## Exemplo de uso

Ligando um `Harness` concreto ao `AgentLoopRunner` (ver também os casos
de uso diretos do runner em `internal/orchestrator/loop_test.go`):

```go
package main

import (
    "github.com/LucasLimaLL/squad-ai-agents/internal/agents"
    "github.com/LucasLimaLL/squad-ai-agents/internal/orchestrator"
)

func runAgent(h agents.Harness, budget *orchestrator.CycleBudget, subPartKey, parentIssue, subIssue string) (orchestrator.ExitResult, error) {
    runner := &orchestrator.AgentLoopRunner{
        AgentSubPartKey: subPartKey,
        Budget:          budget,
        Progress:        &orchestrator.NoProgressDetector{},
        MaxLocalCycles:  h.Config().HardCapCycles,
    }

    ctx := agents.AgentContext{SubIssueID: subIssue, ParentIssueID: parentIssue}

    return runner.Run(func(cycle int) (string, *orchestrator.ExitResult, error) {
        ctx.Cycle = cycle

        action, err := h.Plan(ctx)
        if err != nil {
            return "", nil, err
        }

        obs, err := h.Act(ctx, action)
        if err != nil {
            return "", nil, err
        }
        ctx.PreviousResult = &obs

        result, err := h.Observe(ctx, obs)
        if err != nil {
            return "", nil, err
        }
        if result.Status != "" {
            return obs.Signature, &result, nil
        }
        return obs.Signature, nil, nil
    })
}

func main() {
    budget := orchestrator.NewCycleBudget("pagamentos", "pagamentos-123", 20)

    _ = budget
}
```

Fan-out/fan-in de Dev:

```go
gate := orchestrator.NewFanInGate("pagamentos-123", []string{
    "pagamentos-123-front",
    "pagamentos-123-back",
})

gate.SubParts["pagamentos-123-front"].DevSuccess = true
gate.SubParts["pagamentos-123-front"].SecurityOK = true
gate.SubParts["pagamentos-123-front"].FinopsOK = true

gate.ReadyForIntegratedQA()
```

## Fluxo de branches

- `feature/**` e `fix/**` abrem PR automático para `develop`
  (`.github/workflows/feature-to-develop.yml`).
- Push em `develop` gera bump semântico via conventional commits e cria
  branch/PR de release para `main` (`.github/workflows/develop-to-release.yml`).
  Merge em `main` é sempre humano.
