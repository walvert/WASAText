<template>
	<div class="message-input-container p-3 border-top bg-light">
		<form @submit.prevent="handleSendMessage">
			<!-- Reply Preview -->
			<div v-if="replyingToMessage" class="reply-preview">
				<div class="reply-preview-header">
					<h6 class="reply-preview-title">
						Replying to {{ replyingToMessage.username }}
					</h6>
					<button
						type="button"
						class="reply-preview-close"
						@click="clearReply"
						title="Cancel reply"
					>
						<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
							<path d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z"/>
						</svg>
					</button>
				</div>
				<p class="reply-preview-text">{{ getReplyPreviewText(replyingToMessage) }}</p>
			</div>

			<!-- Show selected image preview if any -->
			<div v-if="selectedMessageImage" class="selected-image-preview mb-3">
				<div class="image-preview-container">
					<img
						:src="messageImagePreviewUrl"
						:alt="selectedMessageImage.name"
						class="preview-message-image"
					>
					<button
						type="button"
						class="btn btn-sm btn-danger remove-image-btn"
						@click="clearMessageImageSelection"
						title="Remove image"
					>
						<svg width="12" height="12" fill="currentColor" viewBox="0 0 16 16">
							<path d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z"/>
						</svg>
					</button>
				</div>
				<div class="image-info">
					<small class="text-muted">{{ selectedMessageImage.name }} ({{ formatFileSize(selectedMessageImage.size) }})</small>
				</div>
			</div>

			<div class="input-group">
				<!-- Add Media Button -->
				<button
					type="button"
					class="btn-add-media"
					@click="openImageModal"
					title="Add image or GIF"
					:disabled="sendingMessage"
				>
					<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
						<path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
					</svg>
				</button>

				<input
					type="text"
					class="form-control message-input"
					:placeholder="getInputPlaceholder()"
					v-model="localNewMessage"
					:disabled="sendingMessage"
					@keypress.enter="handleSendMessage"
					ref="messageInput"
				>
				<button
					class="btn btn-primary send-btn"
					type="submit"
					:disabled="sendingMessage || (!localNewMessage.trim() && !selectedMessageImage)"
				>
					<svg v-if="!sendingMessage" width="16" height="16" fill="currentColor"
						 viewBox="0 0 16 16">
						<path
							d="M15.854.146a.5.5 0 0 1 .11.54l-5.819 14.547a.75.75 0 0 1-1.329.124l-3.178-4.995L.643 7.184a.75.75 0 0 1 .124-1.33L15.314.037a.5.5 0 0 1 .54.11ZM6.636 10.07l2.761 4.338L14.13 2.576 6.636 10.07Zm6.787-8.201L1.591 6.602l4.339 2.76 7.494-7.493Z"/>
					</svg>
					<LoadingSpinner v-else style="width: 16px; height: 16px;"/>
				</button>
			</div>
		</form>
	</div>
</template>

<script>
import LoadingSpinner from './LoadingSpinner.vue'
import {formatFileSize} from "../utils/helpers";

export default {
	name: 'MessageInput',
	components: {
		LoadingSpinner
	},

	props: {
		newMessage: {
			type: String,
			default: ''
		},
		sendingMessage: {
			type: Boolean,
			default: false
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
		}
	},

	emits: [
		'send-message',
		'clear-reply',
		'open-image-modal',
		'clear-message-image-selection',
		'update:new-message',
		'focus-input'
	],

	data() {
		return {
			localNewMessage: this.newMessage
		}
	},

	watch: {
		newMessage(newVal) {
			this.localNewMessage = newVal
		},
		localNewMessage(newVal) {
			this.$emit('update:new-message', newVal)
		}
	},

	methods: {
		handleSendMessage() {
			this.$emit('send-message')
		},

		clearReply() {
			this.$emit('clear-reply')
		},

		openImageModal() {
			this.$emit('open-image-modal')
		},

		clearMessageImageSelection() {
			this.$emit('clear-message-image-selection')
		},

		getInputPlaceholder() {
			if (this.replyingToMessage) {
				return 'Reply to this message...'
			}
			if (this.selectedMessageImage) {
				return 'Add a caption...'
			}
			return 'Type a message...'
		},

		getReplyPreviewText(message) {
			if (!message) return '';

			if (message.type === 'image') {
				return message.text ? '📷 ' + message.text : '📷 Photo';
			} else if (message.type === 'gif') {
				return message.text ? '🎞️ GIF ' + message.text : '🎞️ GIF';
			}

			return message.text || '';
		},

		formatFileSize,
		focusInput() {
			this.$nextTick(() => {
				if (this.$refs.messageInput) {
					this.$refs.messageInput.focus()
				}
			})
		}
	},

	mounted() {
		// Emit focus event when component is ready
		this.$emit('focus-input', this.focusInput)
	}
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";

.message-input-container {
	flex-shrink: 0;
	border-top: 1px solid #dee2e6;
	background-color: #fff;
	padding: 1rem;
}

.message-input-container .input-group {
	display: flex;
	align-items: center;
	width: 100%;
}

.message-input-container .form-control {
	flex: 1;
	border-radius: 1.5rem;
	padding: 0.5rem 1rem;
	border: 1px solid #dee2e6;
	margin-right: 0.5rem;
}

.message-input-container .form-control:focus {
	box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
	border-color: #007bff;
}

.message-input-container .btn {
	flex-shrink: 0;
	border-radius: 50%;
	width: 40px;
	height: 40px;
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 0;
}

</style>
