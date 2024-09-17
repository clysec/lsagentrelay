# Lansweeper Agent Relay/Proxy
An alternative to the cloud service provided by Lansweeper. Listens for requests from the agent, rewrites the configuration on the fly and forwards the request to the actual Lansweeper instance.

To enable sending data to Lansweeper without opening ports for all servers to the Lansweeper instance. This relay can be deployed anywhere, as long as it can reach the Lansweeper instance.

## Usage
### CLI
```
lsagentrelay /path/to/config.yaml
```

### Env
```
# CLI
LSAGENTRELAY_CONFIG=/path/to/config.yaml lsagentrelay

# Docker
-e LSAGENTRELAY_CONFIG=/path/to/config.yaml
```

## Configuration
```yaml
# Lansweeper Server Settings
lansweeper:
  # URL to the Lansweeper Server
  url: lansweephostname.internal.company.net

  # The Agent port for the Lansweeper Server (not the same as the web port), default is 9524
  port: 9524

  # Ignore SSL certificate validation (if using self-signed certificates)
  ignore_ssl: true

# Rewrite settings
# When the agent fetches the configuration from Lansweeper, it needs to rewrite some settings so that
# the agent connects to the proxy instead of the Lansweeper Server.
# This is generally the hostname of the lansweeper server.
# If this is incorrect
rewrite:
  # This is an example, replace with the actual hostname and port(s) you want to rewrite
  # This is case insensitive
  LansweeperHostnameRegex: "(lansweeperhostname|otherhostname)\\.internal\\.company\\.net:(9524|443)"

  # The external name and port of the proxy
  ProxyHostname: "lsagentproxyproxy.company.net:443"


# Server settings
listen:
  # The port to listen on
  port: 8080

  # The host to listen on (default empty, meaning all interfaces)
  host: ""

  # TLS Settings
  tls:
    # Enable TLS. You can also run this behind a proxy with TLS enabled.
    enabled: false

    # Path to the TLS certificate file (PEM format)
    cert: ""

    # Path to the TLS key file (PEM format)
    key: ""

```
