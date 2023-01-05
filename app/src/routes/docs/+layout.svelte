<script lang="ts">
	import Header from './Header.svelte';
	import Sidenav from './Sidenav.svelte';
</script>

<Header />
<div class="relative mx-auto flex max-w-8xl flex-1 justify-center sm:px-2 lg:px-8 xl:px-12">
	<div class="hidden lg:relative lg:block lg:flex-none">
		<div class="absolute inset-y-0 right-0 w-[50vw] bg-slate-50 dark:hidden" />
		<Sidenav />
	</div>
	<main class="flex-1">
		<slot />
	</main>
</div>
