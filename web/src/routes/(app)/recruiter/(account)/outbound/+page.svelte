<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;

	import { getPaginationPages } from '$lib/pagination';

	$: ({ supabase, sequences, pagination, jobs } = data);

	$: showingStart = (pagination.page - 1) * pagination.perPage + 1;
	$: showingEnd = Math.min(pagination.page * pagination.perPage, pagination.numResults);
	$: prevPage = pagination.hasPrev ? pagination.page - 1 : pagination.page;
	$: nextPage = pagination.hasNext ? pagination.page + 1 : pagination.page;

	$: pages = getPaginationPages(pagination.page, pagination.numPages);

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
	// Empty State
	// Smaller Screens!
	// No Job -> Create your first job notification
	// All Sequences w/ Ignore Button & Toggle/Filter
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
		{#if pagination.numPages > 1}
			<div
				class="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3.5 sm:px-6"
			>
				<nav class="flex flex-1 justify-end md:hidden" data-sveltekit-noscroll>
					<a
						href="/account/jobs?page={prevPage}"
						class="relative inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
						class:hover:bg-slate-50={pagination.hasPrev}
						class:cursor-not-allowed={!pagination.hasPrev}>Previous</a
					>
					<a
						href="/account/jobs?page={nextPage}"
						class="relative ml-3 inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
						class:hover:bg-slate-50={pagination.hasNext}
						class:cursor-not-allowed={!pagination.hasNext}>Next</a
					>
				</nav>
				<div class="hidden sm:flex-1 sm:items-center sm:justify-between md:flex">
					<div>
						<p class="text-sm text-slate-700">
							Showing
							<span class="font-semibold">{showingStart}</span>
							to
							<span class="font-semibold">{showingEnd}</span>
							of
							<span class="font-semibold">{pagination.numResults}</span>
							results
						</p>
					</div>
					<div>
						<nav
							class="isolate inline-flex -space-x-px rounded-md shadow-sm"
							aria-label="Pagination"
							data-sveltekit-noscroll
						>
							<a
								href="/account/jobs?page={prevPage}"
								class:focus:z-20={pagination.hasPrev}
								class:hover:bg-slate-50={pagination.hasPrev}
								class:cursor-not-allowed={!pagination.hasPrev}
								class="relative inline-flex items-center rounded-l-md border border-slate-300 bg-white px-2 py-2 text-sm font-semibold text-slate-500"
							>
								<span class="sr-only">Previous</span>
								<!-- Heroicon name: mini/chevron-left -->
								<svg
									class="h-5 w-5"
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 20 20"
									fill="currentColor"
									aria-hidden="true"
								>
									<path
										fill-rule="evenodd"
										d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z"
										clip-rule="evenodd"
									/>
								</svg>
							</a>
							{#each pages as page}
								{@const current = page === `${pagination.page}`}
								{#if page === '...'}
									<span
										class="relative inline-flex items-center border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
										>...</span
									>
								{:else}
									<a
										href="/account/jobs?page={page}"
										class="relative inline-flex items-center border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-500 focus:z-20"
										class:hover:bg-blue-50={!current}
										class:z-10={current}
										class:bg-blue-50={current}
										class:border-blue-500={current}
										class:text-blue-600={current}>{page}</a
									>
								{/if}
							{/each}
							<a
								href="/account/jobs?page={nextPage}"
								disabled={!pagination.hasNext}
								class="relative inline-flex items-center rounded-r-md border border-slate-300 bg-white px-2 py-2 text-sm font-medium text-slate-500"
								class:focus:z-20={pagination.hasNext}
								class:hover:bg-slate-50={pagination.hasNext}
								class:cursor-not-allowed={!pagination.hasNext}
							>
								<span class="sr-only">Next</span>
								<!-- Heroicon name: mini/chevron-right -->
								<svg
									class="h-5 w-5"
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 20 20"
									fill="currentColor"
									aria-hidden="true"
								>
									<path
										fill-rule="evenodd"
										d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z"
										clip-rule="evenodd"
									/>
								</svg>
							</a>
						</nav>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
