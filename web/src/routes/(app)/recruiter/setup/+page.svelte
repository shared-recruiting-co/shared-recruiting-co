<script lang="ts">
	import type { PageData } from './$types';
	import { enhance } from '$app/forms';
	import { supabaseClient } from '$lib/supabase/client';

	export let data: PageData;
	const user = data?.session?.user;
	const email = user?.email;
	// guess company website from email domain
	const emailDomain = email?.split('@')[1];
	const names = user?.user_metadata?.name?.split(' ');

	let firstName = names?.length ? names[0] : '';
	let lastName = names?.length ? names[names.length - 1] : '';
	let companyWebsite = emailDomain && emailDomain !== 'gmail.com' ? `https://${emailDomain}` : '';

	let tos = false;

	export let form: null | {
		success?: boolean;
		errors: Record<string, string>;
	};

	let formError: (field: string) => string | null;
	$: {
		formError = (field: string): string | null => {
			if (!form || !form.errors) return null;

			return (form.errors[field] as string) || null;
		};
	}

	const handleLogout = async () => {
		await supabaseClient.auth.signOut();
	};
</script>

<header
	class="sticky top-0 z-10 flex flex-row flex-nowrap items-center justify-between space-x-4 bg-blue-50  px-4 py-4 text-lg sm:text-xl"
>
	<div class="flex w-96 flex-row items-center space-x-4">
		<img src="/logo.svg" alt="Shared Recruiting Co" class="h-6 md:h-10" />
		<span>Recruiter</span>
	</div>
	<div class="flex w-96 items-center justify-end space-x-4">
		<button class="text-base hover:underline active:underline sm:text-lg" on:click={handleLogout}
			>Log Out</button
		>
	</div>
</header>
<div class="m-2 flex flex-row items-center justify-center sm:mx-auto sm:my-16">
	<div class="flex max-w-xl flex-col space-y-8 rounded-md p-8">
		<h1 class="text-center text-3xl sm:text-4xl">Welcome to SRC Recruiter</h1>
		<p class="mt-2">
			We're excited to have you on board! Before we can create your account, please confirm your
			profile information
		</p>
		<form method="POST" use:enhance class="flex flex-col space-y-6">
			<div>
				<label for="email" class="block text-sm font-medium text-slate-700">Email</label>
				<div class="mt-1">
					<input
						type="email"
						name="email"
						id="email"
						value={email}
						disabled
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 disabled:cursor-not-allowed disabled:border-slate-200 disabled:bg-slate-50 disabled:text-slate-500 sm:text-sm"
						placeholder={email}
					/>
				</div>
			</div>
			<div class="grid grid-cols-2 gap-2">
				<div>
					<label for="firstName" class="block text-sm font-medium text-slate-700">First Name</label>
					<div class="mt-1">
						<input
							required
							type="text"
							name="firstName"
							id="firstName"
							bind:value={firstName}
							class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
							placeholder="Devin"
						/>
					</div>
					{#if formError('firstName')}
						<p class="mt-1 text-xs text-rose-500">{formError('firstName')}</p>
					{/if}
				</div>
				<div>
					<label for="lastName" class="block text-sm font-medium text-slate-700">Last Name</label>
					<div class="mt-1">
						<input
							required
							type="text"
							name="lastName"
							id="lastName"
							bind:value={lastName}
							class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
							placeholder="Stein"
						/>
					</div>
					{#if formError('lastName')}
						<p class="mt-1 text-xs text-rose-500">{formError('lastName')}</p>
					{/if}
				</div>
			</div>
			<div>
				<label for="company" class="block text-sm font-medium text-slate-700">Company</label>
				<div class="mt-1">
					<input
						required
						type="text"
						name="company"
						id="company"
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
						placeholder="ACME Co."
					/>
				</div>
				{#if formError('company')}
					<p class="mt-1 text-xs text-rose-500">{formError('company')}</p>
				{/if}
			</div>
			<div>
				<label for="companyWebsite" class="block text-sm font-medium text-slate-700"
					>Company Website</label
				>
				<div class="mt-1">
					<input
						required
						type="text"
						name="companyWebsite"
						id="companyWebsite"
						bind:value={companyWebsite}
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
						placeholder="https://sharedrecruiting.co"
					/>
				</div>
				{#if formError('companyWebsite')}
					<p class="mt-1 text-xs text-rose-500">{formError('companyWebsite')}</p>
				{/if}
			</div>
			<div>
				<label for="referrer" class="block text-sm font-medium text-slate-700"
					>How did you hear about SRC?</label
				>
				<div class="mt-1">
					<input
						required
						type="text"
						name="referrer"
						id="referrer"
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
						placeholder="A little bird told me"
					/>
				</div>
				{#if formError('referrer')}
					<p class="mt-1 text-xs text-rose-500">{formError('referrer')}</p>
				{/if}
			</div>
			<div>
				<div class="flex justify-between">
					<label for="comment" class="block text-sm font-medium text-slate-700"
						>Anything else you'd like us to know?</label
					>
					<span class="text-xs text-slate-500 sm:text-sm" id="comment-optional">Optional</span>
				</div>
				<div class="mt-1">
					<textarea
						rows="4"
						name="comment"
						id="comment"
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
						placeholder="I'd like to join SRC because..."
					/>
				</div>
			</div>
			<div>
				<!-- Agree to Terms of Service -->
				<div class="mb-4 text-sm">
					<input
						type="checkbox"
						bind:checked={tos}
						required
						id="tos"
						name="tos"
						class="rounded-md"
					/>
					<label for="tos" class="ml-2"
						>I agree to the <a
							href="/legal/terms-of-service"
							class="underline"
							target="_blank"
							rel="noopener noreferrer">SRC Terms of Service</a
						></label
					>
					{#if formError('tos')}
						<p class="mt-1 text-left text-xs text-rose-500">{formError('tos')}</p>
					{/if}
				</div>
			</div>
			<p class="text-sm">
				You can read more about how we use and protect your data in our <a
					href="/legal/privacy-policy"
					target="_blank"
					class="underline">privacy policy</a
				>. If you want an even deeper dive into how we use your data, you can read the
				<a
					href="https://github.com/shared-recruiting-co/shared-recruiting-co"
					class="underline"
					target="_blank"
					rel="noopener noreferrer">SRC code</a
				>.
			</p>

			<div class="flex w-full flex-col">
				<button
					type="submit"
					class="inline-flex justify-center rounded-md border border-transparent bg-slate-900 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
				>
					Get Started
				</button>
				{#if formError('submit')}
					<p class="mt-1 text-center text-xs text-rose-500">{formError('submit')}</p>
				{/if}
			</div>
		</form>
	</div>
</div>
