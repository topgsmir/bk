# bk addons

Package 2 and Package 3 are independently built and versioned here. This module
does not import the BackPack core or use a local module replacement for it.

```sh
go build -o bk-addons .
go test ./... -race -timeout 15m
./bk-addons package2
./bk-addons package3
```

Existing JSON configurations stay under `/etc/bk/external`; service names remain
`bk-ext-<name>.service`. Existing systemd units can continue calling
`/usr/local/bin/bk external run ...`; the launcher routes them to this program.

`internal/externaltunnel` owns the engines and dependency preparation.
`internal/manage` owns setup, optional-package inventory and paired live tests.
`internal/forwarder` is a small, independently maintained TCP/UDP forwarding
library. Test-coordinator, TUI, metrics and forwarding helpers were extracted
from the AGPL bk/BackPack sources into this module; they are not generated or
overwritten during core synchronization. Their attribution and source licence
remain in the distribution's LICENSE and NOTICE.

Core transport tests and core config editing are deliberately owned by the core
program. Addon service actions accept only `bk-ext-` service names.
