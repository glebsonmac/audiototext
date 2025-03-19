# Fluxo Completo da Aplicação de Conversão de Vídeos para HLS com Upload para S3

## **Fase 1
### **1. Recebimento do Vídeo pela API**
- O cliente (frontend ou outra aplicação) faz uma requisição **HTTP multipart/form-data** enviando o arquivo de vídeo para a API.
- A API:
  1. Salva o vídeo temporariamente no servidor (ex: `/tmp/uploads`).
  2. Registra os metadados no **banco de dados PostgreSQL**, incluindo:
     - Nome do arquivo original.
     - Status inicial (`pendente`).
     - Timestamp da solicitação.
     - ID único do vídeo.
  3. Insere a tarefa na **fila de processamento** usando **channels**.

### **2. Processamento Assíncrono da Conversão**
- Um **worker** (goroutine) monitora a fila de processamento e aguarda novos vídeos.
- Quando um novo vídeo entra na fila:
  1. O worker lê o arquivo do disco.
  2. Executa a conversão para **HLS**, gerando:
     - Segmentos `.ts`.
     - Manifesto `.m3u8`.
  3. Atualiza o **banco de dados** para `em processamento`.

### **3. Upload dos Arquivos HLS para S3**
- Após a conversão, uma goroutine inicia o **upload para S3**:
  1. Envia cada segmento `.ts` e o arquivo `.m3u8`.
  2. Após o sucesso do upload, remove os arquivos temporários do servidor.
  3. Atualiza o **banco de dados** com o status `concluído` e a URL do arquivo `.m3u8`.

### **4. Consulta do Status e Acesso ao Vídeo Convertido**
- **API HTTP** permite:
  - Consultar status do vídeo (`pendente`, `em processamento`, `concluído`).
  - Obter a URL do vídeo convertido no S3.

---

## **Contexto Geral do Projeto**
- **Nome do Módulo:** `github.com/devfullcycle/golangtechweek`
- **Versão do Go:** `1.24`
- **Objetivo do Projeto:**
  - Criar um sistema eficiente e escalável para conversão de vídeos em **HLS**, permitindo **upload via API HTTP**, processamento concorrente e armazenamento final na **AWS S3**.
  - A solução permite que aplicações enviem vídeos para serem convertidos automaticamente para streaming, facilitando a distribuição de mídia de forma otimizada.
  - Além da conversão, o sistema fornece um **endpoint HTTP** para consultar o status das conversões e recuperar a URL do vídeo convertido.

- **Principais Tecnologias:**
  - **Linguagem:** Go 1.24
  - **Processamento de Mídia:** `ffmpeg` para conversão de vídeos para HLS
  - **Armazenamento:** AWS S3 para guardar os segmentos e manifestos gerados
  - **Banco de Dados:** PostgreSQL para rastrear status das conversões
  - **Container:** Docker para criar o ambiente de execução tanto para aplicação em Go quanto para o banco de dados PostgreSQL através do docker-compose.yml
  - **Localstack**: Para testar o upload para S3 localmente
