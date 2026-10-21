# Publishing the separated bk distribution

The distribution is maintained at `topgsmir/bk`. Its source keeps the AGPL-3.0
licence, original attribution, and the independently pinned upstream snapshot.

Run `make check` and require the BackPack core, launcher, Additional tunnels and
Tunnel Package 3 workflows to pass. The native addon workflows remain separate
from the core checks. A local loopback pass does not establish public-route
connectivity or native support for every optional method.

`make build` creates a single runnable launcher containing two separately built
component executables. `make release` builds seven Linux variants, creates the
existing `bk_linux_<arch>.tar.gz` archives, computes checksums, and signs them
through the core's existing distribution signer. `release/SBOM.txt` includes
both embedded programs' dependency metadata.

The public signing key is retained in the core overlay and materialized at
`.build/core/internal/app/app.go`. The private key is never source data. Keep
the existing encrypted `.publisher` backup and the repository's
`RELEASE_SIGNING_KEY` secret. `scripts/upload-signing-key.ps1` remains the
explicit account-owned provisioning tool; building does not upload keys.

The root `VERSION`, `upstream.lock.json` distribution version, core app version,
and release tag must agree. The lock also records the independently inherited
core version and addon version. A pushed tag publishes that source; a manual
workflow can use an explicit tag or the declared version, with the workflow's
existing tag/source consistency checks.

Do not tag or publish solely because an upstream release exists. Review its
source changes, sync the pristine snapshot, run the compatibility checks, and
require the appropriate native workflows. Updating core source does not update
or regenerate the addon module.

<div dir="rtl">

## خلاصهٔ فارسی

قبل از انتشار، تست‌های هسته، منوی مشترک و پکیج‌ها باید پاس شوند. خروجی یک
فایل نصب امضاشده است که برنامه‌های مستقل داخل آن قرار دارند. کلید اختصاصی
انتشار، مسیر نصب و قالب آپدیت قبلی حفظ شده‌اند. دریافت آپدیت upstream به معنی
انتشار خودکار نیست؛ نسخه، سازگاری و تست‌های مسیر واقعی باید بررسی شوند.

</div>

*Last verified against bk v1.12.0.*
