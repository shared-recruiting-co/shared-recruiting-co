<script lang="ts">
	// TODO: Type
	export let outboundTemplates: any[];
</script>

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
					<th
						scope="col"
						class="table-cell py-3.5 pl-4 pr-3 font-semibold text-slate-900 lg:px-3"
					/>
				</tr>
			</thead>
			<tbody class="divide-y divide-slate-200 bg-white">
				{#each outboundTemplates as outbound}
					<tr class="text-sm text-slate-700">
						<td class="hidden py-4 pl-4 pr-3 font-medium text-slate-900 sm:pl-6 lg:table-cell">
							{outbound.subject}
						</td>
						<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 lg:px-3">
							{outbound.body.slice(0, 40)}...
						</td>
						<td class="table-cell py-4 pl-4 pr-3 font-medium text-slate-900 lg:px-3">
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

								<span> 2 </span>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
