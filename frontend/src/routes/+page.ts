import { fetchWithNodeAddress } from '$lib/utils/fetch';
import { get } from 'svelte/store';
import { numberOfBlocks } from '$lib/stores/numberOfBlocks';

export const prerender = true;

export async function load() {
    try {
        const blocksToFetch = get(numberOfBlocks);
        const response = await fetchWithNodeAddress(`http://localhost:8080/blocks?number_of_blocks=${blocksToFetch}`);
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