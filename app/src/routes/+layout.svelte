<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';
	import '../app.css';

	import { supabaseClient } from '$lib/supabase/client';

	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	onMount(() => {
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
	<Header />
	<main class="flex-1">
		<slot />
	</main>
	<Footer />
</body>
