#!/usr/bin/env bash
set -euo pipefail
# Destructive fixture tests run ONLY in a disposable Docker container.
[[ -f /.dockerenv && $EUID == 0 ]] || { echo 'Run in a disposable Docker container as root.' >&2; exit 1; }
script="$(cd "$(dirname "$0")/.." && pwd)/uninstall.sh"
export TEST_STATE="$(mktemp -d)"
mkdir -p "$TEST_STATE/bin" /run/systemd/system /etc/systemd/system
export PATH="$TEST_STATE/bin:$PATH"
cat >"$TEST_STATE/bin/crontab" <<'MOCK'
#!/usr/bin/env bash
if [[ "${FAIL_CRON:-0}" == 1 ]]; then echo 'permission denied' >&2; exit 1; fi
if [[ "$1" == -l ]]; then cat "$TEST_STATE/cron"; else cp "$1" "$TEST_STATE/cron"; fi
MOCK
cat >"$TEST_STATE/bin/systemctl" <<'MOCK'
#!/usr/bin/env bash
printf '%s\n' "$*" >> "$TEST_STATE/services"
if [[ "$1" == list-unit-files ]]; then
  for p in /etc/systemd/system/bk-*.service /etc/systemd/system/bk-*.timer; do [[ ! -f "$p" ]] || echo "${p##*/} enabled"; done
fi
if [[ "${FAIL_STOP:-0}" == 1 && "$1" == stop ]]; then exit 1; fi
exit 0
MOCK
chmod +x "$TEST_STATE/bin/"*
fixture() {
  mkdir -p /etc/bk /root/bk/backups /tmp/unrelated-source
  touch /etc/bk/test.toml /root/bk/backups/test /usr/local/bin/bk /tmp/unrelated-source/keep
  echo /tmp/unrelated-source > /etc/bk/install_path
  touch /etc/systemd/system/bk-test.service /etc/systemd/system/bk-test-restart.timer /etc/systemd/system/unrelated.service /etc/systemd/system/unrelated.timer
  mkdir -p /etc/sysctl.d /etc/security/limits.d
  touch /etc/sysctl.d/zz-bk.conf /etc/sysctl.d/99-bk.conf /etc/security/limits.d/99-bk.conf /etc/sysctl.d/unrelated.conf
  printf '%s\n' '0 * * * * /usr/local/bin/bk --restart-all # bk-auto-refresh' '0 * * * * /usr/local/bin/bk --telegram-report # bk-telegram' '0 0 * * * /usr/local/bin/backup # unrelated' > "$TEST_STATE/cron"
}
fixture
bash "$script" --dry-run > "$TEST_STATE/dry-run"
[[ -f /etc/bk/test.toml && -f /root/bk/backups/test && -f /usr/local/bin/bk ]]
[[ $(wc -l < "$TEST_STATE/cron") == 3 ]]
bash "$script" > "$TEST_STATE/cancel" 2>&1
[[ -f /etc/bk/test.toml ]]
if bash "$script" --invalid >/dev/null 2>&1; then echo 'Invalid option accepted' >&2; exit 1; fi
if FAIL_CRON=1 bash "$script" --yes >/dev/null 2>&1; then echo 'Cron failure ignored' >&2; exit 1; fi
[[ -f /etc/bk/test.toml ]]
if FAIL_STOP=1 bash "$script" --yes >/dev/null 2>&1; then echo 'Stop failure ignored' >&2; exit 1; fi
[[ -f /etc/bk/test.toml && -f /usr/local/bin/bk ]]
fixture
bash "$script" --yes
[[ ! -e /etc/bk && ! -e /root/bk && ! -e /usr/local/bin/bk ]]
[[ ! -e /etc/systemd/system/bk-test.service && -f /etc/systemd/system/unrelated.service ]]
[[ -f /tmp/unrelated-source/keep && -f /etc/systemd/system/unrelated.timer && -f /etc/sysctl.d/unrelated.conf ]]
[[ ! -e /etc/systemd/system/bk-test-restart.timer && ! -e /etc/sysctl.d/zz-bk.conf && ! -e /etc/sysctl.d/99-bk.conf && ! -e /etc/security/limits.d/99-bk.conf ]]
[[ $(cat "$TEST_STATE/cron") == '0 0 * * * /usr/local/bin/backup # unrelated' ]]
grep -q 'stop bk-test-restart.timer' "$TEST_STATE/services"
grep -q 'stop bk-test.service' "$TEST_STATE/services"
bash "$script" --yes >/dev/null
fixture
rm -rf -- /root/bk
ln -s /tmp/unrelated-source /root/bk
bash "$script" --yes >/dev/null
[[ ! -L /root/bk && -f /tmp/unrelated-source/keep ]]
echo 'PASS: dry run, cancellation, validation, cron/stop failures, scoped removal, repeat removal and symlink preservation.'
