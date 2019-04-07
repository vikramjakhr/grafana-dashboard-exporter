#!/bin/bash

BIN_DIR=/usr/bin

# Distribution-specific logic
if [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        deb-systemd-invoke stop gde.service
    else
        # Assuming sysv
        invoke-rc.d gde stop
    fi
fi