<template>
	<div class="message-bubble" :class="{ 'has-likes': message.likes && message.likes.length > 0 }">
		<MessageForwardIndicator v-if="message.isForward" />

		<MessageReplySection
			v-if="message.replyTo"
			:reply-message="replyMessage"
			:reply-image-url="replyImageUrl"
			:reply-preview-text="replyPreviewText"
			@image-error="$emit('reply-image-error', message.replyTo)"
		/>

		<MessageContent
			:message="message"
			:message-image-url="messageImageUrl"
			@image-error="$emit('image-error', message)"
		/>

		<MessageMeta
			:message="message"
			:is-current-user="isCurrentUser"
			:is-read="isRead"
		/>

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
		'image-error',
		'reply-image-error',
		'toggle-likes-dropdown',
		'likes-mouse-leave',
		'likes-mouse-enter'
	]
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";

.message-bubble {
	background-color: #e9ecef;
	border-radius: 1rem;
	padding: 0.5rem 0.75rem;
	position: relative;
	word-wrap: break-word;
	max-width: 100%;
}

.message-bubble.has-likes {
	padding-bottom: 2rem;
}

.message-sent .message-bubble {
	background-color: #007bff;
	color: white;
}


</style>
