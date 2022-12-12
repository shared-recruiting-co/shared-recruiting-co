<script lang="ts">
	// import page
  import { page } from '$app/stores';
	// TODOs
	// Clean up this file
	// Create an callback endpoint that
	// - verifies the user is the same
	// - Confirm the X-Requested-With: XmlHttpRequest header is set for popup mode.
	// - fetches access token and refresh token
	// - saves to db
	// - triggers a sync if first time
	// - subscribes to notifications
	// - validate on mobile!!
	// https://developers.google.com/identity/protocols/oauth2/web-server#httprest_3

	// Account Header
	// Client ID env var
	// load script from public folder with promise
	// https://stackoverflow.com/questions/59629947/how-do-i-load-an-external-js-library-in-svelte-sapper

	// check if in browser
	$: console.log('page', $page.data);
</script>

<div class="mx-4 my-8 flex flex-row items-center justify-center sm:mx-auto sm:my-16">
	<div class="">
		<h1>Profile</h1>
		<h1>{$page.data.profile.firstName}</h1>
		<h1>{$page.data.profile.lastName}</h1>
	</div>
</div>
