<script lang="ts">
	// TODO: Type
	export let candidates: any[];
</script>

<div>
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
</div>
