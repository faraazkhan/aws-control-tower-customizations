FROM ubuntu:21.04

MAINTAINER faraaz@samtek-inc.com

ENV DEBIAN_FRONTEND noninteractive

WORKDIR /usr/local/src/ci

COPY source/bin/ bin/

RUN apt-get update && apt-get install -y \
    apt-utils \
    wget \
    bzip2 \
    ca-certificates \
    sudo \
    locales \
    fonts-liberation \
    git \
    libzmq3-dev \
    libcurl4-openssl-dev \
    libssl-dev \
    curl \
    python3-pip \
    python3-dev \
    libtool \
    libffi-dev \
    ruby \
    ruby-dev \
    && cd /usr/local/bin \
    && ln -s /usr/bin/python3 python \
    && pip3 install --upgrade pip
RUN bin/codebuild_scripts/install_stage_dependencies.sh build  \
    && bin/codebuild_scripts/install_stage_dependencies.sh scp \
    && bin/codebuild_scripts/install_stage_dependencies.sh stackset
