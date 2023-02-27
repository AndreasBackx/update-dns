# DigitalOcean Dynamic DNS Tool

A simple tool that retrieves your public IPv4 address and updates the A record of a DigitalOcean domain.

_This script is not a daemon and should be run as a cron job or systemd timer and service._

## Configuration

`CLOUDFLARE_API_TOKEN` is required to be set with your Cloudflare API token.

### Systemd

```ini
# update-dns.service
[Unit]
Description=Update the DNS with the new public IP.

[Service]
Type=simple
ExecStart=/path/to/update-dns --domain ${DOMAIN}

[Install]
WantedBy=multi-user.target

```

```ini
# update-dns.timer
[Unit]
Description=Run the update-dns service regularly.
Requires=update-dns.service

[Timer]
# Time to wait after booting before we run first time
OnBootSec=1min
# Time between running each consecutive time
OnUnitActiveSec=15min
Unit=update-dns.service

[Install]
WantedBy=timers.target
```
