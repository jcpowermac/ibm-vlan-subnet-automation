#!/bin/bash

# Create configuration mode session
LOCATION=$(/usr/bin/https POST "${VYATTA_URL}/rest/conf" --auth vyatta:${VYATTA_PASSWORD} --headers --verify no | grep Location | awk '{print $2}' | tr -d '[:cntrl:]')

PATH="%2Fhome%2Fvyatta%2F${CONFIG_FILE}"

/usr/bin/https POST "${VYATTA_URL}/${LOCATION}/merge/${PATH}" --auth vyatta:${VYATTA_PASSWORD} --verify no --verbose
/usr/bin/https POST "${VYATTA_URL}/${LOCATION}/commit" --auth vyatta:${VYATTA_PASSWORD} --verify no --verbose
/usr/bin/https DELETE "${VYATTA_URL}/${LOCATION}" --auth vyatta:${VYATTA_PASSWORD} --verify no --verbose
