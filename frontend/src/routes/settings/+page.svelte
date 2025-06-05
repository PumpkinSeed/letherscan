<script lang="ts">
	import { nodeAddress } from '$lib/stores/nodeAddress';
	import { numberOfBlocks } from '$lib/stores/numberOfBlocks';
	import { browser } from '$app/environment';

	let localNodeAddress = '';
	let localNumberOfBlocks = 10;

	// Subscribe to the stores
	nodeAddress.subscribe((value) => {
		localNodeAddress = value;
	});

	numberOfBlocks.subscribe((value) => {
		localNumberOfBlocks = value;
	});

	function handleSubmit() {
		nodeAddress.set(localNodeAddress);
		numberOfBlocks.set(localNumberOfBlocks);
		if (browser) {
			alert('Settings saved successfully!');
		}
	}
</script>

<div class="settings-container">
	<h1>Settings</h1>
	
	<form on:submit|preventDefault={handleSubmit} class="settings-form">
		<div class="form-group">
			<label for="nodeAddress">Node Address</label>
			<input
				type="text"
				id="nodeAddress"
				bind:value={localNodeAddress}
				placeholder="Enter node address (e.g., http://localhost:8080)"
				class="input-field"
			/>
		</div>

		<div class="form-group">
			<label for="numberOfBlocks">Number of Blocks to Display</label>
			<input
				type="number"
				id="numberOfBlocks"
				bind:value={localNumberOfBlocks}
				min="1"
				max="100"
				class="input-field"
			/>
			<p class="text-sm text-gray-500 mt-1">Set how many blocks to fetch and display (1-100)</p>
		</div>
		
		<button type="submit" class="save-button">Save Settings</button>
	</form>
</div>

<style>
	.settings-container {
		max-width: 600px;
		margin: 0 auto;
		padding: 2rem;
	}

	h1 {
		margin-bottom: 2rem;
		color: #333;
	}

	.settings-form {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.form-group {
		margin-bottom: 1.5rem;
	}

	label {
		display: block;
		margin-bottom: 0.5rem;
		color: #555;
		font-weight: 500;
	}

	.input-field {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
	}

	.input-field:focus {
		outline: none;
		border-color: #7ea2ee;
		box-shadow: 0 0 0 2px rgba(126, 162, 238, 0.2);
	}

	.save-button {
		background-color: #7ea2ee;
		color: white;
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.save-button:hover {
		background-color: #6b8fd8;
	}
</style> 