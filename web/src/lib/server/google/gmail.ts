import { dev } from '$app/environment';

import { Label } from '$lib/google/labels';
import type { Labels } from '$lib/google/labels';

// Until we have a better local development and testing story,
// We will always return true for dev
const success = new Response('success');

export type WatchRequest = {
	topicName: string;
	labelIds: string[];
	labelFilterAction?: 'include' | 'exclude';
};

// watch for new emails
export const watch = async (accessToken: string, body: WatchRequest) =>
	dev
		? success
		: await fetch('https://gmail.googleapis.com/gmail/v1/users/me/watch', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${accessToken}`
				},
				body: JSON.stringify(body)
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
export const sendMessage = async (
	accessToken: string,
	message: SendMessageArgs
): Promise<Response> => {
	// do nothing in development
	if (dev) return new Response();

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

/**
 * Retrieves the IDs of all Gmail labels that start with "@SRC"
 */
export const getSRCLabels = async (accessToken: string): Promise<Labels> => {
	// The Gmail labels endpoint, gives the full list of labels used in the users Gmail
	const endpoint = 'https://gmail.googleapis.com/gmail/v1/users/me/labels';

	// query the labels endpoint
	const response = await fetch(endpoint, {
		headers: {
			Authorization: `Bearer ${accessToken}`,
			Accept: 'application/json'
		}
	});
	const data = await response.json();

	// get a list of the label IDs that start with the above prefix
	const srcLabels = data.labels
		// all SRC labels start with @SRC, we will use this to identify all the SRC labels
		.filter((label: { name: string }) => label.name.startsWith(Label.SRC))
		.reduce(
			(
				acc: Labels,
				label: {
					name: `${Label}`;
					id: string;
				}
			) => {
				acc[label.name] = label.id;
				return acc;
			},
			{} as Labels
		);

	// return the list of SRC labels (name -> ID)
	return srcLabels;
};

/**
 * Helper function that retrieves the label IDs of a Gmail thread (message).
 *
 * @async
 * @param {string} accessToken - A valid access token for the Gmail API.
 * @param {string} threadId - The ID of the Gmail thread to retrieve label IDs for.
 * @returns {Promise<string[]>} - A promise that resolves with an array of label IDs for the given thread.
 * @throws Throws an error if the API request fails or if the response is invalid.
 */
export const getThreadLabels = async (accessToken: string, threadId: string): Promise<string[]> => {
	// The Gmail thread endpoint, gives the details about a specific Gmail thread
	const endpoint = `https://gmail.googleapis.com/gmail/v1/users/me/threads/${threadId}`;

	// query the thread endpoint
	const response = await fetch(endpoint, {
		headers: {
			Authorization: `Bearer ${accessToken}`,
			Accept: 'application/json'
		}
	});

	// Parse the API response data
	const data = await response.json();

	// if the thread has messages with labels, return them
	if (data.messages) {
		return data.messages[0].labelIds;
	} else {
		return [];
	}
};

/**
 * Remove labels from a Gmail thread using the Gmail API.
 *
 * @returns {Promise<boolean>} A boolean indicating whether the labels were successfully removed from the thread.
 */
export const updateThreadLabels = async ({
	accessToken,
	threadId,
	addLabelIds,
	removeLabelIds
}: {
	accessToken: string;
	threadId: string;
	addLabelIds?: string[];
	removeLabelIds?: string[];
}): Promise<boolean> => {
	// The Gmail modify thread endpoint which will allow us to remove labels from the given thread
	const endpoint = `https://gmail.googleapis.com/gmail/v1/users/me/threads/${threadId}/modify`;

	// POST to the modify endpoint with the array of labels to remove
	const response = await fetch(endpoint, {
		method: 'POST',
		headers: {
			Authorization: `Bearer ${accessToken}`,
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			removeLabelIds,
			addLabelIds
		})
	});

	return response.ok;
};

/**
 * Checks if a given thread ID corresponds to a Gmail thread by attempting to retrieve the thread details using the Gmail API.
 *
 * @async
 * @param {string} accessToken - The access token for the authenticated user.
 * @param {string} threadId - The ID of the thread to check.
 * @returns {Promise<boolean>} - A promise that resolves with a boolean indicating whether or not the thread ID corresponds to a Gmail thread.
 * @throws {Error} - Throws an error if the API request fails or if the response is invalid.
 */
export const isValidThread = async (accessToken: string, threadId: string): Promise<boolean> => {
	// The Gmail thread endpoint, which gives the details of a specific Gmail thread
	const endpoint = `https://gmail.googleapis.com/gmail/v1/users/me/threads/${threadId}`;

	try {
		// Query the Gmail thread endpoint with the specified thread ID
		const response = await fetch(endpoint, {
			headers: {
				Authorization: `Bearer ${accessToken}`,
				Accept: 'application/json'
			}
		});

		// If the API returns a valid response, the thread ID is valid and corresponds to a Gmail thread
		return response.ok;
	} catch (err) {
		// If the API request fails or the response is invalid, throw an error
		throw new Error(
			`Failed to retrieve information for thread ID: ${threadId}, corresponds to a Gmail thread: ${err.message}`
		);
	}
};

