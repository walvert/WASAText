<template>
	<div v-if="show" class="vue-modal" @click.self="handleClose">
		<div class="vue-modal-dialog">
			<div class="vue-modal-content">
				<div class="vue-modal-header">
					<h5 class="vue-modal-title">Add Image or GIF</h5>
					<button type="button" class="vue-modal-close" @click="handleClose">×</button>
				</div>
				<div class="vue-modal-body">
					<!-- File Upload Section -->
					<div class="mb-3">
						<label class="form-label">Choose Image or GIF</label>
						<input
							type="file"
							class="form-control"
							ref="fileInput"
							@change="handleFileSelect"
							accept="image/*,.gif"
						>
						<div class="form-text">Supported formats: JPG, PNG, GIF, WebP. Max size: 10MB</div>
					</div>

					<!-- Image Preview -->
					<div v-if="selectedFile" class="mb-3 text-center">
						<label class="form-label">Preview</label>
						<div class="photo-preview">
							<img
								:src="previewUrl"
								alt="Preview"
								class="preview-message-image"
								style="max-width: 300px; max-height: 300px;"
							>
						</div>
						<div class="mt-2">
							<small class="text-muted">{{ selectedFile.name }} ({{ formatFileSize(selectedFile.size) }})</small>
						</div>
					</div>

					<!-- Error Message -->
					<div v-if="error" class="error-msg">{{ error }}</div>
				</div>
				<div class="vue-modal-footer">
					<button type="button" class="btn btn-secondary" @click="handleClose">Cancel</button>
					<button
						type="button"
						class="btn btn-primary"
						@click="handleSelect"
						:disabled="!selectedFile"
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
	name: 'NewChatImageModal',

	props: {
		show: {
			type: Boolean,
			default: false
		}
	},

	emits: [
		'close',
		'select'
	],

	data() {
		return {
			selectedFile: null,
			previewUrl: null,
			error: null
		}
	},

	watch: {
		show(newVal) {
			if (!newVal) {
				this.resetModal()
			}
		}
	},

	beforeUnmount() {
		this.cleanupPreviewUrl()
	},

	methods: {
		handleClose() {
			this.$emit('close')
		},

		handleFileSelect(event) {
			const file = event.target.files[0]

			if (!file) {
				this.clearSelection()
				return
			}

			const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
			if (!validTypes.includes(file.type)) {
				this.error = 'Please select a valid image file (JPG, PNG, GIF, WebP)'
				this.clearSelection()
				return
			}

			const maxSize = 10 * 1024 * 1024 // 10MB in bytes
			if (file.size > maxSize) {
				this.error = 'File size must be less than 10MB'
				this.clearSelection()
				return
			}

			this.selectedFile = file
			this.error = null

			// Create preview URL
			this.cleanupPreviewUrl()
			this.previewUrl = URL.createObjectURL(file)
		},

		handleSelect() {
			if (!this.selectedFile) return

			this.$emit('select', {
				file: this.selectedFile,
				previewUrl: this.previewUrl
			})

			this.selectedFile = null
			this.previewUrl = null
		},

		clearSelection() {
			this.selectedFile = null
			this.cleanupPreviewUrl()

			if (this.$refs.fileInput) {
				this.$refs.fileInput.value = ''
			}
		},

		cleanupPreviewUrl() {
			if (this.previewUrl) {
				URL.revokeObjectURL(this.previewUrl)
				this.previewUrl = null
			}
		},

		resetModal() {
			this.clearSelection()
			this.error = null
		},

		formatFileSize(bytes) {
			if (bytes === 0) return '0 Bytes'
			const k = 1024
			const sizes = ['Bytes', 'KB', 'MB', 'GB']
			const i = Math.floor(Math.log(bytes) / Math.log(k))
			return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
@import url('../assets/modals.css');
</style>

