<script lang="ts">
	import { supabaseClient } from '$lib/supabase/client';

	const handleLogin = async () => {
		try {
			const { error } = await supabaseClient.auth.signInWithOAuth({
				provider: 'google',
				options: {
					// request refresh token from Google
					queryParams: { access_type: 'offline', prompt: 'consent' },
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

<div class="min-h-screen flex flex-col justify-center items-center">
	<h1 class="text-6xl text-center pb-8">Welcome to Shared Recruiting Co.</h1>
	<button on:click={handleLogin} class="px-4 py-1 bg-gray-900 text-white rounded">Log In</button>
</div>
