<script lang="ts">
	import { supabaseClient } from '$lib/supabase/client';

	let loading = false;
	let email: string;

	const handleLogin = async () => {
		try {
			loading = true;
			const { data, error } = await supabaseClient.auth.signInWithOAuth({
				provider: 'google',
				options: {
					queryParams: { access_type: 'offline' },
					scopes:
						'https://www.googleapis.com/auth/gmail.modify https://www.googleapis.com/auth/gmail.labels'
				}
			});
			console.log(data);
			if (error) throw error;
		} catch (error) {
			if (error instanceof Error) {
				alert(error.message);
			}
		} finally {
			loading = false;
		}
	};
</script>

<div class="min-h-screen flex flex-col justify-center items-center">
	<h1 class="text-6xl text-center pb-8">Welcome to Shared Recruiting Co.</h1>
	<button on:click={handleLogin} class="px-4 py-1 bg-gray-900 text-white rounded">Log In</button>
</div>
