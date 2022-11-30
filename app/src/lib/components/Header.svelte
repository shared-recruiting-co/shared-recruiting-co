<script lang="ts">
	import { page } from '$app/stores';
	import { supabaseClient } from '$lib/supabase/client';

	const handleLogout = async () => {
		await supabaseClient.auth.signOut();
	};
</script>

<header
	class="sticky top-0 z-50 flex flex-row flex-nowrap items-center justify-between space-x-4 bg-white px-4  py-4 text-lg sm:text-xl"
>
	<div class="flex w-96 flex-row items-center space-x-4">
		<img src="/logo.svg" alt="Shared Recruiting Co" class="h-6 md:h-10" />
		<a href="/" class="text-xl text-slate-900 sm:text-2xl md:min-w-[150px]">Shared Recruiting Co.</a
		>
	</div>
	<div class="hidden md:flex md:space-x-4">
		<a href="/#product" class="font-medium text-slate-500 hover:text-slate-900">Product</a>
		<a href="/#process" class="font-medium text-slate-500 hover:text-slate-900">Process</a>
		<a href="#faqs" class="font-medium text-slate-500 hover:text-slate-900">FAQs</a>
	</div>

	<div class="flex w-96 items-center justify-end space-x-4">
		<a
			href="https://github.com/shared-recruiting-co/shared-recruiting-co"
			target="_blank"
			rel="noopener noreferrer"
		>
			<svg fill="currentColor" viewBox="0 0 24 24" class="h-10 w-10 text-slate-900">
				<path
					fill-rule="evenodd"
					d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
					clip-rule="evenodd"
				/>
			</svg>
		</a>
		{#if $page.data.session}
			<button
				class="group inline-flex items-center justify-center rounded-full bg-slate-900 py-2 px-4 text-sm font-semibold text-white hover:bg-slate-700 hover:text-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 active:bg-slate-800 active:text-slate-300"
				on:click={handleLogout}>Log Out</button
			>
		{/if}
	</div>
</header>
