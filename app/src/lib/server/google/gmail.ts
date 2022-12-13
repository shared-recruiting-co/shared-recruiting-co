import { NEW_EMAIL_PUB_SUB_TOPIC } from '$env/static/private';

// watch for new emails
export const watch = async (accessToken: string) =>
	await fetch('https://gmail.googleapis.com/gmail/v1/users/me/watch', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${accessToken}`
		},
		body: JSON.stringify({
			// TODO: Parameterize once we have use cases for multiple topics or labels
			labelIds: ['UNREAD'],
			labelFilterAction: 'include',
			topicName: NEW_EMAIL_PUB_SUB_TOPIC
		})
	});

// stop watching for new emails
export const stop = async (accessToken: string) =>
	await fetch('https://gmail.googleapis.com/gmail/v1/users/me/stop', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${accessToken}`
		}
	});
