FROM ghcr.io/netcracker/qubership-backup-daemon-go:feat_sensitive_data_parsing

ENV ZOOKEEPER_HOME=/opt/zookeeper \
    ZOOKEEPER_BACKUP=/opt/zookeeper/backup-storage

RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.23/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.23/community" >> /etc/apk/repositories

COPY backup-daemon.conf /etc/backup-daemon.conf
COPY docker/requirements.txt ${ZOOKEEPER_HOME}/requirements.txt
COPY scripts/ ${ZOOKEEPER_HOME}/scripts/
COPY sh_scripts/ ${ZOOKEEPER_HOME}/sh_scripts/

RUN set -x \
    && apk add --update --no-cache \
    python3 \
    py3-pip \
    && if pip3 install --help | grep -q -- '--break-system-packages'; then PIP_BREAK_FLAG='--break-system-packages'; else PIP_BREAK_FLAG=''; fi \
    && pip3 install --no-cache-dir ${PIP_BREAK_FLAG} setuptools \
    && pip3 install --no-cache-dir ${PIP_BREAK_FLAG} -r ${ZOOKEEPER_HOME}/requirements.txt \
    && rm -rf /var/cache/apk/*

# Upgrade all tools to avoid vulnerabilities
RUN set -x && apk upgrade --no-cache --available

RUN chmod +x ${ZOOKEEPER_HOME}/scripts/*.py
RUN chmod +x ${ZOOKEEPER_HOME}/sh_scripts/list_instances_in_vault_command.sh

USER 1000:0
VOLUME ["${ZOOKEEPER_BACKUP}"]
