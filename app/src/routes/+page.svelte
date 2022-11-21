<script lang="ts">
	import { supabaseClient } from '$lib/supabase/client';

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
	<h1 class="pb-8 text-center text-6xl">Welcome to Shared Recruiting Co.</h1>
	<button on:click={handleLogin} class="rounded bg-gray-900 px-4 py-1 text-white">Log In</button>
</div>
