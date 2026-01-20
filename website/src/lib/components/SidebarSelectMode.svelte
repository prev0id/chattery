<script>
	import { Users } from '@lucide/svelte';
	import Icon from './Icon.svelte';

	let showTooltip = $state({ users: false, groups: false, info: false });
</script>

<div class="flex w-16 flex-col items-center space-y-8 border-r py-8">
	<div class="size-14">
		<Icon></Icon>
	</div>
	<div
		class="relative inline-block"
		onmouseenter={() => (showTooltip.users = true)}
		onmouseleave={() => (showTooltip.users = false)}
	>
		<Users class="size-10"></Users>
	</div>
	{#if showTooltip.users}
		<div
			id="tooltip-02-top"
			role="tooltip"
			class="bg-foreground absolute bottom-full left-1/2 mb-2 -translate-x-1/2 rounded-md px-3 py-1.5 text-xs whitespace-nowrap text-background shadow-md"
		>
			ChatWithUsers
			<div
				aria-hidden="true"
				class="border-t-foreground absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent"
			></div>
		</div>
	{/if}
</div>
