#!/usr/bin/env sh

if [ "${VIKINGR_TOKEN_FILE}" != "" ]; then
    export VIKINGR_TOKEN=$(cat "${VIKINGR_TOKEN_FILE}")
fi

/usr/bin/vikingr $@
