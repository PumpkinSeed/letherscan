import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Initialize the store with the value from localStorage if available
const storedValue = browser ? localStorage.getItem('numberOfBlocks') : null;
export const numberOfBlocks = writable<number>(storedValue ? parseInt(storedValue) : 10);

// Subscribe to changes and update localStorage
if (browser) {
	numberOfBlocks.subscribe((value) => {
		localStorage.setItem('numberOfBlocks', value.toString());
	});
} 