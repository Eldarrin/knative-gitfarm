#!/bin/bash
buildah unshare
microcontainer=$(buildah from registry.access.redhat.com/ubi8/ubi-micro:8.6-285)
micromount=$(buildah mount $microcontainer)
dnf install \
    --installroot $micromount \
    --releasever 8.6 \
    --setopt install_weak_deps=false \
    --nodocs -y \
    golang \
    git \
    java-1.8.0-openjdk-devel \
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
    userspace-rcu
dnf clean all --installroot $micromount
buildah umount $microcontainer

mkdir tmp
cd tmp
curl -f -L -o runner.tar.gz https://github.com/actions/runner/releases/download/v2.291.1/actions-runner-linux-x64-2.291.1.tar.gz
tar xzf ./runner.tar.gz
rm runner.tar.gz
buildah config --workingdir /runner $microcontainer
buildah copy $microcontainer ./ /runner/
    
buildah run $microcontainer adduser runner
buildah run $microcontainer mkdir /opt/hostedtoolcache
buildah run $microcontainer chmod g+rwx /opt/hostedtoolcache
buildah config --port 8080 $microcontainer

buildah run $microcontainer chown -R runner:runner 
buildsh config --user runner $microcontainer

buildah copy $microcontainer go.mod go.sum webhook.go rand-strings.go ./

buildah run $microcontainer	go get github.com/go-playground/webhooks/v6
buildah run $microcontainer	go get github.com/google/go-github/v45
buildah run $microcontainer	go get golang.org/x/oauth2

buildah run $microcontainer CGO_ENABLED=0 go build -o webhook .

buildah config --entrypoint /runner/webhook $microcontainer

buildah unmount $microcontainer

buildah commit $microcontainer git-main

buildah rm $microcontainer