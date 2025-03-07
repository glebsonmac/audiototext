FROM ubuntu:22.04

# 1. Instalar dependências
RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    ffmpeg \
    python3 \
    g++ \
    cmake \
    libopenblas-dev

# 2. Clonar repositório com versão específica
RUN git clone https://github.com/ggerganov/whisper.cpp && \
    cd whisper.cpp && \
    git checkout v1.5.0  # Versão estável

# 3. Compilar com CMake
RUN cd whisper.cpp && \
    mkdir build && \
    cd build && \
    cmake .. && \
    make -j4

# 4. Baixar modelo
RUN mkdir -p models && \
    curl -L -o models/ggml-large-v3-turbo-q5_0.bin \
    https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v3-turbo-q5_0.bin

# 5. Expor porta e configurar entrypoint
EXPOSE 9083
CMD ["/whisper.cpp/build/bin/main", "-m", "/models/ggml-large-v3-turbo-q5_0.bin", "--server", "--port", "9083"]