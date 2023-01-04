<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
</script>

<body class="flex min-h-screen flex-col text-slate-900">
	<Header />
	<main class="flex-1">
		<slot />
	</main>
	<Footer />
</body>
