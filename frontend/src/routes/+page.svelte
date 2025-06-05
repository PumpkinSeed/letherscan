<script lang="ts">
    import { onMount } from 'svelte';
    import { fetchWithNodeAddress } from '$lib/utils/fetch';

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

    interface BlockHeader {
        parent_hash: string;
        uncle_hash: string;
        miner: string;
        root: string;
        tx_hash: string;
        receipt_hash: string;
        bloom: string;
        difficulty: string;
        number: string;
        gas_limit: number;
        gas_used: number;
        timestamp: number;
        extra: string;
        mix_digest: string;
        nonce: string;
        base_fee: string;
        withdrawals_hash: string;
        blob_gas_used: number;
        excess_blob_gas: number;
        parent_beacon_root: string;
        requests_hash: string;
    }

    interface Block {
        header: BlockHeader;
        transactions: Transaction[];
    }

    export let data;
    let isLoading = false;

    function formatTimestamp(timestamp: number): string {
        return new Date(timestamp * 1000).toLocaleString();
    }

    async function handleReload() {
        try {
            isLoading = true;
            const response = await fetchWithNodeAddress('http://localhost:8080/blocks');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const newData = await response.json();
            data.blocks = newData.blocks;
        } catch (error) {
            console.error('Error fetching blocks:', error);
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="container mx-auto px-4 py-8">
    <div class="flex items-center gap-4 mb-6">
        <h1 class="text-3xl font-bold">Block Explorer</h1>
        <button 
            on:click={handleReload}
            class="p-2 hover:bg-gray-100 rounded-full transition-colors"
            title="Reload blocks"
            disabled={isLoading}
        >
            <svg 
                xmlns="http://www.w3.org/2000/svg" 
                width="24" 
                height="24" 
                viewBox="0 0 24 24" 
                fill="none" 
                stroke="currentColor" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round"
                class="text-gray-600 {isLoading ? 'animate-spin' : ''}"
            >
                <path d="M21 2v6h-6"/>
                <path d="M3 12a9 9 0 0 1 15-6.7L21 8"/>
                <path d="M3 22v-6h6"/>
                <path d="M21 12a9 9 0 0 1-15 6.7L3 16"/>
            </svg>
        </button>
    </div>

    <div class="space-y-6">
        {#each data.blocks as block}
            <div class="bg-white shadow-lg rounded-lg overflow-hidden">
                <div class="p-6">
                    <div class="flex justify-between items-center mb-4">
                        <h2 class="text-xl font-semibold">Block #{block.header.number}</h2>
                        <span class="text-gray-500">{formatTimestamp(block.header.timestamp)}</span>
                    </div>
                    
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                        <div>
                            <p class="text-sm text-gray-600">Miner</p>
                            <p class="font-mono">{block.header.miner}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Gas Used</p>
                            <p>{block.header.gas_used.toLocaleString()} / {block.header.gas_limit.toLocaleString()}</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Base Fee</p>
                            <p>{block.header.base_fee} wei</p>
                        </div>
                        <div>
                            <p class="text-sm text-gray-600">Transactions</p>
                            <p>{block.transactions.length}</p>
                        </div>
                    </div>

                    <div class="mt-4">
                        <h3 class="text-lg font-semibold mb-2">Transactions</h3>
                        <div class="space-y-2">
                            {#each block.transactions as tx}
                                <a href="/transaction/{tx.hash}" class="block bg-gray-50 p-3 rounded hover:bg-gray-100 transition-colors">
                                    <div class="flex justify-between items-center">
                                        <span class="font-mono text-sm">{tx.hash}</span>
                                        <span class="text-sm text-gray-500">Gas: {tx.gas}</span>
                                    </div>
                                    <div class="text-sm text-gray-600 mt-1">
                                        From: {tx.from} → To: {tx.to}
                                    </div>
                                </a>
                            {/each}
                        </div>
                    </div>
                </div>
            </div>
        {/each}
    </div>
</div>

<style>
    /* Add any additional styles here */
</style>
