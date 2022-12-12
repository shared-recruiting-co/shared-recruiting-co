import { PUBLIC_GOOGLE_CLIENT_ID, PUBLIC_GOOGLE_REDIRECT_URI } from '$env/static/public';
import { GOOGLE_CLIENT_SECRET } from '$env/static/private';

export const exchangeCodeForTokens = async (code: string) => {
  return await fetch('https://oauth2.googleapis.com/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded'
		},
		body: `code=${code}&client_id=${PUBLIC_GOOGLE_CLIENT_ID}&client_secret=${GOOGLE_CLIENT_SECRET}&redirect_uri=${PUBLIC_GOOGLE_REDIRECT_URI}&grant_type=authorization_code`
	});
}


// refresh token 
export const refreshAccessToken = async (refreshToken: string) => {
	return await fetch('https://oauth2.googleapis.com/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded'
		},
		body: `client_id=${PUBLIC_GOOGLE_CLIENT_ID}&client_secret=${GOOGLE_CLIENT_SECRET}&refresh_token=${refreshToken}&grant_type=refresh_token`
	});
}
