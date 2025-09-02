<template>
	<div class="message-content-wrapper">
		<!-- Image Message -->
		<div v-if="message.type === 'image' || message.type === 'gif'" class="image-message">
			<!-- Loading state -->
			<div v-if="message.imageLoading" class="image-loading">
				<div class="spinner-border spinner-border-sm" role="status">
					<span class="visually-hidden">Loading image...</span>
				</div>
				<span class="ms-2">Loading image...</span>
			</div>

			<!-- Error state -->
			<div v-else-if="message.imageError" class="image-error">
				<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16" class="me-2">
					<path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zM5.354 4.646a.5.5 0 1 0-.708.708L7.293 8l-2.647 2.646a.5.5 0 0 0 .708.708L8 8.707l2.646 2.647a.5.5 0 0 0 .708-.708L8.707 8l2.647-2.646a.5.5 0 0 0-.708-.708L8 7.293 5.354 4.646z"/>
				</svg>
				Failed to load image
			</div>

			<!-- Actual image -->
			<img
				v-else-if="messageImageUrl"
				:src="messageImageUrl"
				:alt="message.mediaUrl || 'Image'"
				class="message-image"
				@click="handleImageClick"
				@error="$emit('image-error')"
			>

			<!-- Fallback if no image URL -->
			<div v-else class="image-error">
				<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16" class="me-2">
					<path d="M6.002 5.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"/>
					<path d="M2.002 1a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V3a2 2 0 0 0-2-2h-12zm12 1a1 1 0 0 1 1v6.5l-3.777-1.947a.5.5 0 0 0-.577.093l-3.71 3.71-2.66-1.772a.5.5 0 0 0-.63.062L1.002 12V3a1 1 0 0 1 1-1h12z"/>
				</svg>
				Image unavailable
			</div>

			<!-- Caption/text if present -->
			<div v-if="message.text && message.text.trim()" class="message-text image-caption">
				{{ message.text }}
			</div>
		</div>

		<!-- Text Message -->
		<div v-else class="message-text">
			{{ message.text }}
		</div>
	</div>
</template>

<script>
export default {
	name: 'MessageContent',
	props: {
		message: {
			type: Object,
			required: true
		},
		messageImageUrl: {
			type: String,
			default: null
		}
	},
	emits: ['image-error', 'open-image-viewer'],
	methods: {
		handleImageClick() {
			if (this.messageImageUrl) {
				this.$emit('open-image-viewer', {
					url: this.messageImageUrl,
					title: this.message.mediaUrl
				})
			}
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>

