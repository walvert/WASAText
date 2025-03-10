<template>
	<div>
		<LoginView v-if="!isAuthenticated" @login-success="handleLogin" />
		<ConversationView v-else @logout="handleLogout" />
	</div>
</template>

<script>
import LoginView from "@/views/LoginView.vue";
import ConversationView from "@/views/ConversationsView.vue";

export default {
	components: {
		LoginView,
		ConversationView,
	},
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
		},
		handleLogin(token) {
			localStorage.setItem("token", token);
			this.isAuthenticated = true;
		},
		handleLogout() {
			localStorage.removeItem("token");
			this.isAuthenticated = false;
		},
	},
};
</script>

<style>
/* Add your styles here */
</style>
