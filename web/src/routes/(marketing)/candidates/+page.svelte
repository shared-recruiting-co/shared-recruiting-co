<script lang="ts">
	import HowItWorks from '$lib/components/marketing/HowItWorks.svelte';
	import FAQs from '$lib/components/marketing/FAQs.svelte';
	import ProductFeatures from './Features.svelte';
	import Testimonials from './Testimonials.svelte';

	const steps = [
		{
			title: 'Detect & Archive',
			description:
				'SRC detects any inbound job opportunities. Each job opportunity is labeled, replied to, and archived for future reference.'
		},
		{
			title: 'Summarize & Standardize',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardized format. '
		},
		{
			title: 'Review & Engage',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardized format. '
		}
	];

	const faqs = [
		{
			question: 'How do I pronounce SRC?',
			answer: 'SRC is pronounced "source".'
		},
		{
			question: 'How can I join?',
			answer: "Sign up and we'll reach out once you are off the waiting list."
		},
		{
			question: 'How much does it cost to join SRC as a candidate?',
			answer: 'Nothing. It will always be free to join SRC as a candidate.'
		},
		{
			question: 'How does SRC use my data?',
			answer: `We process your emails in order to classify them as recruiting email or not. We do not and will never use your emails for any other purpose. We are take trust and privacy so seriously that all of our code is open source.
					<br/><br/>Read more in our <a class="underline" href="/legal/privacy-policy.md">privacy policy</a>`
		},
		{
			question: 'How do I contribute to SRC?',
			answer:
				"The SRC is open source. You can contribute to SRC on <a class='underline' href='https://github.com/shared-recruiting-co/shared-recruiting-co' target='_blank' rel='noopener noreferrer'>GitHub</a>. All you have to do is make a PR or an issue! <br/><br/> You can also contribute your recruiting emails to improve the SRC's ML model, by forwarding them to <a class='underline' href='mailto:examples@sharedrecruiting.co'>examples@sharederecruiting.co</a>"
		},
		{
			question: 'How does SRC detect job opportunities?',
			answer:
				'Every email you receive is processed by SRCs ML model. If the email is a job opportunity, the SRC will automatically label it as such.'
		},
		{
			question: 'How does the SRC define a job opportunity?',
			answer:
				'We define a job opportunity as an inbound email from a recruiter or hiring manager. This excludes job newsletters or automated notifications.'
		},
		{
			question: 'More Questions?',
			answer:
				"Reach out to us at <a class='underline' href='mailto:team@sharedrecruiting.co'>team@sharederecruiting.co</a>! We'd love to chat."
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
		Your Personal Recruiting<br />
		<span class="font-medium text-blue-600">Assistant</span>
	</h1>
	<p class="mx-auto mt-6 max-w-2xl text-base tracking-tight text-slate-700 sm:text-lg">
		The Shared Recruiting Co. (SRC) is an AI recruiting assistant that lives in your inbox. SRC
		keeps your inbox clean when you aren’t looking for jobs and supercharges your job search once
		you are.
	</p>
	<div class="mt-10 flex justify-center gap-x-6">
		<a
			class="group inline-flex items-center justify-center rounded-md bg-slate-900 py-2 px-4 text-sm font-semibold text-white shadow-lg hover:bg-slate-700 hover:text-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 active:bg-slate-800 active:text-slate-300"
			data-sveltekit-preload-data="hover"
			href="/login">Request an Invite</a
		>
		<a
			class="group inline-flex items-center justify-center rounded-md border border-slate-900 bg-white py-2 px-4 text-sm text-slate-900 shadow-lg hover:border-slate-700 hover:bg-slate-100 focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-slate-900 active:border-slate-800 active:bg-slate-100 active:text-slate-300"
			href="#product">Learn More</a
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
				href="/login">Become a Member</a
			>
		</div>
	</div>
</section>

<FAQs {faqs} />
