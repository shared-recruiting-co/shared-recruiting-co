<script lang="ts">
	export let title = '';
	export let type: 'caution' | 'check' | 'note' | 'warning' = 'note';
</script>

<div
	class="my-8 flex rounded-3xl bg-sky-50 p-4 dark:bg-slate-800/60 dark:ring-1 dark:ring-slate-300/10"
>
	<svg
		xmlns="http://www.w3.org/2000/svg"
		fill="none"
		viewBox="0 0 24 24"
		stroke-width="1.5"
		stroke="currentColor"
		class="h-8 w-8 text-sky-500 dark:text-slate-300"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			d="M12 18v-5.25m0 0a6.01 6.01 0 001.5-.189m-1.5.189a6.01 6.01 0 01-1.5-.189m3.75 7.478a12.06 12.06 0 01-4.5 0m3.75 2.383a14.406 14.406 0 01-3 0M14.25 18v-.192c0-.983.658-1.823 1.508-2.316a7.5 7.5 0 10-7.517 0c.85.493 1.509 1.333 1.509 2.316V18"
		/>
	</svg>

	<div class="ml-4 flex-auto">
		<p class="m-0 text-lg font-semibold text-sky-900 dark:text-sky-400">
			{title}
		</p>
		<div
			class="prose mt-2.5 text-sky-800 [--tw-prose-background:theme(colors.sky.50)] prose-a:text-sky-900 prose-code:text-sky-900 dark:text-slate-300 dark:prose-code:text-slate-300"
		>
			<slot />
		</div>
	</div>
</div>
