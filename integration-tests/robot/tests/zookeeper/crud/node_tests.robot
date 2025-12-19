*** Variables ***
${ZOOKEEPER_CRUD_NODE_PATH}   /zookeeper_crud
${ZOOKEEPER_TESTS_NODE_PATH}  /zookeeper_crud/tests
${CREATION_DATA}              Creation data
${MODIFICATION_DATA}          Modification data

*** Settings ***
Library  String
Library  RetryFailed
Resource  ../../shared/keywords.robot

*** Keywords ***
Setup
    ${zk} =  Connect To Zookeeper
    Set Suite Variable  ${zk}
    Delete Node  ${zk}  ${ZOOKEEPER_CRUD_NODE_PATH}

Check Existence Of Node
    [Arguments]  ${path_to_node}
    ${node} =  Node Exists  ${zk}  ${path_to_node}
    Should Be True  ${node}

Check Absence Of Node
    [Arguments]  ${path_to_node}
    ${node} =  Node Exists  ${zk}  ${path_to_node}
    Should Not Be True  ${node}

Test Node Creation
    Log To Console  \nChecking Node Creation
    Check Absence Of Node  ${ZOOKEEPER_CRUD_NODE_PATH}
    Create Node  ${zk}  ${ZOOKEEPER_CRUD_NODE_PATH}  ${CREATION_DATA}
    Check Existence Of Node  ${ZOOKEEPER_CRUD_NODE_PATH}
    Create Node  ${zk}  ${ZOOKEEPER_TESTS_NODE_PATH}  ${CREATION_DATA}
    Check Existence Of Node  ${ZOOKEEPER_TESTS_NODE_PATH}

Test Reading Data
    Log To Console  Checking Reading Data
    ${data} =  Get Node Value  ${zk}  ${ZOOKEEPER_TESTS_NODE_PATH}
    Should Be Equal As Strings  ${data}  ${CREATION_DATA}

Test Updating Data
    Log To Console  Checking Updating Data
    Update Node Value  ${zk}  ${ZOOKEEPER_TESTS_NODE_PATH}  ${MODIFICATION_DATA}
    ${data} =  Get Node Value  ${zk}  ${ZOOKEEPER_TESTS_NODE_PATH}
    Should Be Equal As Strings  ${data}  ${MODIFICATION_DATA}

Test Node Deletion
    Log To Console  Checking Node Deletion
    Delete Node  ${zk}  ${ZOOKEEPER_CRUD_NODE_PATH}
    Check Absence Of Node  ${ZOOKEEPER_CRUD_NODE_PATH}

Cleanup
    Delete Node  ${zk}  ${ZOOKEEPER_CRUD_NODE_PATH}
    Disconnect From Zookeeper  ${zk}
    ${zk} =  Set Variable  ${None}

*** Test Cases ***
Test Zookeeper CRUD
    [Tags]  zookeeper_crud  zookeeper  test:retry(3)
    [Setup]  Setup
    Test Node Creation
    Test Reading Data
    Test Updating Data
    Test Node Deletion
    [Teardown]  Cleanup
