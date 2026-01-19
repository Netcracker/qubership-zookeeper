FROM ghcr.io/netcracker/qubership-backup-daemon:main_alpine

ENV ZOOKEEPER_HOME=/opt/zookeeper \
    ZOOKEEPER_BACKUP=/opt/zookeeper/backup-storage

ENV PYTHONPATH=/usr/local/lib/python3.7/site-packages/integration_library_builtIn

RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.23/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.23/community" >> /etc/apk/repositories

COPY backup-daemon.conf /etc/backup-daemon.conf
COPY docker/requirements.txt ${ZOOKEEPER_HOME}/requirements.txt
COPY scripts/ ${ZOOKEEPER_HOME}/scripts/
COPY sh_scripts/ ${ZOOKEEPER_HOME}/sh_scripts/

RUN set -x \
    && apk add --upgrade --no-cache curl apk-tools build-base \
    && pip3 install -r ${ZOOKEEPER_HOME}/requirements.txt \
    && apk del build-base \
    && rm -rf /var/cache/apk/*

# Upgrade all tools to avoid vulnerabilities
RUN set -x && apk upgrade --available

# Upgrade tools from edge to avoid vulnerabilities
RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/main" > /etc/apk/repositories  \
    && echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && apk add --update --no-cache --upgrade \
        pcre2 curl gnutls libcurl jq

RUN chmod +x ${ZOOKEEPER_HOME}/scripts/*.py
RUN chmod +x ${ZOOKEEPER_HOME}/sh_scripts/list_instances_in_vault_command.sh

USER 1000:0
WORKDIR ${ZOOKEEPER_HOME}

VOLUME ["${ZOOKEEPER_BACKUP}"]
