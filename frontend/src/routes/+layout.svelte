<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	
	let { children } = $props();
</script>

<div class="layout">
	<aside class="sidebar">
		<div class="logo-container">
			<img src="/logo.png" alt="Leatherscan Logo" class="logo" />
		</div>
		<nav class="nav-content">
			<a href="/" class="menu-item" class:active={$page.url.pathname === '/'}>
				Blocks
			</a>
			<a href="/settings" class="menu-item" class:active={$page.url.pathname === '/settings'}>
				Settings
			</a>
			<a href="/eth-call" class="menu-item" class:active={$page.url.pathname === '/eth-call'}>
				ETH Call
			</a>
			<a href="/send-transaction" class="menu-item" class:active={$page.url.pathname === '/send-transaction'}>
				Send Transaction
			</a>
			<a 
				href="https://github.com/PumpkinSeed/letherscan" 
				target="_blank" 
				rel="noopener noreferrer" 
				class="menu-item github-link"
				title="View on GitHub"
			>
				<svg 
					xmlns="http://www.w3.org/2000/svg" 
					width="20" 
					height="20" 
					viewBox="0 0 24 24" 
					fill="currentColor"
				>
					<path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
				</svg>
				<span>GitHub</span>
			</a>
		</nav>
	</aside>
	<main class="content">
		{@render children()}
	</main>
</div>

<style>
	.layout {
		display: flex;
		min-height: 100vh;
	}

	.sidebar {
		width: 250px;
		background-color: #f3f3f3;
		padding: 1rem;
		display: flex;
		flex-direction: column;
	}

	.logo-container {
		padding: 1rem;
		margin-bottom: 1rem;
		display: flex;
		justify-content: center;
	}

	.logo {
		max-width: 100%;
		height: auto;
	}

	.nav-content {
		flex: 1;
	}

	.content {
		flex: 1;
		padding: 1rem;
	}

	.menu-item {
		display: flex;
		align-items: center;
		padding: 0.75rem 1rem;
		color: #464646;
		text-decoration: none;
		border-radius: 0.5rem;
		margin-bottom: 0.5rem;
		transition: background-color 0.2s;
	}

	.menu-item:hover {
		background-color: #a3a3a3;
	}

	.menu-item.active {
		background-color: #7ea2ee;
	}

	.github-link {
		gap: 0.5rem;
	}

	.github-link:hover {
		background-color: #a3a3a3;
	}
</style>
