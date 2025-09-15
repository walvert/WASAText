<template>
	<div class="messages-list" ref="messagesContainer">
		<Message
			v-for="message in messages"
			:key="message.id"
			:message="message"
			:current-username="currentUsername"
			:selected-chat="selectedChat"
			:last-read-message-id="lastReadMessageId"
			:message-image-url="messageImageUrls[message.id]"
			:message-image-urls="messageImageUrls"
			:messages="messages"
			:show-likes-dropdown="showLikesDropdown"
			:show-reply-dropdown="showReplyDropdown"
			:show-forward-dropdown="showForwardDropdown"
			:show-delete-dropdown="showDeleteDropdown"
			:deleting-message="deletingMessage"
			@image-error="$emit('message-image-error', $event)"
			@reply-image-error="$emit('reply-image-error', $event)"
			@toggle-like="$emit('toggle-message-like', $event)"
			@toggle-likes-dropdown="toggleLikesDropdown"
			@likes-mouse-leave="handleDropdownMouseLeave('likes')"
			@likes-mouse-enter="handleDropdownMouseEnter('likes')"
			@toggle-reply-dropdown="toggleReplyDropdown"
			@toggle-forward-dropdown="toggleForwardDropdown"
			@toggle-delete-dropdown="toggleDeleteDropdown"
			@start-reply="$emit('start-reply', $event)"
			@open-forward-modal="$emit('open-forward-modal', $event)"
			@confirm-delete="$emit('confirm-delete', $event)"
			@reply-mouse-leave="handleDropdownMouseLeave('reply')"
			@reply-mouse-enter="handleDropdownMouseEnter('reply')"
			@forward-mouse-leave="handleDropdownMouseLeave('forward')"
			@forward-mouse-enter="handleDropdownMouseEnter('forward')"
			@delete-mouse-leave="handleDropdownMouseLeave('delete')"
			@delete-mouse-enter="handleDropdownMouseEnter('delete')"
		/>

		<!-- Show loading indicator for pending message -->
		<div v-if="sendingMessage" class="message-wrapper mb-3 message-sent">
			<div class="message-content">
				<div class="message-bubble message-pending">
					<div class="message-text">{{ pendingMessage }}</div>
					<div class="message-meta">
						<span class="message-time">{{ formatMessageTime(new Date()) }}</span>
						<span class="message-status">
              <span class="spinner-border spinner-border-sm message-status-pending" role="status">
                <span class="visually-hidden">Sending...</span>
              </span>
            </span>
					</div>
				</div>
			</div>
		</div>

		<!-- Empty state -->
		<div v-if="messages.length === 0 && !sendingMessage" class="text-center p-4 text-muted">
			<div class="mb-3">
				<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor"
					 stroke-width="1.5">
					<path
						d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
				</svg>
			</div>
			<h6>No messages yet</h6>
			<p class="text-muted small">Start the conversation by sending a message</p>
		</div>
	</div>
</template>

<script>
import Message from './Message.vue'
import { formatMessageTime} from "../utils/helpers";

export default {
	name: 'MessagesSection',
	components: {
		Message,
	},

	props: {
		messages: {
			type: Array,
			default: () => []
		},
		selectedChat: {
			type: Object,
			default: null
		},
		currentUsername: {
			type: String,
			required: true
		},
		lastReadMessageId: {
			type: [Number, String],
			default: null
		},
		messageImageUrls: {
			type: Object,
			default: () => ({})
		},
		sendingMessage: {
			type: Boolean,
			default: false
		},
		pendingMessage: {
			type: String,
			default: ''
		},
		deletingMessage: {
			type: [Number, String],
			default: null
		}
	},

	emits: [
		'message-image-error',
		'reply-image-error',
		'toggle-message-like',
		'start-reply',
		'open-forward-modal',
		'confirm-delete'
	],

	data() {
		return {
			showLikesDropdown: null,
			showReplyDropdown: null,
			showForwardDropdown: null,
			showDeleteDropdown: null
		}
	},

	methods: {
		toggleLikesDropdown(messageId) {
			this.showLikesDropdown = this.showLikesDropdown === messageId ? null : messageId;
			if (this.showLikesDropdown === messageId) {
				this.showReplyDropdown = null;
				this.showDeleteDropdown = null;
				this.showForwardDropdown = null;
			}
		},

		toggleReplyDropdown(messageId) {
			this.showReplyDropdown = this.showReplyDropdown === messageId ? null : messageId;
			if (this.showReplyDropdown === messageId) {
				this.showLikesDropdown = null;
				this.showDeleteDropdown = null;
				this.showForwardDropdown = null;
			}
		},

		toggleForwardDropdown(messageId) {
			this.showForwardDropdown = this.showForwardDropdown === messageId ? null : messageId;
			if (this.showForwardDropdown === messageId) {
				this.showLikesDropdown = null;
				this.showReplyDropdown = null;
				this.showDeleteDropdown = null;
			}
		},

		toggleDeleteDropdown(messageId) {
			this.showDeleteDropdown = this.showDeleteDropdown === messageId ? null : messageId;
			if (this.showDeleteDropdown === messageId) {
				this.showLikesDropdown = null;
				this.showReplyDropdown = null;
				this.showForwardDropdown = null;
			}
		},

		handleDropdownMouseLeave(dropdownType) {
			if (dropdownType === 'likes') {
				this.showLikesDropdown = null;
			} else if (dropdownType === 'reply') {
				this.showReplyDropdown = null;
			} else if (dropdownType === 'delete') {
				this.showDeleteDropdown = null;
			} else if (dropdownType === 'forward') {
				this.showForwardDropdown = null;
			}
		},

		handleDropdownMouseEnter(dropdownType) {
			// Prevents dropdown from closing when hovering
		},

		scrollToBottom() {
			this.$nextTick(() => {
				const messagesContainer = this.$refs.messagesContainer;
				if (messagesContainer) {
					messagesContainer.scrollTop = messagesContainer.scrollHeight;
				}
			});
		},
		formatMessageTime
	},

	watch: {
		messages: {
			handler(newMessages, oldMessages) {
				// Scroll to bottom when new messages are added
				if (newMessages && newMessages.length > 0) {
					// Check if this is a new message (length increased)
					if (!oldMessages || newMessages.length > oldMessages.length) {
						this.scrollToBottom();
					}
				}
			},
			deep: true
		},

		sendingMessage(newVal) {
			if (newVal) {
				// Scroll when a message starts sending
				this.$nextTick(() => {
					this.scrollToBottom();
				});
			}
		}
	},

	mounted() {
		// Close dropdowns when clicking outside
		document.addEventListener('click', () => {
			this.showLikesDropdown = null;
			this.showDeleteDropdown = null;
			this.showReplyDropdown = null;
			this.showForwardDropdown = null;
		});

		// Scroll to bottom on initial load
		this.scrollToBottom();
	}
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";

.messages-list {
	max-width: 100%;
	display: flex;
	flex-direction: column;
	min-height: min-content;
	overflow-y: hidden;
	height: 100%;
	flex: 1;
	padding: 1rem;
}

.messages-list::-webkit-scrollbar {
	display: none;
}
</style>
