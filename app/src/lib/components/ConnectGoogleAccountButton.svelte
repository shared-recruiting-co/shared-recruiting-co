<script lang="ts">
	import { browser } from '$app/environment';
	import { PUBLIC_GOOGLE_CLIENT_ID } from '$env/static/public';

	// props
	export let email: string | undefined;
	export let onConnect: () => void;

	let loaded = Boolean(browser && window.google);
	let error: string;
	let connecting = false;

	const onLoad = () => {
		loaded = true;
	};

	const codeCallback = (response: google.accounts.oauth2.CodeResponse) => {
		connecting = true;
		const xhr = new XMLHttpRequest();
		xhr.open('POST', '/api/account/gmail/connect', true);
		xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
		xhr.setRequestHeader('X-Requested-With', 'XmlHttpRequest');
		// Set custom header for CRSF
		xhr.withCredentials = true;
		xhr.onload = function () {
			connecting = false;
			// check for error status
			if (xhr.status === 200) {
				// success
				onConnect();
			} else {
				const body = JSON.parse(xhr.responseText);
				error = body.message;
			}
		};
		xhr.onerror = function (err) {
			connecting = false;
			// log error
			console.log('onerror', err);
		};
		xhr.send(`code=${response.code}&scope=${response.scope}`);
	};

	const connectAccount = () => {
		if (!loaded) return;
		// reset error
		error = '';

		const client = window.google.accounts.oauth2.initCodeClient({
			client_id: PUBLIC_GOOGLE_CLIENT_ID,
			scope: 'https://www.googleapis.com/auth/gmail.modify',
			// google does magic redirection for us on mobile
			ux_mode: 'popup',
			redirect_uri: 'postmessage',
			hint: email,
			error_callback: (err: unknown) => {
				if (!err) return;
				// ignore popup closed errors
				if (err?.type === 'popup_closed') return;
				error = err?.message;
			},
			callback: codeCallback
		});

		client.requestCode();
	};
</script>

<button
	disabled={!loaded}
	class="flex items-center rounded-md bg-[#1a73e8] py-0.5 pl-0.5 pr-3 text-white shadow hover:bg-[#5a94ee] disabled:cursor-wait"
	on:click|preventDefault={connectAccount}
	><img src="/google.svg" alt="Google logo" class="mr-3 rounded-l-md bg-white p-2" />
	{#if connecting}
		<span class="animate-pulse">Connecting...</span>
	{:else}
		<span>Connect Gmail</span>
	{/if}
</button>
{#if error}
	<p class="mt-2 text-sm text-rose-500">{error}</p>
{/if}
<svelte:head>
	<script src="/scripts/gsi.js" on:load={onLoad} async defer></script>
</svelte:head>
