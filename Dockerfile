FROM ubuntu:20.04@sha256:450e066588f42ebe1551f3b1a535034b6aa46cd936fe7f2c6b0d72997ec61dbd

# fix timezone stalling during build
ENV TZ=America/Toronto
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    libssl-dev \
    libtool \
    make \
    pkg-config \
    wget \
    tmux

WORKDIR /app

# Setup Go programming language
RUN wget https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
RUN chmod +x /usr/local/go/bin/go
RUN chmod +x /usr/local/go/bin/gofmt
RUN cp /usr/local/go/bin/go /usr/bin/go
RUN cp /usr/local/go/bin/gofmt /usr/bin/gofmt

ENV GOLANG_VERSION 1.19.3

# build kaniko
COPY ./kaniko ./kaniko
RUN cd /app/kaniko && \
    make && \
    cp /app/kaniko/out/executor /usr/local/bin/kaniko
RUN rm -rf /app/kaniko

# copy over binutils
COPY ./binutils ./binutils

# copy over sample docker files
COPY ./sampledockerfiles ./sampledockerfiles
