<script lang="ts">
	import { slide } from 'svelte/transition';

	type FAQ = {
		question: string;
		answer: string;
	};
	export let faqs: FAQ[];

	let opened: Record<string, boolean> = {};

	const toggleVisibility = (key: string) => {
		opened[key] = !opened[key];
	};
</script>

<section id="faqs" class="mx-auto max-w-7xl py-12 px-4 sm:py-16 sm:px-6 lg:px-8">
	<div class="mx-auto max-w-3xl divide-y-2 divide-slate-200">
		<div>
			<h2 class="mt-8 mb-4 text-center text-3xl font-extrabold text-slate-900 sm:text-4xl">
				Frequently Asked Questions
			</h2>
			<p class="text-center text-sm text-slate-500">
				For any additional questions, please reach out to us at <a
					class="link"
					href="mailto:team@sharedrecruiting.co">team@sharedrecruiting.co</a
				>. We'd love to chat!
			</p>
		</div>
		<dl class="mt-6 space-y-6 divide-y divide-slate-200">
			{#each faqs as faq}
				<div class="pt-6">
					<dt class="text-lg">
						<!-- Expand/collapse question button -->
						<button
							type="button"
							class="flex w-full items-start justify-between text-left text-slate-400"
							aria-controls={`faq-${faq.question}`}
							aria-expanded={opened[faq.question]}
							on:click={(e) => {
								e.preventDefault();
								toggleVisibility(faq.question);
							}}
						>
							<span class="font-medium text-slate-900">
								{faq.question}
							</span>
							<span class="ml-6 flex h-7 items-center">
								<!--
                  Expand/collapse icon, toggle classes based on question open state.

                  Heroicon name: outline/chevron-down

                  Open: "-rotate-180", Closed: "rotate-0"
                -->
								<svg
									class="h-6 w-6 transform"
									class:-rotate-180={opened[faq.question]}
									class:rotate-0={!opened[faq.question]}
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2"
									stroke="currentColor"
									aria-hidden="true"
								>
									<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
								</svg>
							</span>
						</button>
					</dt>
					{#if opened[faq.question]}
						<dd class="mt-2 pr-12" id={`faq-${faq.question}`} transition:slide>
							<p class="text-base text-slate-500">
								{@html faq.answer}
							</p>
						</dd>
					{/if}
				</div>
			{/each}
		</dl>
	</div>
</section>
