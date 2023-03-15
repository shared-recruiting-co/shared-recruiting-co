<script lang="ts">
	// import { slide, draw, fade } from 'svelte/transition';
	// import type { PageData } from './$types';
	// import { supabaseClient } from '$lib/supabase/client';

	// server page data
	// export let data: PageData;

	// ui state
	//	let profileSaved = false;
	//	let errors: Record<string, string> = {};

	//	let debounceTimeout: NodeJS.Timeout;
	//	const debounceDelay = 1000;
	//	const savedMessageTimeout = 3000;

	//	const formError = (e: typeof errors, field: string) => {
	//		return e[field] || '';
	//	};

	//	const debounce = (func: (...args: any[]) => void, wait: number) => {
	//		return function executedFunction(...args: any[]) {
	//			clearTimeout(debounceTimeout);
	//			debounceTimeout = setTimeout(() => func(...args), wait);
	//		};
	//	};

	// const handleInput = async (e: Event) => {
	// 	const target = e.target as HTMLInputElement;
	// 	const value = target.value;
	// 	const name = target.name;

	// 	// right now all input values are required
	// 	if (!value) {
	// 		errors[name] = 'This field is required';
	// 		return;
	// 	}
	// 	// clear errors
	// 	errors[name] = '';

	// 	const { data: profileData, error } = await supabaseClient
	// 		.from('recruiter')
	// 		.update({ [name]: value })
	// 		.eq('user_id', data.session?.user.id)
	// 		.select()
	// 		.maybeSingle();
	// 	if (!error && profileData && profileData[name as keyof typeof profileData] === value) {
	// 		profileSaved = true;
	// 		setTimeout(() => {
	// 			profileSaved = false;
	// 		}, savedMessageTimeout);
	// 		return;
	// 	}
	// 	errors[name] = 'There was an error saving your changes';
	// };

	// debounce input to limit database writes
	// const debouncedHandleInput = debounce(handleInput, debounceDelay);
</script>

<div>
	<h1 class="text-3xl sm:text-4xl">Welcome to SRC</h1>
	<p class="mt-1 text-sm text-slate-500">
		Right now, we require an onboarding call before you can get started with SRC.
	</p>
</div>
<div class="shadow sm:overflow-hidden sm:rounded-md">
	<div class="space-y-6 bg-sky-50 py-6 px-4 sm:p-6">
		<h3 class="text-xl font-medium leading-6">Onboarding Required</h3>
		<p class="mt-4 max-w-2xl">
			We're excited to have you on board! Before you can get started, we'd love to chat. We want to
			make sure SRC is a good fit for your team and that we can help you achieve your goals.
			<br />
			<br />
			Reach out to us at
			<a class="text-blue-600 underline" href="mailto:team@sharedrecruiting.co"
				>team@sharedrecruiting.co</a
			> to schedule a time to chat.
		</p>
	</div>
</div>
<!--
<div class="relative shadow sm:overflow-hidden sm:rounded-md">
	<div class="space-y-6 bg-white py-6 px-4 sm:p-6">
		<div>
			<h3 class="text-lg font-medium leading-6 text-slate-900">Company</h3>
			<p class="mt-1 text-sm text-slate-500">
				Company information. This information will be display to candidates.
			</p>
		</div>

		<div class="grid grid-cols-6 gap-6">
			<div class="col-span-6 sm:col-span-4">
				<label for="company_name" class="block text-sm font-medium text-slate-700">Name</label>
				<input
					type="text"
					name="company_name"
					id="company_name"
					disabled
					value={data.company?.name}
					class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 disabled:cursor-not-allowed disabled:border-slate-200 disabled:bg-slate-50 disabled:text-slate-500 sm:text-sm"
				/>
			</div>
			<div class="col-span-6 sm:col-span-4">
				<label for="company_website" class="block text-sm font-medium text-slate-700">Website</label
				>
				<input
					type="text"
					name="company_website"
					id="company_website"
					disabled
					value={data.company?.website}
					class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 disabled:cursor-not-allowed disabled:border-slate-200 disabled:bg-slate-50 disabled:text-slate-500 sm:text-sm"
				/>
			</div>
		</div>
	</div>
</div>
<div class="relative shadow sm:overflow-hidden sm:rounded-md">
	{#if profileSaved}
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
					value={data.profile?.email}
					class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 disabled:cursor-not-allowed disabled:border-slate-200 disabled:bg-slate-50 disabled:text-slate-500 sm:text-sm"
				/>
			</div>
			<div class="col-span-6 sm:col-span-3">
				<label for="first_name" class="block text-sm font-medium text-slate-700">First name</label>
				<input
					type="text"
					name="first_name"
					id="first_name"
					autocomplete="given-name"
					on:input={debouncedHandleInput}
					value={data.profile?.firstName}
					class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
				/>
				{#if formError(errors, 'first_name')}
					<p class="mt-1 text-xs text-rose-500">{formError(errors, 'first_name')}</p>
				{/if}
			</div>

			<div class="col-span-6 sm:col-span-3">
				<label for="last_name" class="block text-sm font-medium text-slate-700">Last name</label>
				<input
					type="text"
					name="last_name"
					id="last_name"
					autocomplete="family-name"
					value={data.profile?.lastName}
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

-->
