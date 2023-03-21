<script lang="ts">
	import { slide, draw, fade } from 'svelte/transition';
	import type { PageData } from './$types';

	import ConnectGoogleAccountButton from '$lib/components/ConnectGoogleAccountButton.svelte';
	import GmailIntegration from './GmailIntegration.svelte';

	// server page data
	export let data: PageData;
	$: ({ supabase, profile, session, gmailAccounts } = data);

	// ui state
	$: hasGmailAccount = gmailAccounts.length > 0;
	let profileSaved = false;
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

		const { data: profileData, error } = await supabase
			.from('recruiter')
			.update({ [name]: value })
			.eq('user_id', session?.user.id)
			.select()
			.maybeSingle();
		if (!error && profileData && profileData[name as keyof typeof profileData] === value) {
			profileSaved = true;
			setTimeout(() => {
				profileSaved = false;
			}, savedMessageTimeout);
			return;
		}
		errors[name] = 'There was an error saving your changes';
	};

	// debounce input to limit database writes
	const debouncedHandleInput = debounce(handleInput, debounceDelay);

	const onConnect = async (email?: string) => {
		const resp = await fetch('/api/account/gmail/subscribe', {
			method: 'POST',
			body: JSON.stringify({ email })
		});
		// handle errors
		if (resp.status !== 200) {
			// TODO: Show error!
			errors['activate'] = 'There was an error activating your email. Please try again.';
			return;
		}
	};
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
					value={profile?.email}
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
					value={profile?.firstName}
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
					value={profile?.lastName}
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
	{#if hasGmailAccount}
<div>
	<h2 class="text-xl font-medium leading-6 text-slate-900 sm:text-2xl">Email Integration</h2>
	<p class="mt-1 text-sm text-slate-500">
		Connect the emails you use for candidate outreach. Once connceted, SRC will automatically import
		and sync candidates you reach out to.
	</p>
</div>
{/if}
<!-- 
Two States 
1. No Email Integration -> Explain what it is and show button to connect
2. Email Integration -> Show connected email and button to disconnect + option to connect another email
	- Show connected email 
	- Show last synced date?
	- Show button to disconnect
	- option to add another
-->
<div class="space-y-6">
	{#if hasGmailAccount}
		{#each gmailAccounts as account}
			{@const { email, is_valid } = account}
			{@const settings = profile.emailSettings[email] || {}}
			<GmailIntegration isValid={is_valid} {email} {settings} />
		{/each}
	{:else}
		<div class="bg-white py-6 px-4 shadow sm:overflow-hidden sm:rounded-md sm:p-6">
			<h3 class="text-lg leading-6 font-medium text-slate-900">Setup Your Account</h3>
			<p class="my-2 text-sm">
				We're excited to have you on board! To get started, we'll need to connect your Gmail
				account. Make sure you connect the same account you use for candidate outreach.
			</p>
			<div class="text-sm">
				<p>SRC requires access to your Gmail account to:</p>
				<ul class="mt-4 list-inside space-y-2">
					<li class="flex flex-row items-center">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="currentColor"
							class="h-6 w-6 text-green-600"
						>
							<path
								fill-rule="evenodd"
								d="M19.916 4.626a.75.75 0 01.208 1.04l-9 13.5a.75.75 0 01-1.154.114l-6-6a.75.75 0 011.06-1.06l5.353 5.353 8.493-12.739a.75.75 0 011.04-.208z"
								clip-rule="evenodd"
							/>
						</svg>

						<span class="ml-2">Continuously sync and import candidates you reach out</span>
					</li>
					<li class="flex flex-row items-center">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="currentColor"
							class="h-6 w-6 text-green-600"
						>
							<path
								fill-rule="evenodd"
								d="M19.916 4.626a.75.75 0 01.208 1.04l-9 13.5a.75.75 0 01-1.154.114l-6-6a.75.75 0 011.06-1.06l5.353 5.353 8.493-12.739a.75.75 0 011.04-.208z"
								clip-rule="evenodd"
							/>
						</svg>

						<span class="ml-2"
							>Create and manage
							<span
								class="inline-flex items-center rounded-md bg-blue-500 px-1.5 py-0.5 text-sm font-medium text-white"
								>@SRC</span
							>
							Gmail labels</span
						>
					</li>
					<li class="flex flex-row items-center">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="currentColor"
							class="h-6 w-6 text-green-600"
						>
							<path
								fill-rule="evenodd"
								d="M19.916 4.626a.75.75 0 01.208 1.04l-9 13.5a.75.75 0 01-1.154.114l-6-6a.75.75 0 011.06-1.06l5.353 5.353 8.493-12.739a.75.75 0 011.04-.208z"
								clip-rule="evenodd"
							/>
						</svg>
						<span class="ml-2">Detect and match email sequences to open roles</span>
					</li>
				</ul>
			</div>
			<p class="my-4 text-sm">
				Once connected, SRC will trigger a one-time historic sync to import the last 3 months of
				candidates you've reach out to!
			</p>
			<ConnectGoogleAccountButton {onConnect} />
			{#if formError(errors, 'activate')}
				<p class="mt-2 text-xs text-rose-500">
					{formError(errors, 'activate')}
					<br />
					<span>
						If the error persists, please reach out to <a
							href="mailto:team@sharedrecruiting.co?subject=Error Connceting Gmail"
							class="underline">team@sharedrecruiting.co</a
						>
					</span>
				</p>
			{/if}
		</div>
	{/if}
</div>
