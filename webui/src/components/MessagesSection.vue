<template>
	<div class="messages-list" ref="messagesContainer" @scroll="handleScroll">
		<!-- Enhanced message display with components -->
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
			@jump-to-message="$emit('jump-to-message', $event)"
			@open-image-viewer="$emit('open-image-viewer', $event)"
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
              <div class="spinner-border spinner-border-sm message-status-pending" role="status">
                <span class="visually-hidden">Sending...</span>
              </div>
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

		<!-- New messages indicator -->
		<div
			v-if="hasNewMessages"
			class="new-messages-indicator"
			@click="scrollToNewMessages"
		>
			<div class="indicator-content">
				<span>{{ newMessageCount }} new message{{ newMessageCount !== 1 ? 's' : '' }}</span>
				<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
					<path fill-rule="evenodd" d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"/>
				</svg>
			</div>
		</div>
	</div>
</template>

<script>
import Message from './Message.vue'
import LoadingSpinner from './LoadingSpinner.vue'

export default {
	name: 'MessagesSection',
	components: {
		Message,
		LoadingSpinner
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
		},
		hasNewMessages: {
			type: Boolean,
			default: false
		},
		newMessageCount: {
			type: Number,
			default: 0
		},
		scrollThreshold: {
			type: Number,
			default: 100
		}
	},

	emits: [
		'jump-to-message',
		'open-image-viewer',
		'message-image-error',
		'reply-image-error',
		'toggle-message-like',
		'start-reply',
		'open-forward-modal',
		'confirm-delete',
		'scroll-position-changed',
		'scroll-to-new-messages'
	],

	data() {
		return {
			showLikesDropdown: null,
			showReplyDropdown: null,
			showForwardDropdown: null,
			showDeleteDropdown: null,
			isUserAtBottom: true,
			lastScrollTop: 0
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

		handleScroll() {
			const messagesContainer = this.$refs.messagesContainer;
			if (!messagesContainer) return;

			const { scrollTop, scrollHeight, clientHeight } = messagesContainer;
			const distanceFromBottom = scrollHeight - (scrollTop + clientHeight);

			this.isUserAtBottom = distanceFromBottom <= this.scrollThreshold;
			this.lastScrollTop = scrollTop;

			this.$emit('scroll-position-changed', {
				scrollTop,
				scrollHeight,
				clientHeight,
				distanceFromBottom,
				isUserAtBottom: this.isUserAtBottom
			});
		},

		scrollToBottom() {
			this.$nextTick(() => {
				const messagesContainer = this.$refs.messagesContainer;
				if (messagesContainer) {
					messagesContainer.scrollTop = messagesContainer.scrollHeight;
					this.isUserAtBottom = true;
				}
			});
		},

		scrollToNewMessages() {
			this.scrollToBottom();
			this.$emit('scroll-to-new-messages');
		},

		restoreScrollPosition(savedScrollTop, heightDifference = 0) {
			this.$nextTick(() => {
				const container = this.$refs.messagesContainer;
				if (container) {
					container.scrollTop = savedScrollTop + heightDifference;
				}
			});
		},

		formatMessageTime(timestamp) {
			if (!timestamp) return '';

			const date = new Date(timestamp);
			const now = new Date();

			if (date.toDateString() === now.toDateString()) {
				return date.toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit',
					hour12: false
				});
			}

			const yesterday = new Date(now);
			yesterday.setDate(now.getDate() - 1);
			if (date.toDateString() === yesterday.toDateString()) {
				return 'Yesterday ' + date.toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit',
					hour12: false
				});
			}

			const weekAgo = new Date(now);
			weekAgo.setDate(now.getDate() - 7);
			if (date > weekAgo) {
				return date.toLocaleDateString([], {weekday: 'short'}) + ' ' +
					date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
			}

			if (date.getFullYear() === now.getFullYear()) {
				return date.toLocaleDateString([], {month: 'short', day: 'numeric'}) + ' ' +
					date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
			}

			return date.toLocaleDateString() + ' ' +
				date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
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
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";

.messages-list {
	max-width: 100%;
	display: flex;
	flex-direction: column;
	min-height: min-content;
}

</style>
