#DEPRECATED FOR BUILDAH
FROM registry.access.redhat.com/ubi8/ubi:8.6-754

ARG RUNNER_VERSION=2.291.1

ENV HOME=/home/runner
ENV RUNNER_HOME=/runner
ENV RUNNER_TOOL_CACHE=/opt/hostedtoolcache

WORKDIR /runner

RUN dnf install -y golang \
    git \
    libicu \
    nmap-ncat \
    zip \
    unzip \
    jq \
    tar \
    sudo \
    make && \
    dnf clean all \
    && curl -f -L -o runner.tar.gz https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz \
    && tar xzf ./runner.tar.gz \
    && rm runner.tar.gz \
    && ./bin/installdependencies.sh \
    && adduser runner \
    && mkdir /opt/hostedtoolcache \
    && chmod g+rwx /opt/hostedtoolcache

EXPOSE 8080

RUN chown -R runner:runner "$RUNNER_HOME"

USER runner

COPY go.mod go.sum webhook.go rand-strings.go ./

RUN	go get github.com/go-playground/webhooks/v6
RUN	go get github.com/google/go-github/v45
RUN	go get golang.org/x/oauth2

RUN CGO_ENABLED=0 go build -o webhook .

ENTRYPOINT ["/runner/webhook"]
