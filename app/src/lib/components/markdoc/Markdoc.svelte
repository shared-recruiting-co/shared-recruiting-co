<script lang="ts">
	// Adapted from
	// https://github.com/yuchengkuo/sveltejs-markdoc
	import { parse, transform } from '@markdoc/markdoc';

	import config from './config';
	import Tag from './Tag.svelte';
	import Callout from './Callout.svelte';

	const components = {
		Callout
	};

	// props
	// pass in raw markdown
	export let source: string;

	$: ast = parse(source);
	$: content = transform(ast, config(ast));
</script>

<Tag {content} {components} />
