<template>
	<div>
		<LoadingSpinner v-if="loading" class="d-flex justify-content-center my-5"/>
		<ErrorMsg v-if="error" :message="error" class="mx-3"/>

		<!-- Modern chat list -->
		<div v-if="!loading && !error" class="chat-list">
			<ChatListItem
				v-for="chat in chats"
				:key="chat.id"
				:chat="chat"
				:is-selected="selectedChatId === chat.id"
				:chat-image-url="chatImageUrls[chat.id]"
				:current-username="currentUsername"
				@select="$emit('select-chat', $event)"
				@image-error="$emit('image-error', $event)"
			/>

			<div v-if="chats.length === 0" class="chat-empty-state">
				<div class="chat-empty-icon">
					<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor"
						 stroke-width="1.5">
						<path
							d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
					</svg>
				</div>
				<h6 class="text-muted">No conversations yet</h6>
				<p class="text-muted small">Start chatting by creating a new conversation</p>
			</div>
		</div>
	</div>
</template>

<script>
import ChatListItem from './ChatListItem.vue'

export default {
	name: 'SidebarChatList',
	components: {
		ChatListItem
	},
	props: {
		chats: {
			type: Array,
			required: true
		},
		loading: {
			type: Boolean,
			default: false
		},
		error: String,
		selectedChatId: [String, Number],
		chatImageUrls: {
			type: Object,
			default: () => ({})
		},
		currentUsername: String
	},
	emits: ['select-chat', 'image-error']
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>

