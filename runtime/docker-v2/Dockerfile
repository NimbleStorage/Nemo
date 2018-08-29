FROM ubuntu:bionic as bob
RUN apt-get update
RUN apt-get install -y zfsutils-linux libzfslinux-dev golang git
RUN git clone https://github.com/NimbleStorage/Nemo
WORKDIR Nemo
RUN make

FROM ubuntu:bionic
RUN apt-get update && \
    apt-get install -y zfsutils-linux libzfslinux-dev kmod && \
    mkdir -p /opt/nemo /var/lib/nemo /var/log/nemo
COPY --from=bob Nemo/nemo /usr/sbin/nemo
ADD entrypoint.sh /
