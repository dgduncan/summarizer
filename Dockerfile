# Use the official Golang image as the base image
FROM debian:bookworm

# Set the working directory inside the container
WORKDIR /app

# Install Python and pip to install yt-dlp
RUN apt-get update && \
    apt-get install -y git wget2 cmake python3 ffmpeg g++
# TODO : build GOLANG from source

RUN rm -rf /usr/local/go
RUN ARCH=$(dpkg --print-architecture) && \
    wget2 https://go.dev/dl/go1.23.3.linux-${ARCH}.tar.gz && \
    tar -C /usr/local -xzf go1.23.3.linux-${ARCH}.tar.gz && \
    rm go1.23.3.linux-${ARCH}.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

RUN wget2 https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp  # Make executable
RUN yt-dlp --version

# ------------------------------

RUN git clone https://github.com/ggml-org/whisper.cpp.git --depth 1
WORKDIR /app/whisper.cpp

RUN sh ./models/download-ggml-model.sh small.en-tdrz
RUN sh ./models/download-ggml-model.sh small.en
RUN sh ./models/download-ggml-model.sh base.en
RUN sh ./models/download-ggml-model.sh tiny.en

RUN cmake -B build
RUN cmake --build build -j --config Release
# RUN ./build/bin/quantize models/ggml-small.en.bin models/ggml-small.en-q8_0.bin q8_0
# RUN ./build/bin/quantize models/ggml-base.en.bin models/ggml-base.en-q5_0.bin q5_0
# RUN ./build/bin/quantize models/ggml-tiny.en.bin models/ggml-tiny.en-q4_0.bin q4_0

# ------------------------------

WORKDIR /app

COPY . ./

COPY go.mod ./
RUN go mod download

# ------------------------------

# Default command (can be overridden)
CMD ["go", "run", "cmd/main.go"]


# # ------------------------------
# # Stage 1: Builder
# # ------------------------------
# FROM golang:1.24-bookworm AS builder

# # Install build tools
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     cmake \
#     g++ \
#     git \
#     wget \
#     && rm -rf /var/lib/apt/lists/*

# # 1. Build whisper.cpp
# WORKDIR /build/whisper.cpp
# # Clone only the latest commit to save space
# RUN git clone https://github.com/ggml-org/whisper.cpp.git . --depth 1

# # Download models
# RUN sh ./models/download-ggml-model.sh small.en && \
#     sh ./models/download-ggml-model.sh base.en && \
#     sh ./models/download-ggml-model.sh tiny.en

# # Compile whisper
# RUN cmake -B build -DCMAKE_BUILD_TYPE=Release && \
#     cmake --build build -j --config Release

# # 2. Build Go Application
# WORKDIR /build/app

# # Download deps first for better caching
# COPY go.mod ./
# # If you have go.sum, uncomment the next line
# # COPY go.sum ./ 
# RUN go mod download

# # Copy source and build binary
# COPY . .
# # Build a binary named 'main'
# RUN go build -o main cmd/main.go

# # 3. Download yt-dlp
# RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp && \
#     chmod a+rx /usr/local/bin/yt-dlp

# # ------------------------------
# # Stage 2: Final Runtime
# # ------------------------------
# FROM debian:bookworm-slim as final

# WORKDIR /app

# # Install runtime dependencies.
# # libgomp1 is CRITICAL for whisper.cpp (OpenMP support).
# # libstdc++6 is required for C++ binaries.
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     ffmpeg \
#     python3 \
#     ca-certificates \
#     libgomp1 \
#     libstdc++6 \
#     && rm -rf /var/lib/apt/lists/*

# # Copy yt-dlp from builder
# COPY --from=builder /usr/local/bin/yt-dlp /usr/local/bin/yt-dlp

# # Copy compiled Go binary
# COPY --from=builder /build/app/main ./main

# # Copy whisper models and built binaries
# # Preserving the directory structure likely expected by your app
# RUN mkdir -p whisper.cpp/build/bin
# COPY --from=builder /build/whisper.cpp/models ./whisper.cpp/models
# COPY --from=builder /build/whisper.cpp/build/bin ./whisper.cpp/build/bin

# # Run the binary
# CMD ["./main"]