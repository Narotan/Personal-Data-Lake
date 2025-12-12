#!/bin/bash

set -e

UNAME=$(uname -s)
USER=$(whoami)
PROJECT_DIR=$(pwd)

if [ "$UNAME" = "Linux" ]; then
    echo "Installing aw-client systemd service for Linux..."
    
    # Create systemd service file
    sudo tee /etc/systemd/system/aw-client@.service > /dev/null <<EOF
[Unit]
Description=ActivityWatch Data Collector
After=network.target

[Service]
Type=oneshot
User=%i
WorkingDirectory=${PROJECT_DIR}
ExecStart=${PROJECT_DIR}/bin/aw-client -minutes 5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    # Create systemd timer file
    sudo tee /etc/systemd/system/aw-client@.timer > /dev/null <<EOF
[Unit]
Description=Run ActivityWatch Data Collector every 5 minutes
Requires=aw-client@.service

[Timer]
OnBootSec=2min
OnUnitActiveSec=5min
Unit=aw-client@%i.service

[Install]
WantedBy=timers.target
EOF

    sudo systemctl daemon-reload
    sudo systemctl enable aw-client@$USER.timer
    sudo systemctl start aw-client@$USER.timer
    
    echo "Service installed and started!"
    echo "Check status: sudo systemctl status aw-client@$USER.timer"
    echo "View logs: journalctl -u aw-client@$USER.service -f"

elif [ "$UNAME" = "Darwin" ]; then
    echo "Installing aw-client launchd service for macOS..."
    
    PLIST_FILE="$HOME/Library/LaunchAgents/com.personal-datalake.aw-client.plist"
    
    # Create launchd plist file
    cat > "$PLIST_FILE" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.personal-datalake.aw-client</string>
    <key>ProgramArguments</key>
    <array>
        <string>${PROJECT_DIR}/bin/aw-client</string>
        <string>-minutes</string>
        <string>5</string>
    </array>
    <key>WorkingDirectory</key>
    <string>${PROJECT_DIR}</string>
    <key>StartInterval</key>
    <integer>300</integer>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/aw-client.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/aw-client.error.log</string>
</dict>
</plist>
EOF

    launchctl load "$PLIST_FILE"
    
    echo "Service installed and started!"
    echo "Check status: launchctl list | grep aw-client"
    echo "View logs: tail -f /tmp/aw-client.log"
    echo "Stop service: launchctl unload $PLIST_FILE"

else
    echo "Unsupported OS: $UNAME"
    exit 1
fi
