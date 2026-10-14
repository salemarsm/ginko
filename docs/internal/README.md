# Internal

Notas de trabalho do processo de desenvolvimento deste projeto: contexto para
agentes, registro de decisões e specs.

Estão versionadas de propósito. Decisões sobre arquitetura, escopo e prioridade
fazem parte da história do projeto, e quem chega depois deve conseguir ver
**por que** ele está como está — não só o que ele faz.

## O que é cada coisa

| Arquivo | Para quem | Conteúdo |
|---|---|---|
| [`../roadmap.md`](../roadmap.md) | usuários e contribuidores | o que o produto faz e o que está planejado |
| [`../backlog.md`](../backlog.md) | contribuidores | itens de trabalho priorizados |
| [`project-context.md`](project-context.md) | agentes e mantenedores | como o repositório é de fato, e onde ele é frágil |
| [`decisions/`](decisions/) | quem pergunta "por quê" | a decisão, as opções descartadas e o racional |
| [`cycles/`](cycles/) | mantenedores | orientação e observações de cada ciclo de desenvolvimento |

## Regras

1. **Decisão de projeto é do humano.** Arquitetura, escopo, prioridade,
   trade-offs de design e release. Decisões táticas (como implementar o que já
   foi decidido) são do agente e não passam por aqui.
2. **Registro após o fato não vale.** A decisão é registrada quando é tomada,
   com as opções que estavam na mesa naquele momento — inclusive as descartadas.
3. **Erro também vira registro.** Quando uma hipótese se prova errada, ela fica
   documentada com a correção. Um diretório de decisões que só contém acertos é
   um diretório de marketing.
