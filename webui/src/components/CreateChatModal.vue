<template>
	<div class="vue-modal" v-if="show" @click.self="closeModal">
		<div class="vue-modal-dialog">
			<div class="vue-modal-content">
				<div class="vue-modal-header">
					<h5 class="vue-modal-title">New Conversation</h5>
					<button type="button" class="vue-modal-close" @click="closeModal">×</button>
				</div>
				<div class="vue-modal-body">
					<!-- User Selection Section -->
					<div class="mb-4">
						<label class="form-label">Select Recipients</label>
						<div class="user-search-container mb-3">
							<div class="input-group">
								<input
									type="text"
									class="form-control"
									placeholder="Search users..."
									:value="userSearchQuery"
									@input="handleSearchInput"
								>
							</div>
						</div>

						<!-- Loading state for users -->
						<div v-if="loadingUsers" class="text-center py-3">
							<div class="py-3">
								<LoadingSpinner :loading="loadingUsers" />
							</div>
						</div>

						<!-- Users list -->
						<div v-else-if="filteredUsers.length > 0" class="user-list border rounded p-2" style="max-height: 200px; overflow-y: auto;">
							<div
								v-for="user in filteredUsers"
								:key="user.id"
								class="user-item d-flex align-items-center p-2 rounded hover-bg-light"
								style="cursor: pointer;"
								@click="toggleUserSelection(user)"
							>
								<input
									type="checkbox"
									class="form-check-input me-3"
									:checked="selectedUsers.some(u => u.id === user.id)"
									@click.stop="toggleUserSelection(user)"
								>
								<div class="user-avatar me-3">
									<div class="avatar-circle" v-if="!userImageUrls[user.id]" style="width: 32px; height: 32px; font-size: 12px;">
										<span class="avatar-text">{{ getUserInitials(user.username) }}</span>
									</div>
									<img
										v-else
										:src="userImageUrls[user.id]"
										:alt="user.username"
										class="avatar-image"
										style="width: 32px; height: 32px;"
										@error="() => handleUsersImageError(user.id)"
									>
								</div>
								<div class="user-info">
									<div class="fw-medium">{{ user.username }}</div>
								</div>
							</div>
						</div>

						<!-- No users found -->
						<div v-else-if="users.length === 0 && !loadingUsers" class="text-center py-3 text-muted">
							<p>No users found.</p>
						</div>

						<!-- No filtered results -->
						<div v-else-if="filteredUsers.length === 0 && userSearchQuery && !loadingUsers" class="text-center py-3 text-muted">
							<p>No users found matching "{{ userSearchQuery }}"</p>
						</div>
					</div>

					<!-- Selected Users Display -->
					<div v-if="selectedUsers.length > 0" class="mb-3">
						<label class="form-label">Selected Recipients ({{ selectedUsers.length }})</label>
						<div class="selected-users d-flex flex-wrap gap-2">
							<span
								v-for="user in selectedUsers"
								:key="user.id"
								class="badge bg-primary d-flex align-items-center gap-1"
							>
								{{ user.username }}
								<button
									type="button"
									class="btn-close btn-close-white"
									style="font-size: 0.7em;"
									@click="removeUserSelection(user)"
								></button>
							</span>
						</div>
					</div>

					<!-- Group Name Input -->
					<div v-if="selectedUsers.length > 1" class="mb-3">
						<label class="form-label">Group Name <span class="text-danger">*</span></label>
						<input
							type="text"
							class="form-control"
							:value="newChatName"
							@input="$emit('update:newChatName', $event.target.value)"
							placeholder="Enter group name"
							required
						>
						<div class="form-text">Required when creating a group chat with multiple recipients.</div>
					</div>

					<!-- Media Upload Section -->
					<div v-if="selectedNewChatImage" class="mb-3">
						<label class="form-label">Selected Media</label>
						<div class="image-preview-container">
							<img
								:src="newChatImagePreviewUrl"
								:alt="selectedNewChatImage.name"
								class="preview-message-image"
								style="max-width: 200px; max-height: 150px;"
							>
							<button
								type="button"
								class="btn btn-sm btn-danger remove-image-btn"
								@click="clearNewChatImageSelection"
								title="Remove media"
							>
								<svg width="12" height="12" fill="currentColor" viewBox="0 0 16 16">
									<path d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z"/>
								</svg>
							</button>
						</div>
						<div class="image-info mt-2">
							<small class="text-muted">{{ selectedNewChatImage.name }} ({{ formatFileSize(selectedNewChatImage.size) }})</small>
						</div>
					</div>

					<!-- Initial Message Section -->
					<div class="mb-3">
						<label class="form-label">Initial Message</label>
						<div class="input-group">
							<!-- Add Media Button -->
							<button
								type="button"
								class="btn btn-outline-secondary"
								@click="$emit('open-image-modal')"
								title="Add image or GIF"
							>
								<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
									<path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
								</svg>
							</button>
							<input
								type="text"
								class="form-control"
								:placeholder="selectedNewChatImage ? 'Add a caption...' : 'Type your first message...'"
								:value="initialMessage"
								@keyup.enter="createChat"
								@input="$emit('update:initialMessage', $event.target.value)"
							>
						</div>
						<div class="form-text">Required: Either enter a message or select media to start the conversation.</div>
					</div>

					<ErrorMsg v-if="newChatError" :msg="newChatError"/>
				</div>
				<div class="vue-modal-footer">
					<button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
					<button
						type="button"
						class="btn btn-primary"
						@click="createChat"
						:disabled="!canCreateChat"
					>
						<span v-if="newChatLoading" class="spinner-border spinner-border-sm me-2" role="status">
							<span class="visually-hidden">Creating...</span>
						</span>
						{{ selectedUsers.length > 1 ? 'Create Group' : 'Send Message' }}
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import ErrorMsg from './ErrorMsg.vue';
import { getUserInitials } from '../utils/helpers';
import {formatFileSize} from "../utils/helpers";

