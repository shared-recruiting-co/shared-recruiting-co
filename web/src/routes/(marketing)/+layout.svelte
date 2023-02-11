<script lang="ts">
	import Header from './Header.svelte';
	import Footer from './Footer.svelte';
</script>

<Header />
<main class="flex-1">
	<slot />
</main>
<Footer />
