<script lang="ts">
	import type { Database } from '$lib/supabase/types';

	export let candidates: Database['public']['Tables']['candidate_company_inbound']['Row'][];
</script>

<div>
	<!-- Empty State -->
	{#if candidates.length === 0}
		<div class="flex flex-col items-center justify-center p-8">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="h-12 w-12 text-slate-400"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"
				/>
			</svg>

			<p class="mt-2 font-medium text-slate-900">No Candidates</p>
			<p class="mt-1 max-w-md text-center text-sm text-slate-500">
				SRC automatically imports candidates you reach out. Any candidate you reach out to with an
				outboudn template associated with the job will show up here.
			</p>
		</div>
	{:else}
		<div class="mt-8 overflow-hidden rounded-lg shadow ring-1 ring-black ring-opacity-5 md:mx-0">
			<table class="min-w-full divide-y divide-slate-300">
				<thead class="bg-slate-50 text-left text-sm">
					<tr>
						<th
							scope="col"
							class="hidden py-3.5 pl-4 pr-3 font-semibold text-slate-900 sm:pl-6 lg:table-cell"
							>Candidate</th
						>
						<th scope="col" class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
							>Last Activity</th
						>
						<th
							scope="col"
							class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
						/>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-200 bg-white">
					{#each candidates as candidate}
						<tr class="text-sm text-slate-700">
							<td class="hidden py-4 pl-4 pr-3 font-medium text-slate-900 sm:pl-6 lg:table-cell">
								{candidate.candidate_email}
							</td>
							<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 lg:px-3">
								{new Date(candidate.updated_at).toLocaleDateString()}
							</td>
							<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 lg:px-3"> todo </td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
