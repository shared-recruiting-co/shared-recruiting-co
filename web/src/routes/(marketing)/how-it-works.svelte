<script lang="ts">
	import HowItWorks from '$lib/components/marketing/HowItWorks.svelte';
	const steps = [
		{
			title: 'Detect & Archive',
			description:
				'SRC detects any inbound job opportunities. Each job opportunity is labeled, replied to, and archived for future reference.'
		},
		{
			title: 'Summarize & Standardize',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardized format. '
		},
		{
			title: 'Review & Engage',
			description:
				'SRC does the heavy lifting of summarizing all your job opportunities into a single, standardized format. '
		}
	];
</script>

<HowItWorks {steps} />
