import createAuth0Client from '@auth0/auth0-spa-js';
import type { Auth0Client } from '@auth0/auth0-spa-js';

export let auth0: Auth0Client;

export async function init() {
	auth0 = await createAuth0Client({
		domain: 'evaluate-ams-pro.eu.auth0.com',
		client_id: 'JbCpnPq87PBetrgEtV70JJNNx8aNWXEi',
		audience: 'https://ams-pro.de/apis',
		redirect_uri: window.location.origin,
		cacheLocation: 'localstorage'
	});

	try {
		await auth0.handleRedirectCallback();
		window.location.search = '';
	} catch {
		// NO PROBLEM, login
	}

	if (!(await auth0.isAuthenticated())) {
		await auth0.loginWithRedirect();
	}
}
