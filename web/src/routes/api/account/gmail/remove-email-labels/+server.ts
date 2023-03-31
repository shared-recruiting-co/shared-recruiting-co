import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getRefreshedGoogleAccessToken } from '$lib/supabase/client.server';

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

  // Get the labels before the update
  const labelsBeforeUpdate = await getThreadLabels(accessToken, threadId);

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

  // Get the labels after the update
  const labelsAfterUpdate = await getThreadLabels(accessToken, threadId);

  // Calculate what is expected of the updated labels
  const expectedLabels = labelsBeforeUpdate.filter(label => !labelIdsToRemove.includes(label));

  // Check if the labels retrieved via the Gmail API match what is expected
  const updateSuccessful = (
    expectedLabels.length === labelsAfterUpdate.length &&
    expectedLabels.every((value, index) => value === labelsAfterUpdate[index])
  );

  // Return a boolean indicating whether the update was successful
  return updateSuccessful;
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

  // Get the email and thread ID from the request body (if they exist)
  let email = session.user.email;
  let threadId = 'not present';
  try {
    const { user_email: reqEmail, email_thread_id: reqThreadId } = await request.json();
    email = reqEmail || email;
    threadId = reqThreadId || threadId;
  } catch (err) {
    // Do nothing
  }

  // Get a refreshed Google access token for the current user
  let accessToken = '';
  try {
    accessToken = await getRefreshedGoogleAccessToken(supabase, email);
  } catch (err) {
    throw error(500, err instanceof Error ? err : 'unexpected error occurred');
  }

  // check the thread actually exists in Gmail before attempting the update
  const isValidGmailThread = await isGmailThreadIdValid(accessToken, threadId)
  if (!isValidGmailThread) {
    return new Response({status: 200, statusText: "No corresponding Gmail thread"});
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
