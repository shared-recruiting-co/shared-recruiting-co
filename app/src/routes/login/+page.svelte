<script lang="ts">
	import { supabaseClient } from '$lib/supabase/client';

	const handleLogin = async () => {
		try {
			const { error } = await supabaseClient.auth.signInWithOAuth({
				provider: 'google',
				options: {
					redirectTo: `${window.location.origin}/account/profile`
				}
			});
			if (error) throw error;
		} catch (error) {
			// TODO: handle error
			if (error instanceof Error) {
				alert(error.message);
			}
		}
	};
</script>

<div class="mx-4 my-12 max-w-2xl rounded-md bg-blue-50 p-12 sm:mx-auto">
	<h1 class="text-4xl">Sign In</h1>
	<p class="my-4">Sign into SRC with the email you typically receive recruiting inbound on.</p>
	<div class="mt-12">
		<button
			class="flex items-center rounded-md bg-white px-2 py-2 text-slate-500 shadow hover:bg-slate-50"
			on:click={handleLogin}
			><img src="/google.svg" alt="Google logo" class="mr-3" />Sign in with Google</button
		>
	</div>
</div>
