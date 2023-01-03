<script lang="ts">
	import Markdoc from '@markdoc/markdoc';

	import source from './privacy-policy.md?raw';

	const ast = Markdoc.parse(source);
	const content = Markdoc.transform(ast);

	const html = Markdoc.renderers.html(content);
</script>

<div class="flex h-screen w-full flex-col items-center overflow-auto">
	<h1 class="my-4 text-center text-4xl sm:my-12 sm:text-5xl">Privacy Policy</h1>

	<article class="prose">
		{@html html}
	</article>
</div>
