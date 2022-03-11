import createAuth0Client from '@auth0/auth0-spa-js';

export async function init() {
	const auth0 = await createAuth0Client({
		domain: 'evaluate-ams-pro.eu.auth0.com',
		client_id: 'JbCpnPq87PBetrgEtV70JJNNx8aNWXEi',
		redirect_uri: window.location.origin,
		cacheLocation: 'localstorage'
	});

    try {
        await auth0.handleRedirectCallback()
        window.location.search = ''
    } catch {
        // NO PROBLEM, login
    }

	if (!(await auth0.isAuthenticated())) {
		await auth0.loginWithRedirect();
	}
}
