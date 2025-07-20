<template src="./ChatsView.html"></template>

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
			pendingMessage: '', // Store the message being sent

			// Read status tracking
			lastReadMessageId: null, // Track the last read message ID for current chat

			// User info
			currentUserId: null,
			currentUsername: null,

			// Modals
			newChatType: 'private',
			newChatUsername: '',

			// Rename Group Modal
			showSetGroupNameModal: false,
			editGroupName: '',
			setGroupNameLoading: false,
			setGroupNameError: null,

			// Add member modal
			showAddToGroupModal: false,
			newMemberUsername: '',
			addToGroupLoading: false,
			addToGroupError: null,

			// Profile modal
			showProfileModal: false,
			editUsername: '',
			profileLoading: false,
			profileError: null,

			chatImageUrls: {}, // Store blob URLs for each chat
			imageCache: {}, // Cache blob URLs by image path

			// Polling improvements
			pollingInterval: null,
			isPolling: false,
			pollingRetryCount: 0,
			maxRetries: 3,
			basePollingInterval: 5000, // 5 seconds
			currentPollingInterval: 5000,

			// Request tracking
			activeRequests: new Set(),
			requestController: null,

			// New Chat Modal - Enhanced
			showNewChatModal: false, // Make sure this is initialized
			newChatLoading: false,
			newChatError: null,
			newChatName: '', // For group chats
			initialMessage: '', // Optional initial message

			// User Selection
			users: [],
			filteredUsers: [],
			selectedUsers: [],
			loadingUsers: false,
			userSearchQuery: '',

			// Set Group Photo Modal
			showSetGroupPhotoModal: false,
			selectedPhotoFile: null,
			photoPreviewUrl: null,
			setGroupPhotoLoading: false,
			setGroupPhotoError: null,

			// User profile image
			currentUserImageUrl: null,
			selectedProfileImage: null,
			profileImagePreviewUrl: null,

		}
	},

	computed: {
		selectedChat() {
			return this.chats.find(chat => chat.id === this.selectedChatId);
		}
	},
	async created() {
		document.addEventListener('visibilitychange', this.handleVisibilityChange);
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
				console.log('Loaded user from localStorage:', this.currentUsername);
			} catch (e) {
				console.error('Error parsing user data from localStorage', e);
			}
		}

		// Set default auth header
		this.$axios.defaults.headers.common['Authorization'] = `${token}`;
		console.log('Set Authorization header');

		await this.loadCurrentUserImage();

		// Fetch user chats
		await this.getMyConversations();

		// Start polling for new messages
		this.startPolling();
	},
	beforeUnmount() {
		// Enhanced cleanup
		this.cleanup();

		// Remove visibility change listener
		document.removeEventListener('visibilitychange', this.handleVisibilityChange);

		// Clean up blob URLs to prevent memory leaks
		Object.values(this.imageCache).forEach(blobUrl => {
			URL.revokeObjectURL(blobUrl);
		});
	},

	methods: {
		async getMyConversations() {
			// Create abort controller for this request
			const controller = new AbortController();
			const requestId = 'getMyConversations-' + Date.now();

			try {
				this.loading = true;
				this.error = null;

				// Track active request
				this.activeRequests.add(requestId);

				console.log('Fetching chats with token:', localStorage.getItem('token') ? 'present' : 'missing');

				// Add timeout and abort signal
				const response = await this.$axios.get('/chats', {
					signal: controller.signal,
					timeout: 10000 // 10 second timeout
				});

				console.log('Chats fetched successfully:', response.data);
				this.chats = response.data;

				await this.loadChatImages();

				// Reset retry count on success
				this.pollingRetryCount = 0;
				this.currentPollingInterval = this.basePollingInterval;

			} catch (err) {
				// Don't handle aborted requests
				if (err.name === 'AbortError' || err.code === 'ECONNABORTED') {
					console.log('Request aborted or timed out');
					return;
				}

				console.error('Failed to fetch chats', err);
				this.error = 'Failed to load conversations. Please try again.';

				// Increment retry count and adjust polling interval
				this.pollingRetryCount++;
				this.currentPollingInterval = Math.min(
					this.basePollingInterval * Math.pow(2, this.pollingRetryCount),
					30000 // Max 30 seconds
				);

				// If it's an auth error, logout
				if (err.response?.status === 401) {
					console.log('401 error in getMyConversations, logging out');
					this.$emit('logout');
				}
			} finally {
				this.loading = false;
				this.activeRequests.delete(requestId);
			}
		},

		async getConversation(chatId) {
			const controller = new AbortController();
			const requestId = 'getConversation-' + chatId + '-' + Date.now();

			try {
				this.loadingMessages = true;
				this.messagesError = null;

				this.activeRequests.add(requestId);

				const response = await this.$axios.get(`/chats/${chatId}`, {
					signal: controller.signal,
					timeout: 10000
				});

				// Sort messages by timestamp (oldest first, newest at bottom)
				this.messages = response.data.sort((a, b) => {
					return new Date(a.createdAt) - new Date(b.createdAt);
				});

				// Fetch last read message ID for this chat
				await this.getLastReadMessageId(chatId);

				// Scroll to bottom after messages are rendered - use setTimeout for better reliability
				setTimeout(() => {
					this.scrollToBottom();
				}, 100);

			} catch (err) {
				if (err.name === 'AbortError' || err.code === 'ECONNABORTED') {
					console.log('Messages request aborted or timed out');
					return;
				}

				console.error('Failed to fetch messages', err);
				this.messagesError = 'Failed to load messages. Please try again.';
			} finally {
				this.loadingMessages = false;
				this.activeRequests.delete(requestId);
			}
		},

		async getLastReadMessageId(chatId) {
			const controller = new AbortController();
			const requestId = 'getLastRead-' + chatId + '-' + Date.now();

			try {
				this.activeRequests.add(requestId);

				const response = await this.$axios.get(`/chats/${chatId}/last-read`, {
					signal: controller.signal,
					timeout: 5000
				});

				// Store the previous value to detect changes
				const previousLastRead = this.lastReadMessageId;
				this.lastReadMessageId = response.data.lastReadId;

				// Debug logging
				console.log('Last read message ID updated:', {
					chatId,
					previousLastRead,
					newLastRead: this.lastReadMessageId
				});

				// Force reactivity update if needed (Vue 2)
				this.$forceUpdate();

			} catch (err) {
				if (err.name === 'AbortError' || err.code === 'ECONNABORTED') {
					console.log('Last read request aborted or timed out');
					return;
				}

				console.error('Failed to fetch last read message ID', err);
				// Don't show error to user for this non-critical feature
			} finally {
				this.activeRequests.delete(requestId);
			}
		},

		isMessageRead(message) {
			// Only show read status for messages sent by current user
			if (!this.isCurrentUserMessage(message)) {
				return false;
			}

			// If we don't have lastReadMessageId yet, assume not read
			if (this.lastReadMessageId === null || this.lastReadMessageId === undefined) {
				return false;
			}

			// Ensure both values are numbers for proper comparison
			const messageId = parseInt(message.id);
			const lastReadId = parseInt(this.lastReadMessageId);

			// Message is read if its ID is less than or equal to lastReadMessageId
			return messageId <= lastReadId;
		},

		async sendMessage() {
			if (!this.newMessage.trim() || !this.selectedChatId) return;

			try {
				this.sendingMessage = true;
				this.pendingMessage = this.newMessage.trim();

				// Send message with correct format
				await this.$axios.post(`/chats/${this.selectedChatId}/messages`, {
					type: 'text',
					text: this.newMessage.trim()
				});

				// Clear input
				this.newMessage = '';
				this.pendingMessage = '';

				// Refresh messages and scroll to bottom
				await this.getConversation(this.selectedChatId);

			} catch (err) {
				console.error('Failed to send message', err);
				alert('Failed to send message. Please try again.');
			} finally {
				this.sendingMessage = false;
			}
		},

		scrollToBottom() {
			// Use nextTick to ensure DOM is updated
			this.$nextTick(() => {
				const messagesContainer = this.$refs.messagesContainer;
				if (messagesContainer) {
					// Force scroll to the very bottom
					messagesContainer.scrollTop = messagesContainer.scrollHeight;
				}
			});
		},

		handleNewMessage() {
			// Check if user is at or near the bottom of messages
			const messagesContainer = this.$refs.messagesContainer;
			if (messagesContainer) {
				const { scrollTop, scrollHeight, clientHeight } = messagesContainer;
				const isAtBottom = scrollTop + clientHeight >= scrollHeight - 100; // 100px threshold

				if (isAtBottom) {
					this.scrollToBottom();
				}
			}
		},

		async getUsers() {
			try {
				this.loadingUsers = true;
				this.newChatError = null;

				const response = await this.$axios.get('/users');
				this.users = response.data;
				this.filteredUsers = response.data;

				console.log('Users fetched successfully:', this.users);

			} catch (err) {
				console.error('Failed to fetch users', err);
				this.newChatError = 'Failed to load users. Please try again.';
			} finally {
				this.loadingUsers = false;
			}
		},

		filterUsers() {
			if (!this.userSearchQuery.trim()) {
				this.filteredUsers = this.users;
				return;
			}

			const query = this.userSearchQuery.toLowerCase();
			this.filteredUsers = this.users.filter(user =>
				user.username.toLowerCase().includes(query)
			);
		},

		toggleUserSelection(user) {
			const index = this.selectedUsers.findIndex(u => u.id === user.id);

			if (index > -1) {
				// User is already selected, remove them
				this.selectedUsers.splice(index, 1);
			} else {
				// User is not selected, add them
				this.selectedUsers.push(user);
			}

			// Clear group name if only one user is selected
			if (this.selectedUsers.length <= 1) {
				this.newChatName = '';
			}
		},

		removeUserSelection(user) {
			const index = this.selectedUsers.findIndex(u => u.id === user.id);
			if (index > -1) {
				this.selectedUsers.splice(index, 1);
			}

			// Clear group name if only one user is selected
			if (this.selectedUsers.length <= 1) {
				this.newChatName = '';
			}
		},

		getUserInitials(username) {
			if (!username) return '?';

			const words = username.split(' ');
			if (words.length >= 2) {
				return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase();
			}
			return username.charAt(0).toUpperCase();
		},

		async createNewChat() {
			// Validation
			if (this.selectedUsers.length === 0) {
				this.newChatError = 'Please select at least one recipient';
				return;
			}

			if (this.selectedUsers.length > 1 && !this.newChatName.trim()) {
				this.newChatError = 'Group name is required for group chats';
				return;
			}

			try {
				this.newChatLoading = true;
				this.newChatError = null;

				// Prepare the request body
				const requestBody = {
					type: 'text',
					receivers: this.selectedUsers.map(user => user.username)
				};

				// Add initial message if provided
				if (this.initialMessage.trim()) {
					requestBody.text = this.initialMessage.trim();
				}

				// Add group name if multiple users selected
				if (this.selectedUsers.length > 1) {
					requestBody.chatName = this.newChatName.trim();
				}

				console.log('Creating chat with payload:', requestBody);

				// Create the chat
				const response = await this.$axios.post('/chats', requestBody);

				console.log('Chat created successfully:', response.data);

				// Close modal and reset form
				this.closeNewChatModal();

				// Refresh chats and select the new one
				await this.getMyConversations();

				// Select the new chat if we have an ID from the response
				if (response.data && response.data.id) {
					this.selectChat(response.data.id);
				}

			} catch (err) {
				console.error('Failed to create chat', err);
				this.newChatError = err.response?.data?.message || 'Failed to create conversation. Please try again.';
			} finally {
				this.newChatLoading = false;
			}
		},

		closeNewChatModal() {
			this.showNewChatModal = false;

			// Reset all form data
			this.newChatName = '';
			this.initialMessage = '';
			this.selectedUsers = [];
			this.users = [];
			this.filteredUsers = [];
			this.userSearchQuery = '';
			this.newChatError = null;
			this.newChatLoading = false;
			this.loadingUsers = false;
		},

		openNewChatModal() {
			console.log('Opening new chat modal...'); // Debug log
			this.showNewChatModal = true;

			// Automatically fetch users when modal opens
			this.getUsers();
		},


		async setGroupName() {
			if (!this.editGroupName.trim()) {
				this.setGroupNameError = 'Group name is required';
				return;
			}

			if (!this.selectedChatId) {
				this.setGroupNameError = 'No group selected';
				return;
			}

			try {
				this.setGroupNameLoading = true;
				this.setGroupNameError = null;

				const response = await this.$axios.put(`/chats/${this.selectedChatId}`, {
					chatName: this.editGroupName.trim()
				});

				console.log('Group renamed successfully:', response.data);

				// Update the local chat name immediately for better UX
				const chatIndex = this.chats.findIndex(chat => chat.id === this.selectedChatId);
				if (chatIndex !== -1) {
					this.chats[chatIndex].name = this.editGroupName.trim();
				}

				// Close modal and reset form
				this.showSetGroupNameModal = false;
				this.editGroupName = '';

				// Refresh chats to get latest data from server
				await this.getMyConversations();

				// Show success message
				alert('Group renamed successfully!');

			} catch (err) {
				console.error('Failed to rename group', err);
				this.setGroupNameError = err.response?.data?.message || 'Failed to rename group. Please try again.';
			} finally {
				this.setGroupNameLoading = false;
			}
		},

		async addToGroup() {
			if (!this.newMemberUsername.trim()) {
				this.addToGroupError = 'Username is required';
				return;
			}

			if (!this.selectedChatId) {
				this.addToGroupError = 'No group selected';
				return;
			}

			try {
				this.addToGroupLoading = true;
				this.addToGroupError = null;

				const response = await this.$axios.post(`/chats/${this.selectedChatId}/members`, {
					username: this.newMemberUsername.trim()
				});


				console.log('Member added successfully:', response.data);

				// Close modal and reset form
				this.showAddToGroupModal = false;
				this.newMemberUsername = '';

				// Refresh chats to get latest data
				await this.getMyConversations();

				// Show success message
				alert(`${this.newMemberUsername.trim()} has been added to the group!`);

			} catch (err) {
				console.error('Failed to add member', err);

				// Handle specific error cases
				if (err.response?.status === 404) {
					this.addToGroupError = 'User not found. Please check the username.';
				} else if (err.response?.status === 409) {
					this.addToGroupError = 'User is already a member of this group.';
				} else {
					this.addToGroupError = err.response?.data?.message || 'Failed to add member. Please try again.';
				}
			} finally {
				this.addToGroupLoading = false;
			}
		},

		openSetGroupNameModal() {
			if (this.selectedChat) {
				this.editGroupName = this.selectedChat.name || '';
				this.setGroupNameError = null;
				this.showSetGroupNameModal = true;

				// Focus the input after modal is shown
				this.$nextTick(() => {
					if (this.$refs.groupNameInput) {
						this.$refs.groupNameInput.focus();
						this.$refs.groupNameInput.select();
					}
				});
			}
		},

		openAddToGroupModal() {
			this.newMemberUsername = '';
			this.addToGroupError = null;
			this.showAddToGroupModal = true;

			// Focus the input after modal is shown
			this.$nextTick(() => {
				if (this.$refs.memberUsernameInput) {
					this.$refs.memberUsernameInput.focus();
				}
			});
		},

		async updateProfile() {
			if (!this.editUsername.trim() && !this.selectedProfileImage) {
				this.profileError = 'Please enter a username or select a profile image to update';
				return;
			}

			try {
				this.profileLoading = true;
				this.profileError = null;

				// Update username if provided
				if (this.editUsername.trim()) {
					console.log('Updating username to:', this.editUsername.trim());

					const response = await this.$axios.put('/users', {
						username: this.editUsername.trim(),
					}, {
						headers: {
							'Content-Type': 'application/json'
						}
					});

					console.log('Username update response:', response.data);

					// Update current user info from server response
					if (response.data && response.data.username) {
						this.currentUsername = response.data.username;

						// Update localStorage with new user data
						localStorage.setItem('username', this.currentUsername);

						const userData = {
							id: response.data.id || this.currentUserId,
							username: this.currentUsername,
						};
						localStorage.setItem('user', JSON.stringify(userData));
					}
				}

				// Update profile image if provided
				if (this.selectedProfileImage) {
					console.log('Updating profile image...');

					const formData = new FormData();
					formData.append('image', this.selectedProfileImage);

					const response = await this.$axios.put('/users/image', formData, {
						headers: {
							'Content-Type': 'multipart/form-data'
						}
					});

					console.log('Image update response:', response.data);

					// Clear old cached image URL
					if (this.currentUserImageUrl) {
						URL.revokeObjectURL(this.currentUserImageUrl);
						this.currentUserImageUrl = null;
					}

					// If server returns the new image URL, use it to load the image
					if (response.data && response.data.imageUrl) {
						console.log('Loading new image from:', response.data.imageUrl);
						this.currentUserImageUrl = await this.getChatImageUrl(response.data.imageUrl);
					} else {
						// Fallback: reload user image from server
						await this.loadCurrentUserImage();
					}
				}

				// Close modal and reset
				this.showProfileModal = false;
				this.clearProfileImageSelection();
				this.editUsername = '';

				// Show success message
				alert('Profile updated successfully!');

			} catch (err) {
				console.error('Failed to update profile', err);
				console.error('Error response:', err.response);

				// Check for authentication errors
				if (err.response?.status === 401) {
					console.log('Authentication error during profile update, logging out');
					this.$emit('logout');
					return;
				}

				if (err.response?.status === 413) {
					this.profileError = 'File is too large. Please choose a smaller image.';
				} else if (err.response?.status === 400) {
					this.profileError = 'Invalid file format. Please choose a valid image.';
				} else {
					this.profileError = err.response?.data?.message || 'Failed to update profile. Please try again.';
				}
			} finally {
				this.profileLoading = false;
			}
		},


		// Add method to handle visibility changes (pause polling when tab is hidden)
		handleVisibilityChange() {
			if (document.hidden) {
				this.stopPolling();
			} else {
				this.startPolling();
			}
		},

		async leaveGroup() {
			if (!confirm('Are you sure you want to leave this group?')) return;

			try {
				await this.$axios.delete(`/chats/${this.selectedChatId}/members`);

				// Refresh chats and clear selection
				await this.getMyConversations();
				this.selectedChatId = null;

			} catch (err) {
				console.error('Failed to leave group', err);
				alert('Failed to leave group. Please try again.');
			}
		},

		selectChat(chatId) {
			this.selectedChatId = chatId;
			this.lastReadMessageId = null; // Reset read status when switching chats
			this.getConversation(chatId);
		},

		isCurrentUser(senderId) {
			return senderId === this.currentUserId;
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

		async getChatImageUrl(imagePath) {
			if (!imagePath) return null;

			// Check if we already have a blob URL cached for this image
			if (this.imageCache[imagePath]) {
				return this.imageCache[imagePath];
			}

			try {
				// Get the base URL from your axios configuration
				const baseURL = this.$axios.defaults.baseURL || 'http://localhost:3000';

				// The imagePath should be something like "uploads/chats/image.jpg"
				// We need to construct the URL as: /uploads/chats/image.jpg
				let imageUrl;

				// Remove leading slash if present and ensure it starts with uploads/
				const cleanPath = imagePath.startsWith('/') ? imagePath.slice(1) : imagePath;

				if (cleanPath.startsWith('uploads/')) {
					// Path already includes uploads/, use as is
					imageUrl = `${baseURL}/${cleanPath}`;
				} else {
					// Path doesn't include uploads/, add it
					imageUrl = `${baseURL}/uploads/${cleanPath}`;
				}

				console.log('Fetching image from:', imageUrl); // Debug log

				// Get the token from localStorage
				const token = localStorage.getItem('token');

				// Fetch the image with authorization header
				const response = await this.$axios.get(imageUrl, {
					responseType: 'blob',
					headers: {
						'Authorization': token
					}
				});

				// Convert blob to object URL
				const blobUrl = URL.createObjectURL(response.data);

				// Cache the blob URL
				this.imageCache[imagePath] = blobUrl;

				return blobUrl;

			} catch (error) {
				console.error('Failed to fetch chat image:', error);
				console.error('Image path was:', imagePath);
				return null;
			}
		},

		async loadChatImages() {
			for (const chat of this.chats) {
				if (chat.image && !this.chatImageUrls[chat.id]) {
					try {
						const imageUrl = await this.getChatImageUrl(chat.image);
						if (imageUrl) {
							this.chatImageUrls[chat.id] = imageUrl;
						}
					} catch (error) {
						console.error(`Failed to load image for chat ${chat.id}:`, error);
					}
				}
			}
		},

		async loadCurrentUserImage() {
			try {
				console.log('Loading current user image...');

				const response = await this.$axios.get('/users/image');
				console.log('User image response:', response.data);

				if (response.data && response.data.imageUrl) {
					console.log('Loading image from URL:', response.data.imageUrl);
					this.currentUserImageUrl = await this.getChatImageUrl(response.data.imageUrl);
					console.log('User image loaded successfully');
				} else {
					console.log('No image URL in response, using initials');
				}
			} catch (error) {
				console.error('Failed to load user image:', error);

				// Check if it's a 404 (no image set) vs other errors
				if (error.response?.status === 404) {
					console.log('User has no profile image set');
				} else if (error.response?.status === 401) {
					console.log('Authentication error loading user image');
					// Don't emit logout here as this is called during initialization
				}

				// Don't show error to user, just use initials
				this.currentUserImageUrl = null;
			}
		},

		// Handle user image load error
		handleUserImageError() {
			this.currentUserImageUrl = null;
		},

		// Handle profile image selection
		handleProfileImageSelect(event) {
			const file = event.target.files[0];

			if (!file) {
				this.clearProfileImageSelection();
				return;
			}

			// Validate file type
			if (!file.type.startsWith('image/')) {
				this.profileError = 'Please select a valid image file';
				this.clearProfileImageSelection();
				return;
			}

			// Validate file size (5MB limit)
			const maxSize = 5 * 1024 * 1024; // 5MB in bytes
			if (file.size > maxSize) {
				this.profileError = 'File size must be less than 5MB';
				this.clearProfileImageSelection();
				return;
			}

			this.selectedProfileImage = file;
			this.profileError = null;

			// Create preview URL
			if (this.profileImagePreviewUrl) {
				URL.revokeObjectURL(this.profileImagePreviewUrl);
			}
			this.profileImagePreviewUrl = URL.createObjectURL(file);
		},

		clearProfileImageSelection() {
			this.selectedProfileImage = null;

			if (this.profileImagePreviewUrl) {
				URL.revokeObjectURL(this.profileImagePreviewUrl);
				this.profileImagePreviewUrl = null;
			}

			// Clear the file input
			if (this.$refs.profileImageInput) {
				this.$refs.profileImageInput.value = '';
			}
		},

		openSetGroupPhotoModal() {
			this.selectedPhotoFile = null;
			this.photoPreviewUrl = null;
			this.setGroupPhotoError = null;
			this.showSetGroupPhotoModal = true;
		},

		handlePhotoFileSelect(event) {
			const file = event.target.files[0];

			if (!file) {
				this.clearPhotoSelection();
				return;
			}

			// Validate file type
			if (!file.type.startsWith('image/')) {
				this.setGroupPhotoError = 'Please select a valid image file';
				this.clearPhotoSelection();
				return;
			}

			// Validate file size (5MB limit)
			const maxSize = 5 * 1024 * 1024; // 5MB in bytes
			if (file.size > maxSize) {
				this.setGroupPhotoError = 'File size must be less than 5MB';
				this.clearPhotoSelection();
				return;
			}

			this.selectedPhotoFile = file;
			this.setGroupPhotoError = null;

			// Create preview URL
			if (this.photoPreviewUrl) {
				URL.revokeObjectURL(this.photoPreviewUrl);
			}
			this.photoPreviewUrl = URL.createObjectURL(file);
		},

		clearPhotoSelection() {
			this.selectedPhotoFile = null;

			if (this.photoPreviewUrl) {
				URL.revokeObjectURL(this.photoPreviewUrl);
				this.photoPreviewUrl = null;
			}

			// Clear the file input
			if (this.$refs.photoFileInput) {
				this.$refs.photoFileInput.value = '';
			}
		},

		async setGroupPhoto() {
			if (!this.selectedPhotoFile) {
				this.setGroupPhotoError = 'Please select a photo';
				return;
			}

			if (!this.selectedChatId) {
				this.setGroupPhotoError = 'No group selected';
				return;
			}

			try {
				this.setGroupPhotoLoading = true;
				this.setGroupPhotoError = null;

				// Create FormData for multipart/form-data request
				const formData = new FormData();
				formData.append('image', this.selectedPhotoFile);

				// Make the API request
				const response = await this.$axios.put(`/chats/${this.selectedChatId}/image`, formData, {
					headers: {
						'Content-Type': 'multipart/form-data'
					}
				});

				console.log('Group photo updated successfully:', response.data);

				// Clear the old cached image URL for this chat
				if (this.chatImageUrls[this.selectedChatId]) {
					URL.revokeObjectURL(this.chatImageUrls[this.selectedChatId]);
					delete this.chatImageUrls[this.selectedChatId];
				}

				// Clear image cache entries
				const chat = this.chats.find(c => c.id === this.selectedChatId);
				if (chat && chat.image && this.imageCache[chat.image]) {
					URL.revokeObjectURL(this.imageCache[chat.image]);
					delete this.imageCache[chat.image];
				}

				// Close modal and reset form
				this.showSetGroupPhotoModal = false;
				this.clearPhotoSelection();

				// Refresh chats to get the updated image path from server
				await this.getMyConversations();

				// Show success message
				alert('Group photo updated successfully!');

			} catch (err) {
				console.error('Failed to update group photo', err);

				// Handle specific error cases
				if (err.response?.status === 413) {
					this.setGroupPhotoError = 'File is too large. Please choose a smaller image.';
				} else if (err.response?.status === 400) {
					this.setGroupPhotoError = 'Invalid file format. Please choose a valid image.';
				} else {
					this.setGroupPhotoError = err.response?.data?.message || 'Failed to update group photo. Please try again.';
				}
			} finally {
				this.setGroupPhotoLoading = false;
			}
		},

		getLastMessagePreview(chat) {
			if (!chat.lastMsgText) {
				return 'No messages yet';
			}

			let preview = chat.lastMsgText;

			// For group chats, prepend the username
			if (chat.isGroup && chat.lastMsgUsername) {
				preview = `${chat.lastMsgUsername}: ${chat.lastMsgText}`;
			}

			// Truncate long messages
			const maxLength = 40;
			if (preview.length > maxLength) {
				return preview.substring(0, maxLength) + '...';
			}

			return preview;
		},

		isCurrentUserMessage(message) {
			// Compare using username since that's what you have in the message object
			return message.username === this.currentUsername;
		},

		formatMessageTime(timestamp) {
			if (!timestamp) return '';

			const date = new Date(timestamp);
			const now = new Date();

			// If today, show time only (HH:MM format)
			if (date.toDateString() === now.toDateString()) {
				return date.toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit',
					hour12: false // Use 24-hour format, change to true for 12-hour
				});
			}

			// If yesterday
			const yesterday = new Date(now);
			yesterday.setDate(now.getDate() - 1);
			if (date.toDateString() === yesterday.toDateString()) {
				return 'Yesterday ' + date.toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit',
					hour12: false
				});
			}

			// If this week (within 7 days)
			const weekAgo = new Date(now);
			weekAgo.setDate(now.getDate() - 7);
			if (date > weekAgo) {
				return date.toLocaleDateString([], { weekday: 'short' }) + ' ' +
					date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
			}

			// If this year, show month and day with time
			if (date.getFullYear() === now.getFullYear()) {
				return date.toLocaleDateString([], { month: 'short', day: 'numeric' }) + ' ' +
					date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
			}

			// Otherwise show full date with time
			return date.toLocaleDateString() + ' ' +
				date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
		},

		handleImageError(event) {
			// Get the chat ID from the image element or find it another way
			const imgElement = event.target;
			const chatId = imgElement.getAttribute('data-chat-id');

			if (chatId) {
				// Remove the failed image URL from cache
				delete this.chatImageUrls[chatId];

				// Also remove from image cache
				const chat = this.chats.find(c => c.id == chatId);
				if (chat && chat.image && this.imageCache[chat.image]) {
					URL.revokeObjectURL(this.imageCache[chat.image]);
					delete this.imageCache[chat.image];
				}
			}
		},

		startPolling() {
			this.stopPolling();

			if (this.isPolling) {
				return;
			}

			this.isPolling = true;
			this.pollingRetryCount = 0;
			this.currentPollingInterval = this.basePollingInterval;

			this.schedulePoll();
		},

		schedulePoll() {
			if (!this.isPolling) return;

			this.pollingInterval = setTimeout(async () => {
				if (!this.isPolling) return;

				try {
					// Only poll if we don't have too many active requests
					if (this.activeRequests.size < 2) {
						if (this.selectedChatId) {
							// Fetch messages and read status in parallel for better performance
							await Promise.all([
								this.getConversation(this.selectedChatId),
								this.getLastReadMessageId(this.selectedChatId)
							]);
						}

						// Also refresh chat list occasionally, but less frequently
						if (Math.random() < 0.1) { // 10% chance to refresh chats
							await this.getMyConversations();
						}
					}
				} catch (error) {
					console.error('Polling error:', error);
				} finally {
					// Schedule next poll
					this.schedulePoll();
				}
			}, this.currentPollingInterval);
		},

		stopPolling() {
			this.isPolling = false;

			if (this.pollingInterval) {
				clearTimeout(this.pollingInterval);
				this.pollingInterval = null;
			}

			// Cancel all active requests
			this.cancelAllRequests();
		},

		cancelAllRequests() {
			// If you're using a global axios instance, you might need to implement
			// request cancellation differently. This is a placeholder for that logic.
			console.log('Cancelling active requests:', this.activeRequests.size);
			this.activeRequests.clear();
		},

		async logout() {
			if (confirm('Are you sure you want to logout?')) {
				// Stop polling and clean up first
				this.cleanup();

				// Clear all localStorage data
				localStorage.removeItem('token');
				localStorage.removeItem('user');

				// Clear axios auth header
				delete this.$axios.defaults.headers.common['Authorization'];

				// Emit logout event to App.vue
				this.$emit('logout');
			}
		},

		cleanup() {
			this.stopPolling();
			this.cancelAllRequests();

			// Clean up blob URLs to prevent memory leaks
			Object.values(this.imageCache).forEach(blobUrl => {
				URL.revokeObjectURL(blobUrl);
			});

			// Clean up photo preview URLs
			if (this.photoPreviewUrl) {
				URL.revokeObjectURL(this.photoPreviewUrl);
			}

			if (this.profileImagePreviewUrl) {
				URL.revokeObjectURL(this.profileImagePreviewUrl);
			}

			if (this.currentUserImageUrl) {
				URL.revokeObjectURL(this.currentUserImageUrl);
			}

			this.imageCache = {};
			this.chatImageUrls = {};
			this.photoPreviewUrl = null;
			this.profileImagePreviewUrl = null;
			this.currentUserImageUrl = null;
		}
	}
}
</script>

<style scoped>
@import url('../assets/main.css');
@import "ChatsView.css";
</style>
