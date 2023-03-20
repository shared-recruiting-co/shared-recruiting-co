<script lang="ts">
	import { slide } from 'svelte/transition';

	import AlertModal from '$lib/components/AlertModal.svelte';
	import ConnectGoogleAccountButton from '$lib/components/ConnectGoogleAccountButton.svelte';

	export let isValid: boolean;
	export let email: string;
	// ui state
	// default to closed when valid, open when invalid
	export let isOpen = !isValid;
	let showDeactivateEmailModal = false;
	let errors: Record<string, string> = {};

	const toggle = () => {
		isOpen = !isOpen;
	};

	const formError = (e: typeof errors, field: string) => {
		return e[field] || '';
	};

	// TODO: Add state for
	// TODO: Implement! Add support for recruiter_email_settings
	// - user_id
	// - email
	// - is_active
	const onDeactivateConfirm = async () => {
		showDeactivateEmailModal = false;
		const resp = await fetch('/api/account/gmail/unsubscribe', { method: 'POST' });
		// handle errors
		if (resp.status !== 200) {
			errors['deactivate'] = `There was an error deactivating ${email}. Please try again.`;
			return;
		}
		// isActive = false;
	};

	const onConnect = () => {
		isValid = true;
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
			<!-- Show a green dot if the account is valid, otherwise red-->
			{#if isValid}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="currentColor"
					class="h-6 w-6 text-emerald-600"
				>
					<path
						fill-rule="evenodd"
						d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z"
						clip-rule="evenodd"
					/>
				</svg>
			{:else}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="currentColor"
					class="h-6 w-6 text-rose-600"
				>
					<path
						fill-rule="evenodd"
						d="M12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75 9.75-4.365 9.75-9.75S17.385 2.25 12 2.25zm-1.72 6.97a.75.75 0 10-1.06 1.06L10.94 12l-1.72 1.72a.75.75 0 101.06 1.06L12 13.06l1.72 1.72a.75.75 0 101.06-1.06L13.06 12l1.72-1.72a.75.75 0 10-1.06-1.06L12 10.94l-1.72-1.72z"
						clip-rule="evenodd"
					/>
				</svg>
			{/if}
			<span class="ml-6 flex h-7 items-center">
				<!--
						Expand/collapse icon, toggle classes based on question open state.
						Heroicon name: outline/chevron-down
						Open: "-rotate-180", Closed: "rotate-0"
					-->
				<svg
					class="h-6 w-6 transform"
					class:-rotate-180={isOpen}
					class:rotate-0={!isOpen}
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24 text-slate-600"
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
			{#if isValid}
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
				<AlertModal
					bind:show={showDeactivateEmailModal}
					title="Deactive Inbox Assistant?"
					description="Deactivate the SRC inbox assistant. While disabled, new outbound emails will no longer be automatically imported and synced with SRC. This will not delete your account nor any data. The @SRC labels will remain in your inbox. You can reactivate SRC at anytime."
					cta="Deactivate"
					onConfirm={onDeactivateConfirm}
				/>
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
			{/if}
		</div>
	{/if}
</div>
