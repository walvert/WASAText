<template>
	<div v-if="show" class="forward-modal" @click.self="$emit('close')">
		<div class="forward-modal-dialog">
			<div class="vue-modal-content">
				<div class="vue-modal-header">
					<h5 class="vue-modal-title">Forward Message</h5>
					<button type="button" class="vue-modal-close" @click="$emit('close')">×</button>
				</div>
				<div class="vue-modal-body">
					<!-- Message Preview -->
					<div class="forward-message-preview">
						<div class="forward-preview-header">Forwarding message:</div>
						<div class="forward-preview-content"
							 :class="{ 'has-image': message && (message.type === 'image' || message.type === 'gif') && messageImageUrls[message.id] }">
							<img
								v-if="message && (message.type === 'image' || message.type === 'gif') && messageImageUrls[message.id]"
								:src="messageImageUrls[message.id]"
								:alt="message.mediaUrl || 'Image'"
								class="forward-preview-image"
							>
							<div class="forward-preview-text">
								{{ getForwardPreviewText(message) }}
							</div>
						</div>
					</div>

					<!-- Recipient Selection Tabs -->
					<div class="recipient-tabs">
						<button
							class="recipient-tab"
							:class="{ active: activeTab === 'users' }"
							@click="activeTab = 'users'"
						>
							Users
						</button>
						<button
							class="recipient-tab"
							:class="{ active: activeTab === 'chats' }"
							@click="activeTab = 'chats'"
						>
							Group Chats
						</button>
					</div>

					<!-- Selected Recipients Display -->
					<div v-if="selectedRecipients.length > 0" class="mb-3">
						<label class="form-label">Selected Recipients ({{ selectedRecipients.length }})</label>
						<div class="selected-recipients d-flex flex-wrap">
							<span
								v-for="recipient in selectedRecipients"
								:key="`${recipient.type}-${recipient.id}`"
								class="recipient-chip"
								:class="recipient.type === 'user' ? 'user-chip' : 'chat-chip'"
							>
								{{ recipient.name }}
								<button
									type="button"
									class="recipient-chip-close"
									@click="removeRecipient(recipient)"
								>
									×
								</button>
							</span>
						</div>
					</div>

					<!-- Users Tab Content -->
					<div v-if="activeTab === 'users'">
						<!-- User Search -->
						<div class="mb-3">
							<label class="form-label">Select Users</label>
							<div class="user-search-container mb-3">
								<div class="input-group">
									<input
										type="text"
										class="form-control"
										placeholder="Search users..."
										v-model="userSearchQuery"
										@input="filterUsers"
									>
								</div>
							</div>

							<!-- Loading state for users -->
							<div v-if="loadingUsers" class="forward-loading">
								<div class="forward-loading">
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
										:checked="isRecipientSelected(user.id, 'user')"
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
							<div v-else-if="users.length === 0 && !loadingUsers" class="forward-empty-state">
								<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
									<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
									<circle cx="9" cy="7" r="4"/>
									<path d="M22 21v-2a4 4 0 0 0-3-3.87"/>
									<path d="M16 3.13a4 4 0 0 1 0 7.75"/>
								</svg>
								<p>No users found.</p>
							</div>

							<!-- No filtered results -->
							<div v-else-if="filteredUsers.length === 0 && userSearchQuery && !loadingUsers" class="forward-empty-state">
								<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
									<circle cx="11" cy="11" r="8"/>
									<path d="M21 21l-4.35-4.35"/>
								</svg>
								<p>No users found matching "{{ userSearchQuery }}"</p>
							</div>
						</div>
					</div>

					<!-- Chats Tab Content -->
					<div v-if="activeTab === 'chats'">
						<div class="mb-3">
							<label class="form-label">Select Group Chats</label>

							<!-- Group chats list -->
							<div v-if="groupChats.length > 0" class="forward-chat-list">
								<div
									v-for="chat in groupChats"
									:key="chat.id"
									class="forward-chat-item"
									:class="{ selected: isRecipientSelected(chat.id, 'chat') }"
									@click="toggleChatSelection(chat)"
								>
									<input
										type="checkbox"
										class="form-check-input me-3"
										:checked="isRecipientSelected(chat.id, 'chat')"
										@click.stop="toggleChatSelection(chat)"
									>
									<div class="forward-chat-avatar">
										<img
											v-if="chatImageUrls[chat.id]"
											:src="chatImageUrls[chat.id]"
											:alt="getChatName(chat)"
											style="width: 100%; height: 100%; border-radius: 50%; object-fit: cover;"
										>
										<span v-else>{{ getChatInitials(chat) }}</span>
									</div>
									<div class="forward-chat-info">
										<div class="forward-chat-name">{{ getChatName(chat) }}</div>
										<div class="forward-chat-type">Group Chat</div>
									</div>
								</div>
							</div>

							<!-- No group chats -->
							<div v-else class="forward-empty-state">
								<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
									<path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
								</svg>
								<p>No group chats available for forwarding.</p>
							</div>
						</div>
					</div>

					<!-- Error Message -->
					<div v-if="error" class="error-msg">{{ error }}</div>
				</div>
				<div class="vue-modal-footer">
					<button type="button" class="btn btn-secondary" @click="$emit('close')">Cancel</button>
					<button
						type="button"
						class="btn btn-success"
						@click="submitForward"
						:disabled="loading || selectedRecipients.length === 0"
					>
						<span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status">
							<span class="visually-hidden">Forwarding...</span>
						</span>
						Forward to {{ selectedRecipients.length }} recipient{{ selectedRecipients.length !== 1 ? 's' : '' }}
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { getUserInitials } from '../utils/helpers';

