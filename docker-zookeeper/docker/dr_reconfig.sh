#!/bin/bash

SECRETS_DIR="${ZOOKEEPER_SECRETS_DIR:-/etc/secrets/zookeeper-pod-secrets}"

resolve_secret_value() {
  local secret_key="$1"
  local env_var_name="$2"
  local secret_path="${SECRETS_DIR}/${secret_key}"
  if [[ -r "${secret_path}" ]]; then
    tr -d '\r' < "${secret_path}"
    return 0
  fi
  printf "%s" "${!env_var_name:-}"
}

action=$1
dr_active_side=$2

# Use writable config dir set by entrypoint; fall back to /tmp/zookeeper/config when called standalone
: "${ZK_CONF:=/tmp/zookeeper/config}"
mkdir -p "${ZK_CONF}"

if [[ -z "$SERVER_NAME" ]]; then
  SERVER_NAME="zookeeper"
fi
if [[ -z "$SERVER_ID" ]]; then
  SERVER_ID="1"
fi
if [[ -n "$SERVER_DOMAIN" ]]; then
  SERVER_DOMAIN=".$SERVER_DOMAIN"
fi
if [[ -n "$SERVER_NAMESPACE" ]]; then
  SERVER_NAMESPACE=".$SERVER_NAMESPACE"
fi
if [[ -z "$SERVER_COUNT" ]]; then
  SERVER_COUNT=1
fi

echo "# Server List" > "${ZK_CONF}/zoo.cfg.dynamic"
left_group=""
if [[ -n "$SERVER_DOMAIN" ]]; then
  SERVER_DOMAIN=".left-${DR_SERVER_NAME_SUFFIX}-server"
fi
for (( order_number=1; order_number <= SERVER_COUNT; order_number++ )); do
  echo "server.$order_number=left-$DR_SERVER_NAME_SUFFIX-$order_number$SERVER_DOMAIN$SERVER_NAMESPACE:2888:3888:participant;2181" >> "${ZK_CONF}/zoo.cfg.dynamic"
  if [[ -n ${left_group} ]]; then
      left_group="${left_group}:"
  fi
  left_group="${left_group}${order_number}"
done
right_group=""
if [[ -n "$SERVER_DOMAIN" ]]; then
  SERVER_DOMAIN=".right-${DR_SERVER_NAME_SUFFIX}-server"
fi
for (( order_number=1; order_number <= SERVER_COUNT; order_number++ )); do
  server_number=$(( order_number + SERVER_COUNT ))
  echo "server.$server_number=right-$DR_SERVER_NAME_SUFFIX-$order_number$SERVER_DOMAIN$SERVER_NAMESPACE:2888:3888:participant;2181" >> "${ZK_CONF}/zoo.cfg.dynamic"
  if [[ -n ${right_group} ]]; then
      right_group="${right_group}:"
  fi
  right_group="${right_group}${server_number}"
done
echo "group.1=${left_group}" >> "${ZK_CONF}/zoo.cfg.dynamic"
echo "group.2=${right_group}" >> "${ZK_CONF}/zoo.cfg.dynamic"

left_weight=1
right_weight=1
if [[ ${dr_active_side} == "left" ]]; then
    right_weight=0
fi
if [[ ${dr_active_side} == "right" ]]; then
    left_weight=0
fi
for (( order_number=1; order_number <= SERVER_COUNT; order_number++ )); do
  echo "weight.${order_number}=${left_weight}" >> "${ZK_CONF}/zoo.cfg.dynamic"
done
for (( order_number=SERVER_COUNT + 1; order_number <= SERVER_COUNT + SERVER_COUNT; order_number++ )); do
  echo "weight.${order_number}=${right_weight}" >> "${ZK_CONF}/zoo.cfg.dynamic"
done

if [[ ${action} == "move" ]]; then
  ADMIN_USERNAME="$(resolve_secret_value "admin-username" "ADMIN_USERNAME")"
  ADMIN_PASSWORD="$(resolve_secret_value "admin-password" "ADMIN_PASSWORD")"
  if [[ -n ${ADMIN_USERNAME} ]] && [[ -n ${ADMIN_PASSWORD} ]]; then
    cat > "${ZK_CONF}/client_jaas.conf" << EOL
Client {
           org.apache.zookeeper.server.auth.DigestLoginModule required
           username="${ADMIN_USERNAME}"
           password="${ADMIN_PASSWORD}";
    };
EOL
    export CLIENT_JVMFLAGS="-Djava.security.auth.login.config=${ZK_CONF}/client_jaas.conf" && \
    ./bin/zkCli.sh reconfig -file "${ZK_CONF}/zoo.cfg.dynamic"
  else
    ./bin/zkCli.sh reconfig -file "${ZK_CONF}/zoo.cfg.dynamic"
  fi
fi
