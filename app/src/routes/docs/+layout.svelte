<script lang="ts">
	import Header from './Header.svelte';
	import Sidenav from './Sidenav.svelte';
</script>

<Header />
<div class="flex flex-1">
	<Sidenav />
	<main class="flex-1">
		<slot />
	</main>
</div>
