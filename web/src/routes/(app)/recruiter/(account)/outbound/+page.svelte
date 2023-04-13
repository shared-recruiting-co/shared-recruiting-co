<script lang="ts">
	import type { PageData } from './$types';

	import TableFooter from '$lib/components/TableFooter.svelte';

	export let data: PageData;

	$: ({ supabase, sequences, pagination, jobs } = data);

	const NO_JOB = 'Not Assigned';

	const onSelect = (templateID: string) => async (e: Event) => {
		const target = e.target as HTMLSelectElement;
		let value: string | null = target.value;
		if (value === NO_JOB) {
			value = null;
		}

		await supabase
			.from('recruiter_outbound_template')
			.update({ job_id: value })
			.eq('template_id', templateID);
	};

	// TODO
	// Smaller Screens!
	// No Job -> Create your first job notification
	// All Sequences w/ Ignore Button & Toggle/Filter
	// Dropdown for ignoring? Hover over to show ignore button?
	// View in Gmail button?
	// Use Combobox for Job Selection
</script>

<div>
	<h1 class="text-3xl sm:text-4xl">Outbound</h1>
	<p class="mt-1 text-sm text-slate-500">
		SRC automatically imports and syncs your outbound recruiting emails
	</p>
</div>
<!-- Table -->
<!-- Possible states 1. Empty 2. All Sequences (Default ignore is hidden) 3. All sequences w/ ignore -->
<!-- To start, ignore ignore functionality -->
{#if !sequences || sequences.length === 0}
	<!-- Empty state -->
	<div class="mx-auto flex max-w-lg flex-1 flex-col items-center justify-center py-20">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="h-12 w-12 text-slate-500"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75"
			/>
		</svg>

		<p class="mt-2 text-lg font-medium">No outbound recruiting emails found.</p>
		<p class="mt-4 text-sm font-medium text-slate-700">
			To import your recruiting outbound, connect your Gmail account from your <a
				href="/recruiter/profile"
				class="underline">profile</a
			>. Once conncected, SRC will automatically detect, import, and sync outbound recruiting emails
			and display them here.
		</p>
	</div>
{:else}
	<div>
		<div class="mt-8 overflow-hidden rounded-lg shadow ring-1 ring-black ring-opacity-5 md:mx-0">
			<table class="min-w-full divide-y divide-slate-300">
				<thead class="bg-slate-50 text-left text-sm">
					<tr>
						<th
							scope="col"
							class="hidden py-3.5 pl-4 pr-3 font-semibold text-slate-900 sm:pl-6 lg:table-cell"
							>Subject</th
						>
						<th scope="col" class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
							>Body</th
						>
						<th scope="col" class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
							>Job</th
						>
						<th
							scope="col"
							class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
						/>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-200 bg-white">
					{#each sequences as sequence}
						<tr class="text-sm text-slate-700">
							<td class="hidden py-4 pl-4 pr-3 font-medium text-slate-900 sm:pl-6 lg:table-cell">
								{sequence.subject}
							</td>
							<td
								class="hidden truncate py-4 pl-4 pr-3 font-medium text-slate-700 lg:table-cell lg:px-3"
							>
								{sequence.body.slice(0, 40)}...
							</td>
							<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 lg:px-3">
								<select
									id="job"
									name="job"
									class="mt-2 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-slate-700 ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-blue-600 sm:text-sm sm:leading-6"
									on:change={onSelect(sequence.template_id)}
								>
									<option selected={sequence.job_id === null}>Not Assigned</option>
									{#each jobs as job}
										<option value={job.job_id} selected={job.job_id === sequence.job_id}
											>{job.title}</option
										>
									{/each}
								</select>
							</td>
							<td
								class="hidden truncate py-4 pl-4 pr-3 font-medium text-slate-900 lg:table-cell lg:px-3"
							>
								<div class="flex flex-row items-center text-slate-500">
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 20 20"
										fill="currentColor"
										class="mr-1 h-5 w-5"
									>
										<path
											d="M10 8a3 3 0 100-6 3 3 0 000 6zM3.465 14.493a1.23 1.23 0 00.41 1.412A9.957 9.957 0 0010 18c2.31 0 4.438-.784 6.131-2.1.43-.333.604-.903.408-1.41a7.002 7.002 0 00-13.074.003z"
										/>
									</svg>

									<span>
										{sequence.recipient_count}
									</span>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			<TableFooter pagination={pagination}/>
		</div>
	</div>
{/if}
