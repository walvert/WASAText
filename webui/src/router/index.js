import {createRouter, createWebHashHistory, createWebHistory} from 'vue-router'
import LoginView from "../views/LoginView.vue";
import ConversationsView from "../views/ConversationsView.vue";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', redirect: '/session'},
		{path: '/session', component: LoginView},
		{path: '/chats', component: ConversationsView},
	]
})

export default router
