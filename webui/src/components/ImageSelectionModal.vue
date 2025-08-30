<template>
	<!-- Image Selection Modal -->
	<div v-if="show" class="vue-modal" @click.self="$emit('close')">
		<div class="vue-modal-dialog">
			<div class="vue-modal-content">
				<div class="vue-modal-header">
					<h5 class="vue-modal-title">Add Image or GIF</h5>
					<button type="button" class="vue-modal-close" @click="$emit('close')">×</button>
				</div>
				<div class="vue-modal-body">
					<!-- File Upload Section -->
					<div class="mb-3">
						<label class="form-label">Choose Image or GIF</label>
						<input
							type="file"
							class="form-control"
							ref="messageImageInput"
							@change="handleMessageImageSelect"
							accept="image/*,.gif"
						>
						<div class="form-text">Supported formats: JPG, PNG, GIF, WebP. Max size: 10MB</div>
					</div>

					<!-- Image Preview -->
					<div v-if="tempSelectedImage" class="mb-3 text-center">
						<label class="form-label">Preview</label>
						<div class="photo-preview">
							<img
								:src="tempImagePreviewUrl"
								alt="Preview"
								class="preview-message-image"
								style="max-width: 300px; max-height: 300px;"
							>
						</div>
						<div class="mt-2">
							<small class="text-muted">{{ tempSelectedImage.name }} ({{ formatFileSize(tempSelectedImage.size) }})</small>
						</div>
					</div>

					<!-- Error Message -->
					<div v-if="error" class="error-msg">{{ error }}</div>
				</div>
				<div class="vue-modal-footer">
					<button type="button" class="btn btn-secondary" @click="$emit('close')">Cancel</button>
					<button
						type="button"
						class="btn btn-primary"
						@click="selectImage"
						:disabled="!tempSelectedImage"
					>
						Add to Message
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'ImageSelectionModal',

	props: {
		show: {
			type: Boolean,
			required: true
		}
	},

	emits: ['close', 'select'],

	data() {
		return {
			tempSelectedImage: null,
			tempImagePreviewUrl: null,
			error: null
		}
	},

	beforeUnmount() {
		this.clearTempImageSelection()
	},

	methods: {
		// Handle image file selection in modal
		handleMessageImageSelect(event) {
			const file = event.target.files[0]

			if (!file) {
				this.clearTempImageSelection()
				return
			}

			// Validate file type
			const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
			if (!validTypes.includes(file.type)) {
				this.error = 'Please select a valid image file (JPG, PNG, GIF, WebP)'
				this.clearTempImageSelection()
				return
			}

			// Validate file size (10MB limit)
			const maxSize = 10 * 1024 * 1024 // 10MB in bytes
			if (file.size > maxSize) {
				this.error = 'File size must be less than 10MB'
				this.clearTempImageSelection()
				return
			}

			this.tempSelectedImage = file
			this.error = null

			// Create preview URL
			if (this.tempImagePreviewUrl) {
				URL.revokeObjectURL(this.tempImagePreviewUrl)
			}
			this.tempImagePreviewUrl = URL.createObjectURL(file)
		},

		// Clear temporary image selection in modal
		clearTempImageSelection() {
			this.tempSelectedImage = null

			if (this.tempImagePreviewUrl) {
				URL.revokeObjectURL(this.tempImagePreviewUrl)
				this.tempImagePreviewUrl = null
			}

			// Clear the file input
			if (this.$refs.messageImageInput) {
				this.$refs.messageImageInput.value = ''
			}
		},

		// Select image for message
		selectImage() {
			if (!this.tempSelectedImage) return

			// Emit the selected image data to parent
			this.$emit('select', {
				file: this.tempSelectedImage,
				previewUrl: this.tempImagePreviewUrl
			})

			// Reset modal state but don't revoke URL (parent will handle it)
			this.tempSelectedImage = null
			this.tempImagePreviewUrl = null
			this.error = null

			// Clear file input
			if (this.$refs.messageImageInput) {
				this.$refs.messageImageInput.value = ''
			}
		},

		// Format file size for display
		formatFileSize(bytes) {
			if (bytes === 0) return '0 Bytes'
			const k = 1024
			const sizes = ['Bytes', 'KB', 'MB', 'GB']
			const i = Math.floor(Math.log(bytes) / Math.log(k))
			return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
		}
	},

	watch: {
		show(newValue) {
			if (!newValue) {
				// Clear selection when modal is hidden
				this.clearTempImageSelection()
				this.error = null
			}
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
@import url('../assets/modals.css');
</style>
