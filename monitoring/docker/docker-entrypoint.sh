#!/bin/bash

# Exit immediately if a *pipeline* returns a non-zero status. (Add -x for command tracing)
set -e
if [[ "$DEBUG" == true ]]; then
  set -x
fi

: "${ZOOKEEPER_PROMETHEUS_PORT:=8080}"

#Validate variable name that contains addresses.
#$1 - variable with addresses to validate
function validate_addresses() {
  variable_name=$1
  variable_value=${!variable_name}
  if [[ -z $variable_value ]]; then
    echo >&2 "Error: value of variable [$variable_name] is empty!"
    exit 1
  fi
  if [[ ! (${variable_value} =~ ^((\'.+\')(,*))+$) ]] ; then
    echo >&2 "Error: Value of variable [$variable_name] = [$variable_value] does not match address pattern!"
    exit 1
  fi
}

#Converts comma separated Prometheus addresses to urls.
#
#$1 - comma separated Prometheus addresses
function convert_to_prometheus_urls() {
  echo "${ZOOKEEPER_HOST//2181/$ZOOKEEPER_PROMETHEUS_PORT/metrics}" | sed "s/[[:blank:]]//g;s/^'/'http:\/\//;s/,'/,'http:\/\//g"
}

validate_addresses "ZOOKEEPER_HOST"

if [[ -n "$PROMETHEUS_URLS" ]]; then
  validate_addresses "PROMETHEUS_URLS"
else
  PROMETHEUS_URLS=$(convert_to_prometheus_urls "$ZOOKEEPER_HOST")
fi
export PROMETHEUS_URLS

mkdir -p "${MONITORING_LOGS}"

/sbin/tini -- /entrypoint.sh telegraf