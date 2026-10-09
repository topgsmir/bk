#!/usr/bin/env bash
# Independent optional AWG builder. Never changes bk's source or Go modules.
set -euo pipefail
umask 077
: "${BK_EXTRA_CORES:?missing core directory}"
[[ $EUID -eq 0 ]] || { echo 'Run as root.' >&2; exit 1; }
work="$(mktemp -d /tmp/bk-awg-build.XXXXXX)"
trap 'rm -rf -- "$work"' EXIT
fetch_checked() {
 local url="$1" dest="$2" hash="$3"
 curl -fSL --retry 2 --retry-delay 2 --connect-timeout 15 --max-time 300 "$url" -o "$dest"
 printf '%s  %s\n' "$hash" "$dest" | sha256sum -c -
}
# A complete security patch version is required, including for dependency builds.
go_bin="$(command -v go || true)"
if [[ -z "$go_bin" || $("$go_bin" version | awk '{print $3}') != go1.26.9 ]]; then
 case "$(uname -m)" in
  x86_64) arch=amd64; sum=42d158b4d8f7b61ac0a830567c940a86098fb7aac52e467a5ebec03ef5cc2f8d;;
  aarch64) arch=arm64; sum=4a97373d49fcacdcf3694fea368a500b00ee3e963974f3e7514132717632f052;;
  i?86) arch=386; sum=dea88a548986e02d02f6f24f7ea8a886cd5d6e583419247955a140cd84a74d77;;
  armv6*|armv7*) arch=armv6l; sum=4e7427224d6200800b8c8b9b05b2cb1868df0fb950f57843cf51b1060d91a0c8;;
  *) echo 'No pinned optional AWG toolchain for this architecture.' >&2; exit 1;;
 esac
 fetch_checked "https://go.dev/dl/go1.26.9.linux-$arch.tar.gz" "$work/go.tgz" "$sum"
 tar -xzf "$work/go.tgz" -C "$work"
 go_bin="$work/go/bin/go"
fi
fetch_checked 'https://codeload.github.com/amnezia-vpn/amneziawg-go/tar.gz/1cc94272ca8e9e223a5fe76382f5880f09d3c12d' "$work/awg-go.tgz" '5797a9c6f889bdf91624d68ce13152568b23fc4b5c04b5d175d6580934dfd146'
fetch_checked 'https://codeload.github.com/amnezia-vpn/amneziawg-tools/tar.gz/ee0f0a9aa34ff0a0da4b3433b9512781cfe02843' "$work/awg-tools.tgz" '22438f231d39ea27e4bdc69707ec3103357f3cecc9ad3eecba18a61e5b39c9b9'
mkdir "$work/engine" "$work/tools"
tar -xzf "$work/awg-go.tgz" -C "$work/engine" --strip-components=1
tar -xzf "$work/awg-tools.tgz" -C "$work/tools" --strip-components=1
(
 cd "$work/engine"
 export GOTOOLCHAIN=local
 "$go_bin" get golang.org/x/crypto@v0.57.0 golang.org/x/net@v0.60.0 golang.org/x/sys@v0.48.0
 "$go_bin" build -trimpath -o "$work/amneziawg-go" .
)
make -C "$work/tools/src" -j2
mkdir -p "$BK_EXTRA_CORES"
install -m 700 "$work/amneziawg-go" "$BK_EXTRA_CORES/amneziawg-go"
install -m 700 "$work/tools/src/wg" "$BK_EXTRA_CORES/awg"
cp "$work/engine/LICENSE" "$BK_EXTRA_CORES/AWG-GO-LICENSE"
cp "$work/tools/COPYING" "$BK_EXTRA_CORES/AWG-TOOLS-LICENSE"
echo 'Optional AWG userspace core installed; no kernel module is required.'
