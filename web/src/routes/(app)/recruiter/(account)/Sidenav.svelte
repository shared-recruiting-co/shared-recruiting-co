<script lang="ts">
	import { page } from '$app/stores';
	import { fade } from 'svelte/transition';
	import { afterNavigate } from '$app/navigation';

	$: ({ supabase, profile } = $page.data);
	$: initial = profile?.firstName.charAt(0).toUpperCase();

	const mobileWidth = 768;
	let windowWidth: number;
	let isOpen = false;

	$: {
		// always set is open to false on mobile
		if (isOpen && windowWidth > mobileWidth) {
			isOpen = false;
		}
	}

	// close navigation on mobile when the route changes
	afterNavigate(() => {
		isOpen = false;
	});

	const handleLogout = async () => {
		await supabase.auth.signOut();
	};

	const nav = [
		{
			name: 'Account',
			href: '/recruiter/profile',
			icon: `
  <path fill-rule="evenodd" d="M7.5 6a4.5 4.5 0 119 0 4.5 4.5 0 01-9 0zM3.751 20.105a8.25 8.25 0 0116.498 0 .75.75 0 01-.437.695A18.683 18.683 0 0112 22.5c-2.786 0-5.433-.608-7.812-1.7a.75.75 0 01-.437-.695z" clip-rule="evenodd" />
			`
		},
		{
			name: 'Jobs',
			href: '/recruiter/jobs',
			icon: `
  <path stroke-linecap="round" stroke-linejoin="round" d="M20.25 14.15v4.25c0 1.094-.787 2.036-1.872 2.18-2.087.277-4.216.42-6.378.42s-4.291-.143-6.378-.42c-1.085-.144-1.872-1.086-1.872-2.18v-4.25m16.5 0a2.18 2.18 0 00.75-1.661V8.706c0-1.081-.768-2.015-1.837-2.175a48.114 48.114 0 00-3.413-.387m4.5 8.006c-.194.165-.42.295-.673.38A23.978 23.978 0 0112 15.75c-2.648 0-5.195-.429-7.577-1.22a2.016 2.016 0 01-.673-.38m0 0A2.18 2.18 0 013 12.489V8.706c0-1.081.768-2.015 1.837-2.175a48.111 48.111 0 013.413-.387m7.5 0V5.25A2.25 2.25 0 0013.5 3h-3a2.25 2.25 0 00-2.25 2.25v.894m7.5 0a48.667 48.667 0 00-7.5 0M12 12.75h.008v.008H12v-.008z" />
			`
		},
		{
			name: 'Settings',
			href: '/recruiter/settings',
			icon: `
  <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12a7.5 7.5 0 0015 0m-15 0a7.5 7.5 0 1115 0m-15 0H3m16.5 0H21m-1.5 0H12m-8.457 3.077l1.41-.513m14.095-5.13l1.41-.513M5.106 17.785l1.15-.964m11.49-9.642l1.149-.964M7.501 19.795l.75-1.3m7.5-12.99l.75-1.3m-6.063 16.658l.26-1.477m2.605-14.772l.26-1.477m0 17.726l-.26-1.477M10.698 4.614l-.26-1.477M16.5 19.794l-.75-1.299M7.5 4.205L12 12m6.894 5.785l-1.149-.964M6.256 7.178l-1.15-.964m15.352 8.864l-1.41-.513M4.954 9.435l-1.41-.514M12.002 12l-3.75 6.495" />
			`
		}
	];

	const additionalLinks = [
		{
			name: 'Documentation',
			href: '/docs',
			icon: `
  <path stroke-linecap="round" stroke-linejoin="round" d="M12 7.5h1.5m-1.5 3h1.5m-7.5 3h7.5m-7.5 3h7.5m3-9h3.375c.621 0 1.125.504 1.125 1.125V18a2.25 2.25 0 01-2.25 2.25M16.5 7.5V18a2.25 2.25 0 002.25 2.25M16.5 7.5V4.875c0-.621-.504-1.125-1.125-1.125H4.125C3.504 3.75 3 4.254 3 4.875V18a2.25 2.25 0 002.25 2.25h13.5M6 7.5h3v3H6v-3z" />
			`
		},
		{
			name: 'Help',
			href: 'https://github.com/shared-recruiting-co/shared-recruiting-co/discussions',
			icon: `
  <path stroke-linecap="round" stroke-linejoin="round" d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9 5.25h.008v.008H12v-.008z" />
			`
		}
	];
