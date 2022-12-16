<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import ConnectGoogleAccountButton from '$lib/components/ConnectGoogleAccountButton.svelte';

	const onConnect = () => {
		goto('/account/profile');
	};
	// TODOs
	//
	// Design
	// Create separate layouts for (marketing) and (account) pages
	// Create account header
	// Create account sidebar
	// Design & create account page
	// Active Toggle on settings page
</script>

<div class="m-2 flex flex-row items-center justify-center sm:mx-auto sm:my-16">
	<div class="flex max-w-xl flex-col space-y-8 rounded-md p-8">
		<h1 class="text-center text-3xl">Welcome to the SRC {$page.data.profile.firstName}!</h1>
		<p class="mt-2">
			We're excited to have you on board! Before we can create your account, we need to connect to
			your Google mail account.
		</p>
		<div>
			<p>SRC requires access to your Google mail account to:</p>
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

					<span class="ml-2">Detect new and historic job opportunities</span>
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

					<span class="ml-2">Block automated recruiting email follow-ups</span>
				</li>
			</ul>
		</div>
		<p>
			Once connected, SRC will trigger a one-time historic sync to label your existing inbound job
			opportunities. Aftewards, SRC will automatically detect and label new job opportunities as
			soon as they hit your inbox!
		</p>
		<div>
			<p class="mb-2">Use the button below to connect your account</p>
			<ConnectGoogleAccountButton {onConnect} email={$page.data.session?.user?.email} />
		</div>
		<p class="text-sm">
			You can read more about how we use and protect your data in our <a
				href="/legal/privacy-policy"
				target="_blank"
				class="underline">privacy policy</a
			>. If you want an even deeper dive into how we use your data, you can read the
			<a
				href="https://github.com/shared-recruiting-co/shared-recruiting-co"
				class="underline"
				target="_blank"
				rel="noopener noreferrer">SRC code</a
			>.
		</p>
	</div>
</div>
