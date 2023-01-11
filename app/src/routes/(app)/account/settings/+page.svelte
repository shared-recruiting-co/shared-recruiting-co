<script lang="ts">
</script>
<div class="sm:my-18 my-12 lg:grid lg:grid-cols-12 lg:gap-x-5">
	<!-- Empty space for now -->
	<aside class="block py-6 px-2 sm:px-6 lg:col-span-1 lg:py-0 lg:px-0" />
	<div class="space-y-6 sm:px-6 lg:col-span-9 lg:px-0">
		<div class="px-4 sm:px-0">
			<h1 class="text-3xl sm:text-4xl">Settings</h1>
			<p class="mt-1 text-sm text-slate-500">Private account settings</p>
		</div>
			<div class="shadow sm:overflow-hidden sm:rounded-md">
				<div class="space-y-6 bg-white py-6 px-4 sm:p-6">
					<div class="max-w-2xl">
						<p class="text-red-600 font-semibold text-xs uppercase tracking-tight mb-2">DANGER ZONE</p>
						<h3 class="text-lg font-medium leading-6 text-slate-900">Delete Account</h3>
						<p class="mt-1 text-sm text-slate-500">
							Permenantly delete your SRC account. By deleting your account, all of your data will be removed from SRC's database. Your email data will not be modified during the deletion process, including SRC managed email labels. You can remove the SRC labels manually from your email client UI after your account is deleted.
						</p>
						<p class="mt-2 text-sm text-slate-500">
							Once your account is deleted, you will be logged out.
						</p>
					</div>
					<button
						type="button"
						class="inline-flex w-full justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 sm:w-auto sm:text-sm"
						on:click={() => {
						}}>Delete Account</button
					>
				</div>
			</div>
	</div>
</div>
