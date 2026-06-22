<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/auth.svelte';

	onMount(() => {
		auth.refresh();
	});
</script>

{#if auth.loading}
	<p>Loading…</p>
{:else if auth.user}
	<h1 class="text-2xl font-bold">Hello, {auth.user.name || auth.user.email}</h1>
	<p class="mt-2 text-gray-600">You are authenticated.</p>
	<button
		onclick={() => auth.logout()}
		class="mt-4 rounded bg-gray-200 px-4 py-2 hover:bg-gray-300"
	>
		Log out
	</button>
{:else}
	<h1 class="text-2xl font-bold">Welcome</h1>
	<p class="mt-2 text-gray-600">You are not logged in. Head to the Login page.</p>
{/if}
