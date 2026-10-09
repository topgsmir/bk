import assert from 'node:assert/strict';
import crypto from 'node:crypto';
import fs from 'node:fs';
const source = fs.readFileSync('internal/app/app.go','utf8');
const publicBytes = Buffer.from(source.match(/var ReleasePublicKey = "([^"]+)"/)[1],'base64');
assert.equal(publicBytes.length,32);
const publicKey = crypto.createPublicKey({key:Buffer.concat([Buffer.from('302a300506032b6570032100','hex'),publicBytes]),format:'der',type:'spki'});
const tag = fs.readFileSync('VERSION','utf8').trim();
assert.equal(tag,'v1.8.5.1');
const sums = fs.readFileSync('release/SHA256SUMS');
const signature = Buffer.from(fs.readFileSync('release/SHA256SUMS.sig','utf8').trim(),'base64');
const message = Buffer.concat([Buffer.from(`backpack release ${tag}\n`),sums]);
assert(crypto.verify(null,message,publicKey,signature),'Release signature is invalid');
assert(!crypto.verify(null,Buffer.concat([Buffer.from('backpack release v9.9.9\n'),sums]),publicKey,signature),'Signature accepts the wrong release tag');
let count = 0;
for (const line of sums.toString('utf8').trim().split('\n')) {
  const [expected,name] = line.trim().split(/\s+/);
  assert.match(name,/^backpack_linux_(amd64|arm64|386|s390x|armv[567])\.tar\.gz$/);
  const actual = crypto.createHash('sha256').update(fs.readFileSync(`release/${name}`)).digest('hex');
  assert.equal(actual,expected,`Checksum mismatch: ${name}`);
  count++;
}
assert.equal(count,7);
console.log(`PASS: ${count} release archives, independent publisher signature and tag binding.`);
