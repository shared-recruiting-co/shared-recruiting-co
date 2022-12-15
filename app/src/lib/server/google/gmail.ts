import { dev } from '$app/environment';
import { NEW_EMAIL_PUBSUB_TOPIC } from '$env/static/private';

// Until we have a better local development and testing story,
// We will always return true for dev
const success = new Response('success');

// watch for new emails
export const watch = async (accessToken: string) =>
	dev
		? success
		: await fetch('https://gmail.googleapis.com/gmail/v1/users/me/watch', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${accessToken}`
				},
				body: JSON.stringify({
					// TODO: Parameterize once we have use cases for multiple topics or labels
					labelIds: ['UNREAD'],
					labelFilterAction: 'include',
					topicName: NEW_EMAIL_PUBSUB_TOPIC
				})
		  });

// stop watching for new emails
export const stop = async (accessToken: string) =>
	dev
		? success
		: await fetch('https://gmail.googleapis.com/gmail/v1/users/me/stop', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${accessToken}`
				}
		  });

interface SendMessageArgs {
	from: string;
	to: string;
	subject: string;
	body: string;
}

const formatRawEmailMessage = ({ from, to, subject, body }: SendMessageArgs) =>
	`From: ${from}\r\nTo: ${to}\r\nSubject: ${subject}\r\n\r\n${body}`;

const urlSafeBase64 = (str: string) =>
	Buffer.from(str).toString('base64').replace(/\+/g, '-').replace(/\//g, '_');

// sendMessage  is a simple wrapper around the gmail API to send an email message.
// It does not handle replies, attachments, or other advanced features.
export const sendMessage = async (accessToken: string, message: SendMessageArgs) => {
	const raw = formatRawEmailMessage(message);
	return await fetch('https://gmail.googleapis.com/gmail/v1/users/me/messages/send', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${accessToken}`
		},
		body: JSON.stringify({ raw: urlSafeBase64(raw) })
	});
};
