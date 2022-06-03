#!/bin/bash
# must do this before running "buildah unshare"
# test uid here
microcontainer=$(buildah from registry.access.redhat.com/ubi9-micro:9.0.0-11)
micromount=$(buildah mount "$microcontainer")

yum install \
    --installroot "$micromount" \
    --releasever 9 \
    --setopt install_weak_deps=false \
    --nodocs -y \
    golang \
    git \
    nmap-ncat \
    zip \
    unzip \
    jq \
    tar \
    sudo \
    procps-ng \
    make \
    lttng-ust \
    openssl-libs \
    krb5-libs \
    zlib \
    libicu --nogpgcheck
yum clean all --installroot "$micromount"

mkdir "$micromount/runner"
getrunner.sh "$micromount/runner/"
buildah config --workingdir /runner "$microcontainer"

buildah run "$microcontainer" adduser runner
buildah run "$microcontainer" mkdir /opt/hostedtoolcache
buildah run "$microcontainer" chmod g+rwx /opt/hostedtoolcache
buildah config --port 8080 "$microcontainer"

gobuild.sh "$micromount/runner/runnercontrol"
chmod +x "$micromount/runner/runnercontrol"

buildah run "$microcontainer" chown -R runner:runner /runner
buildah config --user runner "$microcontainer"

buildah config --entrypoint /runner/runnercontrol "$microcontainer"

buildah umount "$microcontainer"

buildah commit "$microcontainer" git-main

buildah rm "$microcontainer"
