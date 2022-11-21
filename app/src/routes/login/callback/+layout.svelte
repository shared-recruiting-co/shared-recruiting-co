<script lang="ts">
	import { onMount } from 'svelte';

	// HACK...hopefully this doesn't stick around too long
	// Supabase's onAuthStateChange isn't reliable enough to use for post-login callback
	// The Supabase session is only created after first page load
	// This triggers a second page load, so the server get the Supabase session, trigger post-login action, and redirect home
	onMount(() => {
		// check if we already have refreshed
		// force refresh on mount
		const url = new URL(window.location.href);
		const hasRefreshed = !!url.searchParams.get('refreshed');

		// if we have already refreshed, something is wrong
		if (hasRefreshed) {
			window.location.replace('/?error=login_callback_failed');
			return;
		}

		// refresh the page with the refresh param to trigger server-side login callback
		url.searchParams.set('refreshed', 'true');
		console.log('redirecting to', url.toString());
		window.location.replace(url.toString());
	});
</script>

<slot />
