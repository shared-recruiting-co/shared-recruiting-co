<script lang="ts">
	import { enhance } from '$app/forms';

	const email = 'devin@sharedrecruiting.co';
	// $page.user.user_metadata.name.split(" ")[0]
	let firstName = 'Devin';
	let lastName = 'Stein';
	let linkedin = '';
	let referrer = '';

	// form validation
	// Show login button if not logged in
	// handle user already on the waitlist (if they are, show a message) (if they are off, do something else)
	//
	// What happens when someone is off the waitlist?
	// Click to active account (we need a lock to prevent people from clicking it multiple times)
</script>

<div class="mx-auto my-20 max-w-2xl rounded-md bg-blue-50 p-12">
	<h1 class="text-4xl">Request Invite</h1>
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
		</div>
		<div>
			<label for="linkedin" class="block text-sm font-medium text-slate-700">LinkedIn Profile</label
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
		<button
			type="submit"
			class="inline-flex justify-center rounded-md border border-transparent bg-slate-900 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
		>
			Request Invite
		</button>
	</form>
</div>
