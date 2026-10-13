import { test } from 'node:test';
import assert from 'node:assert/strict';
import { directMarkup, readDirect } from '../panel/js/ui/directedit.js';
import { fill } from '../panel/js/ui/form.js';

const settings = {
  carrier: 'spoof', holdsPorts: true,
  spoof: { profile: 'udp', uplink: '', downlink: '', srcIPs: '198.51.100.10',
    peerIP: '203.0.113.20', peerSrcIP: '198.51.100.20', interface: 'eth0',
    shufflePort: true, portMin: 49152, portMax: 65535, sockBuf: 8388608,
    mtu: 1500, xdpIface: 'eth0', ttlJitter: true, paddingMax: 80, icmpReply: true },
};
const input = (name, value) => typeof value === 'boolean'
  ? { name, type: 'checkbox', checked: value }
  : { name, type: 'text', value: String(value) };
const root = (...inputs) => ({ querySelectorAll: () => inputs });

test('direct Edit offers source pool, peer addresses and independent rotation on both sides', () => {
  for (const holdsPorts of [true, false]) {
    const markup = directMarkup({ ...settings, holdsPorts });
    for (const name of ['srcIPs', 'peerIP', 'peerSrcIP', 'interface', 'shufflePort', 'portMin', 'portMax', 'sockBuf']) {
      assert.ok(markup.includes(`name="spoof.${name}"`), `missing ${name}`);
    }
    assert.match(markup, /destination remains the peer's real IP/);
  }
  const plain = directMarkup({ carrier: 'udp', holdsPorts: true });
  assert.ok(!plain.includes('spoof.'));
  assert.ok(plain.includes('name="paths"'));
});

test('server settings fill the visible address fields and rotation controls without swapping sources', () => {
  const controls = [...directMarkup(settings).matchAll(/(?:\bname|data-name)="(spoof\.[^"]+)"/g)]
    .map(m => ({ ...input(m[1], m[1] === 'spoof.shufflePort' ? false : ''), tagName: 'INPUT' }));
  const form = { querySelectorAll: selector => selector === '[name]' ? controls : [] };
  const value = { ...settings, spoof: { ...settings.spoof, srcIPs: '198.51.100.10, 198.51.100.11' } };
  fill(form, value);
  const visible = Object.fromEntries(controls.map(n => [n.name, n.type === 'checkbox' ? n.checked : n.value]));
  assert.equal(visible['spoof.srcIPs'], value.spoof.srcIPs);
  assert.equal(visible['spoof.peerIP'], '203.0.113.20');
  assert.equal(visible['spoof.peerSrcIP'], '198.51.100.20');
  assert.equal(visible['spoof.shufflePort'], true);
  assert.equal(visible['spoof.portMin'], '49152');
  assert.deepEqual(readDirect({ querySelectorAll: () => controls }, value).spoof, value.spoof);
});

test('saving an IP change sends a nested carrier, preserves hidden tuning and uses numeric and boolean JSON types', () => {
  const out = readDirect(root(input('spoof.srcIPs', '198.51.100.11, 198.51.100.12'),
    input('spoof.peerIP', settings.spoof.peerIP), input('spoof.peerSrcIP', settings.spoof.peerSrcIP),
    input('spoof.shufflePort', true), input('spoof.portMin', '49152'), input('spoof.portMax', '65535'),
    input('spoof.sockBuf', '8388608'), input('mtu', '1280'), input('autoMtu', false), input('stealth', false)), settings);
  assert.deepEqual(out, { spoof: { ...settings.spoof, srcIPs: '198.51.100.11, 198.51.100.12' },
    mtu: 1280, autoMtu: false, stealth: false });
  assert.equal(settings.spoof.srcIPs, '198.51.100.10');
  assert.ok(!Object.keys(out).some(k => k.startsWith('spoof.')));
});

test('explicitly turning rotation off reaches the API without losing the profile', () => {
  const out = readDirect(root(input('spoof.shufflePort', false)), settings);
  assert.equal(out.spoof.shufflePort, false);
  assert.equal(out.spoof.profile, 'udp');
});

test('non-spoof edits never submit a carrier drawer', () => {
  assert.deepEqual(readDirect(root(input('ports', '443')), { carrier: 'udp' }), { ports: '443' });
});

test('invalid numeric edits fail before saving; profile labels are escaped', () => {
  for (const value of ['oops', 'Infinity', '49152.5']) {
    assert.throws(() => readDirect(root(input('spoof.portMin', value)), settings), /whole number/);
  }
  assert.ok(!directMarkup({ ...settings, spoof: { profile: '<script>alert(1)</script>' } }).includes('<script>'));
});
