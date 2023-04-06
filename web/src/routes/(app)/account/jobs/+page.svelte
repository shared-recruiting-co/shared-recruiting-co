<script lang="ts">
	import type { PageData } from './$types';

	import { getPaginationPages } from '$lib/pagination';

	export let data: PageData;

	$: ({ supabase, profile, jobs } = data);

	$: showingStart = (data.pagination.page - 1) * data.pagination.perPage + 1;
	$: showingEnd = Math.min(
		data.pagination.page * data.pagination.perPage,
		data.pagination.numResults
	);
	$: prevPage = data.pagination.hasPrev ? data.pagination.page - 1 : data.pagination.page;
	$: nextPage = data.pagination.hasNext ? data.pagination.page + 1 : data.pagination.page;

	$: pages = getPaginationPages(data.pagination.page, data.pagination.numPages);

	/**
	 * Function handles removing a job from the job board, given a job, this wil:
	 * 	 - If there is an email associated with the job, remove any SRC gmail lables
	 * 	- Remove the job from Supabase table `user_email_job`
	 *  - Remove the job from the locally scoped `data.jobs` array so it is removed from the UI
	 * @param {string} jobId - The UUID of the job to be removed
	 * @returns {Promise<void>}
	 */
	const handleJobRemoval = async (jobId: string): Promise<void> => {

		// query Supabase to get the needed data from the table `user_email_job` before its deleted
		const { data: user_email_job_data } = await supabase
			.from('user_email_job')
			.select('email_thread_id, user_email')
			.eq('job_id', jobId)
			.maybeSingle()

		// get the email thread ID and user email for this job
		const email_thread_id = user_email_job_data?.email_thread_id
		const email = user_email_job_data?.user_email

		// delete the entry in the user_email_job table the corresponds to the selected job
		const { error: job_deletion_error } = await supabase
			.from('user_email_job')
			.delete()
			.eq('job_id', jobId)

		// if the the delete above is successful, update `jobs`
		if (!job_deletion_error) {

			// find index of the deleted job in the locally scoped jobs and remove it
			const jobIndex = jobs.findIndex((job: { job_id: string }) => job.job_id === jobId);

			// if the index exists, remove that job from the `jobs` array
			if (jobIndex !== -1) {
				jobs.splice(jobIndex, 1);
				jobs = { ...jobs };
			}
		}
		
		// if we can find the specific email assocaiated with this job, remove its SRC email labels
		if (email_thread_id && email) {

			// the remove-email-labels will attempt to remove any SRC labels from the associated email
			const resp = await fetch('/api/account/gmail/remove-email-labels', {
				method: 'POST',
				body: JSON.stringify({ email, email_thread_id })
			});

			// handle errors
			if (resp.status !== 200) {
				return;
			}
		}
		
	}

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
{#if jobs && jobs.length > 0}
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
						<th scope="col" class="px-3 py-3.5" />
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-200 bg-white">
					{#each jobs as job}
						<tr class="text-sm text-slate-700" class:bg-sky-50={job.is_verified}>
							<td class="table-cell py-4 pl-4 pr-3 text-base font-medium text-slate-900 lg:pl-6">
								<div class="flex items-center space-x-4">
									<div
										class="flex h-10 w-10 items-center justify-center rounded-full bg-orange-100"
										alt="Company Intial"
									>
										{job.company_name[0]?.toUpperCase() || '?'}
									</div>
									<div>
										<div class="flex items-center space-x-2">
											<span>
												{job.company_name}
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
								{job.recruiter_name}
								<dl class="text-xs sm:text-sm">
									<dt class="sr-only">Recruiter Email</dt>
									<dd class="mt-1 truncate">{job.recruiter_email}</dd>
								</dl>
							</td>
							<td class="px-3 py-4 text-sm">
								<div
									class="flex items-center justify-end space-x-2 space-y-0 lg:flex-col lg:items-start lg:justify-start lg:space-x-0 lg:space-y-2"
								>
									{#if job.company_website}
										<a
											href={job.company_website}
											target="_blank"
											rel="noopener noreferrer"
											class="flex items-center space-x-2 hover:text-blue-500"
										>
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
											<span class="hidden hover:underline lg:block">
												{new URL(job.company_website).hostname}
											</span>
										</a>
									{/if}
									{#if job.job_description_url}
										<a
											class="flex items-center space-x-2 hover:text-blue-500"
											href={job.job_description_url}
											target="_blank"
											rel="noopener noreferrer"
										>
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
											<span class="hidden hover:underline lg:block">
												{new URL(job.job_description_url).hostname}
											</span>
										</a>
									{/if}
									<div class="flex items-center space-x-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="1.5"
											stroke="currentColor"
											class="hidden h-5 w-5 lg:block"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5m-9-6h.008v.008H12v-.008zM12 15h.008v.008H12V15zm0 2.25h.008v.008H12v-.008zM9.75 15h.008v.008H9.75V15zm0 2.25h.008v.008H9.75v-.008zM7.5 15h.008v.008H7.5V15zm0 2.25h.008v.008H7.5v-.008zm6.75-4.5h.008v.008h-.008v-.008zm0 2.25h.008v.008h-.008V15zm0 2.25h.008v.008h-.008v-.008zm2.25-4.5h.008v.008H16.5v-.008zm0 2.25h.008v.008H16.5V15z"
											/>
										</svg>
										<span>
											{new Date(job.emailed_at).toLocaleDateString()}
										</span>
									</div>
								</div>
							</td>
							<td class="align-top px-3 py-4 flex justify-end">
								{#if !job.is_verified}
								<button 
									class="hover:text-red-500"
									title="Remove this job"
									on:click={handleJobRemoval(job.job_id)}
								>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none" 
											viewBox="0 0 24 24" 
											stroke-width="1.5" 
											stroke="currentColor" 
											class="w-4 h-4">
											<path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
										</svg> 
								</button>
								{/if}
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
