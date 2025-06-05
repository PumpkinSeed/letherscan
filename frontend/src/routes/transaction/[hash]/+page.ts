import { fetchWithNodeAddress } from '$lib/utils/fetch';

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
        const response = await fetchWithNodeAddress(`http://localhost:8080/transaction/${params.hash}`);
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