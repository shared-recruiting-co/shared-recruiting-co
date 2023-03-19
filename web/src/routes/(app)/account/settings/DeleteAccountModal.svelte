<script lang="ts">
	import { page } from '$app/stores';
	import { fly, fade } from 'svelte/transition';
	import { goto } from '$app/navigation';

	$: ({ supabase } = $page.data);

	export let show: boolean = false;
	let reason: string;

	let error: string = '';
	let loading: boolean = false;

	const onConfirm = async () => {
		if (!reason || !reason.trim()) {
			error = 'Please enter a reason for deleting your account.';
			return;
		}
		reason = reason.trim();
		// clear error
		error = '';
		loading = true;

		// wait deletion
		const resp = await fetch('/api/account', {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				reason
			})
		});

		loading = false;
		if (resp.ok) {
			try {
				await supabase.auth.signOut();
			} catch {}
			goto('/join');
		} else {
			// show error
			const { message } = await resp.json();
			error = message;
		}
	};

	const close = () => {
		show = false;
	};
</script>

{#if show}
	<div
		class="relative z-50"
		aria-labelledby="modal-title"
		role="dialog"
		aria-modal="true"
		in:fade={{ duration: 200 }}
		out:fade={{ delay: 150, duration: 200 }}
	>
		<div
			class="fixed inset-0 bg-slate-500 bg-opacity-75 transition-opacity"
			class:ease-out={show}
			class:duration-300={show}
			class:opacity-0={!show}
			class:opacity-100={show}
			class:ease-in={!show}
			class:duration-200={!show}
		/>

		<div class="fixed inset-0 z-10 overflow-y-auto">
			<div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
				<div
					class="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6"
					class:ease-out={show}
					class:duration-300={show}
					class:opacity-0={!show}
					class:opacity-100={show}
					class:translate-y-4={!show}
					class:translate-y-0={show}
					class:scale-95={!show}
					class:scale-100={show}
					class:ease-in={!show}
					class:duration-200={!show}
					in:fly={{ y: -200, duration: 200 }}
				>
					<div class="sm:flex sm:items-start">
						<div
							class="mx-auto flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10"
						>
							<!-- Heroicon name: outline/exclamation-triangle -->
							<svg
								class="h-6 w-6 text-red-600"
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
									d="M12 10.5v3.75m-9.303 3.376C1.83 19.126 2.914 21 4.645 21h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 4.88c-.866-1.501-3.032-1.501-3.898 0L2.697 17.626zM12 17.25h.007v.008H12v-.008z"
								/>
							</svg>
						</div>
						<div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
							<h3 class="text-lg font-medium leading-6 text-slate-900" id="modal-title">
								Confirm Account Deletion
							</h3>
							<div class="mt-2">
								<p class="text-sm text-slate-500">
									Permenantly delete your SRC account. This action cannot be undone.
								</p>
							</div>
							<div class="mt-4">
								<p class="text-left text-xs text-slate-500">
									Before you delete your account, please help us improve SRC and share a few words
									on why you are deleting your account.
								</p>
								<textarea
									class="mt-2 block w-full rounded-md border-slate-300 text-sm shadow-sm"
									class:ring-red-600={error}
									class:ring-1={error}
									rows="3"
									bind:value={reason}
									placeholder="I'm leaving because..."
								/>
								{#if error}
									<p class="mt-2 text-xs text-red-600">{error}</p>
								{/if}
							</div>
						</div>
					</div>
					<div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
						<button
							type="button"
							class="inline-flex w-full justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm"
							disabled={loading}
							class:animate-pulse={loading}
							on:click={onConfirm}
						>
							{#if loading}
								Deleting...
							{:else}
								Delete Account
							{/if}
						</button>
						<button
							type="button"
							class="mt-3 inline-flex w-full justify-center rounded-md border border-slate-300 bg-white px-4 py-2 text-base font-medium text-slate-700 shadow-sm hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:w-auto sm:text-sm"
							on:click={close}>Cancel</button
						>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
