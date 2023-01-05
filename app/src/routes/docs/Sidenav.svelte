<script lang="ts">
	import { page } from '$app/stores';

	const nav = [
		{
			title: 'Introduction',
			links: [
				{ title: 'Welcome', href: '/docs/welcome' },
				{ title: 'Open Source', href: '/docs/open-source' },
				{ title: 'Security & Privacy', href: '/docs/security-privacy' }
			]
		},
		{
			title: 'Core Concepts',
			links: [
				{ title: 'Connect Email', href: '/docs/connect-email' },
				{
					title: 'Email Labels',
					href: '/docs/email-labels'
				},
				{ title: 'Email Settings', href: '/docs/email-settings' }
			]
		},
		{
			title: 'Contributing',
			links: [
				{ title: 'Community', href: '/docs/contributing' },
				{ title: 'Code', href: '/docs/contributing#code' },
				{ title: 'Architecture', href: '/docs/contributing#architecture' },
				{ title: 'Data', href: '/docs/contributing#email' }
			]
		}
	];
</script>

<aside
	class="sticky top-[4.5rem] -ml-0.5 h-[calc(100vh-4.5rem)] overflow-y-auto overflow-x-hidden py-16 pl-0.5"
>
	<nav class="w-64 pr-8 text-base lg:text-sm xl:w-72 xl:pr-16">
		<ul class="space-y-9">
			{#each nav as section}
				<li>
					<h2 class="font-display font-medium text-slate-900 dark:text-white">
						{section.title}
					</h2>
					<ul
						class="mt-2 space-y-2 border-l-2 border-slate-100 dark:border-slate-800 lg:mt-4 lg:space-y-4 lg:border-slate-200"
					>
						{#each section.links as link}
							<li class="relative">
								<a
									href={link.href}
									class={`block w-full pl-3.5 before:pointer-events-none before:absolute before:-left-1 before:top-1/2 before:h-1.5 before:w-1.5 before:-translate-y-1/2 before:rounded-full ${
										link.href === $page.url.pathname
											? 'font-semibold text-sky-500 before:bg-sky-500'
											: 'text-slate-500 before:hidden before:bg-slate-300 hover:text-slate-600 hover:before:block dark:text-slate-400 dark:before:bg-slate-700 dark:hover:text-slate-300'
									}`}
								>
									{link.title}
								</a>
							</li>
						{/each}
					</ul>
				</li>
			{/each}
		</ul>
	</nav>
</aside>
