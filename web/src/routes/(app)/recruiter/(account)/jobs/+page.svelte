<script lang="ts">
	import type { PageData } from './$types';

	import { getPaginationPages } from '$lib/pagination';

	export let data: PageData;

	$: showingStart = (data.pagination.page - 1) * data.pagination.perPage + 1;
	$: showingEnd = Math.min(
		data.pagination.page * data.pagination.perPage,
		data.pagination.numResults
	);
	$: prevPage = data.pagination.hasPrev ? data.pagination.page - 1 : data.pagination.page;
	$: nextPage = data.pagination.hasNext ? data.pagination.page + 1 : data.pagination.page;

	$: pages = getPaginationPages(data.pagination.page, data.pagination.numPages);
</script>

<div>
	<h1 class="text-3xl sm:text-4xl">Jobs</h1>
	<p class="mt-1 max-w-3xl text-sm text-slate-500" />
</div>
{#if data.jobs && data.jobs.length > 0}
	<div>
		<div class="flex flex-row items-center justify-end">
			<a
				href="/recruiter/jobs/new"
				class="mt-4 flex items-center rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="mr-1 h-4 w-4"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m6-6H6" />
				</svg>

				<span>Add Job</span>
			</a>
		</div>
		<div class="mt-4 overflow-hidden rounded-lg shadow ring-1 ring-black ring-opacity-5 md:mx-0">
			<table class="min-w-full divide-y divide-slate-300">
				<thead class="bg-slate-50 text-left text-sm">
					<tr>
						<th
							scope="col"
							class="hidden py-3.5 pl-4 pr-3 font-semibold text-slate-900 sm:pl-6 lg:table-cell"
							>Role</th
						>
						<th scope="col" class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
							>Candidates</th
						>
						<th scope="col" class="px-3 py-3.5 font-semibold text-slate-900">Date</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-200 bg-white">
					{#each data.jobs as job}
						<tr class="text-sm text-slate-700">
							<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 sm:pl-6">
								{job.title}
							</td>
							<td class="px-3 py-4">0</td>
							<td class="px-3 py-4">{new Date(job.updated_at).toLocaleDateString()}</td>
						</tr>
					{/each}
				</tbody>
			</table>
			{#if data.pagination.numPages > 1}
				<div
					class="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3.5 sm:px-6"
				>
					<nav class="flex flex-1 justify-end md:hidden" data-sveltekit-noscroll>
						<a
							href="/account/jobs?page={prevPage}"
							class="relative inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
							class:hover:bg-slate-50={data.pagination.hasPrev}
							class:cursor-not-allowed={!data.pagination.hasPrev}>Previous</a
						>
						<a
							href="/account/jobs?page={nextPage}"
							class="relative ml-3 inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
							class:hover:bg-slate-50={data.pagination.hasNext}
							class:cursor-not-allowed={!data.pagination.hasNext}>Next</a
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
								<span class="font-semibold">{data.pagination.numResults}</span>
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
									class:focus:z-20={data.pagination.hasPrev}
									class:hover:bg-slate-50={data.pagination.hasPrev}
									class:cursor-not-allowed={!data.pagination.hasPrev}
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
			{/if}
		</div>
	</div>
{:else}
	<!-- Empty state -->
	<div class="flex flex-1 flex-col items-center justify-center py-20">
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
				d="M20.25 14.15v4.25c0 1.094-.787 2.036-1.872 2.18-2.087.277-4.216.42-6.378.42s-4.291-.143-6.378-.42c-1.085-.144-1.872-1.086-1.872-2.18v-4.25m16.5 0a2.18 2.18 0 00.75-1.661V8.706c0-1.081-.768-2.015-1.837-2.175a48.114 48.114 0 00-3.413-.387m4.5 8.006c-.194.165-.42.295-.673.38A23.978 23.978 0 0112 15.75c-2.648 0-5.195-.429-7.577-1.22a2.016 2.016 0 01-.673-.38m0 0A2.18 2.18 0 013 12.489V8.706c0-1.081.768-2.015 1.837-2.175a48.111 48.111 0 013.413-.387m7.5 0V5.25A2.25 2.25 0 0013.5 3h-3a2.25 2.25 0 00-2.25 2.25v.894m7.5 0a48.667 48.667 0 00-7.5 0M12 12.75h.008v.008H12v-.008z"
			/>
		</svg>

		<p class="mt-2 text-base font-medium">No jobs found.</p>
		<a
			href="/recruiter/jobs/new"
			class="mt-4 flex items-center rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="mr-1 h-4 w-4"
			>
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m6-6H6" />
			</svg>

			<span>Add Job</span>
		</a>
	</div>
{/if}
