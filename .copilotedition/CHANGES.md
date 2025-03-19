# Checklist de Alterações para WebAssembly + ONNX Runtime

## 1. Limpeza de Arquivos
- [x] Remover todos os arquivos Docker
- [x] Remover configurações gRPC
- [x] Simplificar estrutura de diretórios

## 2. Nova Estrutura
- [x] Criar diretório `web` para arquivos estáticos
- [x] Criar diretório `model` para arquivos ONNX
- [x] Reorganizar pacotes Go para suporte WebAssembly

## 3. Implementações
- [x] Converter servidor para WebAssembly
- [x] Implementar interface de áudio no navegador
- [x] Adicionar suporte ONNX Runtime
- [x] Implementar chat e gravação de áudio

## 4. Novos Recursos
- [x] Interface de chat
- [x] Gravação de áudio em tempo real
- [x] Inferência local usando ONNX Runtime
- [x] Conversão automática de modelos Whisper para ONNX

## 5. Próximos Passos
- [ ] Otimizar modelo ONNX
- [ ] Adicionar suporte offline
- [ ] Melhorar interface do usuário
- [ ] Adicionar testes unitários
