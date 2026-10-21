# Core distribution overlay

upstream/BackPack is an exact, unedited source snapshot of the pinned upstream commit.
Core preparation applies the mechanical bk naming transform, then core.patch.
The patch preserves bk paths, publisher identity, PCK loopback, Spoof defaults,
and the direct-edit fixes. It contains no Package 2 or Package 3 implementation.
The prepared core is disposable build output in .build/core.
Do not edit upstream files or generated .build sources. Change the overlay only
for an intentional core fix. Addon work belongs in addons/.
