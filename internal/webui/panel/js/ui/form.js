/* A <select> keeps its first option when the value it is handed matches none
   of them, which is how "turbo" quietly displayed as "Balanced". The server
   speaks in values ("turbo", "wss"); the preview's options were written with
   labels. Match either, and say so when neither fits. */
function setControl(node, v) {
  if (node.type === 'checkbox') { node.checked = !!v; return; }
  const val = Array.isArray(v) ? v.join(', ') : String(v ?? '');
  if (node.tagName !== 'SELECT') { node.value = val; return; }
  const want = val.trim().toLowerCase();
  const hit = [...node.options].find(o =>
    o.value.trim().toLowerCase() === want || o.textContent.trim().toLowerCase() === want);
  if (hit) { node.value = hit.value; return; }
  /* Nothing matched: show the server's value rather than a wrong one. */
  const opt = new Option(val, val, true, true);
  node.add(opt);
}

const dig = (obj, path) => path.split('.').reduce((o, k) => (o ?? {})[k], obj);


/* The two halves of the API disagree on one name.
 *
 * Reading a tunnel answers with serverHost; writing one takes serverAddr. The
 * field has to be named for the write — a key the edit handler does not know
 * fails the whole save, and that is what it used to do — so the read is
 * bridged here instead. One alias in one place, rather than a form that cannot
 * be filled or cannot be saved. */
const READ_ALIAS = { serverAddr: 'serverHost' };

export function fill(root, values) {
  root.querySelectorAll('[name]').forEach(node => {
    let v = dig(values, node.name);
    if (v === undefined || v === null) v = dig(values, READ_ALIAS[node.name] || '');
    if (v === undefined || v === null) return;
    setControl(node, v);
  });
  root.querySelectorAll('[data-unwired]').forEach(n => {
    n.disabled = true;
    n.title = 'Not settable over the API — use the CLI';
  });
}

