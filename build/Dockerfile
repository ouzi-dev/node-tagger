FROM quay.io/ouzi/go-builder:1.13.7 as builder

WORKDIR /node-tagger

# download modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# add makefile
COPY Makefile .
RUN make setup

COPY . /node-tagger

RUN make test
RUN TARGETS=linux/amd64 make build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/node-tagger \
    USER_UID=1001 \
    USER_NAME=node-tagger

# install operator binary
COPY --from=builder /node-tagger/_dist/linux-amd64/manager/node-tagger ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}