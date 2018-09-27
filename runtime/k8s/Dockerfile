FROM ubuntu:bionic as bob
RUN apt-get update
RUN apt-get install -y zfsutils-linux libzfslinux-dev golang git
RUN git clone https://github.com/NimbleStorage/Nemo
WORKDIR Nemo
RUN make

WORKDIR /
RUN git clone https://github.com/hpe-storage/dory go/src/github.com/hpe-storage/dory
WORKDIR go/src/github.com/hpe-storage/dory
ENV GOPATH /go
RUN make tools && \
    make vendor && \
    make dory

FROM ubuntu:bionic
RUN apt-get update && \
    apt-get install -y zfsutils-linux libzfslinux-dev kmod && \
    mkdir -p /opt/nemo /var/lib/nemo /var/log/nemo && \
    rm -rf /var/cache /var/log
COPY --from=bob Nemo/nemo /usr/sbin/nemo
COPY --from=bob /go/src/github.com/hpe-storage/dory/dory /dory
ADD dory.json /
ADD entrypoint.sh /
ENTRYPOINT /entrypoint.sh
