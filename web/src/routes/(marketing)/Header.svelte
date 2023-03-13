<script lang="ts">
	import { page } from '$app/stores';
	import { afterNavigate } from '$app/navigation';
	import { fly } from 'svelte/transition';
	import { supabaseClient } from '$lib/supabase/client';

	const handleLogout = async () => {
		await supabaseClient.auth.signOut();
	};

	let isMobileMenuOpen = false;
	// close navigation on mobile when the route changes
	afterNavigate(() => {
		isMobileMenuOpen = false;
	});
</script>

<header
	class="sticky top-0 z-10 flex flex-row flex-nowrap items-center justify-between space-x-4 bg-white px-4  py-4 text-lg sm:text-xl"
	aria-label="Main Navigation"
>
	<div class="flex w-96 flex-row items-center space-x-4">
		<a href="/" class="-m-1.5 p-1.5">
			<span class="sr-only">Shared Recruiting Co.</span>
			<img src="/logo.svg" alt="Shared Recruiting Co" class="h-6 md:h-10" />
		</a>
		<a href="/" class="hidden text-xl text-slate-900 sm:text-2xl md:block md:min-w-[150px]"
			>Shared Recruiting Co.</a
		>
		<a href="/" class="text-xl text-slate-900 sm:text-2xl md:hidden md:min-w-[150px]">SRC</a>
	</div>
	<div class="hidden md:space-x-4 lg:flex">
		<a
			data-sveltekit-preload-data="hover"
			href="/candidates"
			class="font-medium text-slate-500 hover:text-slate-900">Candidates</a
		>
		<a
			data-sveltekit-preload-data="hover"
			href="/companies"
			class="font-medium text-slate-500 hover:text-slate-900">Companies</a
		>
		<a
			data-sveltekit-preload-data="hover"
			href="/security"
			class="font-medium text-slate-500 hover:text-slate-900">Security</a
		>
	</div>
	<div class="flex w-96 items-center justify-end space-x-4">
		{#if $page.data.session}
			<a
				href="/account/profile"
				class="focus-visible:outline-blublue-600 ml-auto rounded-md bg-blue-600 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
				data-sveltekit-preload-data="hover"
			>
				Account
			</a>
			<button class="text-base hover:underline active:underline sm:text-lg" on:click={handleLogout}
				>Log Out</button
			>
		{:else}
			<a
				class="hidden text-base hover:underline active:underline sm:text-lg lg:block"
				data-sveltekit-preload-data="hover"
				href="/login">Log in</a
			>
			<a
				class="group inline-flex items-center justify-center rounded-md bg-blue-600 py-2 px-4 text-sm font-semibold text-white shadow-lg hover:bg-blue-700 hover:text-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-900 active:bg-blue-800 active:text-slate-300"
				data-sveltekit-preload-data="hover"
				href="/login">Join</a
			>
			<a
				href="https://github.com/shared-recruiting-co/shared-recruiting-co"
				target="_blank"
				rel="noopener noreferrer"
				class="hidden lg:block"
			>
				<svg
					fill="currentColor"
					viewBox="0 0 24 24"
					class="h-10 w-10 text-slate-900 hover:text-slate-700 active:text-slate-800"
				>
					<path
						fill-rule="evenodd"
						d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
						clip-rule="evenodd"
					/>
				</svg>
			</a>
		{/if}
		<div class="flex lg:hidden">
			<button
				type="button"
				class="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700"
				on:click={() => (isMobileMenuOpen = true)}
			>
				<span class="sr-only">Open main menu</span>
				<svg
					class="h-6 w-6"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					aria-hidden={!isMobileMenuOpen}
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
	<!-- Mobile menu, show/hide based on menu open state. -->
	{#if isMobileMenuOpen}
		<div class="lg:hidden" role="dialog" aria-modal="true">
			<!-- Background backdrop, show/hide based on slide-over state. -->
			<div class="fixed inset-0 z-10 backdrop-blur" transition:fly={{ duration: 100 }} />
			<div
				class="fixed inset-y-0 right-0 z-10 w-full overflow-y-auto bg-white px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-slate-900/10"
				in:fly={{ duration: 100 }}
			>
				<div class="flex items-center gap-x-6">
					<a href="/" class="-m-1.5 p-1.5">
						<span class="sr-only">Shared Recruiting Co.</span>
						<img class="h-8 w-auto" src="/logo.svg" alt="SRC" />
					</a>
					{#if $page.data.session}
						<a
							href="/account/profile"
							class="focus-visible:outline-blublue-600 ml-auto rounded-md bg-blue-600 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
							data-sveltekit-preload-data="hover"
						>
							Account
						</a>
					{:else}
						<a
							href="/login"
							class="focus-visible:outline-blublue-600 ml-auto rounded-md bg-blue-600 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
							data-sveltekit-preload-data="hover">Join</a
						>
					{/if}
					<button
						type="button"
						class="-m-2.5 rounded-md p-2.5 text-slate-700"
						on:click={() => (isMobileMenuOpen = false)}
					>
						<span class="sr-only">Close menu</span>
						<svg
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							aria-hidden={isMobileMenuOpen}
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				</div>
				<div class="mt-6 flow-root">
					<div class="-my-6 divide-y divide-slate-500/10">
						<div class="space-y-2 py-6">
							<a
								href="/candidates"
								class="-mx-3 block rounded-lg py-2 px-3 text-base leading-7 text-slate-900 hover:bg-slate-50"
								data-sveltekit-preload-data="hover">Candidates</a
							>
							<a
								href="/companies"
								class="-mx-3 block rounded-lg py-2 px-3 text-base leading-7 text-slate-900 hover:bg-slate-50"
								data-sveltekit-preload-data="hover">Companies</a
							>
							<a
								href="/security"
								class="-mx-3 block rounded-lg py-2 px-3 text-base leading-7 text-slate-900 hover:bg-slate-50"
								data-sveltekit-preload-data="hover">Security</a
							>
						</div>
						<div class="flex items-center justify-between space-x-2 py-4">
							{#if $page.data.session}
								<button
									class="text-base hover:underline active:underline sm:text-lg"
									on:click={handleLogout}>Log Out</button
								>
							{/if}
							<a
								href="https://github.com/shared-recruiting-co/shared-recruiting-co"
								target="_blank"
								rel="noopener noreferrer"
							>
								<svg
									fill="currentColor"
									viewBox="0 0 24 24"
									class="h-10 w-10 text-slate-900 hover:text-slate-700 active:text-slate-800"
								>
									<path
										fill-rule="evenodd"
										d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
										clip-rule="evenodd"
									/>
								</svg>
							</a>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</header>
