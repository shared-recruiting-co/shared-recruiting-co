import type { HandleServerError } from '@sveltejs/kit';
import '$lib/supabase/client';
import * as Sentry from '@sentry/node';

import { SENTRY_DSN } from '$env/static/private';

Sentry.init({
	dsn: SENTRY_DSN,
	// We recommend adjusting this value in production, or using tracesSampler
	// for finer control
	tracesSampleRate: 1.0
});

export const handleError: HandleServerError = ({ error, event }) => {
	// example integration with https://sentry.io/
	Sentry.captureException(error, { event });

	return {
		message:
			error?.message ??
			'Something went wrong. If this error persists, please contact us at team@sharedrecruiting.co',
		code: error?.code ?? 'UNKNOWN'
	};
};
