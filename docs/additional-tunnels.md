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
| SSH pool | Iran connects to kharej; authenticated SSH channels with checked host keys | TCP |
| SSH reverse | Kharej connects to Iran; verified SSH remote forwarding | TCP |
| RGT reverse TCP | pinned RGT core; kharej connects to Iran | TCP |
| RGT reverse UDP | pinned RGT core; TCP control channel | UDP |
| RGT direct | the upstream manager's VXLAN direct method | TCP and UDP |
| Paqet | pinned upstream KCP/raw TCP core | TCP and UDP |
| Alghadir full stack | GRE -> IPsec ESP -> KCP/FEC -> udp2raw fake TCP -> obfs4 TCP | TCP and UDP |

## Setup

1. Run `bk`, select **11**, then choose the setup method. Setup and connection
   tests automatically prepare missing dependencies on each server. **Install
   dependencies** remains available for advance preparation. Downloads use pinned,
   SHA256-checked cores. AWG's optional
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

SSH setup first tries an existing readable, unencrypted key (`bk_external`,
`id_ed25519`, `id_rsa`, or `id_ecdsa`); you can choose another key path. If login
is not ready, it prepares a dedicated key and runs interactive `ssh-copy-id`.
Verify the peer's fingerprint before accepting it. A password may be requested
on your own terminal to authorize the key; it is never stored in a configuration,
share link, or test report. Encrypted/agent-only keys need a separate unencrypted
key for the unattended service. **11 -> Prepare SSH authentication** also allows
advance preparation. Password login alone does not prove forwarding is allowed.

SSH pool initiates from Iran to kharej and supports up to 64 connections. SSH
reverse initiates from kharej to Iran and carries multiple streams over one SSH
connection, reconnecting after a broken session. Prepare authentication on the
initiating server: Iran for direct SSH, kharej for reverse SSH. Reverse uses the
Iran SSH port/user and a separate unused **Source port** for its remote listener.
The remote listener requests `127.0.0.1`; Iran's bk frontend forwards its public
listen port to that listener, then SSH delivers traffic to kharej's local target.
Keep sshd's normal `GatewayPorts no` setting so the internal listener stays on
loopback. Forwarding policy must permit this path (`AllowTcpForwarding` and,
where configured, `PermitOpen` for direct or `PermitListen` for reverse). bk does
not rewrite the server's global sshd configuration. Failed authorization or
forwarding reports the actual error; it is not counted as a successful test.

Paqet uses one connection and a fixed source port, so NOTRACK/RST rules remain
scoped to this named tunnel instead of disabling tracking for other traffic.

## Real connection test

Iran asks for both assigned IPv4 addresses, lets you choose all methods or one
method (and the SSH port/user/key when applicable), starts throwaway tunnel processes
on unused ports and prints a `bk://et.` link. Paste it on kharej under the same
menu. Both servers prepare missing cores before measurements begin. Direct SSH
prepares the key on Iran; reverse SSH prepares it on kharej, where the terminal
may ask for the Iran host fingerprint/password. Allow up to 30 minutes for peer
preparation. Update **both servers** to v1.10.0 or newer: additional test links
now use plan version 2. The original probes and coordinator are retained; an
additional adapter synchronizes preparation and final peer diagnostics.

The additional test uses the original payload probe: one random echo per second
on a persistent connection for 60 seconds, reconnects after drops, records RTT
and loss, and checks a byte-identical 1 MiB TCP transfer/speed. It also verifies
64, 512 and 1200-byte payloads, and a 16 KiB payload for TCP. Layer-3 methods and
Paqet are exercised separately with TCP and UDP. Missing dependencies are
installed automatically where supported. All 18 selected TCP/UDP variants get
an explicit result. Installation, unsupported architecture/kernel, SSH
credentials or process startup failures produce **SETUP-FAIL** with the reason
from the relevant server, rather than SKIPPED. Actual network failures remain
**DOWN** or **UNSTABLE**; software preparation cannot guarantee a blocked WAN
protocol will work. Small echoes passing does not override failed bulk/payload
checks (for example the reported L2TP-IP timeout).

