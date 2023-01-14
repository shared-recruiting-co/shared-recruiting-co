<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;

	const showingStart = (data.pagination.page - 1) * data.pagination.perPage + 1;
	const showingEnd = Math.min(
		data.pagination.page * data.pagination.perPage,
		data.pagination.numResults
	);
	const prevPage = data.pagination.hasPrev ? data.pagination.page - 1 : data.pagination.page;
	const nextPage = data.pagination.hasNext ? data.pagination.page + 1 : data.pagination.page;

	const getPages = (current: number, total: number): string[] => {
		if (total <= 5) {
			return Array.from({ length: total }, (_, i) => (i + 1).toString());
		}
		switch (current) {
			case 1:
				return ['1', '2', '...', (total - 1).toString(), total.toString()];
			case 2:
				return ['1', '2', '3', '...', (total - 1).toString(), total.toString()];
			case 3:
				// special case for 6
				if (total === 6) {
					return ['1', '2', '3', '4', '5', '6'];
				}
				return ['1', '2', '3', '4', '...', (total - 1).toString(), total.toString()];
			case total - 2:
				// special case for 6
				if (total === 6) {
					return ['1', '2', '3', '4', '5', '6'];
				}
				return [
					'1',
					'2',
					'...',
					(total - 3).toString(),
					(total - 2).toString(),
					(total - 1).toString(),
					total.toString()
				];
			case total - 1:
				return ['1', '2', '...', (total - 2).toString(), (total - 1).toString(), total.toString()];
			case total:
				return ['1', '2', '...', (total - 1).toString(), total.toString()];
			default:
				return [
					'1',
					'2',
					...(current - 1 == 3 ? [] : ['...']),
					(current - 1).toString(),
					current.toString(),
					(current + 1).toString(),
					...(current + 1 == total - 2 ? [] : ['...']),
					(total - 1).toString(),
					total.toString()
				];
		}
	};
	const pages = getPages(data.pagination.page, data.pagination.numPages);

	// TOOD:
	// Description
	// Error state (icon + No jobs + Description)
	// Empty state
	// Realtime?
	// move get pages to utils
</script>

<div>
	<h1 class="text-3xl sm:text-4xl">Jobs</h1>
	<p class="mt-1 text-sm text-slate-500">Coming soon!</p>
</div>
<div>
	<div class="mt-8 overflow-hidden rounded-lg shadow ring-1 ring-black ring-opacity-5 md:mx-0">
		<table class="min-w-full divide-y divide-slate-300">
			<thead class="bg-slate-50 text-left text-sm">
				<tr>
					<th
						scope="col"
						class="hidden py-3.5 pl-4 pr-3 font-semibold text-slate-900 sm:pl-6 lg:table-cell"
						>Company</th
					>
					<th scope="col" class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
						>Role</th
					>
					<th scope="col" class="px-3 py-3.5 font-semibold text-slate-900 sm:table-cell"
						>Recruiter</th
					>
					<th scope="col" class="hidden px-3 py-3.5 font-semibold text-slate-900 lg:table-cell"
						>Email</th
					>
					<th scope="col" class="px-3 py-3.5 font-semibold text-slate-900">Emailed At</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-slate-200 bg-white">
				{#each data.jobs as job}
					<tr class="text-sm text-slate-700">
						<td class="hidden py-4 pl-4 pr-3 font-medium text-slate-900 sm:pl-6 lg:table-cell">
							{job.company}
						</td>
						<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 lg:px-3">
							{job.job_title}
							<dl class="font-normal lg:hidden">
								<dt class="sr-only">Company</dt>
								<dd class="mt-1 truncate text-xs">{job.company}</dd>
							</dl>
						</td>
						<td
							class="w-full max-w-0 py-4 px-3 text-slate-900 sm:w-auto sm:max-w-none lg:text-slate-700"
						>
							{job.recruiter}
							<dl class="font-normal lg:hidden">
								<dt class="sr-only">Email</dt>
								<dd class="mt-1 truncate text-xs">{job.recruiter_email}</dd>
							</dl>
						</td>
						<td class="hidden px-3 py-4 lg:table-cell">{job.recruiter_email} </td>
						<td class="px-3 py-4">{job.emailed_at}</td>
					</tr>
				{/each}
			</tbody>
		</table>
		<div
			class="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3.5 sm:px-6"
		>
			<div class="flex flex-1 justify-end sm:hidden">
				<a
					href="#"
					class="relative inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
					>Previous</a
				>
				<a
					href="#"
					class="relative ml-3 inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50"
					>Next</a
				>
			</div>
			<div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
				<div>
					<p class="text-sm text-slate-700">
						Showing
						<span class="font-medium">{showingStart}</span>
						to
						<span class="font-medium">{showingEnd}</span>
						of
						<span class="font-medium">{data.pagination.numResults}</span>
						results
					</p>
				</div>
				<div>
					<nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
						<a
							href="/account/jobs?page={prevPage}"
							class:focus:z-20={data.pagination.hasPrev}
							class:hover:bg-slate-50={data.pagination.hasPrev}
							class:cursor-not-allowed={!data.pagination.hasPrev}
							class="relative inline-flex items-center rounded-l-md border border-slate-300 bg-white px-2 py-2 text-sm font-medium text-slate-500"
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
							{@const current = page === `${data.pagination.page}`}
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
							disabled={!data.pagination.hasNext}
							class="relative inline-flex items-center rounded-r-md border border-slate-300 bg-white px-2 py-2 text-sm font-medium text-slate-500"
							class:focus:z-20={data.pagination.hasNext}
							class:hover:bg-slate-50={data.pagination.hasNext}
							class:cursor-not-allowed={!data.pagination.hasNext}
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
	</div>
</div>
