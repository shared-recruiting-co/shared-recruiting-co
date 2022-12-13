<script lang="ts">
	import { browser } from '$app/environment';
	import { 
		PUBLIC_GOOGLE_CLIENT_ID, 
		PUBLIC_GOOGLE_REDIRECT_URI,
	} from '$env/static/public';

	type OAuth2CallbackResponse = {
		code: string;
		authuser: string;
		prompt: string;
		scope: string;
	} 

	// props
	export let email: string;
	export let onConnect: () => void;

	let loaded = Boolean(browser && window.google);
	let error: string;

	const onLoad = () => {
		loaded = true;
	}

	const codeCallback = (response: OAuth2CallbackResponse) => {
				const xhr = new XMLHttpRequest();
				xhr.open('POST', '/account/profile/connect', true);
				xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
				xhr.setRequestHeader('X-Requested-With', 'XmlHttpRequest');
				// Set custom header for CRSF
				xhr.withCredentials = true;
				xhr.onload = function () {
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
					// log error
					console.log("onerror",err)
				};
				xhr.send(`code=${response.code}&authuser=${response.authuser}&prompt=${response.prompt}&scope=${response.scope}`);
			}

	const connectAccount = () => {
		if (!loaded) return;
		// reset error
		error = "";

		const client = window.google.accounts.oauth2.initCodeClient({
			client_id: PUBLIC_GOOGLE_CLIENT_ID,
			scope: 'https://www.googleapis.com/auth/gmail.modify',
			// TODO: On mobile, default to redirect
			ux_mode: 'popup',
			redirect_uri: 'postmessage',
			// ux_mode: 'redirect',
			// redirect_uri: PUBLIC_GOOGLE_REDIRECT_URI, access_type: 'offline',
			login_hint: email,
			auto_select: true,
			approval_prompt: 'auto',
			include_granted_scopes: true,
			immediate: true,
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
	class="flex items-center rounded-md bg-[#1a73e8] pl-0.5 py-0.5 pr-3 text-white shadow hover:bg-[#5a94ee] disabled:cursor-wait"
	on:click|preventDefault={connectAccount}
	><img src="/google.svg" alt="Google logo" class="mr-3 p-2 bg-white rounded-l-md" />Connect with Google</button
>
{#if error}
	<p class="mt-2 text-rose-500 text-sm">{error}</p>
{/if}
<svelte:head>
	<script src="/scripts/gsi.js" on:load={onLoad} async defer></script>
</svelte:head>
