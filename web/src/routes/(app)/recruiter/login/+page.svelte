<script lang="ts">
	import type { LayoutData } from './$types';

	export let data: LayoutData;
	let error: string | null = null;

	$: ({ supabase } = data);

	const handleLogin = async () => {
		try {
			const { error } = await supabase.auth.signInWithOAuth({
				provider: 'google',
				options: {
					redirectTo: `${window.location.origin}/recruiter/profile`
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

<header
	class="sticky top-0 z-10 flex flex-row flex-nowrap items-center justify-between space-x-4 bg-blue-50 px-4 py-4 text-lg sm:text-xl"
>
	<div class="flex w-96 flex-row items-center space-x-4">
		<a href="/" class="-m-1.5 p-1.5">
			<span class="sr-only">Shared Recruiting Co.</span>
			<img src="/logo.svg" alt="Shared Recruiting Co" class="h-8 md:h-10" />
		</a>
		<a href="/" class="text-xl text-slate-900 sm:text-2xl md:min-w-[150px]">Recruiter</a>
	</div>
</header>
<div class="mx-4 my-12 rounded-md bg-blue-100 px-12 pt-12 pb-4 sm:mx-auto sm:w-1/3">
	<h1 class="text-4xl">Sign In</h1>
	<p class="my-4">Sign into SRC with your work email.</p>
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
		Not a recruiter? <a href="/login" class="underline">Sign in here</a>
	</p>
</div>
