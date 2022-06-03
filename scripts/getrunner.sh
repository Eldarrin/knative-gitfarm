#!/bin/bash
mkdir tmp && cd tmp || true
if [ ! -f "config.sh" ]; then
  curl -f -L -o runner.tar.gz https://github.com/actions/runner/releases/download/v2.291.1/actions-runner-linux-x64-2.291.1.tar.gz
  tar xzf ./runner.tar.gz
  rm -f runner.tar.gz
fi
cp -r $1