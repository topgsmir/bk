# Tunnel Package 3

Build: **bk → 12 → Build Tunnel Package 3**.
Test: **bk → 0 → Test Tunnel Package 3**. Install the same bk version on both
servers. The original bk engines and original Connection Test remain unchanged.

| Family | Methods | TCP/UDP test cases | Source |
| --- | ---: | ---: | --- |
| Dagger | 65 | 130 | supplied `new2/dagger`, dagger-rs 0.2.1 by ir_spoof |
| Solarpass | 16 | 30 | supplied `new2/solarpass-tun`, core 2.3.0 |
| Backhaul | 16 | 24 | supplied `new2/backhaul`, core 2.0.3 |
| Eylan VPN protocols | 11 | 22 | protocol families in `new2/EylanPanel`; stock engines |
| Total | 108 | 206 | optional independent adapters |

Dagger includes TCP, KCP, HTTP/S, WS/S, DC6, Quantum+, XHTTP/XHTTPS
(stream/packet/auto, both directions), and Quantum/Gaming/TUN profiles
(TCP, UDP, ICMP, GRE, IPIP, BIP, RAW, DCPI). All forward TCP and UDP. DC6
requires assigned IPv6 on both servers; its test asks for these addresses.

Dagger also has 21 dedicated IP Spoof variants: all seven ordinary raw
profiles in Quantum, Gaming and TUN. Source-only, destination-only and combined
outer addresses can be entered separately for each side, in both test and setup.
These fields do not alter the real next hop, add routes, or weaken Noise
authentication. Your actual network must deliver the chosen outer headers.
DCPI and IP spoofing cannot be combined. Quantum+ is the separate UDP/KCP/FEC
carrier; it does not accept the raw spoof options.

Solarpass includes TCP/MUX, WS/S/MUX, QUIC/MUX, Hysteria, Spoof, Gamepass,
Ultimatepass, Speedpass, Ultimategamepass, Speedpassmux and Arenapass.
Spoof and Gamepass are UDP-only; other carriers expose both protocols.
The Spoof return-port repair is checked against complete original and repaired
hashes. Its confirmed peer-filter byte-order workaround is internal; users
enter normal IPv4 addresses. See [binary analysis](../recovered-package3/README.md).

Backhaul includes TCP, TCPMUX, XTCPMUX, WS, WSS, WSMUX, WSSMUX, XWSMUX,
AnyTLS, TUN TCP and TUN IPX over ICMP/IPIP/UDP/TCP/GRE/BIP. TCP and TUN
methods expose UDP too; other streams expose TCP only. Its supplied core
requires non-private public IPv4 actually assigned to the NIC. A NAT-only
public address does not satisfy this requirement; bk reports the cause.

Eylan methods are WireGuard, OpenVPN TCP/UDP, AnyConnect TLS/DTLS,
L2TP/IPsec, and sing-box VLESS/VMess/Trojan/Shadowsocks/Hysteria2. All
expose TCP and UDP. These are independent paired VPNs using stock engines.
The supplied web panel, licensing bypass, signing keys and administrator
accounts are not installed or distributed.

## Setup

1. Choose **12 → Setup Iran**, then the family/method. Enter assigned IPs,
   unused ports and forwarding endpoints. Each instance needs a unique name,
   tunnel ID and unused private subnet. AnyConnect and L2TP/IPsec need Iran
   `.1` and kharej `.2` in one unused `/30`.
2. Copy the secret `bk://e.` link to **12 → Apply setup link** on kharej.
   Start the kharej backend. Treat setup and test links like passwords.
3. On Iran choose **0 → Test Tunnel Package 3 → Iran**, select a family,
   then **Normal Tests** (Dagger/Solarpass) and one numbered group. Choose
   **Test This Group** or one listed method. Each group produces at most ten
   measured TCP/UDP rows, counting each forwarded protocol separately. There
   is no full-matrix or whole-family shortcut. Copy the test link to that menu's
   **Kharej** option. Only the selected dependencies are prepared on both servers.
   Dagger IP Spoof and Solarpass Spoof have a separate **IP Spoof Tests** category;
   normal groups never require spoof addresses. DC6 still needs assigned IPv6.
4. The test starts actual paired cores, exchanges 60 echoes, checks several
   payload sizes and transfers a byte-identical 1 MiB TCP payload. Solarpass
   UDP also checks concurrent clients. DTLS requires a real DTLS handshake;
   TLS fallback cannot pass as DTLS.

At most four forwarding cases run together; fixed-port L2TP/IPsec runs alone.
Run another group with a new test link after the selected group finishes.
Ctrl+C stops temporary instances. Test instances never become permanent units.
A failed install, missing kernel, unavailable address or unsupported architecture
produces SETUP-FAIL with its cause, rather than an unexplained SKIPPED row.
OK means real traffic passed on that pair of servers. No test can make an ISP's
filtered protocol or an unavailable kernel driver work.

## Dependencies and isolation

Supplied Dagger, Solarpass and Backhaul cores support Linux amd64. Other
architectures receive explicit setup errors for these families; original bk
transports remain available. Dagger's MIT source and dependency notices are
preserved in [assets/package3](../assets/package3). Solarpass and Backhaul
remain supplied executables, not recovered original source projects. No license
is inferred for a closed executable from its adapter's open-source license.