export default {
	name: 'ForwardMessageModal',
	props: {
		show: {
			type: Boolean,
			default: false
		},
		message: {
			type: Object,
			default: null
		},
		messageImageUrls: {
			type: Object,
			default: () => ({})
		},
		chatImageUrls: {
			type: Object,
			default: () => ({})
		},
		groupChats: {
			type: Array,
			default: () => []
		},
		loadUsers: {
			type: Function,
			required: true
		},
		forwardMessage: {
			type: Function,
			required: true
		},
		userImageUrls: {
			type: Object,
			default: () => ({})
		},
	},
	emits: ['close', "load-user-images"],
	data() {
		return {
			activeTab: 'users',
			loading: false,
			error: null,
			users: [],
			filteredUsers: [],
			loadingUsers: false,
			userSearchQuery: '',
			selectedRecipients: []
		}
	},
	watch: {
		show(newValue) {
			if (newValue) {
				this.resetModal()
				this.loadUsersData()
			}
		}
	},
	methods: {
		getUserInitials,

		resetModal() {
			this.activeTab = 'users'
			this.selectedRecipients = []
			this.users = []
			this.filteredUsers = []
			this.userSearchQuery = ''
			this.error = null
			this.loading = false
			this.loadingUsers = false
		},

		async loadUsersData() {
			try {
				this.loadingUsers = true
				this.error = null

				const result = await this.loadUsers()

				if (result.success) {
					this.users = result.data
					this.filteredUsers = result.data

					this.$emit('load-user-images', result.data)
				} else {
					this.error = result.error
				}
			} catch (err) {
				this.error = 'Failed to load users. Please try again.'
			} finally {
				this.loadingUsers = false
			}
		},

		handleUsersImageError(userId) {
			this.$emit('handle-users-image-error', userId);
		},

		async submitForward() {
			if (!this.message || this.selectedRecipients.length === 0) {
				this.error = 'Please select at least one recipient'
				return
			}

			try {
				this.loading = true
				this.error = null

				const recipients = this.selectedRecipients.map(recipient => ({
					id: recipient.id,
					type: recipient.type
				}))

				const result = await this.forwardMessage(recipients)

				if (result.success) {
					this.resetModal()
				} else {
					this.error = result.error
				}
			} catch (err) {
				console.error('Error in submitForward:', err)
				this.error = 'An unexpected error occurred. Please try again.'
			} finally {
				this.loading = false
			}
		},

		filterUsers() {
			if (!this.userSearchQuery.trim()) {
				this.filteredUsers = this.users
				return
			}

			const query = this.userSearchQuery.toLowerCase()
			this.filteredUsers = this.users.filter(user =>
				user.username.toLowerCase().includes(query)
			)
		},

		isRecipientSelected(id, type) {
			return this.selectedRecipients.some(recipient =>
				recipient.id === id && recipient.type === type
			)
		},

		toggleUserSelection(user) {
			const existingIndex = this.selectedRecipients.findIndex(recipient =>
				recipient.id === user.id && recipient.type === 'user'
			)

			if (existingIndex > -1) {
				this.selectedRecipients.splice(existingIndex, 1)
			} else {
				this.selectedRecipients.push({
					id: user.id,
					type: 'user',
					name: user.username
				})
			}
		},

		toggleChatSelection(chat) {
			const existingIndex = this.selectedRecipients.findIndex(recipient =>
				recipient.id === chat.id && recipient.type === 'chat'
			)

			if (existingIndex > -1) {
				this.selectedRecipients.splice(existingIndex, 1)
			} else {
				this.selectedRecipients.push({
					id: chat.id,
					type: 'chat',
					name: this.getChatName(chat)
				})
			}
		},

		removeRecipient(recipient) {
			const index = this.selectedRecipients.findIndex(r =>
				r.id === recipient.id && r.type === recipient.type
			)
			if (index > -1) {
				this.selectedRecipients.splice(index, 1)
			}
		},

		getChatName(chat) {
			if (chat.isGroup) {
				return chat.name || 'Unnamed Group'
			} else {
				return chat.name || 'Private Chat'
			}
		},

		getChatInitials(chat) {
			const name = this.getChatName(chat)
			if (chat.isGroup) {
				return name.charAt(0).toUpperCase()
			} else {
				const words = name.split(' ')
				if (words.length >= 2) {
					return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase()
				}
				return name.charAt(0).toUpperCase()
			}
		},

		getForwardPreviewText(message) {
			if (!message) return ''

			if (message.type === 'image') {
				return message.text ? message.text : '📷 Photo'
			} else if (message.type === 'gif') {
				return message.text ? message.text : '🎞️ GIF'
			}

			return message.text || ''
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
@import url('../assets/modals.css');

.forward-modal {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1050;
}

.forward-modal-dialog {
	background: white;
	border-radius: 0.5rem;
	box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
	max-width: 700px;
	width: 95%;
	max-height: 90vh;
	overflow-y: auto;
}

.recipient-tabs {
	display: flex;
	border-bottom: 1px solid #dee2e6;
	margin-bottom: 1rem;
}

.recipient-tab {
	flex: 1;
	padding: 0.75rem 1rem;
	background: none;
	border: none;
	border-bottom: 3px solid transparent;
	cursor: pointer;
	font-weight: 500;
	color: #6c757d;
	transition: all 0.2s ease;
}

.selected-recipients {
	min-height: 40px;
	padding: 0.5rem;
	background-color: #f8f9fa;
	border-radius: 0.375rem;
	border: 1px solid #dee2e6;
	margin-bottom: 1rem;
}

.recipient-tab.active {
	color: #007bff;
	border-bottom-color: #007bff;
	background-color: #f8f9fa;
}

.recipient-tab:hover:not(.active) {
	background-color: #f8f9fa;
	color: #495057;
}

.recipient-chip {
	background-color: #007bff;
	color: white;
	padding: 0.375rem 0.75rem;
	border-radius: 1rem;
	font-size: 0.875rem;
	margin: 0.25rem;
	display: inline-flex;
	align-items: center;
	gap: 0.5rem;
}

.recipient-chip.user-chip {
	background-color: #007bff;
}

.recipient-chip.chat-chip {
	background-color: #28a745;
}

.recipient-chip-close {
	background: none;
	border: none;
	color: white;
	cursor: pointer;
	padding: 0;
	width: 1rem;
	height: 1rem;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	font-size: 0.75rem;
}

.recipient-chip-close:hover {
	background-color: rgba(255, 255, 255, 0.2);
}

.forward-chat-list {
	background-color: #f8f9fa;
	border: 1px solid #dee2e6;
	border-radius: 0.375rem;
	max-height: 200px;
	overflow-y: auto;
}

.forward-chat-item {
	padding: 0.75rem;
	cursor: pointer;
	border-bottom: 1px solid #e9ecef;
	transition: background-color 0.2s ease;
	display: flex;
	align-items: center;
	gap: 0.75rem;
}

.forward-chat-item:last-child {
	border-bottom: none;
}

.forward-chat-item:hover {
	background-color: #e9ecef;
}

.forward-chat-item.selected {
	background-color: #e3f2fd;
}

.forward-chat-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	display: flex;
	align-items: center;
	justify-content: center;
	color: white;
	font-weight: 600;
	font-size: 12px;
	flex-shrink: 0;
}

.forward-chat-info {
	flex: 1;
}

.forward-chat-name {
	font-weight: 500;
	margin-bottom: 0.25rem;
}

.forward-chat-type {
	font-size: 0.75rem;
	color: #6c757d;
}

.forward-message-preview {
	background-color: #f8f9fa;
	border: 1px solid #dee2e6;
	border-radius: 0.5rem;
	padding: 0.75rem;
	margin-bottom: 1rem;
}

.forward-preview-header {
	font-size: 0.875rem;
	font-weight: 600;
	color: #495057;
	margin-bottom: 0.5rem;
}

.forward-preview-content {
	font-size: 0.875rem;
	color: #6c757d;
	font-style: italic;
	max-height: 3em; /* Original compact height for text-only */
	overflow: hidden;
	text-overflow: ellipsis;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
	line-height: 1.5;
}

/* Enhanced styles when image is present */
.forward-preview-content.has-image {
	display: flex;
	align-items: flex-start;
	gap: 0.75rem;
	max-height: 80px; /* More space when image is present */
	-webkit-box-orient: unset; /* Reset flexbox orientation */
}

.forward-preview-image {
	width: 60px;
	height: 60px;
	border-radius: 0.25rem;
	object-fit: cover;
	flex-shrink: 0;
}

.has-image .forward-preview-text {
	flex: 1;
	max-height: 4.5em;
	overflow: hidden;
	text-overflow: ellipsis;
	display: -webkit-box;
	-webkit-line-clamp: 3;
	-webkit-box-orient: vertical;
	line-height: 1.5;
}

.forward-empty-state {
	text-align: center;
	padding: 2rem;
	color: #6c757d;
}

.forward-empty-state svg {
	margin-bottom: 1rem;
	opacity: 0.5;
}

.forward-loading {
	text-align: center;	padding: 2rem;
}

.forward-loading .spinner-border {
	width: 2rem;
	height: 2rem;
}

</style>
