<script lang="ts">
	import { fade } from 'svelte/transition';

	export let title: string;
	export let description: string;
	export let cta: string;
	export let show: boolean = false;
	export let onConfirm: () => void;

	const close = () => {
		show = false;
	};
</script>

{#if show}
	<div
		class="relative z-50"
		aria-labelledby="modal-title"
		role="dialog"
		aria-modal="true"
		in:fade={{ duration: 200 }}
		out:fade={{ delay: 150, duration: 200 }}
	>
		<!--
    Background backdrop, show/hide based on modal state.

    Entering: "ease-out duration-300"
      From: "opacity-0"
      To: "opacity-100"
    Leaving: "ease-in duration-200"
      From: "opacity-100"
      To: "opacity-0"
  -->
		<div
			class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
			class:ease-out={show}
			class:duration-300={show}
			class:opacity-0={!show}
			class:opacity-100={show}
			class:ease-in={!show}
			class:duration-200={!show}
		/>

		<div class="fixed inset-0 z-10 overflow-y-auto">
			<div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
				<!--
        Modal panel, show/hide based on modal state.

        Entering: "ease-out duration-300"
          From: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          To: "opacity-100 translate-y-0 sm:scale-100"
        Leaving: "ease-in duration-200"
          From: "opacity-100 translate-y-0 sm:scale-100"
          To: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
      -->
				<div
					class="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6"
					class:ease-out={show}
					class:duration-300={show}
					class:opacity-0={!show}
					class:opacity-100={show}
					class:translate-y-4={!show}
					class:translate-y-0={show}
					class:scale-95={!show}
					class:scale-100={show}
					class:ease-in={!show}
					class:duration-200={!show}
				>
					<div class="sm:flex sm:items-start">
						<div
							class="mx-auto flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-rose-100 sm:mx-0 sm:h-10 sm:w-10"
						>
							<!-- Heroicon name: outline/exclamation-triangle -->
							<svg
								class="h-6 w-6 text-rose-600"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								aria-hidden="true"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M12 10.5v3.75m-9.303 3.376C1.83 19.126 2.914 21 4.645 21h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 4.88c-.866-1.501-3.032-1.501-3.898 0L2.697 17.626zM12 17.25h.007v.008H12v-.008z"
								/>
							</svg>
						</div>
						<div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
							<h3 class="text-lg font-medium leading-6 text-gray-900" id="modal-title">
								{title}
							</h3>
							<div class="mt-2">
								<p class="text-sm text-gray-500">{description}</p>
							</div>
						</div>
					</div>
					<div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
						<button
							type="button"
							class="inline-flex w-full justify-center rounded-md border border-transparent bg-rose-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-rose-700 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm"
							on:click={onConfirm}>{cta}</button
						>
						<button
							type="button"
							class="mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:w-auto sm:text-sm"
							on:click={close}>Cancel</button
						>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
