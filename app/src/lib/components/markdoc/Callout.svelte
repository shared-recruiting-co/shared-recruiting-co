<script lang="ts">
	export let title = '';
	export let type: 'caution' | 'check' | 'note' | 'warning' = 'note';
	console.log('callout');
</script>

<div class={type}>
	{#if title}
		<h4>{title}</h4>
	{/if}
	<slot />
</div>

<style>
	div {
		padding: 8px;
		background-color: whitesmoke;
	}

	.note {
		background-color: lightblue;
	}
</style>
