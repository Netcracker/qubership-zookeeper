*** Variables ***
${MONITORED_IMAGES}         %{MONITORED_IMAGES}

*** Settings ***
Library  String
Library  Collections
Library  PlatformLibrary  managed_by_operator=%{ZOOKEEPER_IS_MANAGED_BY_OPERATOR}

*** Keywords ***
Get Image Tag
    [Arguments]  ${image}
    ${parts}=  Split String  ${image}  :
    ${length}=  Get Length  ${parts}
    Run Keyword If  ${length} > 1  Return From Keyword  ${parts[2]}  ELSE  Fail  Image has no tag: ${image}

*** Test Cases ***
Test Hardcoded Images
  [Tags]  zookeeper  zookeeper_images
  ${stripped_resources}=  Strip String  ${MONITORED_IMAGES}  characters=,  mode=right
  @{list_resources} =  Split String	${stripped_resources} 	,
  FOR  ${resource}  IN  @{list_resources}
    ${type}  ${name}  ${container_name}  ${image}=  Split String  ${resource}
    ${resource_image}=  Get Resource Image  ${type}  ${name}  %{OS_PROJECT}  ${container_name}

    Log To Console  resource_image: '${resource_image}'
    IF    ${resource_image} == not_found
        Log To Console    Monitored images list: ${MONITORED_IMAGES}
        Fail    Some images are not found, please check .helpers template and description.yaml in delivery
    END

    ${expected_tag}=  Get Image Tag  ${image}
    ${actual_tag}=    Get Image Tag  ${resource_image}

    Log To Console  \n[COMPARE] ${resource}: Expected tag = ${expected_tag}, Actual tag = ${actual_tag}

    Run Keyword And Continue On Failure  Should Be Equal   ${actual_tag}   ${expected_tag}
    
  END
