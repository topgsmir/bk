#!/usr/bin/env bash
# Integration check of the actual release installation and bk entry points.
set -euo pipefail
[[ -f /.dockerenv && $EUID -eq 0 ]] || { echo 'Run only as root in a disposable container.' >&2; exit 1; }
mkdir -p /opt/bk-test /opt/bk-mocks /etc/backpack
printf 'preserve original installation\n' > /etc/backpack/keep
cp install.sh uninstall.sh release/bk_linux_amd64.tar.gz release/SHA256SUMS /opt/bk-test/
cat > /opt/bk-mocks/systemctl <<'MOCK'
#!/bin/sh
case "$1" in is-active) echo active;; is-enabled) echo enabled;; esac
exit 0
MOCK
cat > /opt/bk-mocks/crontab <<'MOCK'
#!/bin/sh
[ "${1:-}" = '-l' ] && { echo "no crontab for root" >&2; exit 1; }
exit 0
MOCK
chmod +x /opt/bk-mocks/*
export PATH="/opt/bk-mocks:$PATH"
bash /opt/bk-test/install.sh </dev/null
[[ $(command -v bk) == /usr/local/bin/bk ]]
bk -v | grep -Fx v1.8.6.1
bk -v | grep -Fx https://github.com/topgsmir/bk
bk help > /opt/bk-test/help-output
grep -q bk /opt/bk-test/help-output
[[ $(cat /etc/bk/install_path) == /root/bk ]]
# The real no-argument entry point writes the panel and monitor service units.
printf '10\n' | bk > /opt/bk-test/menu-output
for unit in bk-monitor bk-webui; do
  grep -q 'ExecStart=/usr/local/bin/bk' "/etc/systemd/system/$unit.service"
done
grep -q 'Select an option' /opt/bk-test/menu-output
grep -q 'https://github.com/topgsmir/bk' /opt/bk-test/menu-output
bash /opt/bk-test/uninstall.sh --yes
[[ ! -e /usr/local/bin/bk && ! -e /etc/bk && ! -e /root/bk ]]
[[ -f /etc/backpack/keep ]]
echo 'PASS: release installation, bk command/menu, bk units, uninstall and original installation preservation.'
