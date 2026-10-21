import hashlib
import importlib.util
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest import mock
import zipfile

SPEC = importlib.util.spec_from_file_location('sync_upstream', Path(__file__).with_name('sync-upstream.py'))
sync = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(sync)

class SyncTests(unittest.TestCase):
    def archive(self, root, members):
        path = root / 'source.zip'
        with zipfile.ZipFile(path, 'w') as package:
            for name, contents in members.items():
                package.writestr(name, contents)
        return path

    def fixture(self, root):
        (root / 'upstream/BackPack').mkdir(parents=True)
        (root / 'upstream/BackPack/VERSION').write_text('old\n')
        (root / 'addons').mkdir()
        (root / 'addons/own.txt').write_bytes(b'own addon source')
        (root / 'integration').mkdir()
        (root / 'integration/upstream-files.json').write_text('{}\n')
        (root / 'tools/preparecore').mkdir(parents=True)
        (root / 'tools/preparecore/main.go').write_text('package main\n')
        (root / '.build').mkdir()
        (root / 'go.mod').write_text('module example.invalid/fixture\n')
        (root / 'upstream.lock.json').write_text('{"commit":"old"}\n')
        files = {'VERSION': b'v1.9.0\n', 'main.go': b'package main\n'}
        tree = {'tree': [{'type': 'blob', 'path': name, 'sha': hashlib.sha1(b'blob ' + str(len(raw)).encode() + b'\0' + raw).hexdigest()} for name, raw in files.items()]}
        archive = self.archive(root, {'BackPack-new/' + name: raw for name, raw in files.items()})
        args = type('Args', (), dict(commit='a' * 40, archive=str(archive), check_only=False))()
        return args, tree

    def test_downloaded_tree_and_overlay_pass_before_replacing_only_upstream(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            args, tree = self.fixture(root)
            with mock.patch.object(sync, 'ROOT', root), mock.patch.object(sync, 'fetch', return_value=json.dumps(tree).encode()), mock.patch.object(sync.subprocess, 'run', return_value=type('Result', (), {'returncode': 0})()):
                sync.run(args)
            self.assertEqual((root / 'upstream/BackPack/VERSION').read_text(), 'v1.9.0\n')
            self.assertEqual((root / 'addons/own.txt').read_bytes(), b'own addon source')
            self.assertEqual(json.loads((root / 'upstream.lock.json').read_text())['commit'], 'a' * 40)
            self.assertEqual(len(list((root / '.build').glob('upstream-before-*'))), 1)

    def test_overlay_conflict_keeps_current_pin_core_and_addons(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            args, tree = self.fixture(root)
            original_pin = (root / 'upstream.lock.json').read_bytes()
            with mock.patch.object(sync, 'ROOT', root), mock.patch.object(sync, 'fetch', return_value=json.dumps(tree).encode()), mock.patch.object(sync.subprocess, 'run', return_value=type('Result', (), {'returncode': 1})()):
                with self.assertRaisesRegex(RuntimeError, 'overlay needs review'):
                    sync.run(args)
            self.assertEqual((root / 'upstream.lock.json').read_bytes(), original_pin)
            self.assertEqual((root / 'upstream/BackPack/VERSION').read_text(), 'old\n')
            self.assertEqual((root / 'addons/own.txt').read_bytes(), b'own addon source')
            self.assertFalse((root / '.build/sync.lock').exists())

    def test_archive_paths_cannot_escape_the_snapshot(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            archive = self.archive(root, {'BackPack-root/../../escape': b'invalid'})
            with self.assertRaisesRegex(ValueError, 'Unsafe archive path'):
                sync.extract(archive, root / 'target')
            self.assertFalse((root / 'escape').exists())

if __name__ == '__main__':
    unittest.main()
