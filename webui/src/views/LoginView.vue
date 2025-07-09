<template>
<div class="container">
<div class="row justify-content-center mt-5">
	<div class="col-md-6 col-lg-4">
		<div class="card shadow">
			<div class="card-header bg-primary text-white text-center py-3">
				<h4>Login</h4>
			</div>
			<div class="card-body p-4">
				<form @submit.prevent="handleLogin">
					<div class="mb-3">
						<label for="username" class="form-label">Username</label>
						<input
							type="text"
							class="form-control"
							id="username"
							v-model="username"
							placeholder="Enter your username"
							required
							autocomplete="username"
							:disabled="loading"
						>
					</div>

					<ErrorMsg v-if="error" :message="error"/>

					<div class="d-grid gap-2 mt-4">
						<button
							type="submit"
							class="btn btn-primary"
							:disabled="loading || !username.trim()"
						>
							<LoadingSpinner v-if="loading" class="me-2"/>
							<span>Enter</span>
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</div>
</div>
</template>

<script>
export default {
	name: 'LoginView',
	data() {
		return {
			username: '',
			loading: false,
			error: ''
		}
	},
	methods: {
		async handleLogin() {
			// Reset error state
			this.error = '';

			// Validate username
			if (!this.username.trim()) {
				this.error = 'Username is required';
				return;
			}

			try {
				this.loading = true;

				// Make API call to log-in
				const response = await this.$axios.post('/session', {
					username: this.username.trim()
				});

				// Extract token from response
				const token = response.data.token || response.data;
				console.log("token:", token);

				if (!token) {
					throw new Error('No token received from server');
				}

				// Store user data in localStorage (but not the token - that's handled in App.vue)
				const userData = {
					username: this.username,
				};
				localStorage.setItem('user', JSON.stringify(userData));

				// Set default auth header for future requests
				this.$axios.defaults.headers.common['Authorization'] = `${token}`;

				// Emit login success event to parent (App.vue) with the token
				this.$emit('login-success', token);

			} catch (err) {
				console.error('Login failed', err);
				this.error = err.response?.data?.message || err.message || 'Login failed. Please try again.';
			} finally {
				this.loading = false;
			}
		}
	}
}
</script>

<style scoped>
.card {
	border-radius: 0.5rem;
	border: none;
}

.card-header {
	border-radius: 0.5rem 0.5rem 0 0 !important;
}

.form-control:focus {
	box-shadow: 0 0 0 0.25rem rgba(13, 110, 253, 0.25);
	border-color: #86b7fe;
}
</style>
