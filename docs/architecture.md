# Independent core and addon programs

The distribution has three separately compiled Go modules:

```text
bk launcher                 (module at repository root)
  ├── bk-core               (prepared from upstream/BackPack + integration)
  └── bk-addons             (module in addons/)
        ├── Package 2
        └── Package 3
```

The launcher communicates with each component through process arguments,
standard input/output, and exit codes. Neither component imports the other.
Service invocations replace the launcher process, retaining the existing PID and
signal behavior. Components are hashed and extracted as streams so startup does
not allocate a whole embedded executable for every tunnel service.
The addon module has no `replace ../core` directive and owns its own test
coordinator, TCP/UDP forwarder, TUI and service-control helpers.

`upstream.lock.json` pins the exact reviewed BackPack commit. Every pristine
source file is hashed in `integration/upstream-files.json`. Preparation refuses
local edits to that snapshot, preserves the previously retained wire identifiers,
performs the mechanical bk naming transform, and applies `integration/core.patch`.
The patch owns the publisher identity, paths, PCK loopback and Spoof fixes; it
contains no Package 2/3 implementation. The result is disposable `.build/core`.

## Existing installations

| Existing command or path | Component after separation |
| --- | --- |
| `bk -c /etc/bk/name.toml` | Core |
| `bk --webui`, `bk --monitor` | Core |
| `bk external run -c /etc/bk/external/name.json` | Addons |
| `/etc/bk/*.toml`, `bk-<name>.service` | Core state |
| `/etc/bk/external/*.json`, `bk-ext-<name>.service` | Addon state |

`bk` opens the shared menu. The core has its existing panel and management menu.
Optional-package management now lists only Packages 2 and 3. No configuration or
systemd unit migration is required for the commands above.

The existing single-file release format is retained: both component executables
are embedded in the launcher, hashed, and extracted into private per-content
directories under the user's cache. A rollback of the launcher also selects its
original component hashes. Installation and signed online updating continue to
replace one binary. SBOM output records the dependencies of both child programs,
not merely the small launcher's standard library.

Distribution removal is owned by the launcher. The existing core-menu Uninstall
action delegates to the same embedded removal script after its confirmation;
addon services are stopped before their configuration and executables are removed.

## Future changes

Addon work is confined to `addons/`, its documentation and its CI workflows.
Upstream work is confined to the pristine snapshot and its pin. Compatibility
fixes are explicit in `integration/`. Core sync checks that overlay before
replacing the snapshot; a conflict leaves the current source and addons intact.

This removes shared Go-package and menu edits between upstream and addons. It
does not make arbitrary future CLI, protocol, kernel or dependency changes
automatically compatible. Core regression tests and the separate native addon
workflows remain release requirements.

<div dir="rtl">

## خلاصهٔ فارسی

منوی `bk` فقط مسئول اجرای بخش انتخاب‌شده است. بک‌پک و پکیج‌های ۲ و ۳ برنامه،
وابستگی و تست مستقل دارند. کد اصلی در پوشهٔ upstream تغییر داده نمی‌شود؛ اصلاحات
اختصاصی هسته در integration مشخص هستند. مسیر تنظیمات و دستور سرویس‌های موجود
حفظ شده و بازگشت به نسخهٔ قبلی، برنامه‌های همان نسخه را هم برمی‌گرداند. آپدیت
بک‌پک به کد پکیج‌ها دست نمی‌زند، ولی تغییر پروتکل یا رابط‌های آینده همچنان بررسی
سازگاری می‌خواهد.

</div>

*Last verified against bk v1.12.0.*
