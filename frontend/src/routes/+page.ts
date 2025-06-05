import { fetchWithNodeAddress } from '$lib/utils/fetch';
import { get } from 'svelte/store';
import { numberOfBlocks } from '$lib/stores/numberOfBlocks';
import { browser } from '$app/environment';

export const prerender = false;

export async function load() {
    try {
        const blocksToFetch = get(numberOfBlocks);
        if (!browser) {
            throw new Error('This code must run in the browser');
        }
        const response = await fetchWithNodeAddress(`${window.location.origin}/blocks?number_of_blocks=${blocksToFetch}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return {
            blocks: data.blocks
        };
    } catch (error) {
        console.error('Error fetching blocks:', error);
        return {
            blocks: []
        };
    }
}