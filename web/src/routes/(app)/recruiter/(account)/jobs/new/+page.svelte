<script lang="ts">
	import { enhance } from '$app/forms';

	let title = '';
	let descriptionURL = '';

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
</script>

<div class="m-2 flex flex-row items-center justify-center sm:mx-auto sm:my-16">
	<div class="flex max-w-xl flex-col space-y-8 rounded-md p-8">
		<h1 class="text-3xl sm:text-4xl">Add a Job</h1>
		<p class="mt-2">Fill out the form below to add a new job to your company's job board.</p>
		<form method="POST" use:enhance class="flex flex-col space-y-6">
			<div>
				<label for="title" class="block text-sm font-medium text-slate-700">Job Title</label>
				<div class="mt-1">
					<input
						required
						type="text"
						name="title"
						id="title"
						bind:value={title}
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
						placeholder="Software Engineer"
					/>
				</div>
				{#if formError('title')}
					<p class="mt-1 text-xs text-rose-500">{formError('title')}</p>
				{/if}
			</div>
			<div>
				<label for="descriptionURL" class="block text-sm font-medium text-slate-700"
					>Job Description URL</label
				>
				<div class="mt-1">
					<input
						required
						type="text"
						name="descriptionURL"
						id="descriptionURL"
						bind:value={descriptionURL}
						class="block w-full rounded-md border-slate-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
						placeholder="https://boards.greenhouse.io/..."
					/>
				</div>
				{#if formError('descriptionURL')}
					<p class="mt-1 text-xs text-rose-500">{formError('descriptionURL')}</p>
				{/if}
			</div>
			<div class="flex w-full flex-col">
				<button
					type="submit"
					class="inline-flex justify-center rounded-md border border-transparent bg-slate-900 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
				>
					Add Job
				</button>
				{#if formError('submit')}
					<p class="mt-1 text-center text-xs text-rose-500">{formError('submit')}</p>
				{/if}
			</div>
		</form>
	</div>
</div>
