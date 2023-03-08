<script lang="ts">
	import HowItWorks from '$lib/components/marketing/HowItWorks.svelte';
	import ProductFeatures from './Features.svelte';
	import Testimonials from '../../testimonials.svelte';
	import FAQs from '../../faqs.svelte';

	const steps = [
		{
			title: 'Sync',
			description:
				'SRC detects any inbound job opportunities. Each job opportunity is labeled, replied to, and archived for future reference.'
		},
		{
			title: 'Import',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardized format. '
		},
		{
			title: 'Review & Engage',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardized format. '
		}
	];
</script>

<div
	id="hero"
	class="mx-auto min-h-screen max-w-7xl px-4 pt-20 pb-16 text-center sm:px-6 lg:px-8 lg:pt-32"
>
	<h1
		class="font-display mx-auto max-w-4xl text-5xl font-medium tracking-tight text-slate-900 sm:text-7xl"
	>
		<!-- <span class="font-medium text-blue-600">Make the Most</span> Out of Your<br /> -->
		<!-- Outbound Recruiting -->
		Stop Wasting Time<br />on the
		<span class="font-medium text-blue-600">Wrong Candidates</span>
	</h1>
	<p class="mx-auto mt-6 max-w-2xl text-lg tracking-tight text-slate-700">
		The Shared Recruiting Co. (SRC) integrates into your exisitng recruiting stack to help you find
		the right candidates at the right time.
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

<ProductFeatures />
<HowItWorks {steps} />
<Testimonials />

<section id="join" class="flex flex-col items-center justify-center py-28">
	<div class="flex flex-col items-center justify-center space-y-6">
		<div class="text-3xl">Better for <span class="text-blue-500 underline">Candidates</span>.</div>
		<div class="text-3xl">Better for <span class="text-blue-500 underline">Companies</span>.</div>
		<div class="text-3xl">
			Better for <span class="text-blue-500 underline">Recruiting</span>.
		</div>
		<div class="pt-4">
			<a
				class="group inline-flex items-center justify-center rounded-md bg-slate-900 py-2 px-4 text-lg font-semibold text-white shadow-lg hover:bg-slate-700 hover:text-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 active:bg-slate-800 active:text-slate-300"
				data-sveltekit-preload-data="hover"
				href="/recruiter/login">Become a Member</a
			>
		</div>
	</div>
</section>

<FAQs />
