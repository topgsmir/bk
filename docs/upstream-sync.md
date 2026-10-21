# Upstream integration, 2026-10-21

Reviewed upstream: `AminMGMT/BackPack` at
`9d58f8449ac04705b2b47f4869dc9f6db0e034af` (v1.9.0).
Original bk source: `topgsmir/bk` at
`4b04ed05b9bf6a1e66ee17b4fa000a76aeb85adc` (v1.11.7).
Both supplied ZIP snapshots were checked file-for-file against those Git trees.
The common source commit was `ff86354bcafbd054e7888e48343145323c9b2207`;
upstream added 58 commits affecting 275 files after that point.

The update brings managed Naive/XHTTP/REALITY carriers; sustained download/upload
measurement; IPv4/IPv6 Connection Test suites; traffic-usage reset; short-HTTPS
relay; bounded cancellation and DNS; TCP pool recovery; QUIC/UDP teardown;
traffic-save integrity; and dependency updates. These are the upstream changes,
not claims that a particular filtered Internet route will carry them.

bk keeps its publisher key, paths and names, original wire identifiers, PCK
loopback routing/receive/reply fixes, Spoof defaults and edit controls. Upstream
layer-4 edit handling and changed-field-only saves are combined with the
existing nested Spoof form. Both servers need updated builds for paired new
helper transports and v2 dual-path tests.

The user selected separate core and addon programs under one shared menu.
Package 2 and 3 were extracted from the bk v1.11.7 sources into their own module.
Their shared test and forwarding support is maintained there, with no import of
the core module. The core is now prepared from a pristine pinned upstream tree.

Validation commands and final outcomes are recorded in `validation.md` after
the integration checks finish. Native external-core and public-route reachability
are distinct from local unit and loopback checks.

<div dir="rtl">

## خلاصهٔ فارسی

۵۸ تغییر upstream از زمان جداشدن فورک بررسی و نسخهٔ ۱٫۹٫۰ آن وارد هسته شده است.
قابلیت‌ها و اصلاحات اختصاصی bk حفظ شده‌اند و پکیج‌های ۲ و ۳ به برنامهٔ مستقل
منتقل شده‌اند. تست محلی به معنی تضمین اتصال روی مسیر واقعی ایران و خارج نیست؛
قابلیت‌های جفت‌شدهٔ جدید به ارتقای هر دو سر تانل نیاز دارند.

</div>

*Last verified against bk v1.12.0.*
