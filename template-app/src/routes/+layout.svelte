<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-gray-50 text-gray-900">
	<header class="border-b bg-white">
		<nav class="mx-auto flex max-w-3xl gap-4 px-4 py-3">
			<a href="/" class="font-medium hover:text-blue-600">Home</a>
			<a href="/login" class="font-medium hover:text-blue-600">Login</a>
		</nav>
	</header>
	<main class="mx-auto max-w-3xl px-4 py-8">
		{@render children()}
	</main>
</div>
