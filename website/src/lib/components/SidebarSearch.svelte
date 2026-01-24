<script>
	import { Search, ArrowRight } from '@lucide/svelte';

	let value = $state('');
</script>

<form class="flex w-full gap-2 border-b p-2">
	<div
		class="bg-background flex h-8 min-w-0 flex-1 items-center gap-2 rounded-md border px-2 text-sm"
	>
		<Search class="h-4 w-4 shrink-0" />
		<input
			type="search"
			bind:value
			placeholder="Search..."
			class="flex-1 bg-transparent outline-none"
		/>
	</div>
	<button
		type="submit"
		class="flex h-8 items-center justify-center rounded-md bg-black px-2 text-white hover:bg-black/90"
	>
		<ArrowRight class="h-4 w-4" />
	</button>
</form>