</script>

<svelte:window bind:innerWidth={windowWidth} />

<aside>
	<!-- Off-canvas menu for mobile, show/hide based on off-canvas menu state. -->
	<div class="relative z-40 md:hidden" role="dialog" aria-modal="true">
		<!--
      Off-canvas menu backdrop, show/hide based on off-canvas menu state.
    -->
		{#if isOpen}
			<div
				class="fixed inset-0 bg-slate-600 bg-opacity-75 transition-opacity duration-300 ease-linear"
				aria-hidden={isOpen ? 'false' : 'true'}
				class:opacity-0={!isOpen}
				class:opacity-100={isOpen}
				transition:fade
			/>
		{/if}

		<div
			class="fixed inset-0 z-40 flex"
			class:-translate-x-full={!isOpen}
			class:translate-x-0={isOpen}
		>
			<!--
        Off-canvas menu, show/hide based on off-canvas menu state.
      -->
			<div
				class="relative flex w-full max-w-xs flex-1 transform flex-col bg-white transition duration-300 ease-in-out"
				class:-translate-x-full={!isOpen}
				class:translate-x-0={isOpen}
			>
				<!--
          Close button, show/hide based on off-canvas menu state.
        -->
				{#if isOpen}
					<div
						class="absolute top-0 right-0 -mr-12 pt-2 duration-300 ease-in-out"
						class:opacity-0={!isOpen}
						class:opacity-100={isOpen}
						transition:fade
					>
						<button
							type="button"
							class="ml-1 flex h-10 w-10 items-center justify-center rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
							on:click={() => (isOpen = false)}
						>
							<span class="sr-only">Close sidebar</span>
							<!-- Heroicon name: outline/x-mark -->
							<svg
								class="h-6 w-6 text-white"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								aria-hidden="true"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
							</svg>
						</button>
					</div>
				{/if}

				<div class="h-0 flex-1 pt-5 pb-4">
					<div class="flex flex-shrink-0 items-center space-x-4 px-4">
						<img class="h-8 w-auto" src="/logo.svg" alt="SRC" />
						<span class="font-medium">Recruiter</span>
					</div>
					<!-- Hack: add padding to bottom to push additionalLinks above the bottom section -->
					<nav class="mt-5 flex h-full flex-col px-2 pb-12">
						<div class="h-full flex-1 space-y-1">
							{#each nav as item}
								{@const current = $page.url.pathname === item.href}
								<a
									href={item.href}
									data-sveltekit-preload-data="hover"
									class="group flex items-center rounded-md px-2 py-2 text-base font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900"
									class:bg-slate-100={current}
									class:text-slate-900={current}
								>
									<svg
										class="mr-4 h-6 w-6 flex-shrink-0 text-slate-400 group-hover:text-slate-500"
										class:text-slate-500={current}
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="1.5"
										stroke="currentColor"
										aria-hidden="true"
									>
										{@html item.icon}
									</svg>
									{item.name}
								</a>
							{/each}
						</div>
						{#each additionalLinks as link}
							<a
								href={link.href}
								class="group flex items-center rounded-md px-2 py-2 text-base font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900"
								target="_blank"
								rel="noopener noreferrer"
							>
								<svg
									class="mr-4 h-6 w-6 flex-shrink-0 text-slate-400 group-hover:text-slate-500"
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="1.5"
									stroke="currentColor"
									aria-hidden="true"
								>
									{@html link.icon}
								</svg>
								{link.name}
							</a>
						{/each}
					</nav>
				</div>
				<div class="flex flex-shrink-0 border-t border-slate-200 p-4">
					<div class="flex w-full items-center">
						<div
							class="flex h-9 w-9 items-center justify-center rounded-full bg-blue-100"
							alt={initial}
						>
							{initial}
						</div>
						<div class="ml-3">
							<p class="text-sm font-medium text-slate-700 group-hover:text-slate-900">
								{profile?.firstName}
								{profile?.lastName}
							</p>
						</div>
						<div class="flex flex-1 justify-end">
							<button on:click={handleLogout} title="Logout">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="1.5"
									stroke="currentColor"
									class="h-6 w-6 text-slate-700 hover:text-slate-900"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75"
									/>
								</svg>
							</button>
						</div>
					</div>
				</div>
			</div>
			<div class="w-14 flex-shrink-0">
				<!-- Force sidebar to shrink to fit close icon -->
			</div>
		</div>
	</div>

	<!-- Static sidebar for desktop -->
	<div class="hidden md:fixed md:inset-y-0 md:flex md:w-64 md:flex-col">
		<!-- Sidebar component, swap this element with another sidebar if you like -->
		<div class="flex min-h-0 flex-1 flex-col border-r border-slate-200 bg-white">
			<div class="flex flex-1 flex-col overflow-y-auto pt-5 pb-4">
				<div class="flex flex-shrink-0 items-center space-x-4 px-4">
					<img class="h-8 w-auto" src="/logo.svg" alt="SRC" />
					<span class="text-lg font-medium">Recruiter</span>
				</div>
				<nav class="mt-5 flex h-full flex-col px-2">
					<div class="flex-1 space-y-1">
						{#each nav as item}
							{@const current = $page.url.pathname === item.href}
							<!-- Current: "bg-slate-100 text-slate-900", Default: "text-slate-600 hover:bg-slate-50 hover:text-slate-900" -->
							<a
								href={item.href}
								class="group flex items-center rounded-md px-2 py-2 text-base font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900"
								class:bg-slate-100={current}
								class:text-slate-900={current}
							>
								<!--
                Heroicon name: outline/home

                Current: "text-slate-500", Default: "text-slate-400 group-hover:text-slate-500"
              -->
								<svg
									class="mr-4 h-6 w-6 flex-shrink-0 text-slate-400 group-hover:text-slate-500"
									class:text-slate-500={current}
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="1.5"
									stroke="currentColor"
									aria-hidden="true"
								>
									{@html item.icon}
								</svg>
								{item.name}
							</a>
						{/each}
					</div>
					{#each additionalLinks as link}
						<a
							href={link.href}
							class="group flex items-center rounded-md px-2 py-2 text-base font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900"
							target="_blank"
							rel="noopener noreferrer"
						>
							<svg
								class="mr-4 h-6 w-6 flex-shrink-0 text-slate-400 group-hover:text-slate-500"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								aria-hidden={!isOpen}
							>
								{@html link.icon}
							</svg>
							{link.name}
						</a>
					{/each}
				</nav>
			</div>
			<div class="flex flex-shrink-0 flex-col border-t border-slate-200 p-4">
				<div class="flex w-full items-center">
					<div
						class="flex h-9 w-9 items-center justify-center rounded-full bg-blue-100"
						alt={initial}
					>
						{initial}
					</div>
					<div class="ml-3">
						<p class="text-sm font-medium text-slate-700 group-hover:text-slate-900">
							{profile?.firstName}
							{profile?.lastName}
						</p>
					</div>
					<div class="flex flex-1 justify-end">
						<button on:click={handleLogout} title="Logout">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								class="h-6 w-6 text-slate-700 hover:text-slate-900"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75"
								/>
							</svg>
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
	<!-- END: Static sidebar for desktop -->

	<!-- Open Sidebar button on Mobile -->
	<div class="relative -mr-10 flex flex-1 flex-col md:mr-0 md:pl-64">
		<div class="sticky top-0 z-10 bg-white pl-1 pt-1 sm:pl-3 sm:pt-3 md:hidden">
			<button
				type="button"
				class="-ml-0.5 -mt-0.5 inline-flex h-12 w-12 items-center justify-center rounded-md text-slate-500 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
				on:click={() => (isOpen = true)}
			>
				<span class="sr-only">Open sidebar</span>
				<!-- Heroicon name: outline/bars-3 -->
				<svg
					class="h-6 w-6"
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
	</div>
</aside>
