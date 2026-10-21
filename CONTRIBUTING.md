# Contributing to separated bk

Work in the owning component:

- Add or change optional tunnels in `addons/` and its native CI workflows.
- Keep `upstream/BackPack` pristine. Use `scripts/sync-upstream.py` for core updates.
- Put intentional bk core fixes in `integration/core.patch`, preserving the
  mechanical normalization and the inherited wire identifiers.
- Change the process boundary and entry menu in the root launcher only.

Run `make check` before publication, plus the appropriate native addon workflows.
The root, core and addon modules are tested independently. A passing unit test
does not replace paired native-core and actual-route checks.

Never edit `.build/core` as the permanent source of a change. It is regenerated.
Never add a module replacement or import coupling addons back to core internals.
The architecture tests enforce the source boundary and pristine snapshot.

Keep the AGPL-3.0 source, original attribution, third-party notices and source
assets. Follow SECURITY.md for private vulnerability reports.
