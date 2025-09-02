<template>
	<div class="message-bubble" :class="{ 'has-likes': message.likes && message.likes.length > 0 }">
		<!-- Forwarded Message Indicator -->
		<MessageForwardIndicator v-if="message.isForward" />

		<!-- Reply Section -->
		<MessageReplySection
			v-if="message.replyTo"
			:reply-message="replyMessage"
			:reply-image-url="replyImageUrl"
			:reply-preview-text="replyPreviewText"
			@jump-to-message="$emit('jump-to-message', message.replyTo)"
			@image-error="$emit('reply-image-error', message.replyTo)"
		/>

		<!-- Message Content -->
		<MessageContent
			:message="message"
			:message-image-url="messageImageUrl"
			@open-image-viewer="$emit('open-image-viewer', $event)"
			@image-error="$emit('image-error', message)"
		/>

		<!-- Message Meta Info -->
		<MessageMeta
			:message="message"
			:is-current-user="isCurrentUser"
			:is-read="isRead"
		/>

		<!-- Like Count Display -->
		<MessageLikes
			v-if="message.likes && message.likes.length > 0"
			:likes="message.likes"
			:show-dropdown="showLikesDropdown"
			@toggle-dropdown="$emit('toggle-likes-dropdown')"
			@mouse-leave="$emit('likes-mouse-leave')"
			@mouse-enter="$emit('likes-mouse-enter')"
		/>
	</div>
</template>

<script>
import MessageForwardIndicator from './MessageForwardIndicator.vue'
import MessageReplySection from './MessageReplySection.vue'
import MessageContent from './MessageContent.vue'
import MessageMeta from './MessageMeta.vue'
import MessageLikes from './MessageLikes.vue'

export default {
	name: 'MessageBubble',
	components: {
		MessageForwardIndicator,
		MessageReplySection,
		MessageContent,
		MessageMeta,
		MessageLikes
	},
	props: {
		message: {
			type: Object,
			required: true
		},
		isCurrentUser: {
			type: Boolean,
			required: true
		},
		isRead: {
			type: Boolean,
			default: false
		},
		messageImageUrl: {
			type: String,
			default: null
		},
		// Props passed from parent for reply handling
		replyMessage: {
			type: Object,
			default: null
		},
		replyImageUrl: {
			type: String,
			default: null
		},
		replyPreviewText: {
			type: String,
			default: ''
		},
		showLikesDropdown: {
			type: Boolean,
			default: false
		}
	},
	emits: [
		'jump-to-message',
		'open-image-viewer',
		'image-error',
		'reply-image-error',
		'toggle-likes-dropdown',
		'likes-mouse-leave',
		'likes-mouse-enter'
	]
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>
