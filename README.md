# DigitalOcean Dynamic DNS Tool

A simple tool that retrieves your public IPv4 address and updates the A record of a DigitalOcean domain.

_This script is not a daemon and should be run as a cron job or systemd timer and service._

## Installation

1. Clone the repository.
2. Create the folder `secret` with inside it `config.json`:

    ```javascript
    {
        // subdomain.example.com will be the A record that is created/edited.
        "domain": "example.com",
        "hostname": "subdomain",
        // Path of file where the last public IP address will be saved.
        // Needs to be writeable, will be created if it does not exist.
        "ip_file_path": "/home/USERNAME/last_ip",
        "token_source": {
            // How to Create a Personal Access Token
            // https://www.digitalocean.com/docs/api/create-personal-access-token/
            "access_token": "DIGITAL_OCEAN_PERSONAL_ACCESS_TOKEN"
        }
    }
    ```


3. Install [packr](https://github.com/gobuffalo/packr):
    ```bash
    go get -u github.com/gobuffalo/packr/...
    ```

4. Install dependencies:
    ```bash
    dep ensure
    ```

5. Create the executable and move the created executable to your server:
    ```bash
    go build
    ```

    Or install it on the current machine if that's the machine where it's going to be used.
    ```bash
    go install
    ```

## Configuration

### Systemd

```toml
# update-dns.service
[Unit]
Description=Update the DNS with the new public IP.

[Service]
Type=simple
ExecStart=/path/to/update-dns

[Install]
WantedBy=multi-user.target

```

```toml
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
