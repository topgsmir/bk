# Integration validation — 2026-10-21

Validated locally with Go 1.26.9 in isolated Linux containers, plus the bundled
Node.js runtime for panel tests.

| Check | Outcome |
| --- | --- |
| Supplied bk and BackPack snapshots against their GitHub Git trees | Exact file-content matches |
| Pristine upstream staged in Git | All 1,115 Git blob hashes match the pinned commit |
| Root launcher tests and vet | Passed, including race detection, routing, cache integrity, bounded executable reads and shared uninstall dry run |
| Independent addon module tests and vet | Passed, including race detection and TCP/UDP forwarding tests |
| Addon staticcheck v0.8.2 | Passed after removing six unused inherited coordinator helpers |
| Core full package run | Every package exercised; final focused reruns described below passed |
| Core WebUI final rerun | Passed with race detection |
| Core menu final rerun and full core vet | Passed |
| Required PCK loopback test | Passed for 127.0.0.1 and 127.0.0.2, with BK_REQUIRE_PCK=1, NET_RAW/NET_ADMIN and race detection |
| Client/server transports, L3, direct, limits and relay handlers | Passed with race detection |
| Panel JavaScript | 41 passed, zero failures or skips |
| Upstream sync transaction tests | 3 passed |
| Actual same-pin sync --check-only | Passed GitHub tree verification and overlay preparation without changing tracked source |
| Single-file bundle | Built and smoke-tested on Linux amd64; core/addon version commands and uninstall dry run passed |
| ARMv7 component/bundle build | Passed; build metadata confirms GOARCH=arm and GOARM=7 |
| Shared removal fixture | Passed in a disposable container; addon services stop before config deletion, unrelated services/source stay intact |
| Workflow YAML and staged whitespace checks | Passed |

The initial full core run found one outdated test fixture that generated a
numeric value for the new string `addr` field. The fixture was corrected; the
final WebUI package passed with `-race`. The final uninstall integration was
then checked through the menu package and full core vet. Other core packages
passed the full run, including the complete local engine end-to-end suite.

The original native Additional tunnels and Package 3 workflows are retained
and now run from the independent addon module. Their full native method matrix,
the complete CI race/coverage/fuzz/audit jobs, all seven release architectures,
and a public Iran/Kharej route were not run locally. Require the CI gates before
merging or publishing. No release tag or signed production release was created.

<div dir="rtl">

## خلاصهٔ فارسی

تست‌های محلی منوی مشترک، پکیج‌های مستقل، پنل، مسیرهای اصلی تانل و PCK پاس شدند.
بستهٔ amd64 اجرا و ساخت ARMv7 هم بررسی شد. مسیر حذف فقط در کانتینر موقت تست
شد. آزمایش همهٔ روش‌های خارجی روی هسته‌های واقعی، همهٔ معماری‌های انتشار و
مسیر عمومی ایران و خارج باید از طریق CI و آزمایش سرورها تکمیل شود.

</div>

*Last verified against bk v1.12.0.*
