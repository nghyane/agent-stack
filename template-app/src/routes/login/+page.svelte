<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';

	let mode = $state<'login' | 'register'>('login');
	let email = $state('');
	let password = $state('');
	let name = $state('');
	let pending = $state(false);
	let errored = $state(false);

	async function onSubmit(e: SubmitEvent) {
		e.preventDefault();
		pending = true;
		errored = false;
		try {
			if (mode === 'login') {
				await auth.login(email, password);
			} else {
				await auth.register(email, password, name);
			}
			await goto('/');
		} catch {
			errored = true;
		} finally {
			pending = false;
		}
	}
</script>

<div class="mx-auto max-w-sm">
	<h1 class="text-2xl font-bold">
		{mode === 'login' ? 'Log in' : 'Create account'}
	</h1>

	<form onsubmit={onSubmit} class="mt-6 space-y-4">
		{#if mode === 'register'}
			<input
				class="w-full rounded border px-3 py-2"
				placeholder="Name"
				bind:value={name}
			/>
		{/if}
		<input
			class="w-full rounded border px-3 py-2"
			type="email"
			placeholder="Email"
			bind:value={email}
			required
		/>
		<input
			class="w-full rounded border px-3 py-2"
			type="password"
			placeholder="Password (min 8 chars)"
			bind:value={password}
			required
		/>

		{#if errored}
			<p class="text-sm text-red-600">
				{mode === 'login'
					? 'Invalid credentials'
					: 'Could not register (email may be taken)'}
			</p>
		{/if}

		<button
			type="submit"
			disabled={pending}
			class="w-full rounded bg-blue-600 px-4 py-2 text-white disabled:opacity-50"
		>
			{pending ? '…' : mode === 'login' ? 'Log in' : 'Sign up'}
		</button>
	</form>

	<button
		onclick={() => (mode = mode === 'login' ? 'register' : 'login')}
		class="mt-4 text-sm text-blue-600"
	>
		{mode === 'login' ? 'Need an account? Register' : 'Have an account? Log in'}
	</button>
</div>
