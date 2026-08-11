# curl-email.sh debug logs

`curl-email.sh` (the Postfix pipe script that forwards mail to the `email-api`)
can optionally log its steps for troubleshooting. **Debug logging is disabled
by default** — the two `echo ... >> /var/log/mail/curl-email.log` lines in the
script are commented out, so nothing is written unless you enable them.

## Enable on a live server (no restart, temporary)

Edit the script directly inside the running `mailserver` container. This takes
effect immediately since `curl-email.sh` is re-read from disk for every
message — no Postfix reload or container restart needed.

```sh
docker compose exec mailserver sed -i \
  -e 's/^# \(echo "\$(date -Iseconds)\)/\1/' \
  /usr/local/bin/curl-email.sh
```

This survives only as long as the current container instance. A recreate/
restart of `mailserver` regenerates `/usr/local/bin/curl-email.sh` from
[docker-data/dms/config/user-patches.sh](../docker-data/dms/config/user-patches.sh),
reverting the file back to disabled.

## Enable persistently (survives mailserver restarts)

Edit [docker-data/dms/config/user-patches.sh](../docker-data/dms/config/user-patches.sh)
(the live copy that's mounted into the container and executed on every
startup) and uncomment the same two debug lines inside the `curl-email.sh`
heredoc. Then recreate the container so the patch script re-runs on a clean
container filesystem:

```sh
docker compose up -d --force-recreate mailserver
```

Use `--force-recreate` rather than `docker compose restart` — the patch
script writes the script with `>>` (append), so restarting the existing
container instead of recreating it will append a duplicate copy of the
script.

## Reading the debug logs

The log file is shared between the `mailserver` and `daemon` containers via
the same host-mounted directory, so it can be read from either container or
directly from the host.

```sh
# from inside the mailserver container
docker compose exec mailserver tail -f /var/log/mail/curl-email.log

# from inside the daemon container
docker compose exec daemon tail -f /var/log/mail/curl-email.log

# from the host
tail -f docker-data/dms/mail-logs/curl-email.log
```

## Log file location

| Where | Path |
|---|---|
| Inside `mailserver` container | `/var/log/mail/curl-email.log` |
| Inside `daemon` container | `/var/log/mail/curl-email.log` |
| On the host | `docker-data/dms/mail-logs/curl-email.log` |

## Clearing cycle

The `daemon` container truncates the file every 24h by default (regardless of
whether debug logging is enabled), so it never grows unbounded. The interval
can be overridden via the `CURL_EMAIL_LOG_CLEAR_INTERVAL` env var on the
`daemon` service (accepts a Go duration string, e.g. `1h`). This is separate
from docker-mailserver's own `LOGROTATE_INTERVAL`/`LOGROTATE_COUNT`, which
does not manage this file.
