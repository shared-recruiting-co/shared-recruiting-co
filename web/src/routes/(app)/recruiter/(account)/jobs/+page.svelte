<script lang="ts">
	import type { PageData } from './$types';

	import TableFooter from '$lib/components/TableFooter.svelte';

	export let data: PageData;
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
						<th scope="col" class="px-3 py-3.5" />
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-200 bg-white">
					{#each data.jobs as job}
						<tr class="text-sm text-slate-700">
							<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 sm:pl-6">
								{job.title}
							</td>
							<td class="px-3 py-4">{job.candidate_count}</td>
							<td class="px-3 py-4">{new Date(job.updated_at).toLocaleDateString()} </td>
							<td class="px-3 py-4">
								<a href="/recruiter/jobs/{job.job_id}" class="text-blue-500 hover:text-blue-600"
									>View</a
								>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			<TableFooter pagination={data.pagination}/>
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
