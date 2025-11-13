# Use the official Golang image as the base image
FROM debian:bookworm

# Set the working directory inside the container
WORKDIR /app

# Install Python and pip to install yt-dlp
RUN apt-get update && \
    apt-get install -y git wget2 cmake python3 ffmpeg g++
# TODO : build GOLANG from source

RUN rm -rf /usr/local/go
RUN wget2 https://go.dev/dl/go1.25.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

RUN wget2 https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp  # Make executable
RUN yt-dlp --version

# ------------------------------

RUN git clone https://github.com/ggml-org/whisper.cpp.git
WORKDIR /app/whisper.cpp

RUN sh ./models/download-ggml-model.sh base.en

RUN cmake -B build
RUN cmake --build build -j --config Release

# ------------------------------

WORKDIR /app

COPY . ./

COPY go.mod ./
RUN go mod tidy
RUN go mod vendor

# ------------------------------

# Default command (can be overridden)
CMD ["go", "run", "cmd/main.go"]