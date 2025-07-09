<template>
	<div>
		<router-view @login-success="handleLogin" @logout="handleLogout" />
	</div>
</template>

<script>
export default {
	data() {
		return {
			isAuthenticated: false,
		};
	},
	created() {
		this.checkAuth();
	},
	methods: {
		checkAuth() {
			const token = localStorage.getItem("token");
			this.isAuthenticated = !!token;

			// Navigate based on auth state
			if (this.isAuthenticated && this.$route.path === '/session') {
				this.$router.push('/chats');
			} else if (!this.isAuthenticated && this.$route.path !== '/session') {
				this.$router.push('/session');
			}
		},
		handleLogin(token) {
			console.log('App.vue received login event with token:', token);

			// Store the token FIRST before updating state
			localStorage.setItem("token", token);

			// Then update the authenticated state
			this.isAuthenticated = true;
			console.log('isAuthenticated set to:', this.isAuthenticated);

			// Navigate to chats after successful login
			// Use nextTick to ensure the localStorage is set before navigation
			this.$nextTick(() => {
				this.$router.push('/chats');
			});
		},
		handleLogout() {
			localStorage.removeItem("token");
			localStorage.removeItem("user");
			this.isAuthenticated = false;

			// Navigate to login after logout
			this.$router.push('/session');
		},
	},
};
</script>

<style>
/* Add your styles here */
</style>
