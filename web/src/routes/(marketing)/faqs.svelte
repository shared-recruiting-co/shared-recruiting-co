<script lang="ts">
	import FAQs from '$lib/components/marketing/FAQs.svelte';

	const faqs = [
		{
			question: 'How do I pronounce SRC?',
			answer: 'SRC is pronounced "source".'
		},
		{
			question: 'How can I join?',
			answer: "Sign up and we'll reach out once you are off the waiting list."
		},
		{
			question: 'How much does it cost to join SRC?',
			answer: 'Nothing. It will always be free to join SRC as a candidate.'
		},
		{
			question: 'How does SRC use my data?',
			answer: `We process your emails in order to classify them as recruiting email or not. We do not and will never use your emails for any other purpose. We are take trust and privacy so seriously that all of our code is open source.
					<br/><br/>Read more in our <a class="underline" href="/legal/privacy-policy.md">privacy policy</a>`
		},
		{
			question: 'How do I contribute to SRC?',
			answer:
				"The SRC is open source. You can contribute to SRC on <a class='underline' href='https://github.com/shared-recruiting-co/shared-recruiting-co' target='_blank' rel='noopener noreferrer'>GitHub</a>. All you have to do is make a PR or an issue! <br/><br/> You can also contribute your recruiting emails to improve the SRC's ML model, by forwarding them to <a class='underline' href='mailto:examples@sharedrecruiting.co'>examples@sharederecruiting.co</a>"
		},
		{
			question: 'How does SRC detect job opportunities?',
			answer:
				'Every email you receive is processed by SRCs ML model. If the email is a job opportunity, the SRC will automatically label it as such.'
		},
		{
			question: 'How does the SRC define a job opportunity?',
			answer:
				'We define a job opportunity as an inbound email from a recruiter or hiring manager. This excludes job newsletters or automated notifications.'
		},
		{
			question: 'More Questions?',
			answer:
				"Reach out to us at <a class='underline' href='mailto:team@sharedrecruiting.co'>team@sharederecruiting.co</a>! We'd love to chat."
		}
	];
</script>

<FAQs {faqs} />
