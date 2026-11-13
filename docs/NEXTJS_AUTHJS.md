# Integrating an existing Next.js/Auth.js application

The Go service is **only a resource server**. Your existing Auth.js/Google login remains in Next.js.

Recommended request flow:
1. Browser calls a protected **server-side** Next.js route handler or server action.
2. Next.js calls `auth()` server-side and rejects missing sessions.
3. Next.js obtains a trusted, stable internal user ID for the session (do **not** trust an arbitrary user ID from the browser).
4. Next.js signs a short-lived Ed25519 **API access JWT**, containing `iss=nextjs-app`, `aud=go-api`, `sub=<trusted-user-id>`, `iat`, `exp` and `scope`.
5. Next.js calls Go with `Authorization: Bearer <JWT>` over HTTPS. The browser never receives this API token.
6. Go validates the JWT and authorizes each endpoint; the item repository restricts ownership by `sub`.

## Example: server-side Next.js only (illustrative, adapt to your Auth.js version)

Use `jose` in the **server runtime**; load the Ed25519 PKCS#8 private key from a secure server-side secret store. This is intentionally documentation, not a frontend embedded into the Go project.

```ts
import 'server-only';
import { SignJWT, importPKCS8 } from 'jose';
import { auth } from '@/auth';

export async function callGoItems() {
  const session = await auth();
  // Implement a trusted user-ID lookup using your Auth.js account/session store.
  const userId = session?.user?.id;
  if (!userId) throw new Error('Unauthorized');
  const privatePem = process.env.GO_API_SIGNING_PRIVATE_KEY_PEM;
  if (!privatePem) throw new Error('Signing key not configured');
  const key = await importPKCS8(privatePem, 'EdDSA');
  const token = await new SignJWT({ scope: 'items:read' })
    .setProtectedHeader({ alg: 'EdDSA', typ: 'JWT' })
    .setIssuer('nextjs-app')
    .setAudience('go-api')
    .setSubject(userId)
    .setIssuedAt()
    .setExpirationTime('5m')
    .sign(key);
  const response = await fetch(`${process.env.GO_API_URL}/api/v1/items`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: 'no-store',
  });
  if (!response.ok) throw new Error(`Go API returned ${response.status}`);
  return response.json();
}
```

**Important:** `session.user.id` is **not guaranteed** to exist with every Auth.js configuration. Add it via authenticated server-side callbacks / database adapter and validate its provenance. `scope` must come from trusted server-side policy, never directly from the browser. For POST, issue `items:write` only when authorized. The issuer must be configured identically on the Go side.

Keep signing keys out of `NEXT_PUBLIC_*`, source control and browser bundles. For production, rotate keys with overlapping key IDs / JWKS or a key-management service; the starter's single static public key deliberately favors simplicity, requiring coordinated rotation. Logout stops future issuance but does not instantly revoke previously minted five-minute tokens.
