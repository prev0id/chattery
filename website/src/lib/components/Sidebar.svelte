<script>
	import SidebarAddChannel from './SidebarAddChannel.svelte';
	import SidebarChannels from './SidebarChannels.svelte';
	import SidebarSearch from './SidebarSearch.svelte';
	import SidebarSelectMode from './SidebarSelectMode.svelte';
</script>

<aside class="m-2 flex w-80 rounded-lg border">
	<SidebarSelectMode></SidebarSelectMode>
	<div class="w-64">
		<SidebarSearch></SidebarSearch>
		<SidebarAddChannel></SidebarAddChannel>
		<SidebarChannels></SidebarChannels>
	</div>
</aside>
