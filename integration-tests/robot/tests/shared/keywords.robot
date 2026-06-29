*** Variables ***
${ZOOKEEPER_HOST}  %{ZOOKEEPER_HOST}
${ZOOKEEPER_PORT}  %{ZOOKEEPER_PORT}
${ZOOKEEPER_OS_PROJECT}  %{ZOOKEEPER_OS_PROJECT}


*** Settings ***
Variables  %{ROBOT_HOME}/SecretData.py
Library  ./lib/ZookeeperLibrary.py  zookeeper_os_project=%{ZOOKEEPER_OS_PROJECT}
...                                 zookeeper_host=${ZOOKEEPER_HOST}
...                                 zookeeper_port=${ZOOKEEPER_PORT}
...                                 zookeeper_enable_ssl=%{ZOOKEEPER_ENABLE_TLS}
Library  PlatformLibrary  managed_by_operator=%{ZOOKEEPER_IS_MANAGED_BY_OPERATOR}


*** Keywords ***
Check Secret
    [Arguments]  ${secret_name}  ${zookeeper_os_project}
    ${response}=  Get Secret  ${secret_name}  ${zookeeper_os_project}
    RETURN  ${response}