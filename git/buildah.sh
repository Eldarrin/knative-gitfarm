#!/bin/bash
#buildah unshare
microcontainer=$(buildah from registry.access.redhat.com/ubi9-micro:9.0.0-11)
micromount=$(buildah mount $microcontainer)

dnf install \
    --installroot $micromount \
    --releasever 9 \
    --setopt install_weak_deps=false \
    --nodocs -y \
    golang \
    git \
    java-17-openjdk \
    libicu \
    nmap-ncat \
    zip \
    unzip \
    jq \
    tar \
    sudo \
    procps-ng \
    make \
    lttng-ust \
    userspace-rcu --nogpgcheck
dnf clean all --installroot $micromount

mkdir $micromount/runner

mkdir tmp
cd tmp
curl -f -L -o runner.tar.gz https://github.com/actions/runner/releases/download/v2.291.1/actions-runner-linux-x64-2.291.1.tar.gz
tar xzf ./runner.tar.gz
rm runner.tar.gz
cp -r * $micromount/runner/
buildah config --workingdir /runner $microcontainer
cd ..
    
buildah run $microcontainer useradd runner
buildah run $microcontainer mkdir /opt/hostedtoolcache
buildah run $microcontainer chmod g+rwx /opt/hostedtoolcache
buildah config --port 8080 $microcontainer

go get github.com/go-playground/webhooks/v6
go get github.com/google/go-github/v45
go get golang.org/x/oauth2

CGO_ENABLED=0 go build -o webhook .
cp webhook $micromount/runner/webhook

buildah run $microcontainer chown -R runner:runner /runner
buildsh config --user runner $microcontainer

buildah config --entrypoint /runner/webhook $microcontainer

buildah umount $microcontainer

buildah commit $microcontainer git-main

buildah rm $microcontainer
