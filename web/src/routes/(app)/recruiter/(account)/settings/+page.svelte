<script lang="ts">
	import type { PageData } from './$types';
	import ConnectGoogleAccountButton from '$lib/components/ConnectGoogleAccountButton.svelte';

	export let data: PageData;

	const onConnect = () => {
		console.log('connected');
		// do nothing for now
	};
</script>

<div>Coming soon.</div>

<ConnectGoogleAccountButton {onConnect} email={data.session?.user?.email} />
