#!/usr/bin/env bash
# bk distribution maintained by topgsmir.
# bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/bk/main/uninstall.sh)
set -euo pipefail

yes=0
dry_run=0
for arg in "$@"; do
  case "$arg" in
    --yes) yes=1 ;;
    --dry-run) dry_run=1 ;;
    -h|--help)
      echo 'Usage: bash uninstall.sh [--yes] [--dry-run]'
      echo 'Removes bk services, settings, tunnels, binary and backups.'
      exit 0 ;;
    *) echo "Unknown argument: $arg" >&2; exit 2 ;;
  esac
done
if (( ! dry_run && EUID != 0 )); then echo 'Run as root (sudo).' >&2; exit 1; fi
echo 'bk uninstall: /etc/bk, /root/bk (including backups), /usr/local/bin/bk'
if (( ! dry_run && ! yes )); then
  answer=''
  read -r -p 'Type DELETE to remove bk and its data: ' answer </dev/tty || true
  if [[ "$answer" != DELETE ]]; then echo 'Cancelled.'; exit 0; fi
fi
run() {
  if (( dry_run )); then printf '[dry-run] '; printf '%q ' "$@"; echo
  else "$@"; fi
}
# A failed cron read must never be treated as an empty crontab.
if command -v crontab >/dev/null 2>&1; then
  cron_dir="$(mktemp -d)"
  trap 'rm -rf -- "$cron_dir"' EXIT
  if crontab -l >"$cron_dir/current" 2>"$cron_dir/error"; then
    awk '!/#[[:space:]]+bk-(auto-refresh|telegram)([[:space:]]|$)/' \
      "$cron_dir/current" >"$cron_dir/kept"
    if ! cmp -s "$cron_dir/current" "$cron_dir/kept"; then
      if (( dry_run )); then echo '[dry-run] remove bk cron entries'
      else crontab "$cron_dir/kept"; fi
    fi
  elif ! grep -qi 'no crontab' "$cron_dir/error"; then
    echo 'Cannot read the current crontab; no bk files were removed.' >&2
    exit 1
  fi
fi
declare -A units=()
shopt -s nullglob
for path in /etc/systemd/system/bk-*.service /etc/systemd/system/bk-*.timer; do units["${path##*/}"]=1; done
if [[ -d /run/systemd/system ]]; then
  command -v systemctl >/dev/null 2>&1 || { echo 'systemctl is required.' >&2; exit 1; }
  unit_list="$(systemctl list-unit-files --no-legend --no-pager 'bk-*.service' 'bk-*.timer')"
  while read -r unit rest; do
    if [[ ( "$unit" == bk-*.service || "$unit" == bk-*.timer ) && "$unit" != */* ]]; then units["$unit"]=1; fi
  done <<<"$unit_list"
  # Stop timers first so they cannot restart a service during removal.
  for kind in timer service; do
    for unit in "${!units[@]}"; do
      if [[ "$unit" == *."$kind" ]]; then
        run systemctl stop "$unit"
        run systemctl disable "$unit"
      fi
    done
  done
fi
for path in /etc/systemd/system/bk-*.service /etc/systemd/system/bk-*.service.d /etc/systemd/system/bk-*.timer /etc/systemd/system/bk-*.timer.d; do
  run rm -rf -- "$path"
done
if [[ -d /run/systemd/system ]]; then run systemctl daemon-reload; fi
# Do not turn an untrusted install_path into a recursive-delete argument.
if [[ -f /etc/bk/install_path ]]; then
  recorded="$(cat /etc/bk/install_path)"
  if [[ -n "$recorded" && "$recorded" != /root/bk ]]; then
    printf 'Custom source checkout preserved: %s\n' "$recorded"
  fi
fi
# Remove only the three tuning files written by this project. Do not reset
# live kernel settings, which could affect other software on the server.
for path in /etc/sysctl.d/zz-bk.conf /etc/sysctl.d/99-bk.conf /etc/security/limits.d/99-bk.conf; do
  if [[ -e "$path" || -L "$path" ]]; then
    run rm -f -- "$path"
    echo 'bk tuning file removed; live kernel settings persist until reboot.'
  fi
done
run rm -f -- /usr/local/bin/bk
run rm -rf -- /etc/bk /root/bk
if (( dry_run )); then echo 'Dry run complete; no bk data was removed.'
else echo 'bk uninstalled.'; fi
