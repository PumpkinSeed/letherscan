<script lang="ts">
    import { page } from '$app/stores';
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

    let transaction: Transaction | null = null;
    let loading = true;
    let error: string | null = null;

    async function fetchTransaction() {
        try {
            const response = await fetch(`http://localhost:8080/transaction/${$page.params.hash}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            transaction = await response.json();
        } catch (e) {
            error = e instanceof Error ? e.message : 'An error occurred while fetching transaction';
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        fetchTransaction();
    });
</script>

<div class="container mx-auto px-4 py-8">
    <div class="mb-6">
        <a href="/" class="text-blue-600 hover:text-blue-800">← Back to Blocks</a>
    </div>

    <h1 class="text-3xl font-bold mb-6">Transaction Details</h1>

    {#if loading}
        <div class="flex justify-center items-center h-32">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
    {:else if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <strong class="font-bold">Error!</strong>
            <span class="block sm:inline"> {error}</span>
        </div>
    {:else if transaction}
        <div class="bg-white shadow-lg rounded-lg overflow-hidden">
            <div class="p-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                        <h2 class="text-xl font-semibold mb-4">Transaction Information</h2>
                        <div class="space-y-4">
                            <div>
                                <p class="text-sm text-gray-600">Hash</p>
                                <p class="font-mono break-all">{transaction.hash}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Nonce</p>
                                <p>{transaction.nonce}</p>
                            </div>
                        </div>
                    </div>

                    <div>
                        <h2 class="text-xl font-semibold mb-4">Transaction Details</h2>
                        <div class="space-y-4">
                            <div>
                                <p class="text-sm text-gray-600">From</p>
                                <p class="font-mono break-all">{transaction.from}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">To</p>
                                <p class="font-mono break-all">{transaction.to}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Value</p>
                                <p>{transaction.value} wei</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Gas Price</p>
                                <p>{transaction.gas_price} wei</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Gas Limit</p>
                                <p>{transaction.gas}</p>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="mt-8">
                    <h2 class="text-xl font-semibold mb-4">Additional Information</h2>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div class="space-y-4">
                            <div>
                                <p class="text-sm text-gray-600">Type</p>
                                <p>{transaction.type}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Method</p>
                                <p>{transaction.method}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Chain ID</p>
                                <p>{transaction.chain_id}</p>
                            </div>
                        </div>
                        <div class="space-y-4">
                            <div>
                                <p class="text-sm text-gray-600">Input Data</p>
                                <p class="font-mono break-all text-sm bg-gray-50 p-2 rounded">{transaction.input}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div> 