<script lang="ts">
	export let label: string;
	export let checked: boolean;
	export let disabled: boolean = false;

	const toggle = () => {
		if (disabled) return;
		checked = !checked;
	};
</script>

<!-- Enabled: "bg-blue-600", Not Enabled: "bg-slate-200" -->
<button
	type="button"
	class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
	class:bg-blue-600={checked}
	class:bg-slate-200={!checked}
	class:cursor-not-allowed={disabled}
	role="switch"
	{disabled}
	aria-checked={checked}
	on:click={toggle}
>
	<span class="sr-only">{label}</span>
	<!-- Enabled: "translate-x-5", Not Enabled: "translate-x-0" -->
	<span
		class="pointer-events-none relative inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
		class:translate-x-5={checked}
		class:translate-x-0={!checked}
	>
		<!-- Enabled: "opacity-0 ease-out duration-100", Not Enabled: "opacity-100 ease-in duration-200" -->
		<span
			class="absolute inset-0 flex h-full w-full items-center justify-center transition-opacity"
			class:opacity-0={checked}
			class:ease-out={checked}
			class:duration-100={checked}
			class:opacity-100={!checked}
			class:ease-in={!checked}
			class:duration-200={!checked}
			aria-hidden={checked}
		>
			<svg class="h-3 w-3 text-slate-400" fill="none" viewBox="0 0 12 12">
				<path
					d="M4 8l2-2m0 0l2-2M6 6L4 4m2 2l2 2"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
			</svg>
		</span>
		<!-- Enabled: "opacity-100 ease-in duration-200", Not Enabled: "opacity-0 ease-out duration-100" -->
		<span
			class="absolute inset-0 flex h-full w-full items-center justify-center transition-opacity"
			class:opacity-100={checked}
			class:ease-in={checked}
			class:duration-200={checked}
			class:opacity-0={!checked}
			class:ease-out={!checked}
			class:duration-100={!checked}
			aria-hidden={!checked}
		>
			<svg class="h-3 w-3 text-indigo-600" fill="currentColor" viewBox="0 0 12 12">
				<path
					d="M3.707 5.293a1 1 0 00-1.414 1.414l1.414-1.414zM5 8l-.707.707a1 1 0 001.414 0L5 8zm4.707-3.293a1 1 0 00-1.414-1.414l1.414 1.414zm-7.414 2l2 2 1.414-1.414-2-2-1.414 1.414zm3.414 2l4-4-1.414-1.414-4 4 1.414 1.414z"
				/>
			</svg>
		</span>
	</span>
</button>
