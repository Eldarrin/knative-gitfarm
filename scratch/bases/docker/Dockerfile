FROM eldarrin/base-main:LATESTTAG

ENV LABELS=docker

COPY yum.repos.d /etc/yum.repos.d/

COPY daemon.json /etc/docker/daemon.json

RUN dnf -y install docker-ce-3:18.09.1-3.el7 net-tools && \
    dnf clean all
