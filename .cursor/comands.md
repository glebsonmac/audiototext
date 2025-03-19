# Comandos Disponíveis

Este documento lista todos os comandos disponíveis para interagir com o projeto através da IA. Os comandos são prefixados com `!` para distingui-los de mensagens normais.

## Comandos de Leitura

### !read [arquivo]
Lê o conteúdo de um arquivo específico.

**Exemplo:** `!read .fcai/state.md`

### !read all
Lê todos os arquivos principais do projeto para obter uma visão geral completa.

**Exemplo:** `!read all`

## Comandos de Tarefas

### !task list
Lista todas as tarefas do projeto, organizadas por status (backlog, em andamento, concluídas).

**Exemplo:** `!task list`

### !task start [número]
Inicia uma nova tarefa do backlog, movendo-a para o status "em andamento".

**Exemplo:** `!task start 5`

### !task update [número]
Atualiza o progresso de uma tarefa em andamento.

**Exemplo:** `!task update 3`

### !task complete [número]
Marca uma tarefa como concluída, movendo-a para o status "concluída".

**Exemplo:** `!task complete 2`

### !task move [número] [status]
Move uma tarefa para um status específico (backlog, in-progress, completed).

**Exemplo:** `!task move 4 backlog`

## Comandos de Atualização

### !update state
Atualiza o arquivo de estado do projeto com as informações mais recentes.

**Exemplo:** `!update state`

### !update all
Atualiza todos os arquivos principais do projeto.

**Exemplo:** `!update all`

## Comandos de Criação

### !create feature [nome]
Cria uma nova feature com a estrutura de pastas padrão.

**Exemplo:** `!create feature user-authentication`

### !create task [feature] [nome]
Cria uma nova tarefa no backlog de uma feature específica.

**Exemplo:** `!create task user-interface login-page`

## Comandos de Análise

### !show progress
Mostra o progresso geral do projeto, incluindo tarefas concluídas e pendentes.

**Exemplo:** `!show progress`

### !show features
Lista todas as features do projeto com uma breve descrição.

**Exemplo:** `!show features`

## Aliases

Alguns comandos têm aliases (atalhos) para facilitar o uso:

- `!r` = `!read`
- `!t` = `!task`
- `!u` = `!update`
- `!c` = `!create`
- `!s` = `!show`

**Exemplo:** `!t list` é equivalente a `!task list`

## Comandos Básicos
- `!help` - Mostra esta lista de comandos
- `!read context` - Lê o contexto geral do projeto
- `!read state` - Lê o estado atual do projeto

## Leitura de Features e Tarefas
- `!read feature <nome>` - Lê a documentação de uma feature
- `!read task <número>` - Lê a documentação de uma tarefa
- `!read result <número>` - Lê os resultados de uma tarefa

## Listagem
- `!list features` - Lista todas as features
- `!list tasks` - Lista todas as tarefas
- `!list tasks active` - Lista tarefas em andamento
- `!list tasks done` - Lista tarefas concluídas
- `!list tasks pending` - Lista tarefas pendentes

## Visualização
- `!show structure` - Mostra a estrutura do projeto
- `!show history` - Mostra o histórico recente de alterações

## Gerenciamento de Tarefas
- `!task start <feature> <número> <nome>` - Inicia uma nova tarefa
- `!task complete <número>` - Marca uma tarefa como concluída
- `!task update <número>` - Atualiza o estado de uma tarefa
- `!task move <número> <status>` - Move uma tarefa para outro status (backlog, in-progress, completed)

## Atualização do Estado
- `!update state` - Atualiza o arquivo de estado do projeto
- `!update all` - Atualiza todos os arquivos de documentação

## Como Usar
Quando precisar recuperar informações sobre o projeto, basta digitar um desses comandos no chat. Por exemplo:

```
!read context
```

Isso fará com que a IA leia o contexto geral do projeto para recuperar informações importantes.

## Fluxo de Trabalho

### Iniciando uma Nova Sessão
1. `!read all` - Para carregar o contexto completo do projeto, lendo todos os arquivos da pasta .fcai
2. `!list tasks active` - Para ver quais tarefas estão em andamento
3. `!read task <número>` - Para ler a tarefa específica que você vai trabalhar

### Trabalhando em uma Tarefa
1. `!read task <número>` - Para entender os requisitos da tarefa
2. Desenvolver a solução
3. `!task update <número>` - Para atualizar o estado da tarefa
4. `!update state` - Para atualizar o estado geral do projeto

### Concluindo uma Tarefa
1. `!task complete <número>` - Para marcar a tarefa como concluída
2. `!update all` - Para atualizar toda a documentação

### Quando Perder o Contexto
1. `!help` - Para ver a lista de comandos disponíveis
2. `!read state` - Para entender o estado atual do projeto
3. `!list tasks active` - Para ver quais tarefas estão em andamento 