export const nav = [
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
			{ title: 'Email Setup', href: '/docs/connect-email' },
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
			{ title: 'Community', href: '/docs/contributing#contributing' },
			{ title: 'Code', href: '/docs/contributing#code' },
			{ title: 'Architecture', href: '/docs/contributing#architecture' },
			{ title: 'Data', href: '/docs/contributing#email' }
		]
	}
];

export const isCurrentPage = (href: string, current: string, hash: string): boolean => {
	// if href contains a hash, try match it
	if (href.includes('#') && hash) {
		return href === current + hash;
	}

	return href === current;
};

export const getSectionTitle = (pathname: string): string => {
	for (const section of nav) {
		for (const link of section.links) {
			if (link.href === pathname) {
				return section.title;
			}
		}
	}
	return '';
};
