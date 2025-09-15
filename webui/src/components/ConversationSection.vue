<template>
	<div v-if="!selectedChatId" class="d-flex justify-content-center align-items-center h-100">
		<div class="text-center text-muted">
			<div class="mb-3">
				<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor"
					 stroke-width="1.5">
					<path
						d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
				</svg>
			</div>
			<h5>Select a conversation to start chatting</h5>
			<p class="mb-4">Or create a new chat to get started</p>
			<button class="btn btn-primary" @click="$emit('open-new-chat-modal')">New Chat</button>
		</div>
	</div>

	<div v-else-if="selectedChat" class="chat-container h-100 d-flex flex-column">
		<ChatHeader
			:selected-chat="selectedChat"
			:chat-image-url="selectedChat && selectedChat.id ? chatImageUrls[selectedChat.id] : null"
			:current-username="currentUsername"
			:chat-members="chatMembers"
			:loading-chat-members="loadingChatMembers"
			:member-image-urls="memberImageUrls"
			@image-error="$emit('handle-image-error', $event)"
			@member-image-error="$emit('handle-member-image-error', $event)"
			@open-rename-group-modal="$emit('open-rename-group-modal')"
			@open-add-to-group-modal="$emit('open-add-to-group-modal')"
			@open-set-group-photo-modal="$emit('open-set-group-photo-modal')"
			@leave-group="$emit('leave-group')"
			@get-group-members="$emit('get-group-members')"
		/>

		<MessagesSection
			v-if="!loadingMessages && !messagesError"
			class="messages-section"
			:messages="messages"
			:selected-chat="selectedChat"
			:current-username="currentUsername"
			:last-read-message-id="lastReadMessageId"
			:message-image-urls="messageImageUrls"
			:sending-message="sendingMessage"
			:pending-message="pendingMessage"
			:deleting-message="deletingMessage"
			ref="messagesSection"
			@message-image-error="$emit('handle-message-image-error', $event)"
			@reply-image-error="$emit('handle-reply-image-error', $event)"
			@toggle-message-like="$emit('toggle-message-like', $event)"
			@start-reply="$emit('start-reply', $event)"
			@open-forward-modal="$emit('open-forward-modal', $event)"
			@confirm-delete="$emit('confirm-delete-message', $event)"
		/>

		<!-- Loading Messages State -->
		<div v-if="loadingMessages" class="d-flex justify-content-center align-items-center flex-grow-1">
			<div class="text-center text-muted">
				<div class="spinner-border mb-3" role="status">
					<span class="visually-hidden">Loading messages...</span>
				</div>
			</div>
		</div>

		<!-- Messages Error State -->
		<div v-if="messagesError" class="d-flex justify-content-center align-items-center flex-grow-1">
			<div class="text-center text-muted">
				<div class="mb-3">
					<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
						<circle cx="12" cy="12" r="10"/>
						<line x1="12" y1="8" x2="12" y2="12"/>
						<line x1="12" y1="16" x2="12.01" y2="16"/>
					</svg>
				</div>
				<h6>Failed to load messages</h6>
				<p class="text-muted small">{{ messagesError }}</p>
				<button class="btn btn-sm btn-outline-primary" @click="$emit('retry-load-messages')">
					Try Again
				</button>
			</div>
		</div>

		<MessageInput
			:new-message="newMessage"
			:sending-message="sendingMessage"
			:replying-to-message="replyingToMessage"
			:selected-message-image="selectedMessageImage"
			:message-image-preview-url="messageImagePreviewUrl"
			ref="messageInput"
			@send-message="$emit('send-message')"
			@clear-reply="$emit('clear-reply')"
			@open-image-modal="$emit('open-image-modal')"
			@clear-message-image-selection="$emit('clear-message-image-selection')"
			@update:new-message="$emit('update:new-message', $event)"
			@focus-input="handleFocusInput"
		/>
	</div>

	<!-- Fallback loading state if selectedChat is still loading -->
	<div v-else class="d-flex justify-content-center align-items-center h-100">
		<div class="text-center text-muted">
			<div class="spinner-border mb-3" role="status">
				<span class="visually-hidden">Loading...</span>
			</div>
			<p>Loading chat...</p>
		</div>
	</div>
</template>

<script>
import ChatHeader from './ChatHeader.vue'
import MessagesSection from './MessagesSection.vue'
import MessageInput from './MessageInput.vue'

export default {
	name: 'ConversationSection',
	components: {
		ChatHeader,
		MessagesSection,
		MessageInput
	},

	props: {
		// Chat selection
		selectedChatId: {
			type: [Number, String],
			default: null
		},
		selectedChat: {
			type: Object,
			default: null
		},
		currentUsername: {
			type: String,
			required: true
		},

		// Messages
		messages: {
			type: Array,
			default: () => []
		},
		loadingMessages: {
			type: Boolean,
			default: false
		},
		messagesError: {
			type: String,
			default: null
		},
		lastReadMessageId: {
			type: [Number, String],
			default: null
		},

		// Message input
		newMessage: {
			type: String,
			default: ''
		},
		sendingMessage: {
			type: Boolean,
			default: false
		},
		pendingMessage: {
			type: String,
			default: ''
		},
		replyingToMessage: {
			type: Object,
			default: null
		},
		selectedMessageImage: {
			type: File,
			default: null
		},
		messageImagePreviewUrl: {
			type: String,
			default: null
		},
		deletingMessage: {
			type: [Number, String],
			default: null
		},

		// Images
		chatImageUrls: {
			type: Object,
			default: () => ({})
		},
		messageImageUrls: {
			type: Object,
			default: () => ({})
		},
		memberImageUrls: {
			type: Object,
			default: () => ({})
		},

		// Chat members
		chatMembers: {
			type: Array,
			default: () => []
		},
		loadingChatMembers: {
			type: Boolean,
			default: false
		}
	},

	emits: [
		// Navigation
		'open-new-chat-modal',
		'select-chat',

		// Message actions
		'send-message',
		'toggle-message-like',
		'start-reply',
		'clear-reply',
		'open-forward-modal',
		'confirm-delete-message',
		'retry-load-messages',

		// Input actions
		'update:new-message',
		'open-image-modal',
		'clear-message-image-selection',

		// Group actions
		'open-rename-group-modal',
		'open-add-to-group-modal',
		'open-set-group-photo-modal',
		'leave-group',
		'get-group-members',

		// Media and error handling
		'handle-image-error',
		'handle-message-image-error',
		'handle-reply-image-error',
		'handle-member-image-error'
	],

	data() {
		return {
			messageInputFocusFunction: null
		}
	},

	methods: {
		handleFocusInput(focusFunction) {
			this.messageInputFocusFunction = focusFunction
		},

		focusMessageInput() {
			if (this.messageInputFocusFunction) {
				this.messageInputFocusFunction()
			}
		},
	},

	mounted() {
		// Focus input when component is mounted and chat is selected
		if (this.selectedChatId && this.selectedChat) {
			this.$nextTick(() => {
				this.focusMessageInput()
			})
		}
	},

	watch: {
		selectedChatId(newChatId) {
			// Focus input when chat changes (only if we have a valid selectedChat)
			if (newChatId && this.selectedChat) {
				this.$nextTick(() => {
					this.focusMessageInput()
				})
			}
		}
	}
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";

.chat-container {
	height: 100%;
	display: flex;
	flex-direction: column;
	min-height: 0;
	overflow: hidden;
}

.messages-section {
	flex: 1;
	overflow: hidden;
	min-height: 0;
	display: flex;
	flex-direction: column;
}
</style>
