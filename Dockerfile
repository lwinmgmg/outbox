FROM lwinmgmg/alpine-git:latest AS gitCloner

ARG gitRepo=https://github.com/lwinmgmg/outbox.git
ARG gitRef=master

WORKDIR /build

RUN git clone --depth 1 --branch ${gitRef} ${gitRepo}

RUN ls -ahl

FROM golang:1.18-alpine

ENV OUTBOX_USER=outbox
ENV OUTBOX_USER_HOME_DIR="/home/${OUTBOX_USER}"
ENV OUTBOX_USER_UID=101
ENV OUTBOX_VERSION=1.18
ENV OUTBOX_INSTALL_DIR="${OUTBOX_USER_HOME_DIR}/${OUTBOX_VERSION}"

RUN addgroup -S ${OUTBOX_USER_UID} && adduser -S ${OUTBOX_USER} -G ${OUTBOX_USER_UID}

COPY --from=gitCloner --chown=${OUTBOX_USER}:${OUTBOX_USER_UID} /build/outbox ${OUTBOX_INSTALL_DIR}/

USER ${OUTBOX_USER}

WORKDIR ${OUTBOX_INSTALL_DIR}

RUN go mod tidy && go mod vendor

RUN go build

RUN touch settings.yaml

CMD [ "./outbox" ]
