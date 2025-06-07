<script lang="ts">
	import { onMount } from 'svelte';
	import { BASE_URL } from '$lib/config';

	interface Method {
		name: string;
		state_mutability: string;
		inputs: string[];
		outputs: string[];
	}

	let address = '';
	let result = '';
	let loading = false;
	let error = '';
	let contractABI = '';
	let stateMutability = '';
	let parsedMethods: Method[] = [];
	let selectedMethod: Method | null = null;
	let showMethodSelection = true;
	let inputValues: string[] = [];
	let privateKey = '';

	const stateMutabilityOptions = [
		{ value: '', label: 'All' },
		{ value: 'view', label: 'View' },
		{ value: 'pure', label: 'Pure' },
		{ value: 'payable', label: 'Payable' }
	];

	function selectMethod(method: Method) {
		selectedMethod = method;
		inputValues = new Array(method?.inputs?.length || 0).fill('');
		showMethodSelection = false;
	}

	function startOver() {
		showMethodSelection = true;
		selectedMethod = null;
		inputValues = [];
		result = '';
	}

	async function handleSubmit() {
		if (!address || !selectedMethod?.name || !contractABI || !privateKey) {
			error = 'Please provide all required fields';
			return;
		}

		loading = true;
		error = '';
		try {
			const response = await fetch(`${BASE_URL}/send-transaction`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					method: selectedMethod.name,
					contract_address: address,
					contract_abi: contractABI,
					private_key: privateKey,
					input: inputValues
				})
			});

			if (!response.ok) {
				throw new Error('Failed to send transaction');
			}

			const responseData = await response.json();
			result = responseData.transaction_hash;
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
	<h1>Send Transaction</h1>
	<p class="description">
		Send a transaction to interact with a smart contract on the blockchain.
	</p>

	{#if !selectedMethod}
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
								class="method-item {selectedMethod?.name === method.name ? 'selected' : ''}"
								on:click={() => selectMethod(method)}
							>
								<div class="method-header">
									<span class="method-name">{method.name}</span>
									<span class="method-mutability">{method.state_mutability}</span>
								</div>
								{#if method.inputs?.length > 0}
									<div class="method-inputs">
										<strong>Inputs:</strong> {method.inputs.join(', ')}
									</div>
								{/if}
								{#if method.outputs?.length > 0}
									<div class="method-outputs">
										<strong>Outputs:</strong> {method.outputs.join(', ')}
									</div>
								{/if}
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

				<div class="form-group">
					<label for="privateKey">Private Key (without 0x prefix)</label>
					<input
						type="password"
						id="privateKey"
						bind:value={privateKey}
						placeholder="Enter your private key..."
						required
					/>
				</div>

				<div class="selected-method">
					<h3>Selected Method</h3>
					<div class="method-details">
						<div class="method-header">
							<span class="method-name">{selectedMethod.name}</span>
							<span class="method-mutability">{selectedMethod.state_mutability}</span>
						</div>
						{#if selectedMethod.inputs?.length > 0}
							<div class="method-inputs">
								<strong>Inputs:</strong> {selectedMethod.inputs.join(', ')}
							</div>
						{/if}
						{#if selectedMethod.outputs?.length > 0}
							<div class="method-outputs">
								<strong>Outputs:</strong> {selectedMethod.outputs.join(', ')}
							</div>
						{/if}
					</div>
				</div>

				{#if selectedMethod.inputs?.length > 0}
					<div class="input-parameters">
						<h3>Input Parameters</h3>
						{#each selectedMethod.inputs as input, i}
							<div class="form-group">
								<label for="input-{i}">{input}</label>
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

				<button type="submit" disabled={loading}>
					{loading ? 'Sending...' : 'Send Transaction'}
				</button>
			</form>
		</div>
	{/if}

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if result}
		<div class="result">
			<h2>Transaction Hash</h2>
			<pre>{result}</pre>
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

	.method-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.method-name {
		font-weight: 600;
		font-size: 1.1rem;
	}

	.method-mutability {
		font-size: 0.8rem;
		padding: 0.25rem 0.5rem;
		background-color: #e9ecef;
		border-radius: 0.25rem;
		color: #495057;
	}

	.method-inputs,
	.method-outputs {
		font-size: 0.9rem;
		color: #666;
		margin-top: 0.5rem;
	}

	.method-details {
		padding: 1rem;
		background-color: #f8f9fa;
		border-radius: 0.25rem;
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
		flex-direction: column;
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

	.method-item.selected .method-mutability {
		background-color: rgba(255, 255, 255, 0.2);
		color: white;
	}

	.method-item.selected .method-inputs,
	.method-item.selected .method-outputs {
		color: rgba(255, 255, 255, 0.9);
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
</style> 