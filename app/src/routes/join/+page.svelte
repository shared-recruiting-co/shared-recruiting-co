<script lang="ts">
	import type { PageData } from './$types';
	import { page } from '$app/stores';
	import { enhance } from '$app/forms';

	const user = $page.data?.session?.user;
	const email = user?.email;
	const names = user?.user_metadata?.name?.split(' ');

	let firstName = names?.length ? names[0] : '';
	let lastName = names?.length ? names[names.length - 1] : '';
	let linkedin = '';
	let referrer = '';

	export let form: null | {
		success?: boolean;
		errors: Record<string, string>;
	};

	export let data: PageData;

	let formError: (field: string) => string | null;
	$: {
		formError = (field: string): string | null => {
			if (!form || !form.errors) return null;

			return (form.errors[field] as string) || null;
		};
	}

	// next steps
	// 3. Redirect to profile creation page if off waitlist/has account
	// 4. Create profile page
	// Confirm first name and last name
	// Verify the granted access to gmail
	// 5. Profile page (edit profile, pause/start email assistant)
	// (Future) sync page
	// (Future) connect LinkedIn
	// 0. Button redesign on homepage
</script>

<div class="mx-4 my-12 max-w-2xl rounded-md bg-blue-100 p-12 sm:mx-auto">
	{#if data?.success || form?.success}
		<div class="text-center">
			<h1 class="text-2xl font-bold">You're on the waitlist!</h1>
			<p class="mt-4">We'll email you once you are off the list.</p>
			<div class="mt-8">
				<a href="/" class="underline hover:text-slate-700">Go home &rarr;</a>
			</div>
		</div>
	{:else}
		<h1 class="text-4xl">Request an Invite</h1>
		<div class="my-6">
			<p>Almost there! SRC is currently invite only.</p>
			<p>Please fill out the form below to get priority on the invite list.</p>
		</div>
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
						placeholder="you@example.com"
					/>
				</div>
			</div>
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
			<div>
				<label for="linkedin" class="block text-sm font-medium text-slate-700"
					>LinkedIn Profile</label
				>
				<div class="mt-1">
					<input
						required
						type="text"
						name="linkedin"
						id="linkedin"
						bind:value={linkedin}
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
						placeholder="https://www.linkedin.com/in/devin-stein-087148107/"
					/>
				</div>
				{#if formError('linkedin')}
					<p class="mt-1 text-xs text-rose-500">{formError('linkedin')}</p>
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
						bind:value={referrer}
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
					<span class="text-sm text-slate-500" id="comment-optional">Optional</span>
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

			<div class="flex w-full flex-col">
				<button
					type="submit"
					class="inline-flex justify-center rounded-md border border-transparent bg-slate-900 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
				>
					Request an Invite
				</button>
				{#if formError('submit')}
					<p class="mt-1 text-center text-xs text-rose-500">{formError('submit')}</p>
				{/if}
			</div>
		</form>
	{/if}
</div>
