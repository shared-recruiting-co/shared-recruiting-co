<script lang="ts">
	// import FAQs from '$lib/components/marketing/FAQs.svelte';
	import FAQs from '../faqs.svelte';

	// screenshot of Github?
	// screenshot of Email?
	// screenshot of CASA?
</script>

<div
	id="hero"
	class="mx-auto min-h-screen max-w-7xl px-4 pt-20 pb-16 text-center sm:px-6 lg:px-8 lg:pt-32"
>
	<h1
		class="font-display mx-auto max-w-4xl text-5xl font-medium tracking-tight text-slate-900 sm:text-7xl"
	>
		The Only <span class="font-medium text-blue-600">Open Source, Privacy First</span>
		<br />Recruiting Platform
	</h1>
	<p class="mx-auto mt-6 max-w-2xl text-lg tracking-tight text-slate-700">
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

<section class="mx-auto min-h-screen max-w-xl px-4 text-left sm:px-0">
	<h2 class="mb-4 text-4xl">CASA Compliant</h2>
	<p class="text-lg">
		The Cloud Application Security Assessment (CASA) is built upon the industry-recognized standards
		of the OWASP's Application Security Verification Standard (ASVS) to provide a consistent set of
		requirements to harden security for any application. CASA provides a uniform way to perform
		trusted assurance assessments of these requirements when such assessments are required for
		applications with potential access to sensitive data. SRC is Tier-2 CASA compliant. SRC's CASA
		compliance has been verified by an independent lab partner as part of the Google OAuth approval
		process.
	</p>
</section>
<section class="mx-auto min-h-screen max-w-xl px-4 text-left sm:px-0">
	<h2 class="mb-4 text-4xl">Open Source</h2>
	<p>
		We believe Linus's law of "given enough eyeballs, all bugs are shallow" applies to security
		issues. The premier example of how open source projects can be more secure than proprietary code
		bases is Bitcoin. In his 2015 talk Andreas M. Antonopoulos describes how closed source banking
		systems have the software equivalent of weak immune systems, because huge security holes can be
		obfuscated for long periods of time, and when eventually exploited can have enormous detrimental
		effects. On the flip side of this is an open source protocol like Bitcoin, where any security
		holes are there for all to see. Exploits are found early and often, and then patched. Remember
		that successful software companies can take more than a decade to build. Over a long time
		period, open source systems will tend towards a more secure state over secretive, proprietary
		systems.
	</p>
	<p class="text-sm italic">Credit: Should I Open Source My Company</p>
</section>
<section class="mx-auto min-h-screen max-w-xl px-4 text-left sm:px-0">
	<h2 class="mb-4 text-4xl">Privacy-First</h2>
	<p>
		At SRC, members trust us with some of their most sensitive information, their emails and their
		job status. Sadly, we've seen firsthand what happens when companies abuse this. We've seen
		employers use LinkedIn to find out which employees are looking for new roles. We've seen
		backlash against recruiting companies for making candidate job profiles public without their
		consent. SRC does not and will never store your emails within our database. We use native email
		functionality, like labels and folders, to securely manage your emails from within your inbox.
		And with SRC's two-way opt-in communication, your job status is confidential. Companies cannot
		see if you are looking for a new role until you choose to start the interview process. By open
		sourcing SRC, we are making more then a promise of privacy. All the ways we use and protect your
		data is transparent and publicly visible.
	</p>
</section>

<FAQs />
