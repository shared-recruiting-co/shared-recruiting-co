import { PUBLIC_GOOGLE_CLIENT_ID, PUBLIC_GOOGLE_REDIRECT_URI } from '$env/static/public';
import { GOOGLE_CLIENT_SECRET } from '$env/static/private';

export const exchangeCodeForTokens = async (code: string, method: 'GET' | 'POST') => {
	const data = new URLSearchParams();
	data.append('code', code);
	data.append('client_id', PUBLIC_GOOGLE_CLIENT_ID);
	data.append('client_secret', GOOGLE_CLIENT_SECRET);
	data.append('grant_type', 'authorization_code');

	if (method === 'GET') {
		data.append('redirect_uri', PUBLIC_GOOGLE_REDIRECT_URI);
	} else {
		data.append('redirect_uri', 'postmessage');
	}

	return await fetch('https://oauth2.googleapis.com/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded'
		},
		body: data
	});
};

// refresh token
export const refreshAccessToken = async (refreshToken: string) => {
	const data = new URLSearchParams();
	data.append('refresh_token', refreshToken);
	data.append('client_id', PUBLIC_GOOGLE_CLIENT_ID);
	data.append('client_secret', GOOGLE_CLIENT_SECRET);
	data.append('grant_type', 'refresh_token');

	return await fetch('https://oauth2.googleapis.com/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded'
		},
		body: data
	});
};
