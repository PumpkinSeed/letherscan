<script lang="ts">
    import { page } from '$app/stores';
    import { fetchWithNodeAddress } from '$lib/utils/fetch';
    import { slide } from 'svelte/transition';

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

    export let data;
    let showDecodeInput = false;
    let contractABI = '';
    let decodedData: { function_name: string; args: Record<string, any> } | null = null;
    let isDecoding = false;
    let decodeError: string | null = null;

    async function handleDecode() {
        if (!contractABI.trim()) {
            decodeError = 'Please enter a contract ABI';
            return;
        }

        try {
            isDecoding = true;
            decodeError = null;
            const response = await fetchWithNodeAddress('http://localhost:8080/decode-contract-call-data', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    contract_abi: contractABI,
                    input_data: data.transaction.input
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            decodedData = await response.json();
        } catch (error) {
            console.error('Error decoding input:', error);
            decodeError = 'Failed to decode input data';
        } finally {
            isDecoding = false;
        }
    }
</script>

<div class="container mx-auto px-4 py-8">
    <div class="mb-6">
        <a href="/" class="text-blue-600 hover:text-blue-800">← Back to Blocks</a>
    </div>

    <h1 class="text-3xl font-bold mb-6">Transaction Details</h1>

    {#if !data.transaction}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
            <strong class="font-bold">Error!</strong>
            <span class="block sm:inline"> Transaction not found</span>
        </div>
    {:else}
        <div class="bg-white shadow-lg rounded-lg overflow-hidden">
            <div class="p-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                        <h2 class="text-xl font-semibold mb-4">Transaction Information</h2>
                        <div class="space-y-4">
                            <div>
                                <p class="text-sm text-gray-600">Hash</p>
                                <p class="font-mono break-all">{data.transaction.hash}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Nonce</p>
                                <p>{data.transaction.nonce}</p>
                            </div>
                        </div>
                    </div>

                    <div>
                        <h2 class="text-xl font-semibold mb-4">Transaction Details</h2>
                        <div class="space-y-4">
                            <div>
                                <p class="text-sm text-gray-600">From</p>
                                <p class="font-mono break-all">{data.transaction.from}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">To</p>
                                <p class="font-mono break-all">{data.transaction.to}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Value</p>
                                <p>{data.transaction.value} wei</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Gas Price</p>
                                <p>{data.transaction.gas_price} wei</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Gas Limit</p>
                                <p>{data.transaction.gas}</p>
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
                                <p>{data.transaction.type}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Method</p>
                                <p>{data.transaction.method}</p>
                            </div>
                            <div>
                                <p class="text-sm text-gray-600">Chain ID</p>
                                <p>{data.transaction.chain_id}</p>
                            </div>
                        </div>
                        <div class="space-y-4">
                            <div>
                                <p class="text-sm text-gray-600">Input Data</p>
                                <p class="font-mono break-all text-sm bg-gray-50 p-2 rounded">{data.transaction.input}</p>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="mt-6">
                    <button
                        on:click={() => showDecodeInput = !showDecodeInput}
                        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors"
                    >
                        {showDecodeInput ? 'Hide Decode Input' : 'Decode Input'}
                    </button>

                    {#if showDecodeInput}
                        <div transition:slide class="mt-4">
                            <div class="mb-4">
                                <label for="contractABI" class="block text-sm font-medium text-gray-700 mb-2">
                                    Contract ABI
                                </label>
                                <textarea
                                    id="contractABI"
                                    bind:value={contractABI}
                                    class="w-full h-32 p-2 border rounded font-mono text-sm"
                                    placeholder="Paste your contract ABI here..."
                                ></textarea>
                            </div>

                            <button
                                on:click={handleDecode}
                                disabled={isDecoding}
                                class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 transition-colors disabled:opacity-50"
                            >
                                {isDecoding ? 'Decoding...' : 'Decode'}
                            </button>

                            {#if decodeError}
                                <p class="mt-2 text-red-500">{decodeError}</p>
                            {/if}

                            {#if decodedData}
                                <div class="mt-4 p-4 bg-gray-50 rounded">
                                    <h3 class="font-semibold mb-2">Decoded Data</h3>
                                    <p class="mb-2"><span class="font-medium">Function:</span> {decodedData.function_name}</p>
                                    <div>
                                        <p class="font-medium mb-1">Arguments:</p>
                                        <pre class="bg-white p-2 rounded overflow-x-auto">{JSON.stringify(decodedData.args, null, 2)}</pre>
                                    </div>
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    {/if}
</div> 