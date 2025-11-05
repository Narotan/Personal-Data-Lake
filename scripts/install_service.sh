#!/bin/bash

set -e

USER=$(whoami)

echo "Installing aw-client systemd service for user $USER..."

sudo cp scripts/systemd/aw-client.service /etc/systemd/system/aw-client@.service
sudo cp scripts/systemd/aw-client.timer /etc/systemd/system/aw-client@.timer

sudo systemctl daemon-reload
sudo systemctl enable aw-client@$USER.timer
sudo systemctl start aw-client@$USER.timer

echo "Service installed and started!"
echo "Check status: sudo systemctl status aw-client@$USER.timer"
echo "View logs: journalctl -u aw-client@$USER.service -f"

