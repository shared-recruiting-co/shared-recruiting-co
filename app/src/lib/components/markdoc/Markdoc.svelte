<script lang="ts">
	// Adapted from
	// https://github.com/yuchengkuo/sveltejs-markdoc
	import markdoc from '@markdoc/markdoc';
	const { parse, transform } = markdoc;

	import config from './config';
	import Tag from './Tag.svelte';

	import Callout from './Callout.svelte';
	import EmailLabel from './EmailLabel.svelte';

	const components = {
		Callout,
		EmailLabel
	};

	// props
	// pass in raw markdown
	export let source: string;

	$: ast = parse(source);
	$: content = transform(ast, config(ast));
</script>

<Tag {content} {components} />
