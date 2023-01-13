<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;
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
							{job.title}
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
	</div>
</div>
