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

			if (this.isAuthenticated && this.$route.path === '/session') {
				this.$router.push('/chats');
			} else if (!this.isAuthenticated && this.$route.path !== '/session') {
				this.$router.push('/session');
			}
		},
		handleLogin(token) {
			console.log('App.vue received login event with token:', token);

			localStorage.setItem("token", token);

			this.isAuthenticated = true;
			console.log('isAuthenticated set to:', this.isAuthenticated);

			this.$nextTick(() => {
				this.$router.push('/chats');
			});
		},
		handleLogout() {
			localStorage.removeItem("token");
			localStorage.removeItem("user");
			this.isAuthenticated = false;

			this.$router.push('/session');
		},
	},
};
</script>

<style>
</style>
