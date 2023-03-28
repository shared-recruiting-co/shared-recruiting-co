<script lang="ts">
	import type { Database } from '$lib/supabase/types';

	export let outboundTemplates: Database['public']['Tables']['recruiter_outbound_template']['Row'] &
		{
			recipient_count: number;
		}[];
</script>

<div>
	<!-- Empty State -->
	{#if outboundTemplates.length === 0}
		<div class="flex flex-col items-center justify-center p-8">
			<svg
				class="h-12 w-12 text-slate-400"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				aria-hidden="true"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
				/>
			</svg>
			<p class="mt-2 font-medium text-slate-900">No Outbound Templates Assigned</p>
			<p class="mt-1 max-w-md text-center text-sm text-slate-500">
				<a href="/recruiter/outbound" class="underline hover:text-blue-500"
					>Assign outbound templates</a
				> to this job to import candidates into SRC.
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
									<span>{outbound.recipient_count || 0}</span>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
