<script>
	import './layout.css';
	let { children } = $props();
	export const prerender = true;
</script>

{@render children()}
