import createClient from 'openapi-fetch';
import { env } from '$env/dynamic/public';
import type { paths } from './schema';

// Typed API client. The base URL is read at RUNTIME from PUBLIC_API_URL
// (`$env/dynamic/public`), so changing it on the host does not require a
// rebuild — set it as a runtime variable on Cloudflare and redeploy config.
// `credentials: 'include'` sends the session cookie cross-site (frontend on
// Pages, API on its own domain). Falls back to the local dev API.
export const api = createClient<paths>({
	baseUrl: env.PUBLIC_API_URL ?? 'http://localhost:8080',
	credentials: 'include'
});
