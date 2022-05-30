FROM registry.access.redhat.com/ubi8/ubi:8.6-754
COPY yum.repos.d /etc/yum.repos.d/

RUN dnf -y install git \
    java-1.8.0-openjdk-devel \
    libicu \
    nmap-ncat \
    zip \
    unzip \
    jq \
    tar \
    sudo \
    procps-ng \
    make && \
    dnf clean all