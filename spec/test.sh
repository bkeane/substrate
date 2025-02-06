#!/bin/bash
set -e

# Globals
TEST_ID=$(date +%s)
HUB_AWS_PROFILE=prod.kaixo.io
SPOKE_AWS_PROFILE=dev.kaixo.io

# Hub Account
if [ -z "${HUB_AWS_PROFILE}" ]; then
    echo "Error: HUB_AWS_PROFILE environment variable must be set"
    exit 1
fi

HUB_ACCOUNT_ID=$(aws --profile $HUB_AWS_PROFILE sts get-caller-identity | jq -r .Account)
HUB_EVENTBUS_NAME=agar-backplane
HUB_EVENTBUS_LOGGING_RULE=debug-logging-route
HUB_EVENTBUS_LOG_GROUP=/aws/events/${HUB_EVENTBUS_NAME}

# Spoke Account
if [ -z "${SPOKE_AWS_PROFILE}" ]; then
    echo "Error: SUBSTRATE_SPOKE_AWS_PROFILE environment variable must be set"
    exit 1
fi

SPOKE_ACCOUNT_ID=$(aws --profile $SPOKE_AWS_PROFILE sts get-caller-identity | jq -r .Account)
SPOKE_EVENTBUS_NAME=agar-backplane
SPOKE_EVENTBUS_LOGGING_RULE=debug-logging-route
SPOKE_EVENTBUS_LOG_GROUP=/aws/events/${SPOKE_EVENTBUS_NAME}

TEST_EVENT_PAYLOAD='[{
    "Source": "'"${HUB_ACCOUNT_ID}"'",
    "DetailType": "TestEvent",
    "Detail": "{\"header\": {\"source\": \"'"${HUB_ACCOUNT_ID}"'\", \"destination\": [\"'"${SPOKE_ACCOUNT_ID}"'\"] }, \"message\": \"Test event '"${TEST_ID}"'\"}",
    "EventBusName": "'"${HUB_EVENTBUS_NAME}"'"
}]'

aws events --profile $HUB_AWS_PROFILE enable-rule --event-bus-name $HUB_EVENTBUS_NAME --name $HUB_EVENTBUS_LOGGING_RULE > /dev/null \
&& echo "✓ Successfully enabled hub debug logging" \
|| echo "✗ Failed to enable hub debug logging"

aws events --profile $SPOKE_AWS_PROFILE enable-rule --event-bus-name $SPOKE_EVENTBUS_NAME --name $SPOKE_EVENTBUS_LOGGING_RULE > /dev/null \
&& echo "✓ Successfully enabled spoke debug logging" \
|| echo "✗ Failed to enable spoke debug logging"

aws events --profile $HUB_AWS_PROFILE put-events --entries "$TEST_EVENT_PAYLOAD" > /dev/null \
&& echo "✓ Test event ${TEST_ID} successfully sent to hub event bus" \
|| echo "✗ Failed to send test event ${TEST_ID} to hub event bus"

echo "sent event..."
echo $TEST_EVENT_PAYLOAD | jq

sleep 5

# aws logs tail --profile $HUB_AWS_PROFILE $HUB_EVENTBUS_LOG_GROUP
aws logs tail --profile $HUB_AWS_PROFILE $HUB_EVENTBUS_LOG_GROUP | grep -q "Test event ${TEST_ID}" \
&& echo "✓ Found test event ${TEST_ID} in hub logs" \
|| echo "✗ Failed to find test event ${TEST_ID} in hub logs"

# aws logs tail --profile $SPOKE_AWS_PROFILE $SPOKE_EVENTBUS_LOG_GROUP
aws logs tail --profile $SPOKE_AWS_PROFILE $SPOKE_EVENTBUS_LOG_GROUP | grep -q "Test event ${TEST_ID}" \
&& echo "✓ Found test event ${TEST_ID} in spoke logs" \
|| echo "✗ Failed to find test event ${TEST_ID} in spoke logs"
