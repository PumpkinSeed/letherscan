<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	
	let { children } = $props();
</script>

<div class="layout">
	<aside class="sidebar">
		<nav>
			<a href="/" class="menu-item" class:active={$page.url.pathname === '/'}>
				Blocks
			</a>
			<a href="/settings" class="menu-item" class:active={$page.url.pathname === '/settings'}>
				Settings
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
	}

	.content {
		flex: 1;
		padding: 1rem;
	}

	.menu-item {
		display: block;
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
</style>
