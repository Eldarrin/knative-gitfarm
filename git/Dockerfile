FROM eldarrin/base-docker:latest

ARG RUNNER_VERSION=2.291.1
ARG DUMB_INIT_VERSION=1.2.5

ENV HOME=/home/runner
ENV RUNNER_ASSETS_DIR=/runnertmp
ENV RUNNER_TOOL_CACHE=/opt/hostedtoolcache

#WORKDIR /actions-runner

RUN ln -s /bin/microdnf /bin/yum || true && \
    mkdir -p "$RUNNER_ASSETS_DIR" \
    && cd "$RUNNER_ASSETS_DIR" \
    && curl -f -L -o runner.tar.gz https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz \
    && tar xzf ./runner.tar.gz \
    && rm runner.tar.gz \
    && ./bin/installdependencies.sh && \
    adduser runner \
    && curl -f -L -o /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v${DUMB_INIT_VERSION}/dumb-init_${DUMB_INIT_VERSION}_x86_64 \
    && chmod +x /usr/local/bin/dumb-init \
    && mkdir /opt/hostedtoolcache \
    && chmod g+rwx /opt/hostedtoolcache

COPY entrypoint.sh logger.bash /usr/bin/

RUN chmod +x /usr/bin/entrypoint.sh && \
    chmod +x /usr/bin/logger.bash

USER runner

ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]
CMD ["entrypoint.sh"]
