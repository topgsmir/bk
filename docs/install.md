# Installing bk

## The normal way

One command as root on the VPS:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/bk/main/install.sh)
```

It downloads the prebuilt release archive for your architecture (amd64/arm64)
into `/root/bk`, **verifies it against the checksum published with the
release**, installs the binary, and opens the menu when it finishes.

Reopen the menu any time with:

```bash
sudo bk
```

Everything lands in a tidy layout — the release bundle in `/root/bk`,
backups in `/root/bk/backups`, tunnel configs in `/etc/bk`. See
[server layout](server-layout.md).

> **Building from source** still works as a fallback: clone the repo and run
> `sudo bash install.sh` inside it. If the release download fails it builds with
> Go, fetching modules **directly first** and via Iran-friendly mirrors
> (RunFlare, goproxy.cn) only when direct access fails.

---

## Offline install (the server cannot reach GitHub)

Download the release on any machine **with** internet, copy it to the server, and
install it there. Nothing is fetched from the VPS.

![Offline install](../img/offline-install.gif)

From the [releases page](https://github.com/topgsmir/bk/releases/latest),
download the archive for the server's architecture — run `uname -m` on it:
`x86_64` → `bk_linux_amd64.tar.gz`, `aarch64` → `bk_linux_arm64.tar.gz`.

### With the installer (recommended)

It also records the layout for the uninstaller. Download `install.sh` and
`SHA256SUMS` alongside the archive, put all three in the **same folder** on the
VPS, and run it. It finds the local archive, verifies it against `SHA256SUMS`,
and never touches the network:

```bash
scp install.sh SHA256SUMS bk_linux_amd64.tar.gz root@SERVER_IP:/root/
ssh root@SERVER_IP "cd /root && sudo bash install.sh"
```

### By hand

Upload the archive to the server, then as root:

```bash
sha256sum bk_linux_amd64.tar.gz        # compare against SHA256SUMS
tar xzf bk_linux_amd64.tar.gz
mkdir -p /etc/bk /root/bk/backups
install -m 0755 bk /usr/local/bin/bk
echo /root/bk > /etc/bk/install_path
sudo bk
```

The `install_path` line is what the built-in uninstaller reads to know what to
remove; skip it and everything still runs, but uninstalling has to be done by
hand. `install -m 0755` already sets the executable bit, so no `chmod` is needed.

### Updating offline

The same way: repeat the steps with the newer archive. `install` replaces the
binary in place, and your tunnels in `/etc/bk` are untouched. Restart them
afterwards with `sudo bk` → **Manage → Restart ALL**. Optional additional
tunnels have their own restart actions under **11 → Manage**; see
[additional tunnels](additional-tunnels.md).

---

## Updating online

**Main menu → 8) Update.** It downloads the release, verifies it against the
published SHA-256, installs it, and **rolls back automatically** if a tunnel does
not come back up. Anything that cannot be verified is refused rather than
installed. [More](updates.md).

## Uninstalling

Run the removal script from your own repository as root:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/bk/main/uninstall.sh)
```

It asks you to type `DELETE`, stops bk timers and services, removes its cron jobs and tuning files,
tunnels, settings, binary and `/root/bk` **including backups**. Unrelated
cron jobs and custom source checkouts are preserved. Kernel tuning already applied
to the running kernel lasts until reboot. `--dry-run` previews it;
`--yes` skips the prompt. The installed menu's **9) Uninstall** also remains available.

## Private repository access

Anonymous raw and release links work only while the repository is public.
For a private repository, authenticate GitHub CLI on a machine with access,
then download the installer and assets from **topgsmir/bk**:

```bash
gh auth login
gh repo clone topgsmir/bk
cd bk
gh release download --repo topgsmir/bk --pattern 'bk_linux_amd64.tar.gz' --pattern SHA256SUMS --dir .
sudo bash install.sh
# Local removal needs no GitHub connection or token:
sudo bash uninstall.sh
```

Choose the asset matching the server architecture. For servers without GitHub
access, copy these files from the authenticated machine and use the offline
steps above. Private releases are updated through this same authenticated,
offline path; the built-in anonymous online updater needs a public repository.
Never paste a token into a public installation command or commit it to source.

---

<div dir="rtl">

## خلاصهٔ فارسی

این توزیع از مخزن topgsmir/bk نصب می‌شود. برای حذف با لینک خودمان،
اسکریپت uninstall.sh را اجرا کن؛ باید DELETE را تأیید کنی و بکاپ‌ها نیز حذف
می‌شوند. --dry-run فقط پیش‌نمایش است. اگر مخزن Private باشد، لینک ناشناس کار
نمی‌کند؛ فایل‌ها را با gh و حساب مجاز بگیر و آفلاین نصب یا حذف کن.

**نصب عادی:** یک دستور با کاربر root روی سرور — آرشیو ریلیز مخصوص معماری سرور را
دانلود می‌کند، با چک‌سام منتشرشده **تأیید** می‌کند، نصب می‌کند و خودش منو را باز
می‌کند. بعداً با `sudo bk` منو را باز کن.

**نصب آفلاین (سروری که به گیت‌هاب دسترسی ندارد):** فایل ریلیز را روی یک ماشین با
اینترنت دانلود کن و به سرور کپی کن. با `uname -m` معماری را ببین: `x86_64` یعنی
amd64 و `aarch64` یعنی arm64. بهترین راه این است که `install.sh` و `SHA256SUMS`
را هم کنار آرشیو در **یک پوشه** بگذاری و اسکریپت را اجرا کنی — خودش فایل محلی را
پیدا و تأیید می‌کند و اصلاً به شبکه دست نمی‌زند. روش دستی هم در بالا آمده؛ فقط
یادت باشد خط `install_path` را بنویسی، چون حذف‌کنندهٔ داخلی از روی آن می‌فهمد چه
چیزی را پاک کند.

**آپدیت:** از منوی اصلی گزینهٔ ۸ — با تأیید SHA-256 و **بازگشت خودکار** اگر تونل
بالا نیامد. آپدیت آفلاین هم همان مراحل نصب با آرشیو جدید است و کانفیگ‌های
`/etc/bk` دست‌نخورده می‌مانند.

</div>

---
[← Back to the docs index](README.md)

---

*Last verified against bk v1.10.0.*
