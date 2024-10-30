# traefik-fail2ban-connector
A Traefik plugin that uses an existing fail2ban setup to block IPs.

## Setup
This plugin needs setup in both Traefik and Fail2Ban.

### Traefik

#### Static

```yaml
experimental:
  plugins:
    fail2ban:
      moduleName: github.com/ClimberJ/traefik-fail2ban-connector
      version: v1.0.0
```

#### Dynamic

```yaml
http:
  middlewares:
    traefik-fail2ban-connector:
      plugin:
        fail2ban:
          {}
```

### Fail2Ban

This assumes you already have a working Fail2Ban configuration.

#### Action
Add this into an action file for Fail2Ban and add the action to the jail you want to use for blocking requests in Traefik.

```
[Definition]

actionban = echo <ip> >> /config/fail2ban/bans.txt


actionunban = grep -v ^<ip>$ /config/fail2ban/bans.txt > /config/fail2ban/temp; cat /config/fail2ban/temp > /config/fail2ban/bans.txt; rm temp
```

### General

The `bans.txt` file from the action will need to be mounted into Traefik at `/etc/fail2ban/bans.txt`.

## Limitations

- This plugin can currently only be used with one jail.
- If one blocked IP (e.g. 192.168.178.100) contains another blocked IP (e.g. 192.168.178.10) unbanning the inner IP will also remove the outer IP from the file.

Feel free to submit PRs for any of these.