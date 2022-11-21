<script lang="ts">
	import { marked } from 'marked';

	import content from './privacy-policy.md?raw';
</script>

<div class="flex h-screen w-full flex-col items-center overflow-scroll">
	<h1 class="my-4 text-center text-4xl sm:my-12 sm:text-5xl">Privacy Policy</h1>

	<article class="prose pt-8">
		{@html marked(content)}
	</article>
</div>
