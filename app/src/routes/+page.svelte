<script lang="ts">
	import { page } from '$app/stores';
	import { supabaseClient } from '$lib/supabase/client';

	import HowItWorks from './how-it-works.svelte';
	import ProductFeatures from './product-features.svelte';
	import Testimonials from './testimonials.svelte';
	import FAQs from './faqs.svelte';

	const handleLogin = async () => {
		try {
			const { error } = await supabaseClient.auth.signInWithOAuth({
				provider: 'google',
				options: {
					redirectTo: `${window.location.origin}/login/callback`,
					// request refresh token from Google
					// prompt: 'consent' forces the consent flow every time the user logs in, so we always regenerate a refresh token
					// Once the app is approved by Google, we can remove the prompt: 'consent' option, so the user only has to consent once
					queryParams: { access_type: 'offline', prompt: 'consent' },
					scopes:
						'https://www.googleapis.com/auth/gmail.modify'
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

<div
	id="hero"
	class="mx-auto min-h-screen max-w-7xl px-4 pt-20 pb-16 text-center sm:px-6 lg:px-8 lg:pt-32"
>
	<h1
		class="font-display mx-auto max-w-4xl text-5xl font-medium tracking-tight text-slate-900 sm:text-7xl"
	>
		Tired of recruiting emails? <br />
		<span class="font-medium text-blue-600">So are we.</span>
	</h1>
	<p class="mx-auto mt-6 max-w-2xl text-lg tracking-tight text-slate-700">
		The Shared Recruiting Co. (SRC) is an AI recruiting assistant that lives in your inbox. SRC keeps
		your inbox clean when you aren’t looking for jobs and supercharges your job search once you are.
	</p>
	<div class="mt-10 flex justify-center gap-x-6">
		<a
			class="group inline-flex items-center justify-center rounded-full bg-slate-900 py-2 px-4 text-sm font-semibold text-white hover:bg-slate-700 hover:text-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 active:bg-slate-800 active:text-slate-300"
			href="https://62g7tcps1dn.typeform.com/to/ZR3hSTZi">Request an Invite</a
		>
		<a
			class="group inline-flex items-center justify-center rounded-full py-2 px-4 text-sm text-slate-700 ring-1 ring-slate-200 hover:text-slate-900 hover:ring-slate-300 focus:outline-none focus-visible:outline-blue-600 focus-visible:ring-slate-300 active:bg-slate-100 active:text-slate-600"
			href="https://github.com/shared-recruiting-co/shared-recruiting-co"
			target="_blank"
			rel="noopener noreferrer"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="h-6 w-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M17.25 6.75L22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3l-4.5 16.5"
				/>
			</svg>
			<span class="ml-3">Read the Code</span></a
		>
	</div>
	<p
		class="mx-auto mt-14 flex max-w-2xl flex-col items-center justify-center text-sm tracking-tight text-slate-700 md:flex-row"
	>
		<span><sup>*</sup>SRC is currently in an invite only beta.</span>
		{#if !$page.data.session}
			<!-- Comply with Google branding requirements -->
			<!-- https://developers.google.com/identity/branding-guidelines#top_of_page -->
			&NonBreakingSpace;If you have an account, then
			<button
				class="ml-2 flex items-center rounded-full px-2 py-2 text-slate-500 shadow"
				on:click={handleLogin}
				><img src="google.svg" alt="Google logo" class="mr-3" />Sign in with Google</button
			>
		{/if}
	</p>
</div>

<ProductFeatures />
<HowItWorks />
<Testimonials />

<FAQs />
