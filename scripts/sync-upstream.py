#!/usr/bin/env python3
"""Update only pristine BackPack; addons are never read, written, or generated."""
from pathlib import Path, PurePosixPath
import argparse, datetime, hashlib, json, os, re, shutil, subprocess, tempfile
import urllib.parse, urllib.request, zipfile

ROOT = Path(__file__).resolve().parents[1]
REPOSITORY = 'AminMGMT/BackPack'

def fetch(url):
    request = urllib.request.Request(url, headers={'User-Agent': 'bk-upstream-sync'})
    with urllib.request.urlopen(request, timeout=90) as response:
        return response.read()

def json_bytes(value):
    return (json.dumps(value, indent=2, ensure_ascii=False) + '\n').encode('utf-8')

def extract(archive, destination):
    with zipfile.ZipFile(archive) as package:
        prefixes = {PurePosixPath(item.filename).parts[0] for item in package.infolist()}
        if len(prefixes) != 1:
            raise ValueError('The archive must contain one repository root')
        for item in package.infolist():
            parts = PurePosixPath(item.filename).parts
            if item.filename.startswith('/') or '..' in parts or '\\' in item.filename or '.git' in parts or any(':' in part for part in parts):
                raise ValueError('Unsafe archive path: ' + item.filename)
            if item.external_attr >> 16 & 0o170000 == 0o120000:
                raise ValueError('Symlinks are not accepted in an upstream snapshot')
            target = destination.joinpath(*parts[1:])
            target.resolve().relative_to(destination.resolve())
            if item.is_dir():
                target.mkdir(parents=True, exist_ok=True)
            else:
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_bytes(package.read(item))

def run(args):
    lock_path = ROOT / '.build/sync.lock'
    lock_path.parent.mkdir(exist_ok=True)
    try:
        descriptor = os.open(lock_path, os.O_CREAT | os.O_EXCL | os.O_WRONLY, 0o600)
    except FileExistsError:
        raise RuntimeError('An upstream sync is already running; inspect .build/sync.lock')
    os.close(descriptor)
    try:
        sync(args)
    finally:
        lock_path.unlink(missing_ok=True)

def sync(args):
    if not re.fullmatch('[0-9a-fA-F]{40}', args.commit):
        raise ValueError('--commit must be the full reviewed 40-character commit SHA')
    commit = args.commit.lower()
    tree = json.loads(fetch(f'https://api.github.com/repos/{REPOSITORY}/git/trees/{commit}?recursive=1'))
    if tree.get('truncated'):
        raise RuntimeError('GitHub returned an incomplete source tree')
    with tempfile.TemporaryDirectory(prefix='sync-', dir=ROOT / '.build') as temporary:
        stage = Path(temporary)
        stage.resolve().relative_to((ROOT / '.build').resolve())
        archive = stage / 'upstream.zip'
        archive.write_bytes(Path(args.archive).read_bytes() if args.archive else fetch(f'https://codeload.github.com/{REPOSITORY}/zip/{commit}'))
        source = stage / 'upstream/BackPack'
        source.mkdir(parents=True)
        extract(archive, source)
        expected = {item['path']: item['sha'] for item in tree['tree'] if item['type'] == 'blob'}
        modes = {item['path']: item['mode'] for item in tree['tree'] if item['type'] == 'blob'}
        actual = {}
        hashes = {}
        for path in source.rglob('*'):
            if not path.is_file():
                continue
            data = path.read_bytes()
            name = path.relative_to(source).as_posix()
            actual[name] = hashlib.sha1(b'blob ' + str(len(data)).encode() + b'\0' + data).hexdigest()
            hashes[name] = hashlib.sha256(data).hexdigest()
        if actual != expected:
            raise RuntimeError('The downloaded archive does not match the reviewed GitHub commit')
        for name, mode in modes.items():
            if mode not in ('100644', '100755'):
                raise RuntimeError('Unsupported upstream file mode: ' + name)
            (source / name).chmod(0o755 if mode == '100755' else 0o644)
        shutil.copytree(ROOT / 'integration', stage / 'integration')
        (stage / 'integration/upstream-files.json').write_bytes(json_bytes(hashes))
        (stage / 'integration/upstream-modes.json').write_bytes(json_bytes(modes))
        shutil.copytree(ROOT / 'tools/preparecore', stage / 'tools/preparecore')
        shutil.copyfile(ROOT / 'go.mod', stage / 'go.mod')
        result = subprocess.run(['go', 'run', './tools/preparecore'], cwd=stage)
        if result.returncode:
            raise RuntimeError('The core overlay needs review. Current upstream and all addons have been kept unchanged.')
        if args.check_only:
            print('Source integrity and the core overlay passed. No tracked source was changed.')
            return
        old_lock = (ROOT / 'upstream.lock.json').read_bytes()
        old_manifest = (ROOT / 'integration/upstream-files.json').read_bytes()
        old_modes = (ROOT / 'integration/upstream-modes.json').read_bytes()
        metadata = json.loads(old_lock)
        metadata.update(commit=commit, version=(source / 'VERSION').read_text(encoding='utf-8').strip(), archive_sha256=hashlib.sha256(archive.read_bytes()).hexdigest())
        target = ROOT / 'upstream/BackPack'
        backup = ROOT / '.build' / ('upstream-before-' + datetime.datetime.now(datetime.timezone.utc).strftime('%Y%m%dT%H%M%S%f'))
        target.resolve().relative_to((ROOT / 'upstream').resolve())
        backup.resolve().relative_to((ROOT / '.build').resolve())
        target.rename(backup)
        try:
            source.rename(target)
            (ROOT / 'integration/upstream-files.json').write_bytes(json_bytes(hashes))
            (ROOT / 'integration/upstream-modes.json').write_bytes(json_bytes(modes))
            (ROOT / 'upstream.lock.json').write_bytes(json_bytes(metadata))
        except Exception:
            if target.exists():
                target.rename(stage / 'failed-source')
            backup.rename(target)
            (ROOT / 'upstream.lock.json').write_bytes(old_lock)
            (ROOT / 'integration/upstream-files.json').write_bytes(old_manifest)
            (ROOT / 'integration/upstream-modes.json').write_bytes(old_modes)
            raise
        print('Pinned upstream to', commit)
        print('Previous source retained in', backup.relative_to(ROOT))
        print('Next: make check; then review the source diff and build a signed release.')

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--commit', required=True, help='full reviewed upstream commit SHA')
    parser.add_argument('--archive', help='optional local ZIP, still verified against GitHub')
    parser.add_argument('--check-only', action='store_true')
    try:
        run(parser.parse_args())
    except Exception as error:
        parser.exit(1, str(error) + '\n')
