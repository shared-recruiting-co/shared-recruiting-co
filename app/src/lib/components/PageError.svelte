<script>
	import { page } from '$app/stores';

	export let home = '/';

	const NOT_FOUND = 404;
	// naive error messages
	let title = $page.status === NOT_FOUND ? 'Page not found' : 'Something went wrong';
	let message =
		$page.status === NOT_FOUND
			? 'Please check the URL in the address bar and try again'
			: $page.error?.message;
</script>

<div class="min-h-full bg-white py-16 px-6 sm:py-24 md:grid md:place-items-center lg:px-8">
	<div class="mx-auto max-w-max">
		<main class="sm:flex">
			<p class="text-4xl font-bold tracking-tight text-blue-600 sm:text-5xl">{$page.status}</p>
			<div class="sm:ml-6">
				<div class="sm:border-l sm:border-gray-200 sm:pl-6">
					<h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">
						{title}
					</h1>
					<p class="mt-1 text-base text-gray-500">
						{message}
					</p>
				</div>
				<div class="mt-10 flex space-x-3 sm:border-l sm:border-transparent sm:pl-6">
					<a
						href={home}
						class="inline-flex items-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
						>Go back home</a
					>
					<a
						href="mailto:team@sharedrecruiting.co"
						class="inline-flex items-center rounded-md border border-transparent bg-blue-100 px-4 py-2 text-sm font-medium text-blue-700 hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
						>Contact us</a
					>
				</div>
			</div>
		</main>
	</div>
</div>
