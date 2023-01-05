<script lang="ts">
	import { page } from '$app/stores';

	const nav = [
		{
			title: 'Introduction',
			links: [
				{ title: 'Welcome', href: '/' },
				{ title: 'Open Source', href: '/docs/installation' },
				{ title: 'Security & Privacy', href: '/docs/security' }
			]
		},
		{
			title: 'Core Concepts',
			links: [
				{ title: 'Connect Gmail Account', href: '/docs/connect-email' },
				{
					title: 'Gmail Labels',
					href: '/docs/labels'
				},
				{ title: 'Email Settings', href: '/docs/email-settings' }
			]
		},
		{
			title: 'Contributing',
			links: [
				{ title: 'Community', href: '/docs/how-to-contribute' },
				{ title: 'Code', href: '/docs/architecture-guide' },
				{ title: 'Architecture', href: '/docs/architecture-guide' },
				{ title: 'Data', href: '/docs/design-principles' }
			]
		}
	];
</script>

<aside>
	<nav class="text-base lg:text-sm">
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
									class={'block w-full pl-3.5 before:pointer-events-none before:absolute before:-left-1 before:top-1/2 before:h-1.5 before:w-1.5 before:-translate-y-1/2 before:rounded-full ' +
										link.href ===
									$page.url.pathname
										? 'font-semibold text-sky-500 before:bg-sky-500'
										: 'text-slate-500 before:hidden before:bg-slate-300 hover:text-slate-600 hover:before:block dark:text-slate-400 dark:before:bg-slate-700 dark:hover:text-slate-300'}
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
