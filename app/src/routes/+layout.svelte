<script lang="ts">
	import '../app.css';

	import { supabaseClient } from '$lib/supabase/client';
	import { invalidate } from '$app/navigation';
	import { onMount } from 'svelte';

	// convert to ms and get iso string
	const expriryFromExpiresAt = (expiresAt: number) => new Date(expiresAt * 1000).toISOString();

	onMount(() => {
		const {
			data: { subscription }
		} = supabaseClient.auth.onAuthStateChange((event, session) => {
			invalidate('supabase:auth');

			if (event === 'SIGNED_IN' && session) {
				// synchronize provider tokens with db
				const { expires_at, provider_token, provider_refresh_token, user } = session;
				const { provider } = user.app_metadata;

				if (!expires_at || !provider_token || !provider_refresh_token || !provider) {
					console.log('missing data. not updating user oauth token');
					return;
				}
				// format tokens for db
				supabaseClient
					.from('user_oauth_token')
					.upsert({
						user_id: user.id,
						provider,
						token: {
							access_token: provider_token,
							refresh_token: provider_refresh_token,
							expiry: expriryFromExpiresAt(expires_at)
						}
					})
					.then(({ data, error }) => {
						if (error) {
							console.log(data, error);
						}
					});
			}
		});

		return () => {
			subscription.unsubscribe();
		};
	});
</script>

<slot />
