FROM eldarrin/base-main:LATESTTAG

RUN dnf -y install yum

RUN useradd build; yum -y update; rpm --restore shadow-utils 2>/dev/null; yum -y install cpp buildah fuse-overlayfs xz --exclude container-selinux; rm -rf /var/cache /var/log/dnf* /var/log/yum.*;

COPY containers.conf /etc/containers/

# Adjust storage.conf to enable Fuse storage.
RUN chmod 644 /etc/containers/containers.conf; sed -i -e 's|^#mount_program|mount_program|g' -e '/additionalimage.*/a "/var/lib/shared",' -e 's|^mountopt[[:space:]]*=.*$|mountopt = "nodev,fsync=0"|g' /etc/containers/storage.conf
RUN mkdir -p /var/lib/shared/overlay-images /var/lib/shared/overlay-layers /var/lib/shared/vfs-images /var/lib/shared/vfs-layers; touch /var/lib/shared/overlay-images/images.lock; touch /var/lib/shared/overlay-layers/layers.lock; touch /var/lib/shared/vfs-images/images.lock; touch /var/lib/shared/vfs-layers/layers.lock

# Define uid/gid ranges for our user https://github.com/containers/buildah/issues/3053
RUN echo -e "build:1:999\nbuild:1001:64535" > /etc/subuid; \
 echo -e "build:1:999\nbuild:1001:64535" > /etc/subgid; \
 mkdir -p /home/build/.local/share/containers; \
 chown -R build:build /home/build

VOLUME /var/lib/containers
VOLUME /home/build/.local/share/containers

# Set an environment variable to default to chroot isolation for RUN
# instructions and "buildah run".
ENV BUILDAH_ISOLATION=chroot

RUN ln -s /bin/buildah /bin/docker && \
    mkdir -p /run/containers/0
