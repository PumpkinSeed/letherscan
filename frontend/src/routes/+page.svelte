<script lang="ts">
    import { onMount } from 'svelte';

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

    interface BlocksResponse {
        blocks: Block[];
    }

    let blocks: Block[] = [];
    let loading = true;
    let error: string | null = null;

    async function fetchBlocks() {
        try {
            const response = await fetch('http://localhost:8080/blocks');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data: BlocksResponse = await response.json();
            blocks = data.blocks;
        } catch (e) {
            error = e instanceof Error ? e.message : 'An error occurred while fetching blocks';
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        fetchBlocks();
    });

    function formatTimestamp(timestamp: number): string {
        return new Date(timestamp * 1000).toLocaleString();
    }
</script>

<div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6">Block Explorer</h1>

    {#if loading}
        <div class="flex justify-center items-center h-32">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
    {:else if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <strong class="font-bold">Error!</strong>
            <span class="block sm:inline"> {error}</span>
        </div>
    {:else}
        <div class="space-y-6">
            {#each blocks as block}
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
                                    <div class="bg-gray-50 p-3 rounded">
                                        <div class="flex justify-between items-center">
                                            <span class="font-mono text-sm">{tx.hash}</span>
                                            <span class="text-sm text-gray-500">Gas: {tx.gas}</span>
                                        </div>
                                        <div class="text-sm text-gray-600 mt-1">
                                            From: {tx.from} → To: {tx.to}
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    /* Add any additional styles here */
</style>
