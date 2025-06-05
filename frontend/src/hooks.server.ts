import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	// The node address will be available in the request headers
	// We'll pass it through to the next request
	const nodeAddress = event.request.headers.get('X-Node-Address');
	
	if (nodeAddress) {
		// Forward the header to the next request
		event.request.headers.set('X-Node-Address', nodeAddress);
	}
	
	return resolve(event);
}; 