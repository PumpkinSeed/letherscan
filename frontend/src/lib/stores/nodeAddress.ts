import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Initialize the store with the value from localStorage if available
const storedAddress = browser ? localStorage.getItem('nodeAddress') : null;
export const nodeAddress = writable<string>(storedAddress || '');

// Subscribe to changes and update localStorage
if (browser) {
	nodeAddress.subscribe((value) => {
		localStorage.setItem('nodeAddress', value);
	});
} 