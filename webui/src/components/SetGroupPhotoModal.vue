<template>
	<div v-if="show" class="vue-modal" @click.self="$emit('close')">
		<div class="vue-modal-dialog">
			<div class="vue-modal-content">
				<div class="vue-modal-header">
					<h5 class="vue-modal-title">Set Group Photo</h5>
					<button class="vue-modal-close" @click="$emit('close')">×</button>
				</div>
				<div class="vue-modal-body">
					<!-- Current Photo Preview -->
					<div v-if="currentPhotoUrl" class="mb-3 text-center">
						<label class="form-label">Current Photo</label>
						<div class="current-photo-preview">
							<img
								:src="currentPhotoUrl"
								:alt="groupName"
								class="current-group-photo"
							>
						</div>
					</div>

					<!-- File Upload Section -->
					<div class="mb-3">
						<label class="form-label">Choose New Photo</label>
						<input
							type="file"
							class="form-control"
							ref="photoFileInput"
							@change="handlePhotoFileSelect"
							accept="image/*"
						>
						<div class="form-text">Supported formats: JPG, PNG, GIF. Max size: 10MB</div>
					</div>

					<!-- Photo Preview -->
					<div v-if="selectedPhotoFile" class="mb-3 text-center">
						<label class="form-label">Preview</label>
						<div class="photo-preview">
							<img
								:src="photoPreviewUrl"
								alt="Preview"
								class="preview-group-photo"
							>
						</div>
						<button
							type="button"
							class="btn btn-sm btn-outline-secondary mt-2"
							@click="clearPhotoSelection"
						>
							Remove Photo
						</button>
					</div>

					<!-- Error Message -->
					<div v-if="error" class="error-msg">{{ error }}</div>
				</div>
				<div class="vue-modal-footer">
					<button class="btn btn-secondary" @click="$emit('close')">Cancel</button>
					<button
						class="btn btn-primary"
						@click="handleSetPhoto"
						:disabled="loading || !selectedPhotoFile"
					>
						<span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status">
							<span class="visually-hidden">Loading...</span>
						</span>
						{{ selectedPhotoFile ? 'Update Photo' : 'Select Photo' }}
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'SetGroupPhotoModal',
	props: {
		show: {
			type: Boolean,
			default: false
		},
		currentPhotoUrl: {
			type: String,
			default: null
		},
		groupName: {
			type: String,
			default: ''
		},
		loading: {
			type: Boolean,
			default: false
		},
		error: {
			type: String,
			default: null
		}
	},

	emits: ['close', 'set-photo'],

	data() {
		return {
			selectedPhotoFile: null,
			photoPreviewUrl: null,
			localError: null
		}
	},

	computed: {
		displayError() {
			return this.localError || this.error
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
		this.clearPhotoSelection()
	},

	methods: {
		handlePhotoFileSelect(event) {
			const file = event.target.files[0]

			if (!file) {
				this.clearPhotoSelection()
				return
			}

			// Validate file type
			if (!file.type.startsWith('image/')) {
				this.localError = 'Please select a valid image file'
				this.clearPhotoSelection()
				return
			}

			// Validate file size (10MB limit)
			const maxSize = 10 * 1024 * 1024 // 10MB
			if (file.size > maxSize) {
				this.localError = 'File size must be less than 10MB'
				this.clearPhotoSelection()
				return
			}

			this.selectedPhotoFile = file
			this.localError = null

			if (this.photoPreviewUrl) {
				URL.revokeObjectURL(this.photoPreviewUrl)
			}
			this.photoPreviewUrl = URL.createObjectURL(file)
		},

		clearPhotoSelection() {
			this.selectedPhotoFile = null

			if (this.photoPreviewUrl) {
				URL.revokeObjectURL(this.photoPreviewUrl)
				this.photoPreviewUrl = null
			}

			if (this.$refs.photoFileInput) {
				this.$refs.photoFileInput.value = ''
			}

			this.localError = null
		},

		handleSetPhoto() {
			if (!this.selectedPhotoFile) {
				this.localError = 'Please select a photo'
				return
			}

			this.$emit('set-photo', this.selectedPhotoFile)
		},

		resetModal() {
			this.clearPhotoSelection()
			this.localError = null
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
@import url('../assets/modals.css');
</style>

