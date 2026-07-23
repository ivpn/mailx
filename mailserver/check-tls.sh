#!/usr/bin/env bash
# check-tls.sh — Check TLS certificates for all domains listed in DOMAINS= in .env
# Usage: ./check-tls.sh [--prefix <host-prefix>] [path/to/.env]
#
# Checks:
#   - STARTTLS connectivity on port 25
#   - Certificate expiry (warns at <14 days)
#   - SAN match: cert must be valid for the MX hostname

PREFIX=""
ENV_FILE=".env"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --prefix)
      PREFIX="$2"
      shift 2
      ;;
    *)
      ENV_FILE="$1"
      shift
      ;;
  esac
done
SMTP_PORT=25
TIMEOUT=10
WARN_DAYS=14

# ── Colours ────────────────────────────────────────────────────────────────────
RED='\033[0;31m'; YELLOW='\033[1;33m'; GREEN='\033[0;32m'
CYAN='\033[0;36m'; BOLD='\033[1m'; RESET='\033[0m'

ok()   { echo -e "  ${GREEN}✓  $*${RESET}"; }
warn() { echo -e "  ${YELLOW}⚠  $*${RESET}"; }
fail() { echo -e "  ${RED}✗  $*${RESET}"; }

# ── Preflight checks ───────────────────────────────────────────────────────────
for cmd in openssl dig awk sed grep date timeout; do
  if ! command -v "$cmd" &>/dev/null; then
    echo -e "${RED}ERROR: required command '$cmd' not found.${RESET}" >&2
    exit 2
  fi
done

if [[ ! -f "$ENV_FILE" ]]; then
  echo -e "${RED}ERROR: env file not found: $ENV_FILE${RESET}" >&2
  exit 2
fi

# ── Extract domains from .env DOMAINS= ────────────────────────────────────────
domains_csv=$(grep -E '^\s*DOMAINS=' "$ENV_FILE" | tail -1 | cut -d= -f2-)
domains=$(echo "$domains_csv" | tr ',' '\n' | sed 's/[[:space:]]//g' | grep -v '^$')

if [[ -z "$domains" ]]; then
  echo -e "${RED}No DOMAINS found in $ENV_FILE${RESET}"
  exit 1
fi

# ── Header ─────────────────────────────────────────────────────────────────────
echo -e "\n${BOLD}TLS Certificate Check${RESET}"
echo -e "Config : ${CYAN}$ENV_FILE${RESET}"
[[ -n "$PREFIX" ]] && echo -e "Prefix : ${CYAN}$PREFIX${RESET}"
echo -e "Date   : $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
echo "════════════════════════════════════════════════════════════════"

all_ok=true
declare -A results  # domain → status line for summary

# ── Per-domain check ───────────────────────────────────────────────────────────
for domain in $domains; do
  echo -e "\n${BOLD}▶ $domain${RESET}"
  domain_ok=true

  # 1. MX lookup (or prefix-based override)
  if [[ -n "$PREFIX" ]]; then
    mx="${PREFIX}.${domain}"
    echo "  MX      : $mx (--prefix)"
  else
    mx=$(dig +short MX "$domain" 2>/dev/null | sort -n | head -1 | awk '{print $2}' | sed 's/\.$//')
    if [[ -z "$mx" ]]; then
      warn "No MX record found — falling back to domain itself"
      mx="$domain"
    else
      echo "  MX      : $mx"
    fi
  fi

  # 2. STARTTLS connect + grab cert
  cert_info=$(echo \
    | timeout "$TIMEOUT" openssl s_client \
        -connect "${mx}:${SMTP_PORT}" \
        -starttls smtp \
        -servername "$mx" \
        2>/dev/null)

  if [[ -z "$cert_info" ]]; then
    fail "Could not connect to ${mx}:${SMTP_PORT} or no TLS offered"
    all_ok=false
    domain_ok=false
    results["$domain"]="NO_CONNECT"
    continue
  fi

  # 3. Parse cert fields
  subject=$(  echo "$cert_info" | openssl x509 -noout -subject      2>/dev/null | sed 's/subject=//')
  issuer=$(   echo "$cert_info" | openssl x509 -noout -issuer       2>/dev/null | sed 's/issuer=//')
  not_after=$( echo "$cert_info" | openssl x509 -noout -enddate     2>/dev/null | cut -d= -f2-)
  san_line=$(  echo "$cert_info" | openssl x509 -noout -ext subjectAltName 2>/dev/null \
               | grep -v 'X509v3' | tr -d ' \n')

  echo "  Subject : $subject"
  echo "  Issuer  : $issuer"
  echo "  SANs    : $san_line"
  echo "  Expires : $not_after"

  # 4. Expiry check
  # BSD date fallback for macOS
  expiry_epoch=$(date -d "$not_after" +%s 2>/dev/null \
              || date -j -f "%b %d %T %Y %Z" "$not_after" +%s 2>/dev/null)
  now_epoch=$(date +%s)
  days_left=$(( (expiry_epoch - now_epoch) / 86400 ))

  if (( days_left < 0 )); then
    fail "Certificate EXPIRED ($days_left days ago)"
    all_ok=false; domain_ok=false
  elif (( days_left < WARN_DAYS )); then
    warn "Expiring in $days_left days — renew now"
    all_ok=false; domain_ok=false
  else
    ok "Expires in $days_left days"
  fi

  # 5. SAN hostname match
  # The cert must cover the MX hostname we connected to.
  # SANs look like: DNS:mail1.mailx.net,DNS:mail1.ambox.net
  # We check for an exact DNS name match or a wildcard match (*.example.com).
  mx_matched=false
  IFS=',' read -ra san_entries <<< "$san_line"
  for entry in "${san_entries[@]}"; do
    name="${entry#DNS:}"
    if [[ "$name" == "$mx" ]]; then
      mx_matched=true
      break
    fi
    # Wildcard match: *.example.com covers mail1.example.com
    if [[ "$name" == \*.* ]]; then
      wildcard_base="${name#\*.}"           # example.com
      mx_suffix="${mx#*.}"                  # example.com (strip first label)
      if [[ "$mx_suffix" == "$wildcard_base" ]]; then
        mx_matched=true
        break
      fi
    fi
  done

  if $mx_matched; then
    ok "SAN matches $mx"
  else
    fail "SAN MISMATCH — cert not valid for '$mx'"
    echo -e "  ${RED}       Cert covers: $san_line${RESET}"
    all_ok=false; domain_ok=false
  fi

  $domain_ok && results["$domain"]="OK" || results["$domain"]="FAIL"
done

# ── Summary ────────────────────────────────────────────────────────────────────
echo -e "\n════════════════════════════════════════════════════════════════"
echo -e "${BOLD}Summary${RESET}"
for domain in $(echo "${!results[@]}" | tr ' ' '\n' | sort); do
  status="${results[$domain]}"
  if [[ "$status" == "OK" ]]; then
    echo -e "  ${GREEN}✓${RESET}  $domain"
  else
    echo -e "  ${RED}✗${RESET}  $domain  ${RED}($status)${RESET}"
  fi
done

echo ""
if $all_ok; then
  echo -e "${GREEN}${BOLD}All certificates OK.${RESET}\n"
  exit 0
else
  echo -e "${RED}${BOLD}One or more certificates need attention.${RESET}\n"
  exit 1
fi