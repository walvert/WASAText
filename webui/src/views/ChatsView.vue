<template>
	<div class="container-fluid">
		<div class="row">
			<!-- Sidebar with conversation list -->
			<div class="col-md-4 col-lg-3 chat-sidebar">
				<div class="sidebar-sticky pt-3">
					<!-- Welcome message -->
					<div class="chat-welcome-section px-3 mb-3">
						<h6 class="text-muted mb-0">Welcome @{{ currentUsername || 'User' }}</h6>
					</div>

					<!-- Header with title and action icons -->
					<div class="d-flex justify-content-between align-items-center px-3 mb-3">
						<h5 class="mb-0">Chats</h5>
						<div class="d-flex chat-gap-2">
							<button
								class="btn btn-sm btn-outline-secondary p-1 chat-btn-icon"
								@click="showNewChatModal = true"
								title="New Chat"
							>
								<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
									<path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
								</svg>
							</button>
							<button
								class="btn btn-sm btn-outline-secondary p-1 chat-btn-icon"
								@click="showProfileModal = true"
								title="Edit Profile"
							>
								<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
									<path d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6Zm2-3a2 2 0 1 1-4 0 2 2 0 0 1 4 0Zm4 8c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4Zm-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10Z"/>
								</svg>
							</button>
							<button
								class="btn btn-sm btn-outline-danger p-1 chat-btn-icon"
								@click="logout"
								title="Logout"
							>
								<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
									<path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0v2z"/>
									<path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3z"/>
								</svg>
							</button>
						</div>
					</div>

					<LoadingSpinner v-if="loading" class="d-flex justify-content-center my-5" />
					<ErrorMsg v-if="error" :message="error" class="mx-3" />

					<!-- Modern chat list -->
					<div v-if="!loading && !error" class="chat-list">
						<div
							v-for="chat in chats"
							:key="chat.id"
							class="chat-item"
							:class="{ active: selectedChatId === chat.id }"
							@click="selectChat(chat.id)"
						>
							<div class="chat-avatar">
								<div class="avatar-circle">
									<span class="avatar-text">{{ getChatInitials(chat) }}</span>
								</div>
							</div>
							<div class="chat-info">
								<div class="chat-header">
									<h6 class="chat-name">{{ getChatName(chat) }}</h6>
									<span class="chat-time">{{ formatTime(chat.lastMsgTime) }}</span>
								</div>
								<div class="chat-preview">
									<p class="last-message">{{ getLastMessagePreview(chat) }}</p>
									<div class="message-type-indicator">
										<span v-if="chat.lastMsgType === 'text'" class="message-type-icon">💬</span>
										<span v-else-if="chat.lastMsgType === 'image'" class="message-type-icon">📷</span>
										<span v-else class="message-type-icon">📎</span>
									</div>
								</div>
							</div>
						</div>

						<div v-if="chats.length === 0" class="chat-empty-state">
							<div class="chat-empty-icon">
								<svg width="48" height="48" fill="currentColor" viewBox="0 0 16 16">
									<path d="M14 1a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H4.414A2 2 0 0 0 3 11.586l-2 2V2a1 1 0 0 1 1-1h12zM2 0a2 2 0 0 0-2 2v12.793a.5.5 0 0 0 .854.353l2.853-2.853A1 1 0 0 1 4.414 12H14a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H2z"/>
								</svg>
							</div>
							<h6 class="text-muted">No conversations yet</h6>
							<p class="text-muted small">Start chatting by creating a new conversation</p>
						</div>
					</div>
				</div>
			</div>

			<!-- Main conversation area -->
			<main class="col-md-8 col-lg-9 ms-sm-auto px-md-4">
				<div v-if="!selectedChatId" class="d-flex justify-content-center align-items-center h-100">
					<div class="text-center text-muted">
						<div class="mb-3">
							<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
								<path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
							</svg>
						</div>
						<h5>Select a conversation to start chatting</h5>
						<p class="mb-4">Or create a new chat to get started</p>
						<button class="btn btn-primary" @click="showNewChatModal = true">New Chat</button>
					</div>
				</div>

				<div v-else class="chat-container h-100 d-flex flex-column">
					<!-- Chat header -->
					<div class="chat-header p-3 border-bottom">
						<div class="d-flex justify-content-between align-items-center">
							<h5 class="mb-0">{{ selectedChat ? getChatName(selectedChat) : '' }}</h5>
							<div v-if="selectedChat && selectedChat.isGroup" class="dropdown">
								<button class="btn btn-sm btn-outline-secondary" data-bs-toggle="dropdown">
									<span class="feather">⋮</span>
								</button>
								<ul class="dropdown-menu dropdown-menu-end">
									<li><button class="dropdown-item" @click="showRenameGroupModal = true">Rename Group</button></li>
									<li><button class="dropdown-item" @click="showAddMemberModal = true">Add Member</button></li>
									<li><hr class="dropdown-divider"></li>
									<li><button class="dropdown-item text-danger" @click="leaveGroup">Leave Group</button></li>
								</ul>
							</div>
						</div>
					</div>

					<!-- Messages list -->
					<div class="chat-messages p-3 flex-grow-1 overflow-auto">
						<LoadingSpinner v-if="loadingMessages" class="d-flex justify-content-center my-5" />
						<ErrorMsg v-if="messagesError" :message="messagesError" />

						<div v-if="!loadingMessages && !messagesError" class="messages-list">
							<div
								v-for="message in messages"
								:key="message.id"
								class="chat-message mb-3"
								:class="{ 'chat-message-sent': message.sender_id === currentUserId }"
							>
								<div class="chat-message-content">
									<div v-if="message.sender_id !== currentUserId" class="message-sender">
										{{ message.sender_username || 'Unknown' }}
									</div>
									<div class="chat-message-bubble p-2">
										{{ message.text }}
									</div>
									<div class="chat-message-time small text-muted">
										{{ formatTime(message.created_at) }}
									</div>
								</div>
							</div>

							<div v-if="messages.length === 0" class="text-center p-4 text-muted">
								No messages yet. Start the conversation!
							</div>
						</div>
					</div>

					<!-- Message input -->
					<div class="chat-input p-3 border-top">
						<form @submit.prevent="sendMessage">
							<div class="input-group">
								<input
									type="text"
									class="form-control"
									placeholder="Type a message..."
									v-model="newMessage"
									:disabled="sendingMessage"
								>
								<button
									class="btn btn-primary"
									type="submit"
									:disabled="sendingMessage || !newMessage.trim()"
								>
									<LoadingSpinner v-if="sendingMessage" class="me-1" style="width: 16px; height: 16px;" />
									<span v-else>Send</span>
								</button>
							</div>
						</form>
					</div>
				</div>
			</main>
		</div>

		<!-- New Chat Modal -->
		<div class="modal fade chat-modal" tabindex="-1" id="newChatModal" v-if="showNewChatModal">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">New Conversation</h5>
						<button type="button" class="btn-close" @click="showNewChatModal = false"></button>
					</div>
					<div class="modal-body">
						<div class="form-check mb-3">
							<input class="form-check-input" type="radio" name="chatType" id="privateChat" value="private" v-model="newChatType">
							<label class="form-check-label" for="privateChat">
								Private Chat
							</label>
						</div>
						<div class="form-check mb-3">
							<input class="form-check-input" type="radio" name="chatType" id="groupChat" value="group" v-model="newChatType">
							<label class="form-check-label" for="groupChat">
								Group Chat
							</label>
						</div>

						<div v-if="newChatType === 'private'" class="mb-3">
							<label class="form-label">Username</label>
							<input type="text" class="form-control" v-model="newChatUsername" placeholder="Enter username">
						</div>

						<div v-if="newChatType === 'group'" class="mb-3">
							<label class="form-label">Group Name</label>
							<input type="text" class="form-control" v-model="newChatName" placeholder="Enter group name">
						</div>

						<ErrorMsg v-if="newChatError" :message="newChatError" />
					</div>
					<div class="modal-footer">
						<button type="button" class="btn btn-secondary" @click="showNewChatModal = false">Cancel</button>
						<button
							type="button"
							class="btn btn-primary"
							@click="createNewChat"
							:disabled="newChatLoading"
						>
							<LoadingSpinner v-if="newChatLoading" class="me-2" />
							Create
						</button>
					</div>
				</div>
			</div>
		</div>

		<!-- Profile Modal -->
		<div class="modal fade chat-modal" tabindex="-1" id="profileModal" v-if="showProfileModal">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">Edit Profile</h5>
						<button type="button" class="btn-close" @click="showProfileModal = false"></button>
					</div>
					<div class="modal-body">
						<div class="mb-3">
							<label class="form-label">Username</label>
							<input type="text" class="form-control" v-model="editUsername" placeholder="Enter new username">
						</div>
						<ErrorMsg v-if="profileError" :message="profileError" />
					</div>
					<div class="modal-footer">
						<button type="button" class="btn btn-secondary" @click="showProfileModal = false">Cancel</button>
						<button
							type="button"
							class="btn btn-primary"
							@click="updateProfile"
							:disabled="profileLoading"
						>
							<LoadingSpinner v-if="profileLoading" class="me-2" />
							Update
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'ChatsView',
	data() {
		return {
			// Chats list
			chats: [],
			loading: true,
			error: null,
			selectedChatId: null,

			// Current chat messages
			messages: [],
			loadingMessages: false,
			messagesError: null,
			newMessage: '',
			sendingMessage: false,

			// User info
			currentUserId: null,
			currentUsername: null,

			// Modals
			showNewChatModal: false,
			newChatType: 'private',
			newChatUsername: '',
			newChatName: '',
			newChatLoading: false,
			newChatError: null,

			showRenameGroupModal: false,
			showAddMemberModal: false,

			// Profile modal
			showProfileModal: false,
			editUsername: '',
			profileLoading: false,
			profileError: null
		}
	},
	computed: {
		selectedChat() {
			return this.chats.find(chat => chat.id === this.selectedChatId);
		}
	},
	async created() {
		// Check if user is authenticated
		const token = localStorage.getItem('token');
		if (!token) {
			console.log('No token found, emitting logout');
			this.$emit('logout');
			return;
		}

		console.log('Token found, setting up ChatsView');

		// Get user data from localStorage
		const userData = localStorage.getItem('user');
		if (userData) {
			try {
				const user = JSON.parse(userData);
				this.currentUsername = user.username;
				console.log('Loaded username from localStorage:', this.currentUsername);
			} catch (e) {
				console.error('Error parsing user data from localStorage', e);
			}
		}

		// Set default auth header
		this.$axios.defaults.headers.common['Authorization'] = `${token}`;
		console.log('Set Authorization header');

		// Fetch user chats
		await this.fetchChats();

		// Start polling for new messages
		this.startPolling();
	},
	beforeUnmount() {
		// Stop polling when component is destroyed
		this.stopPolling();
	},
	methods: {
		async fetchChats() {
			try {
				this.loading = true;
				this.error = null;

				console.log('Fetching chats with token:', localStorage.getItem('token') ? 'present' : 'missing');

				// Get user chats - the backend should identify the user from the token
				const response = await this.$axios.get('/chats');
				console.log('Chats fetched successfully:', response.data);
				this.chats = response.data;

			} catch (err) {
				console.error('Failed to fetch chats', err);
				console.error('Error details:', {
					status: err.response?.status,
					data: err.response?.data,
					headers: err.response?.headers
				});

				this.error = 'Failed to load conversations. Please try again.';

				// If it's an auth error, logout
				if (err.response?.status === 401) {
					console.log('401 error in fetchChats, logging out');
					this.$emit('logout');
				}
			} finally {
				this.loading = false;
			}
		},

		async fetchMessages(chatId) {
			try {
				this.loadingMessages = true;
				this.messagesError = null;

				const response = await this.$axios.get(`/chats/${chatId}`);
				this.messages = response.data;

				// Scroll to bottom of messages
				await this.$nextTick(() => {
					this.scrollToBottom();
				});

			} catch (err) {
				console.error('Failed to fetch messages', err);
				this.messagesError = 'Failed to load messages. Please try again.';
			} finally {
				this.loadingMessages = false;
			}
		},

		async sendMessage() {
			if (!this.newMessage.trim() || !this.selectedChatId) return;

			try {
				this.sendingMessage = true;

				await this.$axios.post(`/chats/${this.selectedChatId}/messages`, {
					text: this.newMessage.trim(),
					msgType: 'text'
				});

				// Clear input and refresh messages
				this.newMessage = '';
				await this.fetchMessages(this.selectedChatId);

			} catch (err) {
				console.error('Failed to send message', err);
				alert('Failed to send message. Please try again.');
			} finally {
				this.sendingMessage = false;
			}
		},

		async createNewChat() {
			if (this.newChatType === 'private' && !this.newChatUsername.trim()) {
				this.newChatError = 'Please enter a username';
				return;
			}

			if (this.newChatType === 'group' && !this.newChatName.trim()) {
				this.newChatError = 'Please enter a group name';
				return;
			}

			try {
				this.newChatLoading = true;
				this.newChatError = null;

				if (this.newChatType === 'private') {
					// Create private chat
					const response = await this.$axios.post('/chats/private', {
						username: this.newChatUsername.trim()
					});

					// Reset form and close modal
					this.newChatUsername = '';
					this.showNewChatModal = false;

					// Refresh chats and select the new one
					await this.fetchChats();
					this.selectChat(response.data.id);

				} else {
					// Create group chat
					const response = await this.$axios.post('/chats/group', {
						name: this.newChatName.trim()
					});

					// Reset form and close modal
					this.newChatName = '';
					this.showNewChatModal = false;

					// Refresh chats and select the new one
					await this.fetchChats();
					this.selectChat(response.data.id);
				}

			} catch (err) {
				console.error('Failed to create chat', err);
				this.newChatError = err.response?.data?.message || 'Failed to create conversation. Please try again.';
			} finally {
				this.newChatLoading = false;
			}
		},

		async updateProfile() {
			if (!this.editUsername.trim()) {
				this.profileError = 'Username is required';
				return;
			}

			try {
				this.profileLoading = true;
				this.profileError = null;

				await this.$axios.put('/user', {
					username: this.editUsername.trim(),
				});

				// Update current user info
				this.currentUsername = this.editUsername.trim();

				// Update localStorage with new user data
				localStorage.setItem('username', this.currentUsername);

				// If you store user data as JSON object, update it
				const userData = {
					id: this.currentUserId,
					username: this.currentUsername,
				};
				localStorage.setItem('user', JSON.stringify(userData));

				// Close modal
				this.showProfileModal = false;

				// Show success message
				alert('Profile updated successfully!');

			} catch (err) {
				console.error('Failed to update profile', err);
				this.profileError = err.response?.data?.message || 'Failed to update profile. Please try again.';
			} finally {
				this.profileLoading = false;
			}
		},

		async logout() {
			if (confirm('Are you sure you want to logout?')) {
				// Clear all localStorage data
				localStorage.removeItem('token');
				localStorage.removeItem('user');

				// Clear axios auth header
				delete this.$axios.defaults.headers.common['Authorization'];

				// Stop polling
				this.stopPolling();

				// Emit logout event to App.vue
				this.$emit('logout');
			}
		},

		async leaveGroup() {
			if (!confirm('Are you sure you want to leave this group?')) return;

			try {
				await this.$axios.delete(`/chats/${this.selectedChatId}/members`);

				// Refresh chats and clear selection
				await this.fetchChats();
				this.selectedChatId = null;

			} catch (err) {
				console.error('Failed to leave group', err);
				alert('Failed to leave group. Please try again.');
			}
		},

		selectChat(chatId) {
			this.selectedChatId = chatId;
			this.fetchMessages(chatId);
		},

		getChatName(chat) {
			if (chat.isGroup) {
				return chat.name || 'Unnamed Group';
			} else {
				// For private chats, use the name field (which should contain the other user's name)
				return chat.name || 'Private Chat';
			}
		},

		getChatInitials(chat) {
			const name = this.getChatName(chat);
			if (chat.isGroup) {
				// For groups, use first letter of group name
				return name.charAt(0).toUpperCase();
			} else {
				// For private chats, use first letter of each word (up to 2 letters)
				const words = name.split(' ');
				if (words.length >= 2) {
					return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase();
				}
				return name.charAt(0).toUpperCase();
			}
		},

		getLastMessagePreview(chat) {
			if (!chat.lastMsgText) {
				return 'No messages yet';
			}

			// Truncate long messages
			const maxLength = 40;
			if (chat.lastMsgText.length > maxLength) {
				return chat.lastMsgText.substring(0, maxLength) + '...';
			}

			return chat.lastMsgText;
		},

		formatTime(timestamp) {
			if (!timestamp) return '';

			const date = new Date(timestamp);
			const now = new Date();

			// If today, show time only
			if (date.toDateString() === now.toDateString()) {
				return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
			}

			// If yesterday
			const yesterday = new Date(now);
			yesterday.setDate(now.getDate() - 1);
			if (date.toDateString() === yesterday.toDateString()) {
				return 'Yesterday';
			}

			// If this week (within 7 days)
			const weekAgo = new Date(now);
			weekAgo.setDate(now.getDate() - 7);
			if (date > weekAgo) {
				return date.toLocaleDateString([], { weekday: 'short' });
			}

			// If this year, show month and day
			if (date.getFullYear() === now.getFullYear()) {
				return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
			}

			// Otherwise show full date
			return date.toLocaleDateString();
		},

		scrollToBottom() {
			const messagesContainer = document.querySelector('.chat-messages');
			if (messagesContainer) {
				messagesContainer.scrollTop = messagesContainer.scrollHeight;
			}
		},

		startPolling() {
			// Poll for new messages every 5 seconds
			this.pollingInterval = setInterval(() => {
				if (this.selectedChatId) {
					this.fetchMessages(this.selectedChatId);
				}

				// Also refresh chat list occasionally
				if (Math.random() < 0.2) { // 20% chance to refresh chats on each poll
					this.fetchChats();
				}
			}, 5000);
		},

		stopPolling() {
			if (this.pollingInterval) {
				clearInterval(this.pollingInterval);
			}
		}
	}
}
</script>

<style scoped>
@import url('../assets/main.css');

/* Component-specific overrides if needed */
.chat-container {
	height: 100vh;
}

.chat-messages {
	height: calc(100vh - 140px);
	overflow-y: auto;
}
</style>
