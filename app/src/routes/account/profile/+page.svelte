<script lang="ts">
	import { page } from '$app/stores';
	import { slide, draw, fade } from 'svelte/transition';

	import { supabaseClient } from '$lib/supabase/client';

	import Toggle from '$lib/components/Toggle.svelte';
	import AlertModal from '$lib/components/AlertModal.svelte';

	let profileSaved = false;
	let settingsSaved = false;
	let errors: Record<string, string> = {};
	let autoContribute = $page.data.profile.autoContribute;
	let autoArchive = $page.data.profile.autoArchive;
	let showDeactivateEmailModal = false;

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
		errors['settings'] = 'There was an error saving your changes';
	};

	const debouncedSaveSettings = debounce(saveSettings, debounceDelay);
	let onSettingsToggle: (checked: boolean) => void;

	$: {
		// keep onSettingsToggle in sync with settings values
		onSettingsToggle = (_checked: boolean) => {
			debouncedSaveSettings();
		};
	}
</script>

<div class="sm:my-18 my-12 lg:grid lg:grid-cols-12 lg:gap-x-5">
	<!-- Empty space for now -->
	<aside class="block py-6 px-2 sm:px-6 lg:col-span-2 lg:py-0 lg:px-0" />
	<div class="space-y-6 sm:px-6 lg:col-span-9 lg:px-0">
		<div class="px-4">
			<h1 class="text-3xl sm:text-4xl">Account</h1>
			<p class="mt-1 text-sm text-slate-500">
				Member since {new Date($page.data.profile.createdAt).toLocaleDateString()}
			</p>
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
						<label for="last_name" class="block text-sm font-medium text-slate-700">Last name</label
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
		<div class="relative shadow sm:overflow-hidden sm:rounded-md">
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
					<h3 class="text-lg font-medium leading-6 text-slate-900">Settings</h3>
					<p class="mt-1 text-sm text-slate-500">Manage your account settings and preferences.</p>
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
						automatically labeled or managed for you. This will not delete your account nor any
						data. The @SRC labels will remain in your inbox. You can reactivate SRC at anytime.
					</p>
				</div>
				<button
					type="button"
					class="inline-flex w-full justify-center rounded-md border border-transparent bg-rose-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-rose-700 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:ring-offset-2 sm:w-auto sm:text-sm"
					on:click={() => {
						showDeactivateEmailModal = true;
					}}>Deactivate</button
				>
			</div>
			<AlertModal
				bind:show={showDeactivateEmailModal}
				title="Deactive Inbox Assistant?"
				description="Are you sure you want to deactivate the SRC inbox assistant? While disabled, recruiting emails will no longer be automatically labeled or managed for you."
				cta="Deactivate"
				onConfirm={() => {
					showDeactivateEmailModal = false;
				}}
			/>
		</div>
	</div>
</div>
