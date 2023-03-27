<script lang="ts">
	import { slide, draw, fade } from 'svelte/transition';
	import type { PageData } from './$types';

	import { UserEmailStats, type UserEmailSettings } from '$lib/supabase/client';
	import { debounce, formError } from '$lib/forms';
	import type { FormErrors } from '$lib/forms';

	import EmailSettings from './EmailSettings.svelte';

	export let data: PageData;
	$: ({ supabase, profile } = data);

	// ui state
	let profileSaved = false;
	let errors: FormErrors = {};
	// settings
	let isActive = data.profile.isActive;
	let isSetup = data.isSetup;
	let autoContribute = data.profile.autoContribute;
	let autoArchive = data.profile.autoArchive;
	// stats
	let lastSyncedAt = data.lastSyncedAt;
	let numInboundJobs = data.numInboundJobs;
	let numOfficialJobs = data.numOfficialJobs;

	supabase
		?.channel('table-db-changes')
		.on(
			'postgres_changes',
			{
				event: '*',
				table: 'user_email_stat',
				schema: 'public'
			},
			(payload) => {
				// rename because of keyword
				const { new: changed } = payload;
				if (!changed) return;
				// TODO: Show a cool animation when the numbers of stats change
				if (changed.stat_id === UserEmailStats.JobsDetected) {
					numInboundJobs = changed.stat_value;
				}
			}
		)
		.on(
			'postgres_changes',
			{
				event: '*',
				table: 'user_email_sync_history',
				schema: 'public'
			},
			(payload) => {
				// rename because of keyword
				const { new: changed } = payload;
				if (!changed || !changed.synced_at) return;
				lastSyncedAt = changed.synced_at;
			}
		)
		.subscribe();

	const debounceDelay = 1000;
	const savedMessageTimeout = 3000;

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
			.from('user_profile')
			.update({ [name]: value })
			.eq('user_id', data.session?.user.id)
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

	const saveEmailSettings = async (_email: string, settings: UserEmailSettings) => {
		const newAutoArchive = settings['auto_archive'] || false;
		const newAutoContribute = settings['auto_contribute'] || false;
		const { error } = await supabase
			.from('user_profile')
			.update({ auto_contribute: newAutoContribute, auto_archive: newAutoArchive })
			.eq('user_id', data.session?.user.id);
		if (error) {
			throw error;
		}
	};

	const timeFormatter = new Intl.DateTimeFormat('en', {
		timeStyle: 'short'
	});
	const numberFormatter = new Intl.NumberFormat('en');
	// written by ChatGPT
	const formatDate = (iso: string): string => {
		const date = new Date(iso);
		const today = new Date();

		if (date.toDateString() === today.toDateString()) {
			// If the date is today, return "Today at" followed by the time
			return `Today at ${timeFormatter.format(date)}`;
		} else {
			// Otherwise, return the month and day followed by "at" and the time
			const options: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' };
			return `${date.toLocaleDateString('en-US', options)} at ${timeFormatter.format(date)}`;
		}
	};
</script>

<div>
	<h1 class="text-3xl sm:text-4xl">Account</h1>
	<p class="mt-1 text-sm text-slate-500">
		Member since {new Date(data.profile.createdAt).toLocaleDateString()}
	</p>
</div>
{#if lastSyncedAt}
	<div>
		<dl class="mt-5 hidden gap-5 sm:grid sm:grid-cols-3">
			<div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
				<dt class="truncate text-sm font-medium text-gray-500">Last Synced</dt>
				<dd class="mt-1 text-2xl tracking-tight text-gray-900">
					{formatDate(lastSyncedAt)}
				</dd>
			</div>

			<div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
				<dt class="truncate text-sm font-medium text-gray-500">Inbound Jobs</dt>
				<dd class="mt-1 text-2xl tracking-tight text-gray-900">
					{numberFormatter.format(numInboundJobs)}
				</dd>
			</div>
			<div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
				<dt class="truncate text-sm font-medium text-gray-500">Verified Jobs</dt>
				<dd class="mt-1 text-2xl tracking-tight text-gray-900">
					{numberFormatter.format(numOfficialJobs)}
				</dd>
			</div>
		</dl>
	</div>
{/if}
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
					value={data.profile.email}
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
					value={data.profile.firstName}
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
					value={data.profile.lastName}
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
<div>
	<h2 class="text-xl font-medium leading-6 text-slate-900 sm:text-2xl">Connected Accounts</h2>
	<p class="mt-2 text-sm text-slate-500">
		Connect and configure your email accounts to get the most of SRC.
	</p>
</div>
<div class="space-y-6">
	<EmailSettings
		isValid={isSetup}
		email={profile.email}
		saveSettings={saveEmailSettings}
		settings={{
			is_active: isActive,
			auto_archive: autoArchive,
			auto_contribute: autoContribute
		}}
	/>
</div>
