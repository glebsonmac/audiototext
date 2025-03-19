 # Guia de Organização e Estrutura

## Estrutura Principal

.fcai/
├── README.md # Este arquivo de orientação
├── commands.md # Lista de comandos disponíveis
├── state.md # Estado atual do projeto
├── structure.md # Estrutura de pastas detalhada
├── features/ # Features do sistema
└── project/ # Documentação geral do projeto

## Instruções de Operação

### Inicialização
1. Ler o [contexto do projeto](.fcai/project/documentation/context.md) para compreender objetivos e escopo
2. Verificar o [estado atual](.fcai/state.md) para identificar progresso e tarefas em andamento
3. Explorar as features existentes na pasta [features/](.fcai/features/)

### Comandos Disponíveis
Utilize comandos prefixados com `!` para interagir com o projeto:

- `!read all` - Lê todos os arquivos do projeto dentro da pasta .fcai
- `!task list` - Lista todas as tarefas do projeto
- `!show progress` - Mostra o progresso atual do projeto

Consulte [commands.md](.fcai/commands.md) para a lista completa.

### Fluxo de Trabalho para Tarefas
1. Verificar tarefas no backlog
2. Iniciar uma tarefa: `!task start <número>`
3. Atualizar progresso: `!task update <número>`
4. Completar tarefa: `!task complete <número>`

### Documentação de Referência
- [Visão geral da arquitetura](.fcai/project/architecture/overview.md)
- [Roadmap do projeto](.fcai/project/planning/roadmap.md)

---

Esta estrutura foi projetada para facilitar a organização de documentação e gerenciamento de tarefas, permitindo uma compreensão clara do estado atual e dos próximos passos do projeto.