# bk component language

**Core** is the BackPack-derived engine, panel and management program, prepared
from the pinned upstream source and the explicit bk distribution overlay.

**Addons** is the independent Package 2/3 program and its own forwarding,
dependency-preparation and test support.

**Launcher** is the small `bk` entry program. It selects a component through
process arguments and standard input/output, without importing its Go code.

**Pristine snapshot** is `upstream/BackPack`. A core update changes this snapshot
and its lock, not addon source.

**Overlay** is the reviewed set of bk naming, publisher, path and core behavior
changes in `integration/`. It is applied when preparing `.build/core`.

The core's domain vocabulary is retained in `upstream/BackPack/CONTEXT.md`.
