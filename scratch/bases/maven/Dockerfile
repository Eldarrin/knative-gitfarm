FROM eldarrin/base-docker:LATESTTAG

ARG javaVer=11.0.15.0.10-2.el8_6

ENV JAVA_HOME_11_X64=/usr/lib/jvm/java-11-openjdk-$javaVer.x86_64 \
    JAVA_HOME=/usr/lib/jvm/java-11-openjdk-$javaVer.x86_64

RUN dnf -y install java-11-openjdk-devel-$javaVer \
    maven-1:3.5.4-5.module+el8+2452+b359bfcd && \
    dnf clean all && \
    alternatives --set java /usr/lib/jvm/java-11-openjdk-$javaVer.x86_64/bin/java