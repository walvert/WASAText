<template>
	<div class="container-fluid">
		<div class="row">
			<!-- Sidebar with conversation list -->
			<div class="col-md-4 col-lg-3 sidebar">
				<div class="sidebar-sticky pt-3">
					<div class="d-flex justify-content-between align-items-center px-3 mb-3">
						<h5 class="mb-0">Conversations</h5>
						<button class="btn btn-sm btn-outline-primary" @click="showNewChatModal = true">
							<span class="feather">+</span> New Chat
						</button>
					</div>

					<LoadingSpinner v-if="loading" class="d-flex justify-content-center my-5" />
					<ErrorMsg v-if="error" :message="error" class="mx-3" />

					<div v-if="!loading && !error" class="list-group list-group-flush">
						<div
							v-for="chat in chats"
							:key="chat.id"
							class="list-group-item list-group-item-action"
							:class="{ active: selectedChatId === chat.id }"
							@click="selectChat(chat.id)"
						>
							<div class="d-flex w-100 justify-content-between">
								<h6 class="mb-1">
									{{ getChatName(chat) }}
								</h6>
								<small v-if="chat.lastMessage" class="text-nowrap">
									{{ formatTime(chat.lastMessage.created_at) }}
								</small>
							</div>
							<p v-if="chat.lastMessage" class="mb-1 text-truncate">
								{{ chat.lastMessage.text }}
							</p>
							<small v-if="chat.is_group" class="text-muted">
								Group · {{ chat.memberCount || 0 }} members
							</small>
						</div>

						<div v-if="chats.length === 0" class="text-center p-4 text-muted">
							No conversations yet.
						</div>
					</div>
				</div>
			</div>

			<!-- Main conversation area -->
			<main class="col-md-8 col-lg-9 ms-sm-auto px-md-4">
				<div v-if="!selectedChatId" class="d-flex justify-content-center align-items-center h-100">
					<div class="text-center text-muted">
						<div class="mb-3">
							<svg width="64" height="64" stroke="currentColor" stroke-width="2" fill="none" class="feather feather-message-square">
								<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
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
							<div v-if="selectedChat && selectedChat.is_group" class="dropdown">
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
								class="message mb-3"
								:class="{ 'message-sent': message.sender_id === currentUserId }"
							>
								<div class="message-content">
									<div v-if="message.sender_id !== currentUserId" class="message-sender">
										{{ message.sender_username || 'Unknown' }}
									</div>
									<div class="message-bubble p-2">
										{{ message.text }}
									</div>
									<div class="message-time small text-muted">
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
		<div class="modal fade" tabindex="-1" id="newChatModal" v-if="showNewChatModal">
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
	</div>
</template>

<script>
export default {
	name: 'ConversationsView',
	data() {
		return {
			// Conversations list
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

			// Modals
			showNewChatModal: false,
			newChatType: 'private',
			newChatUsername: '',
			newChatName: '',
			newChatLoading: false,
			newChatError: null,

			showRenameGroupModal: false,
			showAddMemberModal: false
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
			this.$router.push('/session');
			return;
		}

		// Set default auth header
		this.$axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

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

				// Get user chats
				const response = await this.$axios.get('/chats');
				this.chats = response.data;

				// Get current user ID
				// Note: This depends on how you're storing the user info
				// You might need to modify this based on your API
				const userResponse = await this.$axios.get('/user');
				this.currentUserId = userResponse.data.id;

			} catch (err) {
				console.error('Failed to fetch chats', err);
				this.error = 'Failed to load conversations. Please try again.';
			} finally {
				this.loading = false;
			}
		},

		async selectChat(chatId) {
			this.selectedChatId = chatId;
			await this.fetchMessages(chatId);
		},

		async fetchMessages(chatId) {
			try {
				this.loadingMessages = true;
				this.messagesError = null;

				const response = await this.$axios.get(`/chats/${chatId}/messages`);
				this.messages = response.data;

				// Scroll to bottom of messages
				this.$nextTick(() => {
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
					msgType: 'text' // Based on your backend
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

		getChatName(chat) {
			if (chat.is_group) {
				return chat.chat_name || 'Unnamed Group';
			} else {
				// For private chats, show the other user's name
				// This depends on your API response structure
				return chat.otherUsername || chat.chat_name || 'Private Chat';
			}
		},

		formatTime(timestamp) {
			if (!timestamp) return '';

			const date = new Date(timestamp);
			const now = new Date();

			// If today, show time only
			if (date.toDateString() === now.toDateString()) {
				return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
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
.sidebar {
	height: 100vh;
	border-right: 1px solid rgba(0, 0, 0, 0.1);
	overflow-y: auto;
}

.chat-container {
	height: 100vh;
}

.chat-messages {
	height: calc(100vh - 140px);
	overflow-y: auto;
}

.message {
	display: flex;
	margin-bottom: 12px;
}

.message-sent {
	justify-content: flex-end;
}

.message-content {
	max-width: 75%;
}

.message-bubble {
	background-color: #f0f2f5;
	border-radius: 18px;
	padding: 8px 12px;
	margin-bottom: 2px;
}

.message-sent .message-bubble {
	background-color: #dcf8c6;
}

.message-time {
	font-size: 0.7rem;
	text-align: right;
}

.list-group-item.active {
	background-color: rgba(13, 110, 253, 0.1);
	color: inherit;
	border-color: rgba(0, 0, 0, 0.125);
}

.list-group-item.active:hover {
	background-color: rgba(13, 110, 253, 0.15);
}

.list-group-item {
	cursor: pointer;
}

/* Make modal visible */
.modal {
	display: block !important;
	background-color: rgba(0, 0, 0, 0.5);
}
</style>
