<script>
	export let value = '';
	export let type = 'text';
	export let placeholder = '';
	export let error = null;
	export let icon = null;
	export let containerClass = '';
</script>

<div class="relative flex items-center {containerClass}">
	{#if icon}
		<span class="absolute">
			<svelte:component this={icon} class="mx-3 h-6 w-6 text-gray-300 dark:text-gray-500" />
		</span>
	{/if}

	<input
		{type}
		{placeholder}
		bind:value
		class="focus:ring-opacity-40 block w-full rounded-lg border bg-white py-3 text-gray-700 focus:ring focus:ring-blue-300 focus:outline-none dark:bg-gray-900 dark:text-gray-300"
		class:px-11={!!icon}
		class:px-4={!icon}
		class:border-gray-300={!error}
		class:dark:border-gray-600={!error}
		class:border-red-600={!!error}
		class:focus:border-blue-400={!error}
		class:dark:focus:border-blue-300={!error}
		class:focus:border-red-500={!!error}
		aria-invalid={!!error}
	/>
</div>

{#if error}
	<p class="mt-2 text-sm text-red-600 dark:text-red-400">{error}</p>
{/if}
