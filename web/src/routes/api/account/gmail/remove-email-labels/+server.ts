import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getRefreshedGoogleAccessToken } from '$lib/supabase/client.server';
import { email } from '$lib/components/ConnectGoogleAccountButton.svelte';

/**
 * Retrieves the IDs of all Gmail labels that start with "@SRC"
 *
 * @async
 * @param {string} accessToken - A valid access token for the Gmail API
 * @returns {Array<string>} - An array of the IDs of all Gmail labels that start with "@SRC"
 * @throws Throws an error if the API request fails or if the response is invalid
 */
const getSrcLabelIds = async (accessToken: string): Array<string> => {

    // The Gmail labels endpoint, gives the full list of labels used in the users Gmail
    const labelsEndpoint = `https://gmail.googleapis.com/gmail/v1/users/me/labels`;
  
    // query the labels endpoint
    const response = await fetch(labelsEndpoint, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
        Accept: 'application/json',
      },
    });
    const data = await response.json();

    // all SRC labels start with @SRC, we will use this to identify all the SRC labels
    const labelPrefix = '@SRC';

    // get a list of the label IDs that start with the above prefix
    const srcLabelsIds = data.labels
      .filter(label => label.name.startsWith(labelPrefix))
      .map(label => label.id);
  
    // return the list of SRC label IDs
    return srcLabelsIds;
}

/**
 * Helper function that retrieves the label IDs of a Gmail thread (message).
 *
 * @async
 * @param {string} accessToken - A valid access token for the Gmail API.
 * @param {string} threadId - The ID of the Gmail thread to retrieve label IDs for.
 * @returns {Promise<string[]>} - A promise that resolves with an array of label IDs for the given thread.
 * @throws Throws an error if the API request fails or if the response is invalid.
 */
const getThreadLabels = async (accessToken: string, threadId: string): Promise<string[]> => {

  // The Gmail thread endpoint, gives the details about a specific Gmail thread
  const endpoint = `https://gmail.googleapis.com/gmail/v1/users/me/threads/${threadId}`;
  
  // query the thread endpoint
  const response = await fetch(endpoint, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
      Accept: 'application/json',
    },
  });

  // Parse the API response data
  const data = await response.json();
  
  // if the thread has messages with labels, return them
  if (data.messages) {
    return data.messages[0].labelIds
  } else {
    return [];
  }
}

/**
 * Remove labels from a Gmail thread using the Gmail API.
 * 
 * @param {string} accessToken - The access token for the authenticated user.
 * @param {string} threadId - The ID of the Gmail thread to modify.
 * @param {string[]} labelIdsToRemove - An array of label IDs to remove from the thread.
 * @returns {Promise<boolean>} A boolean indicating whether the labels were successfully removed from the thread.
 */
const removeLabelsFromThread = async (accessToken: string, threadId: string, labelIdsToRemove: string[]): Promise<boolean> => {
    
  // The Gmail modify thread endpoint which will allow us to remove labels from the given thread
  const endpoint = `https://gmail.googleapis.com/gmail/v1/users/me/threads/${threadId}/modify`;

  // POST to the modify endpoint with the array of labels to remove
  const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        removeLabelIds: labelIdsToRemove,
        addLabelIds: [],
    }),
  });

  return response.ok
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
const isGmailThreadIdValid = async (accessToken: string, threadId: string): Promise<boolean> => {
  // The Gmail thread endpoint, which gives the details of a specific Gmail thread
  const threadEndpoint = `https://gmail.googleapis.com/gmail/v1/users/me/threads/${threadId}`;

  try {
    // Query the Gmail thread endpoint with the specified thread ID
    const response = await fetch(threadEndpoint, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
        Accept: 'application/json',
      },
    });

    // If the API returns a valid response, the thread ID is valid and corresponds to a Gmail thread
    return response.ok;
    
  } catch (err) {
    // If the API request fails or the response is invalid, throw an error
    throw new Error(`Failed to retrieve information for thread ID: ${threadId}, corresponds to a Gmail thread: ${err.message}`);
  }
};

/**
 * Main handler function for this POST request route
 *
 * @param request - The incoming request object
 * @param locals.getSession - function used to retrieve the current user's session.
 * @param locals.supabase - the supabase client
 * @returns {Response} - A response indicating whether or not the operation was successful.
 * @throws Throws an error if the user is not authorized or if an unexpected error occurs.
 */
export const POST: RequestHandler = async ({ request, locals: { getSession, supabase } }) => {
  
  // Get the current user's session
  const session = await getSession();
  if (!session) throw error(401, 'unauthorized');

  // retrieve the necessary information from the request 
  const { email, email_thread_id: threadId } = await request.json();

  if (!email || !threadId) {
    throw error(400, 'Could not find both email and email_thread_id from request json');
  }

  // Get a refreshed Google access token for the current user
  let accessToken = '';
  try {
    accessToken = await getRefreshedGoogleAccessToken(supabase, email);
  } catch (err) {
    throw error(500, err instanceof Error ? err : 'Could not retrieve access token');
  }

  // check the thread actually exists in Gmail before attempting the update
  const isValidGmailThread = await isGmailThreadIdValid(accessToken, threadId)
  if (!isValidGmailThread) {
    return new Response({status: 404, statusText: "No corresponding Gmail thread"});
  }

  // Get the label IDs for all "@SRC" labels
  const srcLabelIds = await getSrcLabelIds(accessToken);

  // Remove all "@SRC" labels from the specified Gmail thread
  const success = await removeLabelsFromThread(accessToken, threadId, srcLabelIds);

  if (success) {
    // Return a success response
    return new Response('success');
  } else {
    // Return a failure response with an appropriate error message
    throw error(500, 'Failed to remove SRC labels');
  }
};
