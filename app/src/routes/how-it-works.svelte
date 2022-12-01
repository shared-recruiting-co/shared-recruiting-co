<script lang="ts">
	let data = [
		{
			title: 'Detect & Archive',
			description:
				'SRC detects any inbound job opportunities. Each job opportunity is labeled, replied to, and archived for future reference.'
		},
		{
			title: 'Summarize & Standardize',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardize format. '
		},
		{
			title: 'Review & Engage',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardize format. '
		}
	];
</script>

<section id="process" class="mx-auto flex min-h-screen flex-col items-center py-20 text-slate-900">
	<svg
		xmlns="http://www.w3.org/2000/svg"
		fill="none"
		viewBox="0 0 24 24"
		stroke-width="1.5"
		stroke="currentColor"
		class="h-16 w-16 text-blue-400"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99"
		/>
	</svg>
	<h2 class="text-secondary my-12 max-w-3xl text-center text-4xl font-bold sm:my-16 sm:text-5xl">
		The SRC Process
	</h2>
	<div class="flex max-w-xl flex-col space-y-8 px-2 md:space-y-16">
		{#each data as item, index}
			<div class="flex flex-row items-center space-x-4 md:space-x-6 md:space-y-0">
				<div
					class="flex h-12 w-12 flex-col items-center justify-center rounded-full bg-blue-400 p-5 text-center text-xl font-bold text-white"
				>
					{index + 1}
				</div>
				<div class="flex flex-col">
					<div class="text-xl font-semibold">{item.title}</div>
					<div class="text-base md:text-lg">
						{item.description}
					</div>
				</div>
			</div>
		{/each}
	</div>
</section>
