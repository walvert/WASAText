<template>
	<div class="col-md-4 col-lg-3 chat-sidebar">
		<div class="sidebar-sticky pt-3">
			<SidebarHeader
				:current-user-image-url="currentUserImageUrl"
				:current-username="currentUsername"
				@open-profile="$emit('open-profile')"
				@new-chat="$emit('new-chat')"
				@logout="$emit('logout')"
				@user-image-error="$emit('user-image-error')"
			/>

			<SidebarChatList
				:chats="chats"
				:loading="loading"
				:error="error"
				:selected-chat-id="selectedChatId"
				:chat-image-urls="chatImageUrls"
				:current-username="currentUsername"
				@select-chat="$emit('select-chat', $event)"
				@image-error="$emit('image-error', $event)"
			/>
		</div>
	</div>
</template>

<script>
import SidebarHeader from './SidebarHeader.vue'
import SidebarChatList from './SidebarChatList.vue'

export default {
	name: 'Sidebar',
	components: {
		SidebarHeader,
		SidebarChatList
	},
	props: {
		// User data
		currentUserImageUrl: String,
		currentUsername: String,

		// Chat data
		chats: {
			type: Array,
			required: true
		},
		selectedChatId: [String, Number],
		chatImageUrls: {
			type: Object,
			default: () => ({})
		},

		// UI state
		loading: {
			type: Boolean,
			default: false
		},
		error: String
	},
	emits: [
		'open-profile',
		'new-chat',
		'logout',
		'user-image-error',
		'select-chat',
		'image-error'
	]
}
</script>

<style scoped>
@import "../views/ChatsView.css";

.chat-sidebar {
	height: 100vh;
	border-right: 1px solid #dee2e6;
	overflow-y: auto;
	background-color: #f8f9fa;
}

</style>
