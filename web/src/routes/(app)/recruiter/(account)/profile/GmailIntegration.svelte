<script lang="ts">
	import { slide } from 'svelte/transition';

	import AlertModal from '$lib/components/AlertModal.svelte';
	import ConnectGoogleAccountButton from '$lib/components/ConnectGoogleAccountButton.svelte';

	export let isValid: boolean;
	export let email: string;
	export let settings: Record<string, boolean>;
	// ui state
	// default to closed when valid, open when invalid
	export let isOpen = !isValid;
	let isActive = settings['is_active'] || false;
	let showDeactivateEmailModal = false;
	let errors: Record<string, string> = {};
	let isActivating = false;

	const toggle = () => {
		isOpen = !isOpen;
	};

	const formError = (e: typeof errors, field: string) => {
		return e[field] || '';
	};

	// TODO: Add last synced time
	// TODO: Send welcome email on new email, send welcome back on reactivation
	const onDeactivateConfirm = async () => {
		showDeactivateEmailModal = false;
		const resp = await fetch('/api/account/gmail/unsubscribe', {
			method: 'POST',
			body: JSON.stringify({ email })
		});
		// handle errors
		if (resp.status !== 200) {
			errors['deactivate'] = `There was an error deactivating ${email}. Please try again.`;
			return;
		}
		isActive = false;
		isOpen = false;
	};

	const activateEmail = async () => {
		isActivating = true;
		const resp = await fetch('/api/account/gmail/subscribe', {
			method: 'POST',
			body: JSON.stringify({ email })
		});
		isActivating = false;
		// handle errors
		if (resp.status !== 200) {
			errors['activate'] = 'There was an error activating your inbox assistant. Please try again.';
			return;
		}
		isActive = true;
		isOpen = false;
	};

	const onConnect = async () => {
		await activateEmail();
		isValid = true;
		isOpen = false;
	};
</script>

<div class="bg-white py-6 px-4 shadow sm:overflow-hidden sm:rounded-md sm:p-6">
	<button
		class="flex w-full justify-between"
		on:click={(e) => {
			e.preventDefault();
			toggle();
		}}
	>
		<div class="flex items-center space-x-4">
			<img src="/gmail.svg" alt="Gmail" class="h-6 w-6" />
			<span class="text-slate-600">{email}</span>
		</div>
		<div class="flex items-center space-x-4">
			<!-- Status Badge -->
			{#if isValid && isActive}
				<span
					class="inline-flex items-center rounded-md bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800 md:px-3 md:py-1 md:text-sm"
				>
					Active
				</span>
			{:else if isValid && !isActive}
				<span
					class="inline-flex items-center rounded-md bg-amber-100 px-2.5 py-0.5 text-xs font-medium text-amber-800 md:px-3 md:py-1 md:text-sm"
				>
					Paused
				</span>
			{:else}
				<span
					class="inline-flex items-center rounded-md bg-red-100 px-2.5 py-0.5 text-xs font-medium text-red-800 md:px-3 md:py-1 md:text-sm"
				>
					Connection Lost
				</span>
			{/if}
			<span class="ml-6 flex h-7 items-center">
				<!--
						Expand/collapse icon, toggle classes based on question open state.
						Heroicon name: outline/chevron-down
						Open: "-rotate-180", Closed: "rotate-0"
					-->
				<svg
					class="h-6 w-6 transform text-slate-600"
					class:-rotate-180={isOpen}
					class:rotate-0={!isOpen}
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					aria-hidden="true"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
				</svg>
			</span>
		</div>
	</button>
	{#if isOpen}
		<div transition:slide class="mt-4 border-t pt-4">
			{#if isValid && isActive}
				<!-- Content -->
				<div class="max-w-2xl">
					<h3 class="text-lg font-medium leading-6 text-slate-900">Deactivate</h3>
					<p class="mt-2 mb-4 text-sm text-slate-500">
						Deactivate the SRC inbox assistant. While disabled, new outbound emails will no longer
						be automatically imported and synced with SRC. This will not delete your account nor any
						data. The @SRC labels will remain in your inbox. You can reactivate SRC at anytime.
					</p>
				</div>
				<button
					type="button"
					class="inline-flex w-full justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 sm:w-auto sm:text-sm"
					on:click={() => {
						showDeactivateEmailModal = true;
					}}>Deactivate</button
				>
				{#if formError(errors, 'deactivate')}
					<p class="mt-2 text-xs text-rose-500">
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
				<AlertModal
					bind:show={showDeactivateEmailModal}
					title="Deactivate Inbox Assistant?"
					description="Deactivate the SRC inbox assistant. While disabled, new outbound emails will no longer be automatically imported and synced with SRC. This will not delete your account nor any data. The @SRC labels will remain in your inbox. You can reactivate SRC at anytime."
					cta="Deactivate"
					onConfirm={onDeactivateConfirm}
				/>
			{:else if isValid && !isActive}
				<div class="max-w-2xl">
					<h3 class="text-lg font-medium leading-6 text-slate-900">Activate</h3>
					<p class="mt-2 mb-4 text-sm text-slate-500">
						Your SRC integration is currently disabled. Re-enable it to start automatically syncing
						your outbound emails with SRC. Once re-enabled, SRC will re-sync your inbox between now
						and the last time SRC was active.
					</p>
				</div>
				<button
					type="button"
					class="inline-flex w-full justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 sm:w-auto sm:text-sm"
					class:animate-pulse={isActivating}
					on:click={activateEmail}
				>
					{isActivating ? 'Activating...' : 'Activate'}
				</button>
				{#if formError(errors, 'activate')}
					<p class="mt-2 text-xs text-rose-500">
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
			{:else}
				<div class="mb-4 max-w-2xl">
					<h3 class="text-lg font-medium leading-6 text-slate-900">Connection Lost</h3>
					<p class="mt-2 text-sm text-slate-500">
						We lost connection to your Gmail account. Please reconnect to continue using SRC. Once
						re-enabled, SRC will re-sync your inbox between now and the last time the connection was
						active.
					</p>
				</div>
				<ConnectGoogleAccountButton {onConnect} {email} />
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
			{/if}
		</div>
	{/if}
</div>
