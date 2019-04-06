#!/bin/bash

BIN_DIR=/usr/bin
DATA_DIR=/var/lib/gde
LOG_DIR=/var/log/gde
SCRIPT_DIR=/usr/lib/gde/scripts

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/gde
    chmod +x /etc/init.d/gde
}

function install_systemd {
    cp -f $SCRIPT_DIR/gde.service /lib/systemd/system/gde.service
}

function enable_systemd {
    systemctl enable gde
}

function enable_update_rcd {
    update-rc.d gde defaults
}

function enable_chkconfig {
    chkconfig --add gde
}

if ! id gde >/dev/null 2>&1; then
    useradd --system -U -M gde -s /bin/false -d $DATA_DIR
fi
chmod a+rX $BIN_DIR/gde*

mkdir -p $LOG_DIR
chown -R -L gde:gde $LOG_DIR
mkdir -p $DATA_DIR
chown -R -L gde:gde $DATA_DIR

test -f /etc/default/gde || touch /etc/default/gde

# Distribution-specific logic
if [[ -f /etc/redhat-release ]]; then
    # RHEL-variant logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd
        # Do not enable service
    else
        # Assuming SysV
        install_init
        # Do not enable service
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd
        enable_systemd
    else
        # Assuming SysV
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            enable_update_rcd
        else
            enable_chkconfig
        fi
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
        # Amazon Linux logic
        install_init
        # Do not enable service
    fi
fi
