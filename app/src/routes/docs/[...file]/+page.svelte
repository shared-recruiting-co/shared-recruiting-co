<script lang="ts">
	import { page } from '$app/stores';

	import Markdoc from '@markdoc/markdoc';

	const render = (source: string) => {
		const ast = Markdoc.parse(source);
		const content = Markdoc.transform(ast);

		return Markdoc.renderers.html(content);
	};
</script>

<div class="flex w-full flex-col items-center">
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
		{@html render($page.data.raw)}
	</article>
</div>

<style lang="postcss">
	.content {
		@apply prose prose-slate max-w-none dark:text-slate-400 dark:prose-invert;
		/* headings */
		/* omitted: prose-headings:font-display */
		@apply prose-headings:scroll-mt-28 prose-headings:font-normal lg:prose-headings:scroll-mt-[8.5rem];
		/* lead */
		@apply prose-lead:text-slate-500 dark:prose-lead:text-slate-400;
		/* links */
		@apply prose-a:font-semibold dark:prose-a:text-sky-400;
		/* link underline */
		@apply prose-a:no-underline prose-a:shadow-[inset_0_-2px_0_0_var(--tw-prose-background,#fff),inset_0_calc(-1*(var(--tw-prose-underline-size,4px)+2px))_0_0_var(--tw-prose-underline,theme(colors.sky.300))] hover:prose-a:[--tw-prose-underline-size:6px] dark:[--tw-prose-background:theme(colors.slate.900)] dark:prose-a:shadow-[inset_0_calc(-1*var(--tw-prose-underline-size,2px))_0_0_var(--tw-prose-underline,theme(colors.sky.800))] dark:hover:prose-a:[--tw-prose-underline-size:6px];
		/* pre */
		@apply prose-pre:rounded-xl prose-pre:bg-slate-900 prose-pre:shadow-lg dark:prose-pre:bg-slate-800/60 dark:prose-pre:shadow-none dark:prose-pre:ring-1 dark:prose-pre:ring-slate-300/10;
		/* hr */
		@apply dark:prose-hr:border-slate-800;
	}
</style>
