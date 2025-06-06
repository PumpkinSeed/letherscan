import { fetchWithNodeAddress } from '$lib/utils/fetch';
import { browser } from '$app/environment';
import { BASE_URL } from '$lib/config';

export const prerender = false;

interface Transaction {
    hash: string;
    nonce: number;
    block_hash: string | null;
    block_number: string;
    transaction_index: number;
    from: string;
    to: string;
    value: string;
    gas_price: string;
    gas: number;
    input: string;
    v: string;
    r: string;
    s: string;
    chain_id: string;
    type: string;
    method: string;
    isPending: boolean;
}

export async function load({ params }: { params: { hash: string } }) {
    try {
        if (!browser) {
            throw new Error('This code must run in the browser');
        }
        const response = await fetchWithNodeAddress(`${BASE_URL}/transaction/${params.hash}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const transaction = await response.json();
        return {
            transaction
        };
    } catch (error) {
        console.error('Error fetching transaction:', error);
        return {
            transaction: null
        };
    }
} 