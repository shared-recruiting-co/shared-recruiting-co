<script lang="ts">
	import type { LayoutData } from './$types';

	export let data: LayoutData;

	$: ({ supabase } = data);

	let error: string | null = null;

	const handleLogin = async () => {
		try {
			const { error } = await supabase.auth.signInWithOAuth({
				provider: 'google',
				options: {
					redirectTo: `${window.location.origin}/account/profile`
				}
			});
			if (error) throw error;
		} catch (error: unknown) {
			if (error instanceof Error) {
				error = error.message || 'Something went wrong. Please try again.';
			}
		}
	};
</script>

<div class="mx-4 my-12 max-w-2xl rounded-md bg-blue-100 px-12 pt-12 pb-4 sm:mx-auto">
	<h1 class="text-4xl">Sign In</h1>
	<p class="my-4">Sign into SRC with the email you typically receive recruiting inbound on.</p>
	<div class="mt-12">
		<button
			class="flex items-center rounded-md bg-white px-2 py-2 text-slate-500 shadow hover:bg-slate-50"
			on:click={handleLogin}
			><img src="/google.svg" alt="Google logo" class="mr-3" />Sign in with Google</button
		>
		{#if error}
			<div class="mt-2 text-sm text-rose-500">{error}</div>
		{/if}
	</div>
	<p class="mt-12 mb-4 text-sm">
		Are you a recruiter? If so, <a href="/recruiter/login" class="underline">sign in here</a>
	</p>
</div>
