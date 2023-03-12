<script lang="ts">
	import FAQs from './FAQs.svelte';
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
		The Shared Recruiting Co. (SRC) is an open-source, candidate-centric recruiting platform. SRC
		makes recruiting more effecient and less painful by promoting transparent, opt-in communication
		between candidates and companies.
	</p>
	<div class="mt-10 flex justify-center gap-x-6">
		<a
			class="group inline-flex items-center justify-center rounded-md bg-slate-900 py-2 px-4 text-sm font-semibold text-white shadow-lg hover:bg-slate-700 hover:text-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 active:bg-slate-800 active:text-slate-300"
			data-sveltekit-preload-data="hover"
			href="/login">Request an Invite</a
		>
		<a
			class="group inline-flex items-center justify-center rounded-md border border-slate-900 bg-white py-2 px-4 text-sm text-slate-900 shadow-lg hover:border-slate-700 hover:bg-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 active:border-slate-800 active:bg-slate-100 active:text-slate-300"
			href="/#product">Learn More</a
		>
	</div>
	<p
		class="mx-auto mt-14 flex max-w-2xl flex-col items-center justify-center text-sm tracking-tight text-slate-700 md:flex-row"
	>
		<span
			><sup>*</sup>SRC is currently in an invite only beta. If you already have an account,
			<a data-sveltekit-preload-data="hover" href="/login" class="underline">log in</a></span
		>
	</p>
</div>

<FAQs />
