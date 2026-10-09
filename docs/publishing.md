# Publishing the topgsmir distribution

This checkout is a **local preparation only**. Nothing has been pushed,
detached, made private or released. The repository and product name remain
BackPack at the owner's request. Upstream's `NOTICE` and `TRADEMARK.md` require
a different name for a modified distribution unless its owner permits this
use. Obtain that permission before publishing under the current name.

## Identity and attribution

- Repository: https://github.com/topgsmir/BackPack
- Maintainer of this distribution: topgsmir
- Source: https://github.com/AminMGMT/BackPack
- Original author: Amin Mohammadi (AminMGMT)
- License and original notices are retained, including the required attribution.
- `.github/description.txt` is the prepared GitHub About description.

## Make the repository independent

The existing GitHub repository is a public fork. A fork cannot independently
change visibility. After reviewing GitHub's warnings, use **Settings → General
→ Danger Zone → Leave fork network**. This keeps Git commit history and the
repository name but permanently detaches the fork and can discard GitHub
metadata such as issues, pull requests, stars and watchers. Do not delete and
recreate the repository as a shortcut.

After detaching, **Change repository visibility** in that same settings area
controls Public/Private. Leave it public if you want anonymous installation
links. Private downloads require authentication; see [install.md](install.md).
Private storage does not replace the upstream license or notices.

References: [Detach a fork](https://docs.github.com/en/pull-requests/how-tos/work-with-forks/detaching-a-fork)
and [visibility](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility).

## Signing key

This distribution pins its own Ed25519 public key in `internal/app/app.go`.
The matching private key is kept locally in the ignored
`.publisher/release-signing-key.dpapi`, encrypted for the Windows account that
prepared it. It is never committed or included in a source archive. Retain a
secure backup or put it in an appropriate secret store before replacing this
machine/account.

Once publication is authorized, sign GitHub CLI in as an account that can
administer `topgsmir/BackPack` and run:

```powershell
gh auth login
./scripts/upload-signing-key.ps1
```

This uploads only the `RELEASE_SIGNING_KEY` Actions secret. It does not push
source or publish releases. Losing the private key requires planned key
rotation; do not replace the public key casually on installed servers.

## Build and release

Use the Go version from `go.mod`. Verify before publishing:

```bash
go test ./... -timeout 20m
# Destructive removal tests belong in a disposable container only:
docker run --rm -v "$PWD:/work" -w /work golang:1.26.6-bookworm bash tests/uninstall.sh
RELEASE_TAG=v1.8.5.1 make release
```

The release workflow requires the signing secret and the signer rejects a key
that does not match the public key in the source. After checking local changes,
push source and the `v1.8.5.1` tag. The Release workflow builds seven Linux
architectures and publishes checksums, their signature, an SBOM, `install.sh`
and `uninstall.sh` under **your** repository. Its manual run requires an
existing version tag and checks out that tag.

When authorized, update the About description with:

```bash
gh repo edit topgsmir/BackPack --description "$(cat .github/description.txt)"
```

Only after a successful release do these public links install/remove this
prepared version:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/BackPack/main/install.sh)
bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/BackPack/main/uninstall.sh)
```

The removal script also works offline and needs no installation or download.
It confirms data removal and supports `--dry-run`. It leaves custom source
checkouts outside `/root/BackPack` alone.


---

<div dir="rtl">

## خلاصهٔ فارسی

این نسخه فعلاً محلی است و در GitHub منتشر نشده. نام مخزن BackPack مانده و
نگه‌دارندهٔ این توزیع topgsmir است؛ سازندهٔ اصلی AminMGMT و مجوز AGPL-3.0 حفظ
شده‌اند. برای انتشار نسخهٔ تغییریافته با همین نام، ابتدا اجازهٔ صاحب نام را بگیر.

پس از اجازه، مخزن را از Settings و Leave fork network مستقل کن؛ این کار دائمی
است و ممکن است اطلاعاتی مثل Issues و ستاره‌ها را حذف کند. سپس امکان انتخاب
Public یا Private داری. لینک عمومی نصب فقط در حالت Public بدون ورود کار می‌کند؛
برای Private از GitHub CLI با حساب مجاز فایل‌ها را بگیر و آفلاین نصب کن.

کلید خصوصی انتشار محلی و رمزگذاری‌شده است و وارد Git یا فایل زیپ نمی‌شود.
اسکریپت upload-signing-key.ps1 بعد از ورود به GitHub آن را در secret مخزن ثبت
می‌کند. انتشار واقعی نیازمند push کد، تگ نسخه و موفق‌شدن تست و ساخت ریلیز است.
اسکریپت uninstall.sh با تأیید تو تنظیمات و بکاپ‌ها را پاک می‌کند؛ --dry-run فقط
پیش‌نمایش است.

</div>

*Last verified against Backpack v1.8.5.1.*
