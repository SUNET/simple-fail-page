FROM debian:stable

MAINTAINER lundberg <lundberg@sunet.se>

ENV DEBIAN_FRONTEND noninteractive

RUN /bin/sed -i s/deb.debian.org/ftp.se.debian.org/g /etc/apt/sources.list

RUN apt-get update && \
    apt-get -y dist-upgrade && \
    apt-get install -y \
      sed \
      git \
      iputils-ping \
    && apt-get clean

RUN rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/simple-fail-page
COPY index.html /opt/simple-fail-page/index.html
COPY simple-fail-page.yaml /opt/simple-fail-page/simple-fail-page.yaml
COPY assets/ /opt/simple-fail-page/assets/

COPY simple-fail-page.linux.amd64 /usr/local/sbin/simple-fail-page

VOLUME ["/opt/simple-fail-page"]

WORKDIR /opt/simple-fail-page
CMD /usr/local/sbin/simple-fail-page
