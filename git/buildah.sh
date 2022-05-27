microcontainer=$(buildah from registry.access.redhat.com/ubi8/ubi-micro)
micromount=$(buildah mount $microcontainer)
yum install \
    --installroot $micromount \
    --releasever 8 \
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
    make
yum clean all \
    --installroot $micromount
buildah umount $microcontainer
buildah commit $microcontainer git-main