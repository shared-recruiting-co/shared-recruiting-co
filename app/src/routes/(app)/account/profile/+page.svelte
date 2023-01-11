<script lang="ts">
	import { page } from '$app/stores';
	import { slide, draw, fade } from 'svelte/transition';

	import { supabaseClient, UserEmailStats } from '$lib/supabase/client';

	import Toggle from '$lib/components/Toggle.svelte';
	import AlertModal from '$lib/components/AlertModal.svelte';
	import ConnectGoogleAccountButton from '$lib/components/ConnectGoogleAccountButton.svelte';

	// ui state
	let profileSaved = false;
	let settingsSaved = false;
	let errors: Record<string, string> = {};
	// settings
	let isActive = $page.data.profile.isActive;
	let isSetup = $page.data.isSetup;
	let autoContribute = $page.data.profile.autoContribute;
	let autoArchive = $page.data.profile.autoArchive;
	let showDeactivateEmailModal = false;
	// stats
	let lastSyncedAt = $page.data.lastSyncedAt;
	let numEmailsProcessed = $page.data.numEmailsProcessed;
	let numJobsDetected = $page.data.numJobsDetected;

	supabaseClient
		.channel('table-db-changes')
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
				if (changed.stat_id === UserEmailStats.EmailsProcessed) {
					numEmailsProcessed = changed.stat_value;
				} else if (changed.stat_id === UserEmailStats.JobsDetected) {
					numJobsDetected = changed.stat_value;
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

	const onConnect = () => {
		isSetup = true;
		isActive = true;
	};

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

	const saveSettings = async () => {
		settingsSaved = false;
		const { data, error } = await supabaseClient
			.from('user_profile')
			.update({ auto_contribute: autoContribute, auto_archive: autoArchive })
			.eq('user_id', $page.data.session?.user.id)
			.select()
			.maybeSingle();
		if (!error && data) {
			settingsSaved = true;
			setTimeout(() => {
				settingsSaved = false;
			}, savedMessageTimeout);
			return;
		}
		errors['settings'] = 'There was an error saving your changes. Please try again.';
	};

	const debouncedSaveSettings = debounce(saveSettings, debounceDelay);
	let onSettingsToggle: (checked: boolean) => void;

	$: {
		// keep onSettingsToggle in sync with settings values
		onSettingsToggle = (_checked: boolean) => {
			debouncedSaveSettings();
		};
	}

	const onDeactivateConfirm = async () => {
		showDeactivateEmailModal = false;
		const resp = await fetch('/api/account/gmail/unsubscribe', { method: 'POST' });
		// handle errors
		if (resp.status !== 200) {
			errors['deactivate'] =
				'There was an error deactivating the inbox assistant. Please try again.';
			return;
		}
		isActive = false;
	};

	const activateEmail = async () => {
		const resp = await fetch('/api/account/gmail/subscribe', { method: 'POST' });
		// handle errors
		if (resp.status !== 200) {
			errors['activate'] = 'There was an error activating your inbox assistant. Please try again.';
			return;
		}
		isActive = true;
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
		Member since {new Date($page.data.profile.createdAt).toLocaleDateString()}
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
				<dt class="truncate text-sm font-medium text-gray-500">Emails Analyzed</dt>
				<dd class="mt-1 text-2xl tracking-tight text-gray-900">
					{numberFormatter.format(numEmailsProcessed)}
				</dd>
			</div>

			<div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
				<dt class="truncate text-sm font-medium text-gray-500">Jobs Identified</dt>
				<dd class="mt-1 text-2xl tracking-tight text-gray-900">
					{numberFormatter.format(numJobsDetected)}
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
					value={$page.data.profile.email}
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
					value={$page.data.profile.firstName}
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
{#if isActive && isSetup}
	<div class="relative shadow sm:overflow-hidden sm:rounded-md">
		{#if formError(errors, 'settings')}
			<div class="absolute top-0 right-0 mt-6 mr-8 flex items-center space-x-2 text-green-600">
				<p class="mt-1 text-xs text-rose-500">{formError(errors, 'settings')}</p>
			</div>
		{/if}
		{#if settingsSaved}
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
				<h3 class="text-lg font-medium leading-6 text-slate-900">Email Settings</h3>
				<p class="mt-1 text-sm text-slate-500">Manage your email settings and preferences.</p>
			</div>
			<ul class="mt-2 divide-y divide-slate-200">
				<li class="flex items-center justify-between py-4">
					<div class="flex flex-col pr-4 sm:pr-8">
						<p class="text-sm font-medium text-slate-900" id="privacy-option-2-label">
							Hide Recruiting Emails from Inbox
						</p>
						<p class="mt-1 text-xs text-slate-500 sm:text-sm" id="privacy-option-2-description">
							Keep your inbox distraction free when you aren't actively looking for a new role.
							<br />
							Recruiting emails are always accessible under the @SRC folder.
						</p>
					</div>
					<Toggle
						bind:checked={autoArchive}
						label="Recruiting Assistant"
						onToggle={onSettingsToggle}
					/>
				</li>
				<li class="flex items-center justify-between py-4">
					<div class="flex flex-col pr-4 sm:pr-8">
						<p class="text-sm font-medium text-slate-900" id="privacy-option-1-label">
							Auto-Contribute Recruiting Emails
						</p>
						<p class="mt-1 text-xs text-slate-500 sm:text-sm" id="privacy-option-1-description">
							Help us build the best product for candidates by automatically contributing your
							inbound recruiting emails to our recruiting email dataset.
							<br />
							You can always manually contribute emails by forwarding them to
							<a class="underline" href="mailto:examples@sharedrecruiting.co"
								>examples@sharedrecruiting.co</a
							>.
						</p>
					</div>
					<Toggle
						bind:checked={autoContribute}
						label="Auto Contribute"
						onToggle={onSettingsToggle}
					/>
				</li>
				<li class="flex items-center justify-between py-4">
					<div class="flex flex-col pr-4 sm:pr-8">
						<p class="text-sm font-medium text-slate-900" id="privacy-option-3-label">
							Block Automated Email Sequences&NonBreakingSpace;
							<span class="text-slate-500">(Coming Soon)</span>
						</p>
						<p class="mt-1 text-xs text-slate-500 sm:text-sm" id="privacy-option-3-description">
							Block automated recruiting sequences by automatically replying to recruiters with a
							standard message.
						</p>
					</div>
					<Toggle checked={false} disabled label="Auto-Reply" />
				</li>
			</ul>
		</div>
	</div>
	<div class="shadow sm:overflow-hidden sm:rounded-md">
		<div class="space-y-6 bg-white py-6 px-4 sm:p-6">
			<div class="max-w-2xl">
				<h3 class="text-lg font-medium leading-6 text-slate-900">Deactivate Inbox Assistant</h3>
				<p class="mt-1 text-sm text-slate-500">
					Deactivate the SRC inbox assistant. While disabled, recruiting emails will no longer be
					automatically labeled or managed for you. This will not delete your account nor any data.
					The @SRC labels will remain in your inbox. You can reactivate SRC at anytime.
				</p>
			</div>
			<button
				type="button"
				class="inline-flex w-full justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 sm:w-auto sm:text-sm"
				on:click={() => {
					showDeactivateEmailModal = true;
				}}>Deactivate</button
			>
			{#if formError(errors, 'deactivate')}
				<p class="mt-1 text-xs text-rose-500">
					{formError(errors, 'deactivate')}
					<br />
					<span>
						If the error persists, please reach out to <a
							href="mailto:team@sharedrecruiting.co?subject=Error Deactivating Inbox Assistant"
							class="underline">team@sharedrecruiting.co</a
						>
					</span>
				</p>
			{/if}
		</div>
		<AlertModal
			bind:show={showDeactivateEmailModal}
			title="Deactive Inbox Assistant?"
			description="Are you sure you want to deactivate the SRC inbox assistant? While disabled, recruiting emails will no longer be automatically labeled or managed for you."
			cta="Deactivate"
			onConfirm={onDeactivateConfirm}
		/>
	</div>
{:else if !isSetup}
	<div class="shadow sm:overflow-hidden sm:rounded-md">
		<div class="space-y-6 bg-sky-50 py-6 px-4 sm:p-6">
			<div class="max-w-2xl">
				<h3 class="text-lg font-medium leading-6 text-slate-900">Connect Gmail Account</h3>
				<p class="mt-1 text-sm text-slate-500">
					We lost connection to your Gmail account. Please reconnect to continue using SRC. Once
					re-enabled, SRC will re-sync your inbox between now and the last time the connection was
					active.
				</p>
			</div>
			<ConnectGoogleAccountButton {onConnect} email={$page.data.profile?.email} />
		</div>
	</div>
	<!-- else not active -->
{:else}
	<div class="shadow sm:overflow-hidden sm:rounded-md">
		<div class="space-y-6 bg-white py-6 px-4 sm:p-6">
			<div class="max-w-2xl">
				<h3 class="text-lg font-medium leading-6 text-slate-900">Activate Inbox Assistant</h3>
				<p class="mt-1 text-sm text-slate-500">
					Your SRC Inbox Assistant is currently disabled. Re-enable it to start monitoring your
					inbox for job opportunities. Once re-enabled, SRC will re-sync your inbox between now and
					the last time SRC was active.
				</p>
			</div>
			<button
				type="button"
				class="inline-flex w-full justify-center rounded-md border border-transparent bg-green-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 sm:w-auto sm:text-sm"
				on:click={activateEmail}>Activate</button
			>
			{#if formError(errors, 'activate')}
				<p class="mt-1 text-xs text-rose-500">
					{formError(errors, 'activate')}
					<br />
					<span>
						If the error persists, please reach out to <a
							href="mailto:team@sharedrecruiting.co?subject=Error Activating Inbox Assistant"
							class="underline">team@sharedrecruiting.co</a
						>
					</span>
				</p>
			{/if}
		</div>
	</div>
{/if}
