<script lang="ts" context="module">
	// on mount
	import { browser } from '$app/environment';
	import { PUBLIC_GOOGLE_CLIENT_ID } from '$env/static/public';

	type OAuth2CallbackResponse = {
		code: string;
		authuser: string;
		prompt: string;
		scope: string;
	} 

	const connectAccount = () => {
		if (!browser || !window.google) return;

		const client = window.google.accounts.oauth2.initCodeClient({
			client_id: PUBLIC_GOOGLE_CLIENT_ID,
			scope: 'https://www.googleapis.com/auth/gmail.modify',
			// what about mobile??
			ux_mode: 'popup',
			access_type: 'offline',
			auto_select: true,
			approval_prompt: 'auto',
			include_granted_scopes: true,
			immediate: true,
			// error_callback: myErrorCallback,
			callback: (response: OAuth2CallbackResponse) => {
				const xhr = new XMLHttpRequest();
				xhr.open('POST', '/account/profile/connect', true);
				xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
				// Set custom header for CRSF
				xhr.setRequestHeader('X-Requested-With', 'XmlHttpRequest');
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
		});

		client.requestCode();
	};
</script>

<script lang="ts">
	// TODOs
	// Clean up this file
	// Create an callback endpoint that
	// - verifies the user is the same
	// - Confirm the X-Requested-With: XmlHttpRequest header is set for popup mode.
	// - fetches access token and refresh token
	// - saves to db
	// - triggers a sync if first time
	// - subscribes to notifications
	// - validate on mobile!!
	// https://developers.google.com/identity/protocols/oauth2/web-server#httprest_3

	// Account Header
	// Client ID env var
	// load script from public folder with promise
	// https://stackoverflow.com/questions/59629947/how-do-i-load-an-external-js-library-in-svelte-sapper

	// check if in browser
</script>

<div class="mx-4 my-8 flex flex-row items-center justify-center sm:mx-auto sm:my-16">
	<div class="">
		<h1 class="text-2xl">Welcome to the SRC, Devin!</h1>
		<p class="text-gray-500">Let's get started by connecting your email.</p>
		<button
			class="flex items-center rounded-md bg-white px-2 py-2 text-slate-500 shadow hover:bg-slate-50"
			on:click|preventDefault={connectAccount}
			><img src="/google.svg" alt="Google logo" class="mr-3" />Connect Account</button
		>
		<script src="/scripts/gsi.js" async defer></script>
	</div>
</div>
