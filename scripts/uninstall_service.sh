#!/bin/bash

set -e

USER=$(whoami)

echo "Uninstalling aw-client systemd service for user $USER..."

sudo systemctl stop aw-client@$USER.timer || true
sudo systemctl disable aw-client@$USER.timer || true
sudo rm -f /etc/systemd/system/aw-client@.service
sudo rm -f /etc/systemd/system/aw-client@.timer
sudo systemctl daemon-reload

echo "Service uninstalled!"

