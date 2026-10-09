# Server layout (file locations)

Everything lives in a tidy, predictable layout. You can also see this any time
from **Manage → File Locations** in the CLI.

| Path | What |
|------|------|
| `/root/bk` | The release bundle and downloaded archives. |
| `/root/bk/backups` | [Backup](backup-restore.md) `.tar.gz` files. |
| `/etc/bk` | Tunnel configs (one `.toml` per tunnel) and runtime state. |
| `/usr/local/bin/bk` | The binary itself. |
| `bk-<name>.service` | A systemd unit per tunnel. |
| `bk-monitor.service` | The [monitor service](monitor-service.md). |

The install directory is recorded in `/etc/bk/install_path`, which is what
the uninstaller reads to know what to remove.

---

<div dir="rtl">

## خلاصهٔ فارسی

همه‌چیز در یک ساختار مرتب و قابل‌پیش‌بینی است. همین را هر وقت خواستی از
`Manage → File Locations` در CLI هم می‌بینی.

`‎/root/bk` بستهٔ ریلیز و آرشیوهای دانلودشده ·
`‎/root/bk/backups` فایل‌های [پشتیبان](backup-restore.md) ·
`‎/etc/bk` کانفیگ تونل‌ها (برای هر تونل یک فایل `.toml`) و وضعیت اجرا ·
`‎/usr/local/bin/bk` خودِ باینری ·
`bk-<name>.service` یک یونیت systemd به‌ازای هر تونل ·
`bk-monitor.service` [سرویس مانیتور](monitor-service.md).

مسیر نصب در `/etc/bk/install_path` ثبت می‌شود و حذف‌کنندهٔ داخلی از روی
همان می‌فهمد چه چیزی را پاک کند.

</div>

---
[← Back to the docs index](README.md)

---

*Last verified against Backpack v1.8.5.*
