<script lang="ts" >
	import ConnectGoogleAccountButton from '$lib/components/ConnectGoogleAccountButton.svelte';
	// TODOs
	// 
	// Refactor Connect button to be a component
	// Validate on mobile!! Use redirect flow if necessary
	// Trigger welcome email 
	// Update cloud functions logic (sync to date OR sync to history id)
	// Update watch cloud function to use filter action
	// Create separate layouts for (marketing) and (account) pages 
	// Create account header
	// Create account sidebar
	// Design & create account page 
	// Watch/Stop endpoints
</script>

<div class="mx-4 my-8 flex flex-row items-center justify-center sm:mx-auto sm:my-16">
	<div class="">
		<h1 class="text-2xl">Welcome to the SRC, Devin!</h1>
		<p class="text-gray-500">Let's get started by connecting your email.</p>
		<ConnectGoogleAccountButton />
	</div>
</div>
