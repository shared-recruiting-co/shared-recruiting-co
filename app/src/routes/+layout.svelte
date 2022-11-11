<script>
	import '../app.css';

	import { supabaseClient } from '$lib/supabase/client';
	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		const {
			data: { subscription }
		} = supabaseClient.auth.onAuthStateChange((event, session) => {
			invalidate('supabase:auth');

			if (event === "SIGNED_IN" && session) {
				// synchronize provider tokens with db
				const { provider_token, provider_refresh_token } = session;
				// do something
			}
		});

		return () => {
			subscription.unsubscribe();
		};
	});
</script>

<slot />
