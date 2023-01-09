import type { HandleServerError } from '@sveltejs/kit';
import '$lib/supabase/client';
// import * as Sentry from '@sentry/node';

// Sentry.init({/*...*/})

export const handleError = (({ error, event }) => {
	// example integration with https://sentry.io/
	// Sentry.captureException(error, { event });

	return {
		message:
			error?.message ??
			'Something went wrong. If this error persists, please contact us at team@sharedrecruiting.co',
		code: error?.code ?? 'UNKNOWN'
	};
}) satisfies HandleServerError;
