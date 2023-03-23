<script lang="ts">
	import type { PageData } from ',/$types';
	import { page } from '$app/stores';

	export let data: PageData;

	$: ({ job } = data);
	$: hash = $page.url.hash;

	//

	// TODO
	// Breadcrumbs
	// Edit Job
</script>

<div>
	<h1 class="text-3xl sm:text-4xl">{job.title}</h1>
	<p class="mt-2 text-sm text-slate-500">
		<a
			href={job.description_url}
			target="_blank"
			rel="noopener noreferrer"
			class="group flex items-center text-slate-500 hover:text-blue-600 hover:underline"
		>
			<span>
				{job.description_url}
			</span>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1"
				stroke="currentColor"
				class="ml-2 h-4 w-4 group-hover:text-blue-600"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25"
				/>
			</svg>
		</a>
	</p>
</div>
<div>
	<div>
		<div class="border-b border-slate-200 md:text-lg">
			<nav class="-mb-px flex space-x-8" aria-label="Tabs">
				<!-- Current: "border-blue-500 text-blue-600", Default: "border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700" -->
				<a
					href="#candidates"
					class="group inline-flex items-center border-b-2 border-transparent py-4 px-1 text-sm font-medium md:text-base"
					class:text-slate-500={hash !== '#candidates' && hash !== ''}
					class:hover:border-slate-300={hash !== '#candidates' && hash !== ''}
					class:hover:text-slate-700={hash !== '#candidates' && hash !== ''}
					class:border-blue-500={hash === '#candidates' || hash === ''}
					class:text-blue-600={hash === '#candidates' || hash === ''}
				>
					<!-- Current: "text-blue-500", Default: "text-slate-400 group-hover:text-slate-500" -->
					<svg
						class="-ml-0.5 mr-2 h-6 w-6"
						class:text-slate-400={hash !== '#candidates' && hash !== ''}
						class:group-hover:text-slate-500={hash !== '#candidates' && hash !== ''}
						class:text-blue-500={hash === '#candidates' || hash === ''}
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
					>
						<path
							d="M10 8a3 3 0 100-6 3 3 0 000 6zM3.465 14.493a1.23 1.23 0 00.41 1.412A9.957 9.957 0 0010 18c2.31 0 4.438-.784 6.131-2.1.43-.333.604-.903.408-1.41a7.002 7.002 0 00-13.074.003z"
						/>
					</svg>
					<span>Candidates</span>
				</a>

				<a
					href="#outbound"
					class="group inline-flex items-center border-b-2 border-transparent py-4 px-1 text-sm font-medium md:text-base"
					class:text-slate-500={hash !== '#outbound' && hash !== ''}
					class:hover:border-slate-300={hash !== '#outbound' && hash !== ''}
					class:hover:text-slate-700={hash !== '#outbound' && hash !== ''}
					class:border-blue-500={hash === '#outbound' || hash === ''}
					class:text-blue-600={hash === '#outbound' || hash === ''}
				>
					<svg
						class="-ml-0.5 mr-2 h-6 w-6"
						class:text-slate-400={hash !== '#outbound' && hash !== ''}
						class:group-hover:text-slate-500={hash !== '#outbound' && hash !== ''}
						class:text-blue-500={hash === '#outbound' || hash === ''}
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5"
						/>
						/>
					</svg>
					<span>Outbound Templates</span>
				</a>
			</nav>
		</div>
	</div>
</div>
