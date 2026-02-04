# Email Pipeline

This document contains an overview of how email moves through the system: **SMTP ingress → Postfix → HTTPS handoff to API → SMTP egress**.

## Components

- **Mailserver**: docker-mailserver (Postfix) accepts inbound SMTP and hands messages to the API.
- **Email API**: HTTPS service that receives raw messages and triggers outbound delivery.
- **DB + Redis**: dependencies required by the API.
- **SMTP relay**: where the API delivers outbound email (can be an external provider, or the mailserver itself).

## Ports

Mailserver (Postfix):
- `25/tcp` SMTP
- `465/tcp` ESMTP (implicit TLS)
- `587/tcp` ESMTP (STARTTLS)

API:
- The mailserver calls the API at **`https://email-api/v1/email`**.

## Inbound flow (SMTP → Postfix → API)

1. **Inbound SMTP** arrives at Postfix on ports `25/465/587`.
2. Postfix routes recipient domains to a local transport named `curl_email` (via postfix maps in `postfix-virtual.cf`).
3. `curl_email` is implemented as a pipe to `/usr/local/bin/curl-email.sh` (created via `user-patches.sh`).
4. The pipe script POSTs the full RFC822 message bytes to the API:
   - URL: `https://email-api/v1/email`
   - Header: `Authorization: Bearer <PSK>`
   - Body: raw message (`curl --data-binary @-`)

## What the API does with inbound mail (high level)

- The API parses the message and identifies the inbound **alias address** (the original recipient).
- It looks up which **real mailbox recipient(s)** are configured for that alias.
- It then **forwards** the message to those real recipient(s) via the configured SMTP relay (see “Outbound flow”).
- If the message is not deliverable (temporary errors), the API returns a non-2xx status so Postfix will defer/retry.

## Handoff auth (PSK)

- The API endpoint `POST /v1/email` requires a shared secret via `Authorization: Bearer <token>`.
- **The `PSK` value must match on both sides**:
  - mailserver: `PSK` in [mailserver/.env.sample](mailserver/.env.sample)
  - api: `PSK` in [api/.env.sample](api/.env.sample)

## Retry semantics

- If the API returns a non-2xx response, the mailserver pipe uses `curl --fail` and exits with **code 75**.
- Exit code 75 signals a **temporary failure** to Postfix, so Postfix **defers and retries** later.

What this means operationally:
- **API down / network issues / wrong PSK / API cannot submit outbound email** ⇒ Postfix queue growth and repeated deferrals until fixed.
- Default Postfix retry intervals apply.

## Outbound flow (API → SMTP relay)

- The API connects as an SMTP client to the configured relay (`SMTP_CLIENT_HOST` / `SMTP_CLIENT_PORT`).
  - `SMTP_CLIENT_PORT=587` ⇒ STARTTLS

## Minimal configuration checklist

Mailserver:
- Route the domain(s) to `curl_email` in `postfix-virtual.cf` and run `postmap /etc/postfix/virtual`.
- Ensure `curl_email` alias exists (pipe to `/usr/local/bin/curl-email.sh`) and run `newaliases`.
- Set `PSK`, `DOMAIN`, `SMTP_USER`, `SMTP_PASS` (used by `user-patches.sh`).

API:
- Set `PSK` (must match mailserver).
- Set `SMTP_CLIENT_*` for outbound delivery.
- Ensure DB/Redis are reachable (as defined in [api/compose.yml](api/compose.yml)).
