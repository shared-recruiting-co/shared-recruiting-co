<script lang="ts">
	export let hint: string;

	let isVisible = false;
</script>

<div on:mouseenter={() => (isVisible = true)} on:mouseleave={() => (isVisible = false)}>
	<svg
		xmlns="http://www.w3.org/2000/svg"
		fill="none"
		viewBox="0 0 24 24"
		stroke-width="1.5"
		stroke="currentColor"
		class="h-5 w-5 cursor-pointer text-slate-400 hover:text-slate-500"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z"
		/>
	</svg>
	{#if isVisible}
		<div
			role="tooltip"
			class="absolute z-50 inline-block max-w-xl rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-900 shadow-sm"
		>
			{@html hint}
		</div>
	{/if}
</div>
