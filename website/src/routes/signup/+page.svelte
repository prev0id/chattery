<script>
	import AuthHeader from '$lib/components/AuthHeader.svelte';
	import AuthInput from '$lib/components/AuthInput.svelte';
	import { User, Mail, Lock } from '@lucide/svelte';

	let username = '';
	let email = '';
	let password = '';
	let confirmPassword = '';

	let errors = {
		username: null,
		email: null,
		password: null,
		confirmPassword: null
	};

	let submitting = false;

	const validateForm = () => {
		errors = {
			username: null,
			email: null,
			password: null,
			confirmPassword: null
		};

		let isValid = true;

		if (!username.trim()) {
			errors.username = 'Username is required.';
			isValid = false;
		}

		if (!email.trim()) {
			errors.email = 'Email address is required.';
			isValid = false;
		} else if (!/^\S+@\S+\.\S+$/.test(email)) {
			errors.email = 'Please enter a valid email address.';
			isValid = false;
		}

		if (!password) {
			errors.password = 'Password is required.';
			isValid = false;
		} else if (password.length < 8) {
			errors.password = 'Password must be at least 8 characters long.';
			isValid = false;
		}

		if (confirmPassword !== password) {
			errors.confirmPassword = 'Passwords do not match.';
			isValid = false;
		}

		return isValid;
	};

	const handleSubmit = async (e) => {
		e.preventDefault();

		if (validateForm()) {
			submitting = true;

			// TODO: Replace with real API call
			console.log('Signup successful:', { username, email, password });
			alert('Sign up successful!');

			submitting = false;
		}
	};
</script>

<section class="bg-white dark:bg-gray-900">
	<div class="container mx-auto flex min-h-screen items-center justify-center px-6">
		<form class="w-full max-w-md" on:submit={handleSubmit} novalidate>
			<AuthHeader active="signup" />

			<AuthInput
				containerClass="mt-8"
				icon={User}
				bind:value={username}
				placeholder="Username"
				error={errors.username}
			/>

			<AuthInput
				containerClass="mt-4"
				icon={Mail}
				type="email"
				bind:value={email}
				placeholder="Email address"
				error={errors.email}
			/>

			<AuthInput
				containerClass="mt-4"
				icon={Lock}
				type="password"
				bind:value={password}
				placeholder="Password"
				error={errors.password}
			/>

			<AuthInput
				containerClass="mt-4"
				icon={Lock}
				type="password"
				bind:value={confirmPassword}
				placeholder="Confirm Password"
				error={errors.confirmPassword}
			/>

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
					<a href="/signin" class="text-sm text-blue-500 hover:underline dark:text-blue-400">
						Already have an account?
					</a>
				</div>
			</div>
		</form>
	</div>
</section>
