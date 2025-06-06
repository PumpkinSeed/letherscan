<script lang="ts">
	import { onMount } from 'svelte';
	import { BASE_URL } from '$lib/config';

	let address = '';
	let data = '';
	let result = '';
	let loading = false;
	let error = '';
	let contractABI = '';
	let stateMutability = '';
	let parsedMethods: string[] = [];
	let selectedMethod: string | null = null;
	let showMethodSelection = true;
	let inputValues: string[] = [];
	let decodedResult: Record<string, any> | null = null;

	const stateMutabilityOptions = [
		{ value: '', label: 'All' },
		{ value: 'view', label: 'View' },
		{ value: 'pure', label: 'Pure' },
		{ value: 'payable', label: 'Payable' }
	];

	function selectMethod(method: string) {
		selectedMethod = method;
		// Extract input parameters from method signature
		const match = method.match(/\((.*?)\)/);
		if (match && match[1]) {
			const params = match[1].split(',').map(p => p.trim());
			inputValues = new Array(params.length).fill('');
		} else {
			inputValues = [];
		}
		showMethodSelection = false;
	}

	function startOver() {
		showMethodSelection = true;
		selectedMethod = null;
		inputValues = [];
		result = '';
		decodedResult = null;
	}

	async function handleSubmit() {
		if (!address || !selectedMethod || !contractABI) {
			error = 'Please provide all required fields';
			return;
		}

		loading = true;
		error = '';
		try {
			const response = await fetch(`${BASE_URL}/eth-call`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					method: selectedMethod,
					contract_address: address,
					contract_abi: contractABI,
					input: inputValues
				})
			});

			if (!response.ok) {
				throw new Error('Failed to execute eth_call');
			}

			const responseData = await response.json();
			result = responseData.raw_response;
			decodedResult = responseData.decoded;
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'An unknown error occurred';
		} finally {
			loading = false;
		}
	}

	async function handleParseABI() {
		if (!contractABI) {
			error = 'Please provide a contract ABI';
			return;
		}

		loading = true;
		error = '';
		selectedMethod = null;
		try {
			const response = await fetch(`${BASE_URL}/parse-contract-abi`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					contract_abi: contractABI,
					state_mutability_filter: stateMutability
				})
			});

			if (!response.ok) {
				throw new Error('Failed to parse contract ABI');
			}

			const responseData = await response.json();
			parsedMethods = responseData.methods;
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'An unknown error occurred';
		} finally {
			loading = false;
		}
	}
</script>

