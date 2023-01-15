<script lang="ts">
	import Sidenav from './Sidenav.svelte';
</script>

<div class="flex h-full flex-1">
	<Sidenav />
	<main class="flex-1 py-16 lg:grid lg:grid-cols-12 lg:gap-x-5">
		<div class="space-y-6 px-2 sm:pl-12 sm:pr-4 lg:col-span-10 lg:col-start-2 lg:px-0">
			<slot />
		</div>
	</main>
</div>
