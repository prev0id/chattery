<script lang="ts">
	import Hash from '@lucide/svelte/icons/hash';
	import Volume2 from '@lucide/svelte/icons/volume-2';

	export let channel;
	export let selected = false;
	export let clickable = true;

	function handleClick() {
		if (clickable) {
			dispatch('select', channel);
		}
	}

	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
</script>

<button
	disabled={!clickable}
	on:click={handleClick}
	class="
		flex w-full items-center gap-1.5 rounded-md px-2 py-1 text-sm
		hover:bg-muted
		{selected ? 'bg-muted-soft text-info' : ''}
		{!clickable ? 'cursor-default' : ''}
	"
>
	{#if channel.type === 'voice'}
		<Volume2 class="h-4 w-4 shrink-0" />
	{:else}
		<Hash class="h-4 w-4 shrink-0" />
	{/if}

	<span class="flex-1 truncate text-left">{channel.name}</span>

	{#if channel.unread}
		<span
			class="flex h-4 min-w-4 items-center justify-center rounded-full bg-danger px-1 text-xs font-bold text-background"
		>
			{channel.unread}
		</span>
	{/if}

	{#if channel.users}
		<span class="text-xs">{channel.users}</span>
	{/if}
</button>
