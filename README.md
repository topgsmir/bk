# bk — independent tunnel components

bk maintained by **topgsmir** combines one entry menu with two independent
programs. BackPack core and Package 2/3 have separate source, Go modules,
dependencies, tests, and processes.

> Based on BackPack by Amin Mohammadi (AminMGMT)
> https://github.com/AminMGMT/BackPack

## Install

Run this command as **root on both servers** (Iran and kharej):

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/bk/main/install.sh)
```

The installer downloads the latest release for your server's architecture,
checks its published SHA-256 checksum, installs `bk`, and opens the shared menu.
Choose **1) BackPack core**, **2) Package 2**, or **3) Package 3**.

To open the menu again:

```bash
sudo bk
```

[Full installation and offline setup](docs/install.md) ·
[Latest release](https://github.com/topgsmir/bk/releases/latest) ·
[راهنمای فارسی](README_FA.md)

## Project structure

| Directory | Owner and purpose |
| --- | --- |
| `upstream/BackPack/` | Exact, unedited upstream snapshot, pinned by commit |
| `integration/` | bk naming, protocol-compatibility literals, and distribution fixes |
| `addons/` | Your Package 2 and Package 3 code; builds without core source |
| `main.go`, `tools/` | Shared entry menu and component build tools |
| `.build/core/` | Generated core source; never edit or commit it |

Run `bk` to choose BackPack core, Package 2, or Package 3. The BackPack panel
belongs to the core. Optional-package management stays in its independent CLI.
`bk core ...`, `bk package2`, `bk package3`, and `bk manage-addons` also work.
Existing `bk -c ...`, `bk --webui`, and `bk external run -c ...` service
commands retain their meanings and config/service paths.

## Build and check

Use Go 1.26.10 or newer, Git, and Make on Linux. Python 3 runs sync checks;
Node.js runs the core panel tests.

```sh
make build
make check
./bk components
```

The launcher embeds the independently compiled `bk-core` and `bk-addons`.
It verifies their hashes and extracts private executables on demand. The release
archive still contains one `bk` binary, so the existing signed installer,
updater, backup, and rollback format remains usable. This is one distribution
bundle containing separate programs; an upstream sync never rewrites addon code.

## Update the core

Review the upstream commit and use its full SHA:

```sh
python3 scripts/sync-upstream.py --commit REVIEWED_40_CHARACTER_SHA --check-only
python3 scripts/sync-upstream.py --commit REVIEWED_40_CHARACTER_SHA
make check
```

The sync verifies the downloaded files against GitHub's commit tree and checks
the distribution overlay before replacing the pristine snapshot. A conflict
leaves the current source and all addons intact. Review and adjust the overlay
when an upstream change affects a retained bk fix. Independent programs remove
addon merge conflicts; core compatibility fixes still need review.

See [architecture](docs/architecture.md), [release procedure](docs/publishing.md),
[Package 2](docs/additional-tunnels.md), and [Package 3](docs/package3.md).

<div dir="rtl">

## خلاصهٔ فارسی

کد اصلی بک‌پک در `upstream/BackPack` دست‌نخورده نگه داشته می‌شود. پکیج‌های ۲ و ۳
داخل `addons` برنامهٔ مستقل دارند و برای ساخت یا اجرا به کد داخلی بک‌پک وابسته
نیستند. دستور `bk` منوی مشترک است؛ مسیر تنظیمات و سرویس‌های قبلی حفظ شده‌اند.
آپدیت هسته کد پکیج‌های شما را تغییر نمی‌دهد. اصلاحات اختصاصی هسته مثل PCK و
هویت انتشار در لایهٔ `integration` مشخص هستند و هنگام تغییر upstream بررسی
می‌شوند. هر دو سر تانل باید پیش از استفاده از قابلیت‌های جدید ارتقا پیدا کنند.

</div>
