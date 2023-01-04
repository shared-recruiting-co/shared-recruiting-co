<script lang="ts">
	import Markdoc from '@markdoc/markdoc';

	import source from './docs.md?raw';

	const ast = Markdoc.parse(source);
	const content = Markdoc.transform(ast);

	const html = Markdoc.renderers.html(content);
</script>

<div class="flex h-screen w-full flex-col items-center overflow-auto">
	<article class="prose">
		{@html html}
	</article>
</div>
