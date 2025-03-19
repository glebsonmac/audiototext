# Estrutura de Pastas do Projeto

Este documento descreve a estrutura de pastas recomendada para projetos que utilizam este template. A estrutura foi projetada para facilitar a organização de documentação, tarefas e features do projeto.

## Visão Geral

```
.fcai/
├── README.md                     # Visão geral do projeto e instruções iniciais
├── commands.md                   # Lista de comandos disponíveis para interação
├── state.md                      # Estado atual do projeto e progresso
├── structure.md                  # Este arquivo (estrutura de pastas)
├── features/                     # Pasta para features do sistema
│   └── [feature-name]/           # Um diretório para cada feature
│       ├── documentation/        # Documentação da feature
│       │   └── overview.md       # Visão geral da feature
│       ├── backlog/              # Tarefas planejadas para o futuro
│       │   └── [number]-[name]/  # Pasta para cada tarefa no backlog
│       │       └── [number]-[name].md       # Descrição da tarefa
│       ├── in-progress/          # Tarefas em andamento
│       │   └── [number]-[name]/  # Pasta para cada tarefa em andamento
│       │       ├── [number]-[name].md       # Descrição da tarefa
│       │       └── [number]-[name].result.md     # Resultados parciais da tarefa
│       └── completed/            # Tarefas concluídas
│           └── [number]-[name]/  # Pasta para cada tarefa concluída
│               ├── [number]-[name].md       # Descrição da tarefa
│               └── [number]-[name].result.md     # Resultados da tarefa
└── project/                      # Documentação geral do projeto
    ├── documentation/            # Documentação principal
    │   └── context.md            # Contexto geral do projeto
    ├── planning/                 # Planejamento do projeto
    │   └── roadmap.md            # Roadmap e cronograma
    ├── architecture/             # Documentação de arquitetura
    │   └── overview.md           # Visão geral da arquitetura
    └── analysis/                 # Análises e estudos
        └── [analysis-name].md    # Documentos de análise específicos
```

## Detalhamento das Pastas

### Arquivos Principais

- **README.md**: Contém uma visão geral do projeto, instruções de como começar e referências para documentação mais detalhada.
- **commands.md**: Lista todos os comandos disponíveis para interagir com o projeto através da IA.
- **state.md**: Mantém o estado atual do projeto, incluindo progresso, features ativas e tarefas em andamento.
- **structure.md**: Este arquivo, que descreve a estrutura de pastas do projeto.

### Pasta `features/`

Esta pasta contém um diretório para cada feature do sistema. Cada feature segue a mesma estrutura interna:

- **documentation/**: Contém a documentação específica da feature.
  - **overview.md**: Visão geral da feature, incluindo propósito, funcionalidades e arquitetura interna.

- **backlog/**: Contém tarefas planejadas para o futuro.
  - **[number]-[name]/**: Uma pasta para cada tarefa no backlog.
    - **task.md**: Descrição detalhada da tarefa.

- **in-progress/**: Contém tarefas que estão sendo trabalhadas atualmente.
  - **[number]-[name]/**: Uma pasta para cada tarefa em andamento.
    - **task.md**: Descrição detalhada da tarefa.
    - **result.md**: Resultados parciais da tarefa.

- **completed/**: Contém tarefas que foram concluídas.
  - **[number]-[name]/**: Uma pasta para cada tarefa concluída.
    - **task.md**: Descrição detalhada da tarefa.
    - **result.md**: Resultados finais da tarefa.

### Pasta `project/`

Esta pasta contém a documentação geral do projeto:

- **documentation/**: Contém a documentação principal do projeto.
  - **context.md**: Descreve o contexto geral do projeto, incluindo objetivos, escopo e visão geral.

- **planning/**: Contém documentos relacionados ao planejamento do projeto.
  - **roadmap.md**: Descreve o roadmap e cronograma do projeto.

- **architecture/**: Contém documentação relacionada à arquitetura do sistema.
  - **overview.md**: Visão geral da arquitetura do sistema.

- **analysis/**: Contém análises e estudos relacionados ao projeto.
  - **[analysis-name].md**: Documentos específicos de análise.

## Convenções de Nomenclatura

1. **Pastas de Tarefas**: Sempre use o formato `[number]-[task-name]`, onde:
   - `[number]` é um identificador numérico único para a tarefa
   - `[task-name]` é um nome curto e descritivo, usando hífens para separar palavras

2. **Arquivos Markdown**: Use nomes em minúsculas, com hífens para separar palavras (ex: `command-implementation.md`).

3. **Features**: Use nomes em minúsculas, com hífens para separar palavras (ex: `user-interface`).

## Gerenciamento de Tarefas

As tarefas seguem um fluxo de trabalho específico:

1. Inicialmente, as tarefas são criadas no `backlog/` da feature relevante.
2. Quando o trabalho começa, a tarefa é movida para `in-progress/`.
3. Quando a tarefa é concluída, ela é movida para `completed/`.

Este fluxo pode ser gerenciado usando os comandos definidos em `commands.md`, como `!task start`, `!task update` e `!task complete`. 