<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import '../app.css';

	import { supabaseClient } from '$lib/supabase/client';

	onMount(() => {
		// Add note to developers that open the console
		console.log('Greetings Developer 🖖');
		console.log(
			'Something wrong? Report it at https://github.com/shared-recruiting-co/shared-recruiting-co/issues'
		);

		const {
			data: { subscription }
		} = supabaseClient.auth.onAuthStateChange((_event, _session) => {
			invalidate('supabase:auth');
		});

		return () => {
			subscription.unsubscribe();
		};
	});
</script>

<body class="flex min-h-screen flex-col text-slate-900">
	<slot />
</body>
