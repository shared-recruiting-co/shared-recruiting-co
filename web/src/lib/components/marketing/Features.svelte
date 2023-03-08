<script lang="ts">
	type Feature = {
		icon: string;
		title: string;
		description: string;
	};
	export let features: Feature[];
	export let title: string;
	export let description: string;
</script>

<section
	id="product"
	class="flex min-h-screen flex-col items-center bg-slate-900 py-20 text-blue-300"
>
	<h2 class="max-w-4xl text-center text-4xl font-bold sm:text-5xl sm:leading-tight">
		{@html title}
	</h2>
	<h3 class="mt-8 mb-16 max-w-2xl px-4 text-center text-slate-100">
		{description}
	</h3>
	<div class="flex w-2/3 flex-wrap justify-between">
		{#each features as feature}
			<div
				class="mb-12 flex max-w-xs flex-col items-center space-y-4 text-center sm:items-start sm:text-left"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="text-secondary-focus h-12 w-12"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="1.5"
				>
					{@html feature.icon}
				</svg>
				<p class="text-secondary-focus text-xl font-semibold">{feature.title}</p>
				<p class="text-slate-100">
					{feature.description}
				</p>
			</div>
		{/each}
	</div>
</section>
