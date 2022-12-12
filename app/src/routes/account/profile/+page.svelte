<script lang="ts">
	// import page
  import { page } from '$app/stores';
	// TODOs
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
	// load script from public folder with promise
	// https://stackoverflow.com/questions/59629947/how-do-i-load-an-external-js-library-in-svelte-sapper

	// check if in browser
	$: console.log('page', $page.data);
</script>

<div class="my-12 lg:grid lg:grid-cols-12 lg:gap-x-5">
  <aside class="py-6 px-2 hidden lg:block sm:px-6 lg:col-span-3 lg:py-0 lg:px-0">
    <nav class="space-y-1">
      <!-- Current: "bg-slate-50 text-blue-700 hover:text-indigo-700 hover:bg-white", Default: "text-slate-900 hover:text-slate-900 hover:bg-slate-50" -->
      <a href="#" class="bg-slate-50 text-blue-700 hover:text-indigo-700 hover:bg-white group rounded-md px-3 py-2 flex items-center text-sm font-medium" aria-current="page">
        <!--
          Heroicon name: outline/user-circle

          Current: "text-blue-500 group-hover:text-indigo-500", Default: "text-slate-400 group-hover:text-slate-500"
        -->
        <svg class="text-blue-500 group-hover:text-indigo-500 flex-shrink-0 -ml-1 mr-3 h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0012 15.75a7.488 7.488 0 00-5.982 2.975m11.963 0a9 9 0 10-11.963 0m11.963 0A8.966 8.966 0 0112 21a8.966 8.966 0 01-5.982-2.275M15 9.75a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        <span class="truncate">Profile</span>
      </a>
    </nav>
  </aside>

  <div class="space-y-6 sm:px-6 lg:col-span-9 lg:px-0">
    <form action="#" method="POST">
      <div class="shadow sm:overflow-hidden sm:rounded-md">
        <div class="space-y-6 bg-white py-6 px-4 sm:p-6">
          <div>
            <h3 class="text-lg font-medium leading-6 text-slate-900">Profile</h3>
            <p class="mt-1 text-sm text-slate-500">Let us know how you'd like us to address you in emails</p>
          </div>

          <div class="grid grid-cols-6 gap-6">
            <div class="col-span-6 sm:col-span-3">
              <label for="first-name" class="block text-sm font-medium text-slate-700">First name</label>
              <input type="text" name="first-name" id="first-name" autocomplete="given-name" value={$page.data.profile.firstName} class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 sm:text-sm">
            </div>

            <div class="col-span-6 sm:col-span-3">
              <label for="last-name" class="block text-sm font-medium text-slate-700">Last name</label>
              <input type="text" name="last-name" id="last-name" autocomplete="family-name" value={$page.data.profile.lastName} class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 sm:text-sm">
            </div>

            <div class="col-span-6 sm:col-span-4">
              <label for="email-address" class="block text-sm font-medium text-slate-700">Email address</label>
              <input type="text" name="email-address" id="email-address" autocomplete="email" disabled
						value={$page.data.profile.email}
            class="mt-1 block w-full rounded-md border border-slate-300 py-2 px-3 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-indigo-500 sm:text-sm disabled:cursor-not-allowed disabled:border-slate-200 disabled:bg-slate-50 disabled:text-slate-500">
            </div>
        </div>
      </div>
    </form>
  </div>
</div>
