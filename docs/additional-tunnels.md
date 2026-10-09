# Additional tunnels (optional)

The original bk engines and setup wizards are unchanged. Main menu **11** opens
an independent manager. Main menu **0** now offers the original Connection Test
and a separate Additional Tunnels test.

## Supported methods

| Method | Path | Forwarded traffic |
| --- | --- | --- |
| GRE IPv4 | IPv4 protocol 47, keyed GRE | TCP and UDP |
| L2TPv3 IP | IPv4 protocol 115 | TCP and UDP |
| L2TPv3 UDP | paired UDP port | TCP and UDP |
| AmneziaWG | obfuscated encrypted WireGuard, kernel or pinned userspace core | TCP and UDP |
| SSH pool | authenticated SSH channels with checked host keys | TCP |
| RGT reverse TCP | pinned RGT core; kharej connects to Iran | TCP |
| RGT reverse UDP | pinned RGT core; TCP control channel | UDP |
| RGT direct | the upstream manager's VXLAN direct method | TCP and UDP |
| Paqet | pinned upstream KCP/raw TCP core | TCP and UDP |
| Alghadir full stack | GRE -> IPsec ESP -> KCP/FEC -> udp2raw fake TCP -> obfs4 TCP | TCP and UDP |

## Setup

1. Run `bk`, select **11**, then **Install dependencies** for the chosen method
   on both servers. This downloads pinned, SHA256-checked cores. AWG's optional
   userspace builder uses verified source/toolchain archives and does not change
   bk's Go modules. RGT's upstream binary supports amd64 only. Alghadir's pinned
   udp2raw binary currently supports amd64, 386 and 32-bit ARM; native ARM64 is
   explicitly refused instead of installing an incompatible binary. L2TP loads
   the required kernel drivers; Ubuntu's matching extra-module package is
   installed when needed. Kernels without L2TP support are reported as errors.
2. Select Setup Iran or Setup Kharej, choose the method and enter the assigned
   IPv4 addresses, port, paired ID/secret, tunnel subnet (where applicable), MTU
   and forwarding endpoints. Each additional tunnel has its own named unit.
3. Copy the resulting `bk://e.` link to **11 -> Apply setup link** on the other
   server. The link contains the shared secret; keep it private. Machine-local
   SSH keys and physical-interface/MAC overrides are never copied in the link.
4. Start the backend service on kharej, then use **0 -> Additional tunnels** on
   both servers to verify the real path. A running systemd unit alone is not
   evidence of successful traffic. Network/cloud firewalls must permit the
   method's required protocol/port.

Use different tunnel names, IDs and unused private/benchmark tunnel subnets
(/24 through /31) for simultaneous instances. Overlapping local interfaces are
refused; temporary test plans are restricted to /30 subnets and loopback echoes.
The IPv4 entered as local must be assigned to a NIC. This initial implementation
uses IPv4; NAT-only public addresses and IPv6 need separate handling.

SSH: prepare the dedicated key through **11 -> Prepare SSH authentication**.
Check the host fingerprint against kharej before accepting it. Password login is
used only interactively by ssh-copy-id; the tunnel uses key authentication.
SSH supports arbitrary local/remote ports and up to 64 parallel connections.
Paqet uses one connection and a fixed source port, so NOTRACK/RST rules remain
scoped to this named tunnel instead of disabling tracking for other traffic.

## Real connection test

Iran asks for both assigned IPv4 addresses, lets you choose all methods or one
method (and the SSH port/user/key when applicable), starts throwaway tunnel processes
on unused ports and prints a `bk://et.` link. Paste it on kharej under the same
menu. The original short-exchange TCP/UDP coordinator is reused unchanged.

The additional test uses the original payload probe: one random echo per second
on a persistent connection for 60 seconds, reconnects after drops, records RTT
and loss, and checks a byte-identical 1 MiB TCP transfer/speed. It also verifies
64, 512 and 1200-byte payloads, and a 16 KiB payload for TCP. Layer-3 methods and
Paqet are exercised separately with TCP and UDP. Missing dependencies are
reported explicitly; they are never counted as successful tests.

Test resources are transient: no enabled service, permanent config, rc.local
edit or cron entry is created. Stop with Ctrl+C. Only the test's own interfaces,
firewall rules, processes and files are removed.

## Alghadir

The supplied script started disconnected GRE/IPsec/KCP/udp2raw/obfs4 processes.
The optional implementation actually connects every layer. Its private network
namespace contains GRE, transport-mode AES-GCM IPsec, KCP/FEC, udp2raw and two
TUN packet boundaries. The outer raw packets are framed over an authenticated
obfs4 stream between the two VPS. IPsec uses fresh directional keys derived
from authenticated random nonces for each outer session, replay protection and
an enforced packet lifetime; restart never reuses AES-GCM nonces. It does not
change /etc/ipsec.conf, restart a shared strongSwan daemon or flush global rules.
Only obfs4 TCP is exposed on the public path; the inner protocol ports stay in
the private namespace. This is a new composition, not a claim that the supplied
script already had working layered protection.