export default {
	name: 'CreateChatModal',
	components: {
		ErrorMsg
	},
	props: {
		show: {
			type: Boolean,
			default: false
		},
		users: {
			type: Array,
			default: () => []
		},
		filteredUsers: {
			type: Array,
			default: () => []
		},
		selectedUsers: {
			type: Array,
			default: () => []
		},
		loadingUsers: {
			type: Boolean,
			default: false
		},
		newChatName: {
			type: String,
			default: ''
		},
		selectedNewChatImage: {
			type: File,
			default: null
		},
		newChatImagePreviewUrl: {
			type: String,
			default: ''
		},
		initialMessage: {
			type: String,
			default: ''
		},
		newChatError: {
			type: String,
			default: null
		},
		newChatLoading: {
			type: Boolean,
			default: false
		},
		userImageUrls: {
			type: Object,
			default: () => ({})
		},
		canCreateChat: {
			type: Boolean,
			default: false
		},
		userSearchQuery: {
			type: String,
			default: ''
		},
	},

	methods: {
		getUserInitials,

		closeModal() {
			this.$emit('close');
		},

		handleSearchInput(event) {
			const searchQuery = event.target.value;
			this.$emit('update:userSearchQuery', searchQuery);
			this.$emit('filter-users', searchQuery);
		},

		toggleUserSelection(user) {
			this.$emit('toggle-user-selection', user);
		},

		removeUserSelection(user) {
			this.$emit('remove-user-selection', user);
		},

		handleUsersImageError(userId) {
			this.$emit('handle-users-image-error', userId);
		},

		clearNewChatImageSelection() {
			this.$emit('clear-new-chat-image-selection');
		},

		formatFileSize,
		createChat() {
			this.$emit('create-new-chat');
		}
	},

	watch: {
		show(newVal) {
			if (newVal) {
				if (this.users.length === 0) {
					this.$emit('get-users');
				}
			} else {
				this.$emit('update:userSearchQuery', '');
			}
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
@import url('../assets/modals.css');
</style>
