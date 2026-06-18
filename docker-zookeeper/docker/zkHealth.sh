#!/bin/bash

need_to_run_backup_server() {
  local is_process_present
  is_process_present=$(pgrep -a -f "/zookeeper-assistant -c backup" | grep "${ZOOKEEPER_BACKUP_SOURCE_DIR}")
  echo "$is_process_present"
  if [[ -z "${is_process_present}" ]]; then
    return 0
  fi
  return 1
}

run_backup_server() {
  mkdir -p /tmp/zookeeper
  nohup /zookeeper-assistant -c backup -s "${ZOOKEEPER_BACKUP_SOURCE_DIR}" -d "${ZOOKEEPER_BACKUP_DESTINATION_DIR}" &> /tmp/zookeeper/output_backup_server.log &
  while true; do
    curl -s -X GET "http://localhost:8081/" 2>/tmp/zookeeper/server_errors.log
    if grep -q "Failed to connect to localhost port 8081" /tmp/zookeeper/server_errors.log; then
      sleep 1s
    else
      echo "" > /tmp/zookeeper/server_errors.log
      break
    fi
  done
}

liveness() {
  if [[ ! -d ${ZOOKEEPER_DATA}/logs ]];then
    mkdir "${ZOOKEEPER_DATA}"/logs
  fi

  local output
  output="Check health for localhost. time: $(date)"

  local success=false
  for attempt in {1..5}; do
    local result
    if [[ "${ENABLE_SSL}" == "true" ]]; then
      if [[ "${ENABLE_2WAY_SSL}" == "true" ]]; then
        result=$(echo "ruok" | openssl s_client -crlf -quiet -connect localhost:2181 -cert /opt/zookeeper/tls/tls.crt -key /opt/zookeeper/tls/tls.key -CAfile /opt/zookeeper/tls/ca.crt 2>/dev/null)
      else
        result=$(echo "ruok" | openssl s_client -crlf -quiet -connect localhost:2181 -CAfile /opt/zookeeper/tls/ca.crt 2>/dev/null)
      fi
    else
      result=$(echo ruok | nc -w 2 -q 1 localhost 2181)
    fi
    output="${output}\n- ruok attempt ${attempt} - time: $(date) - result: $result"
    if echo "${result}" | grep -q "imok"; then
      success=true
      break
    fi
  done

  if [[ ${success} == true ]]; then
    for attempt in {1..5}; do
      local result
      if [[ "${ENABLE_SSL}" == "true" ]]; then
        if [[ "${ENABLE_2WAY_SSL}" == "true" ]]; then
          result=$(echo "srvr" | openssl s_client -crlf -quiet -connect localhost:2181 -cert /opt/zookeeper/tls/tls.crt -key /opt/zookeeper/tls/tls.key -CAfile /opt/zookeeper/tls/ca.crt 2>/dev/null)
        else
          result=$(echo "srvr" | openssl s_client -crlf -quiet -connect localhost:2181 -CAfile /opt/zookeeper/tls/ca.crt 2>/dev/null)
        fi
      else
        result=$(echo srvr | nc -w 2 -q 1 localhost 2181)
      fi
      output="${output}\n- srvr attempt ${attempt} - time: $(date) - result: $result"
      if echo "${result}" | grep -P -q "Mode: (follower|leader|standalone)"; then
        return 0
      fi
    done
  fi

  output="${output}\n- failed"
  echo -e "${output}" > "${ZOOKEEPER_DATA}"/logs/health.log
  return 1
}

readiness() {
  if [[ -n ${ADMIN_USERNAME} && -n ${ADMIN_PASSWORD} ]]; then
    if /zookeeper-assistant -c health -u "${ADMIN_USERNAME}:${ADMIN_PASSWORD}"; then
      return 0
    fi
  else
    if /zookeeper-assistant -c health; then
      return 0
    fi
  fi
  echo -e  "ZooKeeper connection failed. time: $(date) - result: ZooKeeper Quorum does not work" >> "${ZOOKEEPER_DATA}"/logs/health.log
  return 1
}

if need_to_run_backup_server; then
  run_backup_server
fi

case $1 in
health)
    liveness
    exit $?
    ;;
liveness-probe)
    liveness
    exit $?
    ;;
readiness-probe)
    readiness
    exit $?
    ;;
esac
