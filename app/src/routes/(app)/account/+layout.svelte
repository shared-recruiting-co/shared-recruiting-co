<script lang="ts">
	import Sidenav from './Sidenav.svelte';
</script>

<div class="flex h-full flex-1">
	<Sidenav />
	<main class="flex-1">
		<slot />
	</main>
</div>
