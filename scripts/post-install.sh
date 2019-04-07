#!/bin/bash

BIN_DIR=/usr/bin
LOG_DIR=/var/log/gde
SCRIPT_DIR=/usr/lib/gde/scripts
LOGROTATE_DIR=/etc/logrotate.d

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/gde
    chmod +x /etc/init.d/gde
}

function install_systemd {
    cp -f $SCRIPT_DIR/gde.service $1
    systemctl enable gde || true
    systemctl daemon-reload || true
}

function install_update_rcd {
    update-rc.d gde defaults
}

function install_chkconfig {
    chkconfig --add gde
}

if ! grep "^gde:" /etc/group &>/dev/null; then
    groupadd -r gde
fi

if ! id gde &>/dev/null; then
    useradd -r -M gde -s /bin/false -d /etc/gde -g gde
fi

test -d $LOG_DIR || mkdir -p $LOG_DIR
chown -R -L gde:gde $LOG_DIR
chmod 755 $LOG_DIR

# Remove legacy symlink, if it exists
if [[ -L /etc/init.d/gde ]]; then
    rm -f /etc/init.d/gde
fi
# Remove legacy symlink, if it exists
if [[ -L /etc/systemd/system/gde.service ]]; then
    rm -f /etc/systemd/system/gde.service
fi

# Add defaults file, if it doesn't exist
if [[ ! -f /etc/default/gde ]]; then
    touch /etc/default/gde
fi

# Add .d configuration directory
if [[ ! -d /etc/gde/gde.d ]]; then
    mkdir -p /etc/gde/gde.d
fi

# Distribution-specific logic
if [[ -f /etc/redhat-release ]] || [[ -f /etc/SuSE-release ]]; then
    # RHEL-variant logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd /usr/lib/systemd/system/gde.service
    else
        # Assuming SysVinit
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd /lib/systemd/system/gde.service
        deb-systemd-invoke restart gde.service || echo "WARNING: systemd not running."
    else
        # Assuming SysVinit
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
        invoke-rc.d gde restart
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ "$NAME" = "Amazon Linux" ]]; then
        # Amazon Linux 2+ logic
        install_systemd /usr/lib/systemd/system/gde.service
    elif [[ "$NAME" = "Amazon Linux AMI" ]]; then
        # Amazon Linux logic
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
    fi
fi