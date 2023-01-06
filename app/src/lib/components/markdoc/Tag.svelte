<script lang="ts">
	import type { RenderableTreeNodes } from '@markdoc/markdoc';
	import type { ComponentsMap } from './types';
	import type { Tag } from '@markdoc/markdoc';

	const isTag = (tag: any): tag is Tag => {
		return !!(tag?.$$mdtype === 'Tag');
	};

	export let content: RenderableTreeNodes;
	export let components: ComponentsMap = {};

	if (
		!Array.isArray(content) &&
		isTag(content) &&
		typeof content.name === 'string' &&
		content.name.at(0)?.match(/[A-Z]/) &&
		!components[content.name]
	)
		console.warn(
			`${content.name} seems like a Svelte component but not provided in component props.`
		);
</script>

{#if content !== null && typeof content === 'object'}
	<!-- RenderableTreeNodes[] -->
	{#if Array.isArray(content)}
		{#each content as node}
			<svelte:self content={node} {components} />
		{/each}

		<!-- RenderableTreeNode -->
	{:else if isTag(content)}
		{@const { name, attributes = {}, children = [] } = content}

		<!-- Svelte components -->
		{#if components[name] || typeof name !== 'string'}
			{@const component = typeof name === 'function' ? name : components[name]}
			{#if children.length === 0}
				<svelte:component this={component} {...attributes} />
			{:else}
				<svelte:component this={component} {...attributes}>
					{#each children as child}
						{#if typeof child === 'string' || typeof child === 'number'}
							{child}
						{:else}
							<svelte:self content={child} {components} />
						{/if}
					{/each}
				</svelte:component>
			{/if}

			<!-- other tags -->
		{:else if children.length === 0}
			<svelte:element this={name} {...attributes} />
		{:else}
			<svelte:element this={name} {...attributes}>
				{#each children as child}
					{#if typeof child === 'string' || typeof child === 'number'}
						{child}
					{:else}
						<svelte:self content={child} {components} />
					{/if}
				{/each}
			</svelte:element>
		{/if}
	{/if}
{/if}