## Storage and removal

Configurations: `/etc/bk/external/<name>.json`, mode 0600.
Cores: `/etc/bk/external/cores/`.
Units: `bk-ext-<name>.service`; command remains `/usr/local/bin/bk`.
Manage/start/stop/restart/log/delete through **11 -> Manage**. Deletion stops the
supervised process first, allowing it to undo its own kernel resources. The
existing bk uninstaller also stops `bk-*.service` before removing /etc/bk.
The original backup/panel/watchdog currently manage the original engines;
additional tunnels are managed through menu 11 and their own systemd units.
Back up `/etc/bk/external` separately until panel/backup integration is added.

## Sources

- The five supplied local files under `new/` are retained unchanged as
  references. Their concepts are implemented here; the original scripts are
  neither sourced, executed nor redistributed by this manager.
- RGT manager: https://github.com/black-sec/RGT
  (reference commit d71d5e18a7e783f80c2961db86fc4fe84bdb8abf, MIT).
- Paqet manager: https://github.com/Ramin-Setoodehnia/Paqet-Tunnel
  (configuration reference only; its script is not redistributed).
- Paqet core: https://github.com/hanselime/paqet, v1.0.0-alpha.21.
- AmneziaWG: https://github.com/amnezia-vpn/amneziawg-go (v0.2.19) and
  https://github.com/amnezia-vpn/amneziawg-tools
  (ee0f0a9aa34ff0a0da4b3433b9512781cfe02843).
- udp2raw: https://github.com/wangyu-/udp2raw, 20230206.0.
- obfs4proxy: the distribution's signed package, following Tor's managed
  transport specification https://spec.torproject.org/pt-spec/.

bk remains based on BackPack by Amin Mohammadi (AminMGMT).

<div dir="rtl">

## خلاصهٔ فارسی

موتورهای اصلی bk و تنظیماتشان حفظ شده‌اند. امکانات تازه از **گزینهٔ ۱۱**
مدیریت می‌شوند: GRE، دو حالت L2TPv3، AWG، SSH، سه حالت RGT، Paqet و Alghadir.
برای تست از **گزینهٔ ۰ ← Additional tunnels** استفاده کن؛ می‌توانی همهٔ روش‌ها
یا یک روش را انتخاب کنی. آزمون روی تونل واقعی، ۶۰ ثانیه انتقال داده، تأخیر،
قطع و وصل، صحت داده و اندازه‌های مختلف بسته را بررسی می‌کند. نبودن وابستگی
یا خطای راه‌اندازی نتیجهٔ موفق محسوب نمی‌شود.

ابتدا در هر دو سرور از گزینهٔ ۱۱ وابستگی روش انتخابی را نصب کن. سپس یک سمت
را بساز و لینک `bk://e.` را در گزینهٔ Apply setup link سمت دیگر وارد کن.
این لینک کلید مشترک دارد؛ آن را عمومی نکن. برای تست، ایران لینک موقت
`bk://et.` می‌دهد و خارج آن را در منوی تست وارد می‌کند. آدرس محلی باید IPv4
واقعاً اختصاص‌یافته به کارت شبکه باشد؛ حالت NAT و IPv6 در این افزونه پوشش
داده نشده است. فایروال هر دو سرور نیز باید پروتکل و درگاه روش انتخابی را
اجازه دهد. RGT به هستهٔ amd64 نیاز دارد و Alghadir فعلاً روی ARM64 پشتیبانی
نمی‌شود. SSH فقط TCP است و قبل از استفاده باید کلید و اثرانگشت سرور خارج
را تأیید کنی؛ در تست نیز می‌توانی درگاه و کاربر SSH را مشخص کنی.

Alghadir تازه، همهٔ لایه‌های GRE، IPsec، KCP، udp2raw و obfs4 را در یک مسیر
داده به هم وصل می‌کند. کلید نشست جدید است و منابع داخل فضای شبکهٔ جداگانه
ساخته می‌شوند. اسکریپت‌های اولیهٔ پوشهٔ `new` تغییر نکرده‌اند و اجرا نمی‌شوند.

سرویس‌ها، تنظیمات و حذف هر تونل در گزینهٔ **۱۱ ← Manage** مستقل هستند.
پنل، بکاپ و بازنشانی گروهی موتورهای اصلی فعلاً این تونل‌های افزوده را مدیریت
نمی‌کنند؛ از `/etc/bk/external` جداگانه بکاپ بگیر و بعد از به‌روزرسانی برنامه،
سرویس‌های افزوده را از همین منو دوباره راه‌اندازی کن. حذف کامل bk همهٔ
سرویس‌های آن را متوقف و تنظیماتش را پاک می‌کند.

</div>

*Last verified against bk v1.9.0.*
