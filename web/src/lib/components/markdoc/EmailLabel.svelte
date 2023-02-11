<script lang="ts">
	export let title = '';
	export let color: 'blue' | 'green' | 'red' = 'blue';
</script>

<span
	class="inline-flex items-center rounded-md px-1.5 py-0.5 text-sm font-medium text-white"
	class:bg-blue-600={color === 'blue'}
	class:bg-red-600={color === 'red'}
	class:bg-green-600={color === 'green'}
>
	{title}
</span>
