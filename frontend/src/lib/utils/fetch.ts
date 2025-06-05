import { get } from 'svelte/store';
import { nodeAddress } from '../stores/nodeAddress';

export async function fetchWithNodeAddress(input: RequestInfo | URL, init?: RequestInit) {
	const address = get(nodeAddress);
	
	if (!init) {
		init = {};
	}
	if (!init.headers) {
		init.headers = {};
	}
	
	if (address) {
		(init.headers as Record<string, string>)['X-Node-Address'] = address;
	}
	
	return fetch(input, init);
} 