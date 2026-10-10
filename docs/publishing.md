# Publishing bk

bk is the distribution maintained by topgsmir, based on BackPack by Amin
Mohammadi (AminMGMT): https://github.com/AminMGMT/BackPack.
The original AGPL-3.0 license, NOTICE and trademark policy remain in the source.
The modified product uses its own name and icon.

## Repository and visibility

https://github.com/topgsmir/bk is a standalone repository. It preserves the
source history and can be Public or Private through Settings → General →
Danger Zone → Change repository visibility. The original BackPack fork is
separate. Anonymous installation and online updates require the bk repository
to be public; authenticated private installation is described in [install.md](install.md).

See [GitHub visibility documentation](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility).

## Signed releases

The pinned Ed25519 public key is in internal/app/app.go. The matching private
key is locally encrypted for the Windows account in the ignored
.publisher/release-signing-key.dpapi file. It is never part of Git or source
archives. Keep a secure backup before replacing the machine/account.
The script scripts/upload-signing-key.ps1 uploads it to the bk repository's
RELEASE_SIGNING_KEY Actions secret after gh auth login. The signer rejects a
key that does not match the public key in the source.

## Build and publish

Use the Go version from go.mod. Before tagging, run:

```bash
go test ./... -timeout 20m
go test ./internal/utils/network -race -timeout 5m
# Run removal and privileged PCK tests in a disposable container only.
docker run --rm --cap-add=NET_RAW --cap-add=NET_ADMIN -e BK_REQUIRE_PCK=1 -v "$PWD:/work" -w /work golang:1.26.9-bookworm bash -c 'bash tests/uninstall.sh && go test ./internal/utils/network -run TestPckLoopback -count=1 -v'
RELEASE_TAG=v1.10.0 make release
```

Require the Additional tunnels workflow to pass its real 60-second tests for
every optional method before tagging. Push reviewed source to main and its
version tag. The Release workflow tests,
builds seven Linux architectures and publishes checksums, a signature, SBOM,
install.sh and uninstall.sh. Manual workflow runs require an existing tag.
The release tag, VERSION and app.Version must agree.

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/bk/main/install.sh)
bk
bash <(curl -fsSL https://raw.githubusercontent.com/topgsmir/bk/main/uninstall.sh)
```

The removal script works offline, asks for confirmation, and supports
--dry-run. It removes only bk's standard paths and units. Custom source
checkouts outside /root/bk are left alone.

---

<div dir="rtl">

## خلاصهٔ فارسی

مخزن bk متعلق به topgsmir و مستقل است؛ منبع اصلی، مجوز و سابقهٔ کد حفظ شده‌اند.
از Settings و Change repository visibility می‌توان حالت عمومی یا خصوصی را
تغییر داد. نصب و به‌روزرسانی با لینک عمومی نیاز به مخزن عمومی دارد؛ برای حالت
خصوصی راهنمای نصب آفلاین را بخوان. دستور bk منوی مدیریت را باز می‌کند.
کلید خصوصی انتشار رمزگذاری شده است و وارد Git نمی‌شود؛ نسخه‌های منتشرشده با
کلید اختصاصی امضا می‌شوند. حذف با تأیید کاربر انجام می‌شود و --dry-run پیش‌نمایش است.

</div>

*Last verified against bk v1.11.0.*
