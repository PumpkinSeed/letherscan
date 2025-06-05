import { fetchWithNodeAddress } from '$lib/utils/fetch';

export const prerender = true;

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

interface Block {
    header: {
        number: string;
        // ... other header fields
    };
    transactions: Transaction[];
}

interface BlocksResponse {
    blocks: Block[];
}

// This function tells SvelteKit which transaction hashes to prerender
export async function entries() {
    // Fetch all blocks to get transaction hashes
    const response = await fetchWithNodeAddress('http://localhost:8080/blocks');
    if (!response.ok) {
        return [];
    }
    const data: BlocksResponse = await response.json();
    
    // Extract all transaction hashes from all blocks
    const hashes = data.blocks.flatMap((block: Block) => 
        block.transactions.map((tx: Transaction) => tx.hash)
    );
    
    return hashes.map((hash: string) => ({ hash }));
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