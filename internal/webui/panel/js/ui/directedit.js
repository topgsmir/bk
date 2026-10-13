import { esc } from '../lib/dom.js';

const sw = (name, title, text) => `<div class="tg"><div class="tx"><b>${title}</b><span>${text}</span></div>` +
  `<div class="sww" data-name="${name}"><i></i></div></div>`;
const field = (name, label, hint = '') => `<div class="f"><label>${label}</label>` +
  `<input name="${name}" type="text">${hint ? `<div class="hint">${hint}</div>` : ''}</div>`;

export function directMarkup(set) {
  const ports = !!set.holdsPorts;
  const spoof = set.carrier === 'spoof';
  const profile = set.spoof || {};
  return `<div class="pane" data-tab="Direct">` +
    (ports ? field('ports', 'Forwarded ports', 'Comma separated. Blank leaves a plain TUN tunnel.') +
      sw('acceptUdp', 'Accept UDP', 'Carry UDP over the forwarded ports.') : '') +
    `<div class="f"><label>Performance preset</label><div class="sel" data-name="preset">Turbo` +
    `<span class="sp"></span><svg class="x" viewBox="0 0 24 24"><path d="M6 9l6 6 6-6"/></svg></div></div>` +
    `<div class="two">${field('mtu', 'MTU')}` +
    (set.carrier === 'udp' ? field('paths', 'Sockets', '1 is one socket; 2–8 spread it.') : '<div></div>') + `</div>` +
    sw('autoMtu', 'Auto MTU', 'Measure the path and correct the MTU once up.') +
    sw('fec', 'Error correction (FEC)', 'Spare packets repair loss. Both ends must match.') +
    (spoof ? `<div class="f"><label>IP spoofing · ${esc(profile.profile || 'udp')}</label>` +
      `<div class="hint">Forged addresses change the packet source. The destination remains the peer's real IP.</div></div>` +
      field('spoof.srcIPs', 'This server’s forged source IPv4', 'Comma separated for a source pool. Blank uses the real source.') +
      field('spoof.peerIP', 'Other server’s real IPv4', 'The real destination for this server’s packets.') +
      field('spoof.peerSrcIP', 'Expected forged IPv4 from the other server', 'Blank accepts any authenticated source; leave blank for a rotating pool. Update the peer’s expected source if you change your forged source.') +
      field('spoof.interface', 'Network interface', 'Blank selects the route automatically.') +
      sw('spoof.shufflePort', 'Rotate source port', 'Independent of Stealth. The destination port stays fixed.') +
      `<div class="two">${field('spoof.portMin', 'Source port minimum', '0 uses the default range.')}` +
      `${field('spoof.portMax', 'Source port maximum', '0 uses the default range.')}</div>` +
      field('spoof.sockBuf', 'Packet socket buffer (bytes)', '0 uses the engine default.') +
      sw('stealth', 'Stealth', 'Padding and header cosmetics. Both ends must match.') : '') +
    (ports ? `<div class="two">${field('maxConnections', 'Max connections', '0 is no limit.')}` +
      `${field('bandwidthMbps', 'Bandwidth (Mbit/s)', '0 is no limit.')}</div>` : '') +
    `</div>`;
}

const DIRECT_NUMBERS = new Set(['mtu', 'paths', 'maxConnections', 'bandwidthMbps',
  'spoof.portMin', 'spoof.portMax', 'spoof.sockBuf']);

export function readDirect(root, settings) {
  // SpoofTune replaces the whole carrier table. Retain settings this form does
  // not edit (direction profiles, outer MTU, XDP, ICMP replies and cosmetics).
  const out = settings.carrier === 'spoof' ? { spoof: { ...settings.spoof } } : {};
  root.querySelectorAll('.panes input[name]').forEach(n => {
    const v = n.type === 'checkbox' ? n.checked
            : DIRECT_NUMBERS.has(n.name) ? Number(n.value.trim())
            : n.value;
    if (typeof v === 'number' && (!Number.isFinite(v) || !Number.isInteger(v))) {
      throw new Error(`${n.name} must be a whole number.`);
    }
    if (n.name.startsWith('spoof.')) {
      if (settings.carrier === 'spoof') out.spoof[n.name.slice(6)] = v;
    } else {
      out[n.name] = v;
    }
  });
  return out;
}
