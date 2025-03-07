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
    curl -L -o models/large-v3-turbo-q5_0.bin \
    https://huggingface.co/ggerganov/whisper.cpp/resolve/main/large-v3-turbo-q5_0.bin

EXPOSE 9083
CMD ["/whisper.cpp/server", "-m", "/models/large-v3-turbo-q5_0.bin", "--port", "9083", "--convert"]