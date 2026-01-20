<script>
	export let active = 'signup';
</script>

<div class="mx-auto flex justify-center">
	<img class="h-7 w-auto sm:h-8" src="https://merakiui.com/images/logo.svg" alt="Logo" />
</div>

<div class="mt-6 flex items-center justify-center">
	<a
		href="/signin"
		data-sveltekit-reload
		class="w-1/3 border-b pb-4 text-center font-medium capitalize transition-all
			{active === 'signin'
			? 'border-b-2 border-blue-500 text-gray-800 dark:border-blue-400 dark:text-white'
			: 'border-transparent text-gray-500 dark:text-gray-300'}"
	>
		sign in
	</a>

	<a
		href="/signup"
		data-sveltekit-reload
		class="w-1/3 border-b pb-4 text-center font-medium capitalize transition-all
			{active === 'signup'
			? 'border-b-2 border-blue-500 text-gray-800 dark:border-blue-400 dark:text-white'
			: 'border-transparent text-gray-500 dark:text-gray-300'}"
	>
		sign up
	</a>
</div>