<div class="container">
	<h1>ETH Call</h1>
	<p class="description">
		Execute an eth_call to read data from the blockchain without creating a transaction.
	</p>

	{#if showMethodSelection}
		<div class="abi-section">
			<div class="form-group">
				<label for="contractABI">Contract ABI</label>
				<textarea
					id="contractABI"
					bind:value={contractABI}
					placeholder="Paste your contract ABI here..."
					rows="10"
					required
				></textarea>
			</div>

			<div class="form-group">
				<label for="stateMutability">State Mutability</label>
				<select id="stateMutability" bind:value={stateMutability}>
					{#each stateMutabilityOptions as option}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
			</div>

			<button type="button" on:click={handleParseABI} disabled={loading} class="parse-button">
				{loading ? 'Parsing...' : 'Parse ABI'}
			</button>

			{#if parsedMethods.length > 0}
				<div class="methods-container">
					<h3>Available Methods</h3>
					<div class="methods-list">
						{#each parsedMethods as method}
							<button
								class="method-item {selectedMethod === method ? 'selected' : ''}"
								on:click={() => selectMethod(method)}
							>
								{method}
							</button>
						{/each}
					</div>
				</div>
			{/if}
		</div>
	{:else}
		<div class="input-section">
			<button type="button" on:click={startOver} class="start-over-button">
				Start Over
			</button>

			<form on:submit|preventDefault={handleSubmit} class="form">
				<div class="form-group">
					<label for="address">Contract Address</label>
					<input
						type="text"
						id="address"
						bind:value={address}
						placeholder="0x..."
						required
					/>
				</div>

				{#if selectedMethod}
					<div class="selected-method">
						<h3>Selected Method</h3>
						<pre>{selectedMethod}</pre>
					</div>

					{#if inputValues.length > 0}
						<div class="input-parameters">
							<h3>Input Parameters</h3>
							{#each inputValues as value, i}
								<div class="form-group">
									<label for="input-{i}">Parameter {i + 1}</label>
									<input
										type="text"
										id="input-{i}"
										bind:value={inputValues[i]}
										placeholder="Enter parameter value..."
									/>
								</div>
							{/each}
						</div>
					{/if}
				{/if}

				<button type="submit" disabled={loading}>
					{loading ? 'Executing...' : 'Execute Call'}
				</button>
			</form>
		</div>
	{/if}

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if result}
		<div class="result">
			<h2>Result</h2>
			{#if decodedResult}
				<div class="decoded-result">
					<h3>Decoded Response</h3>
					<pre>{JSON.stringify(decodedResult, null, 2)}</pre>
				</div>
			{/if}
			<div class="raw-result">
				<h3>Raw Response</h3>
				<pre>{result}</pre>
			</div>
		</div>
	{/if}
</div>

<style>
	.container {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	.description {
		color: #666;
		margin-bottom: 2rem;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	label {
		font-weight: 500;
	}

	input {
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 0.5rem;
		font-family: monospace;
	}

	button {
		padding: 0.75rem 1.5rem;
		background-color: #7ea2ee;
		color: white;
		border: none;
		border-radius: 0.5rem;
		cursor: pointer;
		font-weight: 500;
		transition: background-color 0.2s;
	}

	button:hover:not(:disabled) {
		background-color: #6b8fd8;
	}

	button:disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}

	.error {
		margin-top: 1rem;
		padding: 1rem;
		background-color: #fee;
		color: #c00;
		border-radius: 0.5rem;
	}

	.result {
		margin-top: 2rem;
		padding: 1rem;
		background-color: #f5f5f5;
		border-radius: 0.5rem;
	}

	.result h2 {
		margin-top: 0;
		margin-bottom: 1rem;
	}

	pre {
		margin: 0;
		white-space: pre-wrap;
		word-break: break-all;
	}

	.abi-section {
		margin-bottom: 2rem;
		padding: 1.5rem;
		background-color: #f8f9fa;
		border-radius: 0.5rem;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	textarea {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 0.5rem;
		font-family: monospace;
		resize: vertical;
		background-color: white;
	}

	select {
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 0.5rem;
		background-color: white;
	}

	.parse-button {
		margin-top: 1.5rem;
	}

	.methods-container {
		margin-top: 2rem;
		padding-top: 1.5rem;
		border-top: 1px solid #eee;
		max-width: 100%;
		overflow-x: hidden;
	}

	.methods-container h3 {
		margin: 0 0 1rem 0;
		color: #333;
	}

	.methods-list {
		display: grid;
		grid-template-columns: 1fr;
		gap: 0.75rem;
	}

	.method-item {
		text-align: left;
		padding: 1rem;
		background-color: white;
		border: 1px solid #ddd;
		color: #333;
		transition: all 0.2s;
		font-family: monospace;
		font-size: 0.9rem;
		line-height: 1.5;
		white-space: pre-wrap;
		word-break: break-word;
		height: 100%;
		display: flex;
		align-items: center;
	}

	.method-item:hover {
		background-color: #f0f4ff;
		border-color: #7ea2ee;
		transform: translateY(-1px);
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.method-item.selected {
		background-color: #7ea2ee;
		border-color: #7ea2ee;
		color: white;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.input-section {
		margin-bottom: 2rem;
		padding: 1.5rem;
		background-color: #f8f9fa;
		border-radius: 0.5rem;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	.start-over-button {
		margin-bottom: 1.5rem;
		background-color: #6c757d;
	}

	.start-over-button:hover {
		background-color: #5a6268;
	}

	.selected-method {
		margin-bottom: 1.5rem;
		padding: 1rem;
		background-color: white;
		border-radius: 0.5rem;
		border: 1px solid #ddd;
	}

	.selected-method h3 {
		margin: 0 0 0.5rem 0;
		color: #333;
	}

	.input-parameters {
		margin-bottom: 1.5rem;
	}

	.input-parameters h3 {
		margin: 0 0 1rem 0;
		color: #333;
	}

	.decoded-result {
		margin-bottom: 1.5rem;
		padding: 1rem;
		background-color: white;
		border-radius: 0.5rem;
		border: 1px solid #ddd;
	}

	.raw-result {
		padding: 1rem;
		background-color: white;
		border-radius: 0.5rem;
		border: 1px solid #ddd;
	}
</style> 