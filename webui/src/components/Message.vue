<template>
	<div
		class="message-wrapper mb-3"
		:class="{ 'message-sent': isCurrentUser }"
		:data-message-id="message.id"
	>
		<div class="message-content" style="position: relative;">
			<!-- Show sender name for group chats and incoming messages -->
			<div
				v-if="selectedChat?.isGroup && !isCurrentUser"
				class="message-sender"
			>
				{{ message.username }}
			</div>

			<!-- Message Bubble -->
			<MessageBubble
				:message="message"
				:is-current-user="isCurrentUser"
				:is-read="isRead"
				:message-image-url="messageImageUrl"
				:reply-message="replyMessage"
				:reply-image-url="replyImageUrl"
				:reply-preview-text="replyPreviewText"
				:show-likes-dropdown="showLikesDropdown === message.id"
				@jump-to-message="$emit('jump-to-message', $event)"
				@open-image-viewer="$emit('open-image-viewer', $event)"
				@image-error="$emit('image-error', $event)"
				@reply-image-error="$emit('reply-image-error', $event)"
				@toggle-likes-dropdown="$emit('toggle-likes-dropdown', message.id)"
				@likes-mouse-leave="$emit('likes-mouse-leave')"
				@likes-mouse-enter="$emit('likes-mouse-enter')"
			/>

			<!-- Hover Actions -->
			<MessageHoverActions
				:message="message"
				:is-current-user="isCurrentUser"
				:is-liked="isLiked"
				:is-deleting="deletingMessage === message.id"
				:show-reply-dropdown="showReplyDropdown === message.id"
				:show-forward-dropdown="showForwardDropdown === message.id"
				:show-delete-dropdown="showDeleteDropdown === message.id"
				@toggle-like="$emit('toggle-like', message)"
				@toggle-reply-dropdown="$emit('toggle-reply-dropdown', message.id)"
				@toggle-forward-dropdown="$emit('toggle-forward-dropdown', message.id)"
				@toggle-delete-dropdown="$emit('toggle-delete-dropdown', message.id)"
				@start-reply="$emit('start-reply', message)"
				@open-forward-modal="$emit('open-forward-modal', message)"
				@confirm-delete="$emit('confirm-delete', message)"
				@reply-mouse-leave="$emit('reply-mouse-leave')"
				@reply-mouse-enter="$emit('reply-mouse-enter')"
				@forward-mouse-leave="$emit('forward-mouse-leave')"
				@forward-mouse-enter="$emit('forward-mouse-enter')"
				@delete-mouse-leave="$emit('delete-mouse-leave')"
				@delete-mouse-enter="$emit('delete-mouse-enter')"
			/>
		</div>
	</div>
</template>

<script>
import MessageBubble from './MessageBubble.vue'
import MessageHoverActions from './MessageHoverActions.vue'

export default {
	name: 'SingleMessage',
	components: {
		MessageBubble,
		MessageHoverActions
	},
	props: {
		message: {
			type: Object,
			required: true
		},
		currentUsername: {
			type: String,
			required: true
		},
		selectedChat: {
			type: Object,
			default: null
		},
		lastReadMessageId: {
			type: [Number, String],
			default: null
		},
		messageImageUrl: {
			type: String,
			default: null
		},
		messageImageUrls: {
			type: Object,
			default: () => ({})
		},
		messages: {
			type: Array,
			default: () => []
		},
		// Dropdown states
		showLikesDropdown: {
			type: [Number, String, null],
			default: null
		},
		showReplyDropdown: {
			type: [Number, String, null],
			default: null
		},
		showForwardDropdown: {
			type: [Number, String, null],
			default: null
		},
		showDeleteDropdown: {
			type: [Number, String, null],
			default: null
		},
		deletingMessage: {
			type: [Number, String, null],
			default: null
		}
	},
	emits: [
		'jump-to-message',
		'open-image-viewer',
		'image-error',
		'reply-image-error',
		'toggle-like',
		'toggle-likes-dropdown',
		'likes-mouse-leave',
		'likes-mouse-enter',
		'toggle-reply-dropdown',
		'toggle-forward-dropdown',
		'toggle-delete-dropdown',
		'start-reply',
		'open-forward-modal',
		'confirm-delete',
		'reply-mouse-leave',
		'reply-mouse-enter',
		'forward-mouse-leave',
		'forward-mouse-enter',
		'delete-mouse-leave',
		'delete-mouse-enter'
	],
	computed: {
		isCurrentUser() {
			return this.message.username === this.currentUsername
		},
		isRead() {
			if (!this.isCurrentUser || this.lastReadMessageId === null) {
				return false
			}
			const messageId = parseInt(this.message.id)
			const lastReadId = parseInt(this.lastReadMessageId)
			return messageId <= lastReadId
		},
		isLiked() {
			return this.message.likes && this.message.likes.includes(this.currentUsername)
		},
		replyMessage() {
			if (!this.message.replyTo) return null
			return this.messages.find(msg => msg.id === this.message.replyTo)
		},
		replyImageUrl() {
			if (!this.replyMessage || (this.replyMessage.type !== 'image' && this.replyMessage.type !== 'gif')) {
				return null
			}
			return this.messageImageUrls[this.message.replyTo] || null
		},
		replyPreviewText() {
			if (!this.replyMessage) return ''

			if (this.replyMessage.type === 'image') {
				return this.replyMessage.text ? this.replyMessage.text : '📷 Photo'
			} else if (this.replyMessage.type === 'gif') {
				return this.replyMessage.text ? this.replyMessage.text : '🎞️ GIF'
			}

			return this.replyMessage.text || ''
		}
	}
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";

.message-wrapper {
	display: flex;
	width: 100%;
	margin-bottom: 0.5rem;
	flex-shrink: 0;
}

.message-wrapper.message-sent {
	justify-content: flex-end;
}

.message-wrapper:not(.message-sent) {
	justify-content: flex-start;
}

.message-content {
	max-width: 70%;
	display: flex;
	flex-direction: column;
}

.message-sent .message-content {
	align-items: flex-end;
}

.message-wrapper:not(.message-sent) .message-content {
	align-items: flex-start;
}

.message-sender {
	font-size: 0.75rem;
	color: #6c757d;
	margin-bottom: 0.25rem;
	padding-left: 0.75rem;
	font-weight: 500;
}

</style>
