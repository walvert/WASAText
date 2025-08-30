<template>
	<div v-if="show" class="vue-modal" @click.self="handleClose">
		<div class="vue-modal-dialog">
			<div class="vue-modal-content">
				<div class="vue-modal-header">
					<h5 class="vue-modal-title">Edit Profile</h5>
					<button type="button" class="vue-modal-close" @click="handleClose">×</button>
				</div>
				<div class="vue-modal-body">
					<!-- Current Profile Image -->
					<div v-if="currentUserImageUrl" class="mb-3 text-center">
						<label class="form-label">Current Profile Picture</label>
						<div class="current-photo-preview">
							<img
								:src="currentUserImageUrl"
								:alt="currentUsername || 'User'"
								class="current-group-photo"
								style="max-width: 150px; max-height: 150px;"
							>
						</div>
					</div>

					<!-- Username Field -->
					<div class="mb-3">
						<label class="form-label">Username</label>
						<input
							type="text"
							class="form-control"
							v-model="username"
							placeholder="Enter new username"
							@keyup.enter="validateAndSubmit"
						>
					</div>

					<!-- Profile Image Upload -->
					<div class="mb-3">
						<label class="form-label">Profile Picture</label>
						<input
							type="file"
							class="form-control"
							ref="profileImageInput"
							@change="handleProfileImageSelect"
							accept="image/*"
						>
						<div class="form-text">Supported formats: JPG, PNG, GIF. Max size: 5MB</div>
					</div>

					<!-- Profile Image Preview -->
					<div v-if="selectedProfileImage" class="mb-3 text-center">
						<label class="form-label">Preview</label>
						<div class="photo-preview">
							<img
								:src="profileImagePreviewUrl"
								alt="Preview"
								class="preview-group-photo"
								style="max-width: 150px; max-height: 150px;"
							>
						</div>
						<button
							type="button"
							class="btn btn-sm btn-outline-secondary mt-2"
							@click="clearProfileImageSelection"
						>
							Remove Photo
						</button>
					</div>

					<ErrorMsg v-if="displayError" :msg="displayError"/>
				</div>
				<div class="vue-modal-footer">
					<button type="button" class="btn btn-secondary" @click="handleClose">Cancel</button>
					<button
						type="button"
						class="btn btn-primary"
						@click="validateAndSubmit"
						:disabled="loading"
					>
						<span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status">
							<span class="visually-hidden">Updating...</span>
						</span>
						Update Profile
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import ErrorMsg from './ErrorMsg.vue' // Assuming you have this component

export default {
	name: 'ProfileModal',
	components: {
		ErrorMsg
	},
	props: {
		show: {
			type: Boolean,
			required: true
		},
		currentUserImageUrl: {
			type: String,
			default: ''
		},
		currentUsername: {
			type: String,
			default: ''
		},
		loading: {
			type: Boolean,
			default: false
		},
		error: {
			type: String,
			default: ''
		}
	},
	emits: ['close', 'update-profile'],
	data() {
		return {
			username: '',
			selectedProfileImage: null,
			profileImagePreviewUrl: '',
			validationError: ''
		}
	},
	computed: {
		displayError() {
			return this.validationError || this.error
		}
	},
	watch: {
		currentUsername: {
			handler(newUsername) {
				this.username = newUsername || ''
			},
			immediate: true
		},
		show(newShow) {
			if (newShow) {
				// Reset form when opening
				this.username = this.currentUsername || ''
				this.clearProfileImageSelection()
			}
		}
	},
	methods: {
		handleProfileImageSelect(event) {
			const file = event.target.files[0]
			this.validationError = ''

			if (file) {
				if (file.size > 5 * 1024 * 1024) {
					this.validationError = 'File size must be less than 5MB'
					return
				}

				if (!file.type.startsWith('image/')) {
					this.validationError = 'Please select a valid image file'
					return
				}

				this.selectedProfileImage = file

				// Create preview URL
				const reader = new FileReader()
				reader.onload = (e) => {
					this.profileImagePreviewUrl = e.target.result
				}
				reader.readAsDataURL(file)
			}
		},

		clearProfileImageSelection() {
			this.selectedProfileImage = null
			this.profileImagePreviewUrl = ''
			this.validationError = ''
			if (this.$refs.profileImageInput) {
				this.$refs.profileImageInput.value = ''
			}
		},

		validateAndSubmit() {
			this.validationError = ''

			// Check if at least one field has changes
			const hasUsernameChange = this.username.trim() && this.username.trim() !== this.currentUsername
			const hasImageChange = this.selectedProfileImage !== null

			if (!hasUsernameChange && !hasImageChange) {
				this.validationError = 'Please enter a different username or select a profile image to update'
				return
			}

			if (this.username.trim().length > 0 && this.username.trim().length < 2) {
				this.validationError = 'Username must be at least 2 characters'
				return
			}

			// Emit the update event with the form data
			this.$emit('update-profile', {
				username: hasUsernameChange ? this.username.trim() : null,
				profileImage: this.selectedProfileImage
			})
		},

		handleClose() {
			this.$emit('close')
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
@import url('../assets/modals.css');
</style>
