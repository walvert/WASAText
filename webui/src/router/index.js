import {createRouter, createWebHistory} from 'vue-router'
import LoginView from "../views/LoginView.vue";
import ChatsView from "../views/ChatsView.vue";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', redirect: '/session'},
		{
			path: '/session',
			component: LoginView,
			meta: { requiresGuest: true }
		},
		{
			path: '/chats',
			component: ChatsView,
			meta: { requiresAuth: true }
		},
	]
})

// Navigation guard to protect routes
router.beforeEach((to, from, next) => {
	const token = localStorage.getItem('token');
	const isAuthenticated = !!token;

	if (to.meta.requiresAuth && !isAuthenticated) {
		// Redirect to login if trying to access protected route without auth
		next('/session');
	} else if (to.meta.requiresGuest && isAuthenticated) {
		// Redirect to chats if trying to access login while authenticated
		next('/chats');
	} else {
		next();
	}
});

export default router
