<script>
	import AuthHeader from '$lib/components/AuthHeader.svelte';
	import AuthInput from '$lib/components/AuthInput.svelte';
	import { Mail, Lock } from '@lucide/svelte';

	let submitting = false;

	const handleSubmit = async () => {
		submitting = true;

		console.log('Signin successful:');
		alert('Sign in successful!');

		submitting = false;
	};
</script>

<section class="bg-white dark:bg-gray-900">
	<div class="container mx-auto flex min-h-screen items-center justify-center px-6">
		<form class="w-full max-w-md" on:submit={handleSubmit} novalidate>
			<AuthHeader active="signup" />

			<AuthInput containerClass="mt-4" icon={Mail} type="email" placeholder="Email address" />

			<AuthInput containerClass="mt-4" icon={Lock} type="password" placeholder="Password" />

			<AuthInput containerClass="mt-4" icon={Lock} type="password" placeholder="Confirm Password" />

			<div class="mt-6">
				<button
					type="submit"
					disabled={submitting}
					class="focus:ring-opacity-50 w-full transform rounded-lg bg-blue-500 px-6 py-3 text-sm font-medium tracking-wide text-white capitalize transition-colors duration-300 hover:bg-blue-400 focus:ring focus:ring-blue-300 focus:outline-none {submitting
						? 'cursor-not-allowed opacity-70'
						: ''}"
				>
					{submitting ? 'Signing up...' : 'Sign Up'}
				</button>

				<div class="mt-6 text-center">
					<a href="/signup" class="text-sm text-blue-500 hover:underline dark:text-blue-400">
						Don't have an account?
					</a>
				</div>
			</div>
		</form>
	</div>
</section>
