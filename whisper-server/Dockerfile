FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    ffmpeg \
    python3 \
    g++ \
    libssl-dev

# Build whisper.cpp
RUN git clone https://github.com/ggerganov/whisper.cpp && \
    cd whisper.cpp && \
    make server

# Download model
RUN mkdir -p models && \
    curl -L -o models/ggml-large-v3.bin \
    https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-large-v3.bin

EXPOSE 9083
CMD ["/whisper.cpp/server", "-m", "/models/ggml-large-v3.bin", "--port", "9083", "--convert"]