sing-box 1.14.3 comes from the official SagerNet release, checked against
pinned archive and installed-core hashes. Other Eylan methods use included
WireGuard userspace code or distribution packages. Server executables for
AnyConnect and L2TP are extracted without deploying a global VPN daemon.

Each instance owns its named unit, TUN, configuration and temporary firewall
rules. Stop removes owned state; shared firewall tables and XFRM state are
never flushed. Conflicting interfaces, subnets, VPN ports or IPsec owners are
refused. Solarpass/Backhaul have protected private kernel-setting paths.
L2TP uses private runtime and PPP directories without replacing existing VPN
configuration.

L2TP/IPsec requires free UDP **500, 4500, 1701**, `/dev/ppp`, and ESP/XFRM
kernel support. Its paired IKEv2 authentication uses AES-256-GCM; an owned
rule blocks L2TP lacking incoming IPsec protection. AnyConnect verifies a
paired certificate and public-key pin and honors the negotiated MTU.
Backhaul TCP forwarding adds paired mutual TLS. Native UDP maps add an
authenticated encrypted session envelope between the paired adapters.

## Source and verification

Based on BackPack by Amin Mohammadi (AminMGMT)
https://github.com/AminMGMT/BackPack

Adapters maintained by [topgsmir](https://github.com/topgsmir/bk).
Official references: [sing-box](https://github.com/SagerNet/sing-box),
[OpenVPN](https://openvpn.net/community-docs/community-articles/openvpn-2-6-manual.html),
[OpenConnect](https://www.infradead.org/openconnect/manual.html),
[ocserv](https://ocserv.gitlab.io/www/manual.html),
[strongSwan](https://docs.strongswan.org/docs/5.9/swanctl/swanctlConf.html).

The mandatory Tunnel Package 3 workflow runs native 60-echo matrices for all
families and verifies cleanup. Local native testing uses two Linux network
namespaces and the production bk binary. These tests prove paired setup and
traffic in the lab; test the actual Iran/kharej route before selecting a method.

## خلاصهٔ فارسی

<div dir="rtl">

پکیج تانل ۳ در گزینهٔ **۱۲** برای ساخت و در **۰ ← Test Tunnel Package 3**
برای تست قرار دارد. این بسته شامل ۱۰۸ روش از Dagger، Solarpass، Backhaul
و پروتکل‌های استاندارد Eylan است؛ مجموع آزمون‌های TCP و UDP آن ۲۰۶ حالت است.
کد موتورهای اصلی bk تغییر نکرده و فایل‌های اولیهٔ `new2` محفوظ هستند.

روی هر دو سرور یک نسخهٔ یکسان نصب کن. ابتدا سمت ایران را بساز و لینک
محرمانهٔ تنظیمات را در خارج وارد کن. آدرس محلی باید روی کارت شبکهٔ سرور
وجود داشته باشد. برای هر تونل نام، شناسه، درگاه و زیرشبکهٔ آزاد انتخاب کن.
برای AnyConnect و L2TP/IPsec آدرس ایران باید .۱ و آدرس خارج .۲ در یک /۳۰ باشد.

تست ابتدا وابستگی‌ها را آماده می‌کند و سپس تونل واقعی را می‌سازد. هر حالت
۶۰ پاسخ، اندازه‌های مختلف بسته و برای TCP انتقال یک مگابایت داده را بررسی
می‌کند. کل تست حدود ۶۰ تا ۹۰ دقیقه طول می‌کشد و با Ctrl+C متوقف می‌شود.
خطای نصب، نبودن درایور یا معماری پشتیبانی‌نشده، همراه دلیل با SETUP-FAIL
نمایش داده می‌شود. سبز شدن در آزمایش محلی، تضمین باز بودن مسیر شرکت‌های
اینترنتی نیست؛ مسیر ایران و خارج خودت را با همین منو اندازه بگیر.

هسته‌های ارائه‌شدهٔ Dagger، Solarpass و Backhaul مخصوص Linux amd64 هستند.
DC6 به IPv6 واقعی در هر دو سمت نیاز دارد و Backhaul به IPv4 عمومیِ واقعاً
نصب‌شده روی کارت شبکه. L2TP/IPsec از UDP ۵۰۰، ۴۵۰۰ و ۱۷۰۱ استفاده می‌کند؛
اگر VPN دیگری این درگاه‌ها را گرفته باشد، تنظیماتش دست‌کاری نمی‌شود.
MTU توافق‌شدهٔ AnyConnect رعایت می‌شود و گزینهٔ DTLS تنها با اتصال DTLS
واقعی پاس می‌شود. پنل وب Eylan، کد دور زدن مجوز و کلیدهای آن نصب نمی‌شوند.

مدیریت، توقف، لاگ و حذف سرویس‌های پکیج در **۱۲ ← Manage** است. فایل‌های
تنظیمات در `/etc/bk/external` هستند و فعلاً بکاپ موتور اصلی آن‌ها را شامل
نمی‌شود؛ جداگانه بکاپ بگیر. لینک‌های تست و ساخت را مانند رمز نگه دار.

</div>

*Last verified against bk v1.11.2.*
