<script>
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import Settings from '@lucide/svelte/icons/settings';
	import Plus from '@lucide/svelte/icons/plus';
	import SidebarChannelButton from './SidebarChannelButton.svelte';

	const servers = [
		{
			id: 1123,
			name: 'MyServer',
			channels: [
				{ id: 6, name: 'general', type: 'text', unread: 3 },
				{ id: 7, name: 'off-topic', type: 'text', unread: 0 },
				{ id: 8, name: 'General Voice', type: 'voice', users: 4 }
			]
		},
		{
			id: 312312,
			name: 'other server',
			channels: [
				{ id: 2, name: 'welcome', type: 'text' },
				{ id: 3, name: 'rules', type: 'text' },
				{ id: 4, name: 'announcements', type: 'text' }
			]
		}
	];

	let expanded = $state(new Set([1123, 312312]));
	let selected = $state();

	function toggle(id) {
		if (expanded.has(id)) {
			expanded.delete(id);
		} else {
			expanded.add(id);
		}
		expanded = new Set(expanded);
	}
</script>

<div>
	{#each servers as server (server.id)}
		<div class="flex items-center justify-between border-b bg-muted-soft p-2">
			<button
				class="flex w-full items-center rounded-md hover:bg-muted"
				onclick={() => toggle(server.id)}
			>
				{#if expanded.has(server.id)}
					<ChevronDown class="mr-1 h-3 w-3" />
				{:else}
					<ChevronRight class="mr-1 h-3 w-3" />
				{/if}
				<span class="text-sm font-semibold">{server.name}</span>
			</button>

			<button class="rounded-md p-1 hover:bg-muted">
				<Settings class="h-4 w-4" />
			</button>
		</div>

		{#if expanded.has(server.id)}
			<div class="space-y-1 border-b p-2">
				{#each server.channels as channel (channel.id)}
					<SidebarChannelButton
						{channel}
						selected={selected?.id === channel.id}
						on:select={(e) => (selected = e.detail)}
					/>
				{/each}
			</div>
		{:else if server.channels.some((c) => c.id === selected?.id)}
			<div class="space-y-1 border-b p-2">
				<SidebarChannelButton channel={selected} selected={true} clickable={false} />
			</div>
		{/if}
	{/each}
	<div class="h-2"></div>
</div>
