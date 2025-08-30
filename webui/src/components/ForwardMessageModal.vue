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
						<div class="forward-preview-content">
							<img
								v-if="message && (message.type === 'image' || message.type === 'gif') && messageImageUrls[message.id]"
								:src="messageImageUrls[message.id]"
								:alt="message.mediaUrl || 'Image'"
								class="forward-preview-image"
							>
							{{ getForwardPreviewText(message) }}
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
									<button class="btn btn-outline-secondary" type="button" @click="loadUsersData">
										<svg v-if="!loadingUsers" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
											<path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"/>
										</svg>
										<span v-else class="spinner-border spinner-border-sm" role="status">
											<span class="visually-hidden">Loading...</span>
										</span>
									</button>
								</div>
							</div>

							<!-- Loading state for users -->
							<div v-if="loadingUsers" class="forward-loading">
								<div class="spinner-border" role="status">
									<span class="visually-hidden">Loading users...</span>
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
										<div class="avatar-circle" style="width: 32px; height: 32px; font-size: 12px;">
											<span class="avatar-text">{{ getUserInitials(user.username) }}</span>
										</div>
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
		}
	},
	emits: ['close'],
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
					console.log('Forward users fetched successfully:', this.users)
				} else {
					this.error = result.error
				}
			} catch (err) {
				console.error('Error loading users:', err)
				this.error = 'Failed to load users. Please try again.'
			} finally {
				this.loadingUsers = false
			}
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

		getUserInitials(username) {
			if (!username) return '?'
			const words = username.split(' ')
			if (words.length >= 2) {
				return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase()
			}
			return username.charAt(0).toUpperCase()
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
</style>
