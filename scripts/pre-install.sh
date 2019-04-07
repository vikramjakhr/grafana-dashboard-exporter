#!/bin/bash

if [[ -f /etc/gde/gde.conf ]]; then
    backup_name="gde.conf.$(date +%s).backup"
    echo "A backup of your current configuration can be found at: /etc/gde/${backup_name}"
    cp -a "/etc/gde/gde.conf" "/etc/gde/${backup_name}"
fi