For non-interactive software preparation, run as root:

```sh
bk external prepare --kind all
# or only selected methods:
bk external prepare --kind paqet,rgt-tcp,awg,alghadir
```

This command does not authorize SSH keys. Use the interactive menu for SSH.
Kernel modules, verified downloaded cores and SSH authorization persist for
later runs; temporary tunnel resources described below are still removed.

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
مدیریت می‌شوند: GRE، دو حالت L2TPv3، AWG، SSH مستقیم و ریورس، سه حالت RGT، Paqet و Alghadir.
برای تست از **گزینهٔ ۰ ← Additional tunnels** استفاده کن؛ می‌توانی همهٔ روش‌ها
یا یک روش را انتخاب کنی. آزمون روی تونل واقعی، ۶۰ ثانیه انتقال داده، تأخیر،
قطع و وصل، صحت داده و اندازه‌های مختلف بسته را بررسی می‌کند. نبودن وابستگی
یا خطای راه‌اندازی نتیجهٔ موفق محسوب نمی‌شود.

در نسخهٔ جدید، وابستگی‌های روش انتخابی هنگام ساخت و تست خودکار آماده می‌شوند.
هر دو سرور را به نسخهٔ v1.10.0 به‌روزرسانی کن؛ سپس یک سمت
را بساز و لینک `bk://e.` را در گزینهٔ Apply setup link سمت دیگر وارد کن.
این لینک کلید مشترک دارد؛ آن را عمومی نکن. برای تست، ایران لینک موقت
`bk://et.` می‌دهد و خارج آن را در منوی تست وارد می‌کند. آدرس محلی باید IPv4
واقعاً اختصاص‌یافته به کارت شبکه باشد؛ حالت NAT و IPv6 در این افزونه پوشش
داده نشده است. فایروال هر دو سرور نیز باید پروتکل و درگاه روش انتخابی را
اجازه دهد. RGT به هستهٔ amd64 نیاز دارد و Alghadir فعلاً روی ARM64 پشتیبانی
نمی‌شود. SSH فقط TCP است. در حالت مستقیم ایران به خارج وصل می‌شود و در
حالت ریورس خارج به ایران. اثرانگشت سرور مقصد را بررسی کن؛ ممکن است برای
ثبت کلید، همان‌جا رمز ورود درخواست شود. رمز در تنظیمات ذخیره نمی‌شود.
خطای نصب یا آماده‌سازی با SETUP-FAIL و دلیل دقیق نمایش داده می‌شود؛
قطع واقعی مسیر با DOWN یا UNSTABLE نمایش داده می‌شود.

Alghadir تازه، همهٔ لایه‌های GRE، IPsec، KCP، udp2raw و obfs4 را در یک مسیر
داده به هم وصل می‌کند. کلید نشست جدید است و منابع داخل فضای شبکهٔ جداگانه
ساخته می‌شوند. اسکریپت‌های اولیهٔ پوشهٔ `new` تغییر نکرده‌اند و اجرا نمی‌شوند.

سرویس‌ها، تنظیمات و حذف هر تونل در گزینهٔ **۱۱ ← Manage** مستقل هستند.
پنل، بکاپ و بازنشانی گروهی موتورهای اصلی فعلاً این تونل‌های افزوده را مدیریت
نمی‌کنند؛ از `/etc/bk/external` جداگانه بکاپ بگیر و بعد از به‌روزرسانی برنامه،
سرویس‌های افزوده را از همین منو دوباره راه‌اندازی کن. حذف کامل bk همهٔ
سرویس‌های آن را متوقف و تنظیماتش را پاک می‌کند.

</div>

*Last verified against bk v1.10.0.*
