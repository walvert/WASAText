<template>
	<div
		class="message-reply-section"
		:class="{ 'reply-deleted': !replyMessage }"
		:style="{ cursor: replyMessage ? 'pointer' : 'not-allowed' }"
	>
		<div v-if="replyMessage" class="reply-section-header">
			Reply to {{ replyMessage.username }}
		</div>
		<div v-else class="reply-section-header">
			Reply to message
		</div>

		<div v-if="replyMessage" class="reply-section-content">
			<!-- Image preview for image/gif messages -->
			<div
				v-if="(replyMessage.type === 'image' || replyMessage.type === 'gif') && replyImageUrl"
				class="reply-image-preview"
			>
				<img
					:src="replyImageUrl"
					:alt="replyPreviewText"
					class="reply-preview-image"
					@error="handleImageError"
				>
			</div>

			<!-- Image placeholder if image failed to load -->
			<div
				v-else-if="(replyMessage.type === 'image' || replyMessage.type === 'gif') && !replyImageUrl"
				class="reply-image-placeholder"
			>
				<svg width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
					<path d="M6.002 5.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"/>
					<path d="M2.002 1a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V3a2 2 0 0 0-2-2h-12zm12 1a1 1 0 0 1 1v6.5l-3.777-1.947a.5.5 0 0 0-.577.093l-3.71 3.71-2.66-1.772a.5.5 0 0 0-.63.062L1.002 12V3a1 1 0 0 1 1-1h12z"/>
				</svg>
			</div>

			<div class="reply-section-text">
				{{ replyPreviewText }}
			</div>
		</div>

		<div v-else class="reply-section-content">
			<div class="reply-section-text reply-section-deleted">
				This message has been deleted
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'MessageReplySection',
	props: {
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
		}
	},
	emits: ['image-error'],
	methods: {
		handleImageError() {
			this.$emit('image-error')
		}
	}
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";

.reply-image-preview {
	width: 40px;
	height: 40px;
	border-radius: 4px;
	overflow: hidden;
	flex-shrink: 0;
	margin-right: 8px;
}

.reply-preview-image {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 4px;
}

.reply-image-placeholder {
	width: 40px;
	height: 40px;
	background-color: #f8f9fa;
	border: 1px solid #dee2e6;
	border-radius: 4px;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 8px;
	flex-shrink: 0;
}

.reply-section-content {
	display: flex;
	align-items: center;
	gap: 8px;
}

.reply-section-text {
	flex: 1;
	min-width: 0;
}
</style>
