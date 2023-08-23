<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { dev } from '$app/environment';
	import { onMount } from 'svelte';
	import { inject } from '@vercel/analytics';

	import type { LayoutData } from './$types';

	import '../app.css';

	export let data: LayoutData;

	$: ({ supabase, session } = data);

	// use Vercel analytics
	inject({ mode: dev ? 'development' : 'production' });

	onMount(() => {
		// Add note to developers that open the console
		console.log('Greetings Developer 🖖');
		console.log(
			'Something wrong? Report it at https://github.com/shared-recruiting-co/shared-recruiting-co/issues'
		);

		const {
			data: { subscription }
		} = supabase.auth.onAuthStateChange((event, _session) => {
			if (_session?.expires_at !== session?.expires_at) {
        invalidate('supabase:auth');
      }
		});

		return () => {
			subscription.unsubscribe();
		};
	});
</script>

<body class="flex min-h-screen flex-col text-slate-900">
	<slot />
</body>
