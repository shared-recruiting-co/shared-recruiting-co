<script lang="ts">
	import { page } from '$app/stores';
	import type { PageData } from './$types';

	import Markdoc from '$lib/components/markdoc/Markdoc.svelte';

	import { getSectionTitle, getLinkTitle } from '../navigation';

	export let data: PageData;

	$: section = getSectionTitle($page.url.pathname);
	$: title = getLinkTitle($page.url.pathname);

	// TODO: Truncate document title if too long for SEO
</script>

<svelte:head>
	<title>{title} | {section}</title>
</svelte:head>

<div class="min-w-0 max-w-2xl flex-auto px-4 py-16 lg:max-w-none lg:pr-0 lg:pl-8 xl:px-16">
	<!-- omitted: prose-headings:font-display -->
	<article
		class="
dark:prose-pre:ring-slate-300/10; prose prose-slate max-w-none prose-headings:scroll-mt-28
prose-headings:font-normal prose-a:font-semibold prose-a:no-underline
prose-a:shadow-[inset_0_-2px_0_0_var(--tw-prose-background,#fff),inset_0_calc(-1*(var(--tw-prose-underline-size,4px)+2px))_0_0_var(--tw-prose-underline,theme(colors.sky.300))] hover:prose-a:[--tw-prose-underline-size:6px]
prose-pre:rounded-xl prose-pre:bg-slate-900
prose-pre:shadow-lg prose-lead:text-slate-500 dark:text-slate-400 dark:prose-invert dark:[--tw-prose-background:theme(colors.slate.900)] dark:prose-a:text-sky-400
dark:prose-a:shadow-[inset_0_calc(-1*var(--tw-prose-underline-size,2px))_0_0_var(--tw-prose-underline,theme(colors.sky.800))] dark:hover:prose-a:[--tw-prose-underline-size:6px] dark:prose-pre:bg-slate-800/60 dark:prose-pre:shadow-none dark:prose-pre:ring-1 dark:prose-hr:border-slate-800 dark:prose-lead:text-slate-400
lg:prose-headings:scroll-mt-[8.5rem]
		"
	>
		{#if section}
			<p class="text-sm font-semibold text-sky-500">{section}</p>
		{/if}
		<Markdoc source={data.markdown.source} />
	</article>
</div>
