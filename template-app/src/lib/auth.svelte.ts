import { api } from './api/client';
import type { components } from './api/schema';

export type User = components['schemas']['UserResponse'];

// Auth state as a singleton rune store. Components read `auth.user` reactively;
// mutations update it in place. Equivalent role to TanStack Query's auth hooks,
// but with Svelte 5 runes and no extra dependency.
class AuthStore {
	user = $state<User | null>(null);
	loading = $state(true);

	// Fetch the current user. A 401 means "logged out" (null), not an error.
	async refresh(): Promise<void> {
		this.loading = true;
		const { data, error } = await api.GET('/auth/me');
		this.user = error ? null : data;
		this.loading = false;
	}

	async login(email: string, password: string): Promise<void> {
		const { data, error } = await api.POST('/auth/login', {
			body: { email, password }
		});
		if (error) throw error;
		this.user = data;
	}

	async register(email: string, password: string, name: string): Promise<void> {
		const { data, error } = await api.POST('/auth/register', {
			body: { email, password, name }
		});
		if (error) throw error;
		this.user = data;
	}

	async logout(): Promise<void> {
		await api.POST('/auth/logout');
		this.user = null;
	}
}

export const auth = new AuthStore();
