<script lang="ts">
	import { page } from '$app/stores'
	import { supabaseClient } from '$lib/supabase/client';

	const handleLogout = async () => {
		await supabaseClient.auth.signOut()
	}

	const handleLogin = async () => {
		try {
			const { error } = await supabaseClient.auth.signInWithOAuth({
				provider: 'google',
				options: {
					redirectTo: `${window.location.origin}/login/callback`,
					// request refresh token from Google
					queryParams: { access_type: 'offline' },
					scopes:
						'https://www.googleapis.com/auth/gmail.modify https://www.googleapis.com/auth/gmail.labels'
				}
			});
			if (error) throw error;
		} catch (error) {
			if (error instanceof Error) {
				alert(error.message);
			}
		}
	};
</script>

<div class="flex min-h-screen flex-col items-center justify-center">
	<h1 class="pb-8 text-center text-6xl">Welcome to the S.R.C</h1>
	{#if $page.data.session}
		<p class="py-8 text-center text-2xl">You are logged in!</p>
	<button on:click={handleLogout} class="rounded bg-gray-900 px-4 py-1 text-white">Log Out</button>
	{:else}
	<button on:click={handleLogin} class="rounded bg-gray-900 px-4 py-1 text-white">Log In</button>
	{/if}
</div>
