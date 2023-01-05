<script lang="ts">
	import Header from './Header.svelte';
	import Sidenav from './Sidenav.svelte';
</script>

<Header />
<div class="relative mx-auto flex max-w-8xl flex-1 justify-center py-16 sm:px-2 lg:px-8 xl:px-12">
	<Sidenav />
	<main class="flex-1">
		<slot />
	</main>
</div>
