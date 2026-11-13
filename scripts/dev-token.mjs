// Development-only signing helper. Node.js 20+ recommended.
import { readFileSync } from 'node:fs';
import { sign } from 'node:crypto';
const priv = readFileSync('.local/keys/api-private.pem');
const now = Math.floor(Date.now()/1000);
const claims = { iss:'nextjs-app', aud:'go-api', sub:'dev-user', iat:now, exp:now+300, scope:'items:read items:write' };
const encoded = (input) => Buffer.from(JSON.stringify(input)).toString('base64url');
const data = `${encoded({alg:'EdDSA',typ:'JWT'})}.${encoded(claims)}`;
console.log(`${data}.${sign(null,Buffer.from(data),priv).toString('base64url')}`);
