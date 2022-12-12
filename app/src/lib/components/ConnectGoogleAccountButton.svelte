<script lang="ts">
	import { browser } from '$app/environment';
	import { PUBLIC_GOOGLE_CLIENT_ID, PUBLIC_GOOGLE_REDIRECT_URI } from '$env/static/public';

	type OAuth2CallbackResponse = {
		code: string;
		authuser: string;
		prompt: string;
		scope: string;
	} 

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
				xhr.withCredentials = true;
				xhr.onload = function () {
					// check for error status
					if (xhr.status === 200) {
						// success
						console.log('success');
						// TODO: redirect
					} else {
						console.log('error');
						// log error
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

		const client = window.google.accounts.oauth2.initCodeClient({
			client_id: PUBLIC_GOOGLE_CLIENT_ID,
			scope: 'https://www.googleapis.com/auth/gmail.modify',
			// what about mobile??
			// ux_mode: 'popup',
			ux_mode: 'redirect',
			redirect_uri: PUBLIC_GOOGLE_REDIRECT_URI,
			access_type: 'offline',
			auto_select: true,
			approval_prompt: 'auto',
			include_granted_scopes: true,
			immediate: true,
			error_callback: (err) => {
				// type === 'popup_closed' display appropriate error
				error = err.message;
			},
			callback: codeCallback
		});

		client.requestCode();
	};
</script>

<button
	disabled={!loaded}
	class="flex items-center rounded-md bg-white px-2 py-2 text-slate-500 shadow hover:bg-slate-50 disabled:cursor-wait"
	on:click|preventDefault={connectAccount}
	><img src="/google.svg" alt="Google logo" class="mr-3" />Connect Account</button
>
<svelte:head>
	<script src="/scripts/gsi.js" on:load={onLoad} async defer></script>
</svelte:head>
