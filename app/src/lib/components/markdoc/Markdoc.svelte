<script lang="ts">
	// Adapted from
	// https://github.com/yuchengkuo/sveltejs-markdoc
	import { parse, transform } from '@markdoc/markdoc';

	import Tag from './Tag.svelte';
	import Callout from './Callout.svelte';
	import type { Config } from '@markdoc/markdoc';

	const components = {
		Callout
	};
	const config: Config = {
		tags: {
			callout: {
				render: Callout,
				attributes: {
					type: {
						type: String,
						default: 'note',
						matches: ['caution', 'check', 'note', 'warning']
					}
				}
			}
		}
	};

	// props
	// pass in raw markdown
	export let source: string;

	$: ast = parse(source);
	$: content = transform(ast, config);
</script>

<Tag {content} {components} />
