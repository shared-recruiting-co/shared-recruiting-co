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

	// TOOD:
	// Error state (icon + No jobs + Description)
	// Realtime
	// Show Actions (Interested, Not Interested, Remove)
	// Show Company Favicon/Logo
</script>

<div>
	<h1 class="text-3xl sm:text-4xl">Jobs</h1>
	<p class="mt-1 max-w-3xl text-sm text-slate-500">Review your inbound job opportunity history.</p>
	<p class="mt-2 max-w-3xl text-sm text-slate-500">
		Inbound recruiting emails are automatically parsed and added to your job board. Only emails that
		describe the company and role will be added to your job board.
	</p>
</div>
{#if data.jobs && data.jobs.length > 0}
	<div>
		<div class="mt-8 overflow-hidden rounded-lg shadow ring-1 ring-black ring-opacity-5 md:mx-0">
			<table class="min-w-full divide-y divide-slate-300">
				<thead class="bg-slate-50 text-left text-sm">
					<tr>
						<th
							scope="col"
							class="py-3.5 pl-4 pr-3 font-semibold text-slate-900 sm:pl-6 lg:table-cell">Role</th
						>
						<th scope="col" class="px-3 py-3.5 font-semibold text-slate-900 sm:table-cell"
							>Recruiter</th
						>
						<th scope="col" class="px-3 py-3.5 font-semibold text-slate-900" />
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-200 bg-white">
					{#each data.jobs as job}
						<tr class="text-sm text-slate-700" class:bg-blue-50={job.is_verified}>
							<td class="table-cell py-4 pl-4 pr-3 text-base font-medium text-slate-900 lg:pl-6">
								<div class="flex items-center space-x-4">
									<div
										class="flex h-10 w-10 items-center justify-center rounded-full bg-indigo-100"
										alt="Company Intial"
									>
										A
									</div>
									<div>
										<div class="flex items-center space-x-2">
											<span>
												{job.company}
											</span>
											{#if job.is_verified}
												<svg
													xmlns="http://www.w3.org/2000/svg"
													viewBox="0 0 24 24"
													fill="currentColor"
													class="h-5 w-5 text-blue-500"
												>
													<path
														fill-rule="evenodd"
														d="M8.603 3.799A4.49 4.49 0 0112 2.25c1.357 0 2.573.6 3.397 1.549a4.49 4.49 0 013.498 1.307 4.491 4.491 0 011.307 3.497A4.49 4.49 0 0121.75 12a4.49 4.49 0 01-1.549 3.397 4.491 4.491 0 01-1.307 3.497 4.491 4.491 0 01-3.497 1.307A4.49 4.49 0 0112 21.75a4.49 4.49 0 01-3.397-1.549 4.49 4.49 0 01-3.498-1.306 4.491 4.491 0 01-1.307-3.498A4.49 4.49 0 012.25 12c0-1.357.6-2.573 1.549-3.397a4.49 4.49 0 011.307-3.497 4.49 4.49 0 013.497-1.307zm7.007 6.387a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z"
														clip-rule="evenodd"
													/>
												</svg>
											{/if}
										</div>
										<dl class="text-xs sm:text-sm">
											<dt class="sr-only">Job Title</dt>
											<dd class="mt-1 truncate">{job.job_title}</dd>
										</dl>
									</div>
								</div>
							</td>
							<td class="w-full max-w-0 py-4 px-3 text-base text-slate-900 sm:w-auto sm:max-w-none">
								{job.recruiter}
								<dl class="text-xs sm:text-sm">
									<dt class="sr-only">Recruiter Email</dt>
									<dd class="mt-1 truncate">{job.recruiter_email}</dd>
								</dl>
							</td>
							<td class="px-3 py-4 text-sm">
								<div class="space-y-1">
									{#if job.company_website}
										<div class="flex items-center space-x-2">
											<svg
												xmlns="http://www.w3.org/2000/svg"
												fill="none"
												viewBox="0 0 24 24"
												stroke-width="1.5"
												stroke="currentColor"
												class="h-5 w-5"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418"
												/>
											</svg>
											<span> Company Website </span>
										</div>
									{/if}
									{#if job.job_description_url}
										<div class="flex items-center space-x-2">
											<svg
												xmlns="http://www.w3.org/2000/svg"
												fill="none"
												viewBox="0 0 24 24"
												stroke-width="1.5"
												stroke="currentColor"
												class="h-5 w-5"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													d="M20.25 14.15v4.25c0 1.094-.787 2.036-1.872 2.18-2.087.277-4.216.42-6.378.42s-4.291-.143-6.378-.42c-1.085-.144-1.872-1.086-1.872-2.18v-4.25m16.5 0a2.18 2.18 0 00.75-1.661V8.706c0-1.081-.768-2.015-1.837-2.175a48.114 48.114 0 00-3.413-.387m4.5 8.006c-.194.165-.42.295-.673.38A23.978 23.978 0 0112 15.75c-2.648 0-5.195-.429-7.577-1.22a2.016 2.016 0 01-.673-.38m0 0A2.18 2.18 0 013 12.489V8.706c0-1.081.768-2.015 1.837-2.175a48.111 48.111 0 013.413-.387m7.5 0V5.25A2.25 2.25 0 0013.5 3h-3a2.25 2.25 0 00-2.25 2.25v.894m7.5 0a48.667 48.667 0 00-7.5 0M12 12.75h.008v.008H12v-.008z"
												/>
											</svg>
											<span> Job Description</span>
										</div>
									{/if}
									<div class="flex items-center space-x-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="1.5"
											stroke="currentColor"
											class="hidden h-5 w-5 sm:block"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5m-9-6h.008v.008H12v-.008zM12 15h.008v.008H12V15zm0 2.25h.008v.008H12v-.008zM9.75 15h.008v.008H9.75V15zm0 2.25h.008v.008H9.75v-.008zM7.5 15h.008v.008H7.5V15zm0 2.25h.008v.008H7.5v-.008zm6.75-4.5h.008v.008h-.008v-.008zm0 2.25h.008v.008h-.008V15zm0 2.25h.008v.008h-.008v-.008zm2.25-4.5h.008v.008H16.5v-.008zm0 2.25h.008v.008H16.5V15z"
											/>
										</svg>
										<span>
											{job.emailed_at}
										</span>
									</div>
								</div>
							</td>
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
		<p class="slate-500 mt-4 max-w-xs text-center text-sm">
			Once your job emails are processed they will appear in your job board.
		</p>
	</div>
{/if}
