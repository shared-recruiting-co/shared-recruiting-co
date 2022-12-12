<script lang="ts" context="module">
	import { browser } from '$app/environment';
	import { PUBLIC_GOOGLE_CLIENT_ID } from '$env/static/public';

	// TODOs
	// 
	// Refactor Connect button to be a component
	// Validate on mobile!! Use redirect flow if necessary
	// Create separate layouts for (marketing) and (account) pages 
	// Create account header
	// Create account sidebar
	// Design & create account page 
	// Trigger welcome email 
	// Update cloud functions logic (sync to date OR sync to history id)
	// Update watch cloud function to use filter action
	// Watch/Stop endpoints
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
		});

		client.requestCode();
	};
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
