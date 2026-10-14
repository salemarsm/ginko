# Decision 001 — Verdade do repositório antes de funcionalidade nova

- **Data:** 2026-10-14
- **Ciclo:** 1
- **Tipo:** decisão de projeto (humano)
- **Decidido por:** Rafael Salema Marques

## Contexto

O projeto estava parado desde 2026-05-06, com a tag `v1.0` publicada e o código
implementado, mas a documentação descrevendo um produto mais antigo. Detalhes e
evidências em [`../cycles/cycle-001.md`](../cycles/cycle-001.md).

## Opções consideradas

1. **Verdade do repositório primeiro** (higiene de release) — whitepaper, CI,
   `.gitattributes`, reconciliação de docs, issues a partir do backlog.
   *Escolhida.*
2. **Dogfooding** — usar o ginko como camada de memória da própria expertise,
   instalando o MCP e ingerindo conhecimento real.
   *Adiada:* é o passo que materializa a intenção do projeto, mas construir sobre
   um repositório que mente sobre si mesmo gera backlog baseado em premissa errada.
3. **Robustez de retrieval** — alocador de token budget, cache de contexto,
   benchmarks em 1k/10k/100k, adapter de embeddings.
   *Descartada por ora:* otimizar o coração técnico antes de validar o uso real.
4. **Coordenação multi-agente** — leases, handoffs e conflitos ligados ao harness
   de desenvolvimento.
   *Descartada por ora:* escopo grande, e depende do item 2 ter validado o uso básico.

## Decisão tomada

Corrigir primeiro todas as afirmações falsas do repositório sobre o próprio
estado, restaurar o CI do whitepaper e migrar o backlog para issues rastreáveis.
Funcionalidade nova vem depois.

## Rationale

O diferencial declarado do ginko é ser auditável e honesto sobre o que existe —
"a LLM não é o banco de dados", "embeddings são índices, não verdade". Um
repositório cujo whitepaper descreve o estado de três releases atrás contradiz a
própria tese, e isso é mais caro do que qualquer feature faltando: contamina a
confiança em todo o resto.

Os itens 2, 3 e 4 não foram rejeitados, apenas ordenados. O item 2 gera backlog a
partir de uso real em vez de suposição, e por isso deve vir antes de otimizar.

**Trade-off aceito:** este ciclo não entrega nenhuma capacidade nova ao usuário
final. Em troca, os ciclos seguintes partem de um repositório que descreve a
realidade.

## Ações derivadas

- [x] Atualizar o whitepaper para o estado v1.0 e recompilar o PDF
- [x] Corrigir o workflow `whitepaper.yml` (pacote `lmodern` e permissão de escrita)
- [x] Adicionar `.gitattributes` fixando LF
- [x] Corrigir os arquivos com formatação real fora do `gofmt`
- [x] Reconciliar README, `docs/roadmap.md`, `docs/install.md` e CHANGELOG
- [x] Criar issues a partir dos itens P0 do backlog
- [ ] Revisar e mergear o PR #10 *(decisão do humano)*
- [ ] Ciclo 2: decidir entre dogfooding (item 2) e robustness de retrieval (item 3)

## Revisão futura

Revisitar se o ciclo 2 mostrar que o problema real não era documentação, mas
alguma lacuna funcional que o dogfooding exponha de imediato — nesse caso a
ordem foi conservadora demais e o custo de oportunidade aparece.
