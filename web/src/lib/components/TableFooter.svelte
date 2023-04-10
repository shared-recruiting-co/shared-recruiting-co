<script lang="ts">
    import type { Pagination } from '$lib/pagination';

    export let pagination: Pagination;
</script>

{#if pagination.pagesCount > 1}
<div
    class="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3.5 sm:px-6"
>
    <nav class="flex flex-1 justify-end md:hidden" data-sveltekit-noscroll>
        <a
            href="/account/jobs?page={pagination.prevPage}"
            class="relative inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
            class:hover:bg-slate-50={pagination.prevPageValid}
            class:cursor-not-allowed={pagination.prevPageValid}>Previous</a
        >
        <a
            href="/account/jobs?page={pagination.nextPage}"
            class="relative ml-3 inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
            class:hover:bg-slate-50={pagination.nextPageValid}
            class:cursor-not-allowed={!pagination.nextPageValid}>Next</a
        >
    </nav>
    <div class="hidden sm:flex-1 sm:items-center sm:justify-between md:flex">
        <div>
            <p class="text-sm text-slate-700">
                Showing
                <span class="font-semibold">{pagination.resultShowingFirst}</span>
                to
                <span class="font-semibold">{pagination.resultShowingLast}</span>
                of
                <span class="font-semibold">{pagination.resultsCount}</span>
                results
            </p>
        </div>
        <div>
            <nav
                class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                aria-label="Pagination"
                data-sveltekit-noscroll
            >
                <a
                    href="/account/jobs?page={pagination.prevPage}"
                    class:focus:z-20={pagination.prevPageValid}
                    class:hover:bg-slate-50={pagination.prevPageValid}
                    class:cursor-not-allowed={!pagination.prevPageValid}
                    class="relative inline-flex items-center rounded-l-md border border-slate-300 bg-white px-2 py-2 text-sm font-semibold text-slate-500"
                >
                    <span class="sr-only">Previous</span>
                    <!-- Heroicon name: mini/chevron-left -->
                    <svg
                        class="h-5 w-5"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                        aria-hidden="true"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z"
                            clip-rule="evenodd"
                        />
                    </svg>
                </a>
                {#each pagination.pagesDisplayArray as page}
                    {@const current = page === `${pagination.currentResultsPage}`}
                    {#if page === '...'}
                        <span
                            class="relative inline-flex items-center border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700"
                            >...</span
                        >
                    {:else}
                        <a
                            href="/account/jobs?page={page}"
                            class="relative inline-flex items-center border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-500 focus:z-20"
                            class:hover:bg-blue-50={!current}
                            class:z-10={current}
                            class:bg-blue-50={current}
                            class:border-blue-500={current}
                            class:text-blue-600={current}>{page}</a
                        >
                    {/if}
                {/each}
                <a
                    href="/account/jobs?page={pagination.nextPage}"
                    class="relative inline-flex items-center rounded-r-md border border-slate-300 bg-white px-2 py-2 text-sm font-medium text-slate-500"
                    class:focus:z-20={pagination.nextPageValid}
                    class:hover:bg-slate-50={pagination.nextPageValid}
                    class:cursor-not-allowed={!pagination.nextPageValid}
                >
                    <span class="sr-only">Next</span>
                    <!-- Heroicon name: mini/chevron-right -->
                    <svg
                        class="h-5 w-5"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                        aria-hidden="true"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z"
                            clip-rule="evenodd"
                        />
                    </svg>
                </a>
            </nav>
        </div>
    </div>
</div>
{/if}