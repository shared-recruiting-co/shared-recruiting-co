<script lang="ts">
	import { fly } from 'svelte/transition';
	import { afterNavigate } from '$app/navigation';

	import Navigation from './Navigation.svelte';

	let isOpen = false;

	// close navigation on mobile when the route changes
	afterNavigate(() => {
		isOpen = false;
	});
</script>

<!-- Desktop Menu -->
<aside
	class="sticky top-[4.5rem] -ml-0.5 hidden h-[calc(100vh-4.5rem)] overflow-y-auto overflow-x-hidden py-16 pl-0.5 lg:block"
>
	<Navigation />
</aside>
<!-- Mobile Menu -->
<aside class="flex lg:hidden">
	<div class="relative">
		<button
			type="button"
			class="flex items-center justify-center rounded-md text-slate-500 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-sky-500"
			on:click={() => (isOpen = true)}
		>
			<span class="sr-only">Open sidebar</span>
			<!-- Heroicon name: outline/bars-3 -->
			<svg
				class="h-8 w-8"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				aria-hidden={isOpen ? 'true' : 'false'}
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
				/>
			</svg>
		</button>
	</div>
	{#if isOpen}
		<div
			class="fixed inset-0 z-50 flex items-start overflow-y-auto bg-slate-900/50 pr-10 backdrop-blur lg:hidden"
			transition:fly={{ x: -100, duration: 200 }}
		>
			<div class="min-h-full w-full max-w-xs bg-white px-4 pt-4 pb-12 dark:bg-slate-900 sm:px-6">
				<div class="flex items-center space-x-4">
					<button
						type="button"
						class="mr-4 flex items-center justify-center rounded-md text-slate-500 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-sky-500"
						on:click={() => (isOpen = false)}
					>
						<span class="sr-only">Close sidebar</span>
						<!-- Heroicon name: outline/x-mark -->
						<svg
							class="h-8 w-8 text-slate-500"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							aria-hidden={!isOpen}
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
					<a href="/" class="block">
						<img src="/logo.svg" alt="Shared Recruiting Co" class="h-8 w-auto" />
					</a>
				</div>
				<div class="mt-8 px-1">
					<Navigation />
				</div>
			</div>
		</div>
	{/if}
	<!-- </Dialog.Panel> -->
	<!-- </Dialog> -->
</aside>
