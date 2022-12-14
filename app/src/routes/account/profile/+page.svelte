<script lang="ts">
	import { page } from '$app/stores';
	import { slide, draw, fade } from 'svelte/transition';

	import { supabaseClient } from '$lib/supabase/client';

	let saved = false;
	let errors: Record<string, string> = {};

	let debounceTimeout: NodeJS.Timeout;
	const debounceDelay = 1000;
	const savedMessageTimeout = 3000;

	const formError = (e: typeof errors, field: string) => {
		return e[field] || '';
	};

	const debounce = (func: (...args: any[]) => void, wait: number) => {
		return function executedFunction(...args: any[]) {
			clearTimeout(debounceTimeout);
			debounceTimeout = setTimeout(() => func(...args), wait);
		};
	};

	const handleInput = async (e: Event) => {
		const target = e.target as HTMLInputElement;
		const value = target.value;
		const name = target.name;

		// right now all input values are required
		if (!value) {
			errors[name] = 'This field is required';
			return;
		}
		// clear errors
		errors[name] = '';

		const { data, error } = await supabaseClient
			.from('user_profile')
			.update({ [name]: value })
			.eq('user_id', $page.data.session?.user.id)
			.select()
			.maybeSingle();
		if (!error && data && data[name as keyof typeof data] === value) {
			saved = true;
			setTimeout(() => {
				saved = false;
			}, savedMessageTimeout);
			return;
		}
		errors[name] = 'There was an error saving your changes';
	};

	// debounce input to limit database writes
	const debouncedHandleInput = debounce(handleInput, debounceDelay);
</script>

<div class="my-12 lg:grid lg:grid-cols-12 lg:gap-x-5">
	<aside class="hidden py-6 px-2 sm:px-6 lg:col-span-3 lg:block lg:py-0 lg:px-0">
		<nav class="space-y-1">
			<!-- Current: "bg-slate-50 text-blue-700 hover:text-indigo-700 hover:bg-white", Default: "text-slate-900 hover:text-slate-900 hover:bg-slate-50" -->
			<a
				href="/account/profile"
				class="group flex items-center rounded-md bg-slate-50 px-3 py-2 text-sm font-medium text-blue-700 hover:bg-white hover:text-indigo-700"
				aria-current="page"
			>
				<!--
          Heroicon name: outline/user-circle

          Current: "text-blue-500 group-hover:text-indigo-500", Default: "text-slate-400 group-hover:text-slate-500"
        -->
				<svg
					class="-ml-1 mr-3 h-6 w-6 flex-shrink-0 text-blue-500 group-hover:text-indigo-500"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					aria-hidden="true"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M17.982 18.725A7.488 7.488 0 0012 15.75a7.488 7.488 0 00-5.982 2.975m11.963 0a9 9 0 10-11.963 0m11.963 0A8.966 8.966 0 0112 21a8.966 8.966 0 01-5.982-2.275M15 9.75a3 3 0 11-6 0 3 3 0 016 0z"
					/>
				</svg>
				<span class="truncate">Profile</span>
			</a>
		</nav>
	</aside>

	<div class="space-y-6 sm:px-6 lg:col-span-9 lg:px-0">
		<h1 class="px-4 text-3xl sm:text-4xl">Account</h1>
		<form action="#" method="POST">
			<div class="relative shadow sm:overflow-hidden sm:rounded-md">
				{#if saved}
					<div
						class="absolute top-0 right-0 mt-6 mr-8 flex items-center space-x-2 text-green-600"
						in:slide
						out:fade
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="h-6 w-6"
						>
							<path
								transition:draw
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
						<span>Saved</span>
					</div>
				{/if}
				<div class="space-y-6 bg-white py-6 px-4 sm:p-6">
					<div>
						<h3 class="text-lg font-medium leading-6 text-slate-900">Profile</h3>
						<p class="mt-1 text-sm text-slate-500">
							Let us know how people typically address you in emails.
						</p>
					</div>

					<div class="grid grid-cols-6 gap-6">
						<div class="col-span-6 sm:col-span-4">
							<label for="email-address" class="block text-sm font-medium text-slate-700"
								>Email address</label
							>
							<input
								type="text"
								name="email-address"
								id="email-address"
								autocomplete="email"
								disabled
								value={$page.data.profile.email}
								class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 disabled:cursor-not-allowed disabled:border-slate-200 disabled:bg-slate-50 disabled:text-slate-500 sm:text-sm"
							/>
						</div>
						<div class="col-span-6 sm:col-span-3">
							<label for="first_name" class="block text-sm font-medium text-slate-700"
								>First name</label
							>
							<input
								type="text"
								name="first_name"
								id="first_name"
								autocomplete="given-name"
								on:input={debouncedHandleInput}
								value={$page.data.profile.firstName}
								class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
							/>
							{#if formError(errors, 'first_name')}
								<p class="mt-1 text-xs text-rose-500">{formError(errors, 'first_name')}</p>
							{/if}
						</div>

						<div class="col-span-6 sm:col-span-3">
							<label for="last_name" class="block text-sm font-medium text-slate-700"
								>Last name</label
							>
							<input
								type="text"
								name="last_name"
								id="last_name"
								autocomplete="family-name"
								value={$page.data.profile.lastName}
								on:input={debouncedHandleInput}
								class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
							/>
							{#if formError(errors, 'first_name')}
								<p class="mt-1 text-xs text-rose-500">{formError(errors, 'first_name')}</p>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</form>
	</div>
</div>
