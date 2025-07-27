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
			basePollingInterval: 2000, // 5 seconds
			currentPollingInterval: 2000,

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

			// Image message support
			showImageModal: false,
			selectedMessageImage: null,
			messageImagePreviewUrl: null,
			tempSelectedImage: null,
			tempImagePreviewUrl: null,
			imageModalError: null,
			messageImageUrls: {}, // Store blob URLs for message images

			// Image viewer
			showImageViewer: false,
			imageViewerUrl: null,
			imageViewerTitle: null,

			// New scroll management properties
			lastMessageCount: 0,
			shouldScrollToBottom: false,
			lastScrollTop: 0,
			isUserAtBottom: true,
			hasNewMessages: false,
			newMessageCount: 0,
			scrollThreshold: 100, // pixels from bottom to consider "at bottom"
			preserveScrollPosition: false, // Flag to preserve scroll during polling
			savedScrollPosition: 0, // Store scroll position during updates

			// Flags for different scenarios
			isFirstChatLoad: false,
			shouldScrollAfterSend: false,

			// New Chat Image Support
			showNewChatImageModal: false,
			selectedNewChatImage: null,
			newChatImagePreviewUrl: null,
			tempNewChatImage: null,
			tempNewChatImagePreviewUrl: null,
			newChatImageModalError: null,

			showLikesDropdown: null,

			showDeleteDropdown: null,
			deletingMessage: null, // Track which message is being deleted

			showReplyDropdown: null,
			replyingToMessage: null,
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
		// Remove scroll listener
		const messagesContainer = this.$refs.messagesContainer;
		if (messagesContainer) {
			messagesContainer.removeEventListener('scroll', this.handleScroll);
		}
		// Enhanced cleanup
		this.cleanup();

		// Remove visibility change listener
		document.removeEventListener('visibilitychange', this.handleVisibilityChange);

		// Clean up blob URLs to prevent memory leaks
		Object.values(this.imageCache).forEach(blobUrl => {
			URL.revokeObjectURL(blobUrl);
		});
	},

	// Enhanced scroll monitoring method
	mounted() {
		this.$nextTick(() => {
			const messagesContainer = this.$refs.messagesContainer;
			if (messagesContainer) {
				messagesContainer.addEventListener('scroll', this.handleScroll);
			}
		});

		// Close dropdowns when clicking outside
		document.addEventListener('click', () => {
			this.showLikesDropdown = null;
			this.showDeleteDropdown = null; // Add this line
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

				// Sort chats by lastMsgTime (most recent first)
				this.chats = response.data.sort((a, b) => {
					const timeA = new Date(a.lastMsgTime || 0);
					const timeB = new Date(b.lastMsgTime || 0);
					return timeB - timeA; // Descending order (newest first)
				});

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

		// Replace your getConversation method with this fixed version:
		async getConversation(chatId, isFirstLoad = false) {
			const controller = new AbortController();
			const requestId = 'getConversation-' + chatId + '-' + Date.now();

			try {
				this.loadingMessages = true;
				this.messagesError = null;
				this.activeRequests.add(requestId);

				// Store scroll position ONLY for polling (not first load)
				let savedScrollTop = 0;
				let savedScrollHeight = 0;
				const messagesContainer = this.$refs.messagesContainer;

				if (messagesContainer && !isFirstLoad && !this.shouldScrollAfterSend) {
					savedScrollTop = messagesContainer.scrollTop;
					savedScrollHeight = messagesContainer.scrollHeight;

					// Update user position status
					const { scrollTop, scrollHeight, clientHeight } = messagesContainer;
					const distanceFromBottom = scrollHeight - (scrollTop + clientHeight);
					this.isUserAtBottom = distanceFromBottom <= this.scrollThreshold;
				}

				const response = await this.$axios.get(`/chats/${chatId}`, {
					signal: controller.signal,
					timeout: 10000
				});

				// Store previous message count to detect new messages
				const previousMessageCount = this.messages.length;
				const previousMessageIds = new Set(this.messages.map(m => m.id));

				// Sort messages by timestamp (oldest first, newest at bottom)
				const newMessages = response.data.sort((a, b) => {
					return new Date(a.createdAt) - new Date(b.createdAt);
				});

				// Detect new messages
				const hasNewMessages = newMessages.some(msg => !previousMessageIds.has(msg.id));
				const newMessageCount = newMessages.filter(msg => !previousMessageIds.has(msg.id)).length;

				// Update messages
				this.messages = newMessages;

				// Load images for image messages
				await this.loadMessageImages();

				// Load likes for all messages
				await this.getComments();

				// Fetch last read message ID
				await this.getLastReadMessageId(chatId);

				// Handle scrolling after DOM update
				this.$nextTick(() => {
					const container = this.$refs.messagesContainer;
					if (!container) return;

					if (isFirstLoad) {
						// Case 1: First chat opening - scroll to bottom
						console.log('First load - scrolling to bottom');
						this.scrollToBottom();

					} else if (this.shouldScrollAfterSend) {
						// Case 4: New message sent - scroll to bottom
						console.log('Message sent - scrolling to bottom');
						this.scrollToBottom();
						this.shouldScrollAfterSend = false;

					} else {
						// Case 2: Polling refresh - restore exact position
						console.log('Polling refresh - restoring position');
						const heightDifference = container.scrollHeight - savedScrollHeight;
						container.scrollTop = savedScrollTop + heightDifference;

						// Case 3: New messages received while user is mid-chat
						if (hasNewMessages && !this.isUserAtBottom) {
							console.log('New messages detected, showing indicator');
							this.hasNewMessages = true;
							this.newMessageCount += newMessageCount;
						}
					}
				});

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

		async getComments() {
			for (const message of this.messages) {
				if (!message.likes) {
					try {
						const response = await this.$axios.get(`/messages/${message.id}/comments`);
						message.likes = response.data || [];
					} catch (error) {
						console.error(`Failed to load likes for message ${message.id}:`, error);
						message.likes = [];
					}
				}
			}
		},

		async commentMessage(message) {
			try {
				const response = await this.$axios.put(`/messages/${message.id}/comments`);
				message.likes = response.data || [];
				this.$forceUpdate();
			} catch (error) {
				console.error('Failed to add like:', error);
			}
		},

		async deleteComment(message) {
			try {
				const response = await this.$axios.delete(`/messages/${message.id}/comments`);
				message.likes = response.data || [];
				this.$forceUpdate();
			} catch (error) {
				console.error('Failed to remove like:', error);
			}
		},

		async toggleMessageLike(message) {
			const isLiked = this.isMessageLikedByUser(message);

			if (isLiked) {
				await this.deleteComment(message);
			} else {
				await this.commentMessage(message);
			}
		},

		isMessageLikedByUser(message) {
			return message.likes && message.likes.includes(this.currentUsername);
		},

		toggleLikesDropdown(messageId) {
			this.showLikesDropdown = this.showLikesDropdown === messageId ? null : messageId;
			// Close other dropdowns when opening likes dropdown
			if (this.showLikesDropdown === messageId) {
				this.showReplyDropdown = null;
				this.showDeleteDropdown = null;
			}
		},

		// Add this new method to track user scroll position
		handleScroll() {
			const messagesContainer = this.$refs.messagesContainer;
			if (!messagesContainer) return;

			const { scrollTop, scrollHeight, clientHeight } = messagesContainer;

			// Update bottom detection
			const distanceFromBottom = scrollHeight - (scrollTop + clientHeight);
			this.isUserAtBottom = distanceFromBottom <= this.scrollThreshold;

			// Hide indicator if user scrolls to bottom
			if (this.isUserAtBottom && this.hasNewMessages) {
				this.hasNewMessages = false;
				this.newMessageCount = 0;
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

		// Reply functionality
		toggleReplyDropdown(messageId) {
			this.showReplyDropdown = this.showReplyDropdown === messageId ? null : messageId;
			// Close other dropdowns when opening reply dropdown
			if (this.showReplyDropdown === messageId) {
				this.showLikesDropdown = null;
				this.showDeleteDropdown = null;
			}
		},

		startReply(message) {
			this.replyingToMessage = message;
			this.showReplyDropdown = null; // Close dropdown

			// Focus the message input
			this.$nextTick(() => {
				if (this.$refs.messageInput) {
					this.$refs.messageInput.focus();
				}
			});
		},

		clearReply() {
			this.replyingToMessage = null;
		},

		getReplyMessage(messageId) {
			return this.messages.find(msg => msg.id === messageId);
		},

		getReplyPreviewText(message) {
			if (!message) return '';

			if (message.type === 'image') {
				return message.text ? message.text : '📷 Photo';
			} else if (message.type === 'gif') {
				return message.text ? message.text : '🎞️ GIF';
			}

			return message.text || '';
		},

		jumpToMessage(messageId) {
			console.log('Attempting to jump to message:', messageId);

			// Close any open dropdowns first
			this.showReplyDropdown = null;
			this.showDeleteDropdown = null;
			this.showLikesDropdown = null;

			this.$nextTick(() => {
				// Find the message element using the data attribute
				const messageElement = document.querySelector(`[data-message-id="${messageId}"]`);

				if (messageElement) {
					console.log('Found message element, scrolling to it');

					// Scroll to the message with smooth animation
					messageElement.scrollIntoView({
						behavior: 'smooth',
						block: 'center',
						inline: 'nearest'
					});

					// Add highlight effect to make the message stand out
					const messageBubble = messageElement.querySelector('.message-bubble');
					if (messageBubble) {
						// Add highlight class or inline style
						messageBubble.style.transition = 'background-color 0.3s ease';
						messageBubble.style.backgroundColor = 'rgba(0, 123, 255, 0.2)';
						messageBubble.style.borderRadius = '1rem';

						// Remove highlight after 3 seconds
						setTimeout(() => {
							messageBubble.style.backgroundColor = '';
							messageBubble.style.transition = '';
						}, 3000);
					}
				} else {
					console.warn('Message element not found for ID:', messageId);
					console.log('Available message elements:',
						Array.from(document.querySelectorAll('[data-message-id]'))
							.map(el => el.getAttribute('data-message-id'))
					);
				}
			});
		},

		// Updated sendMessage method with better preview handling
		async sendMessage() {
			if ((!this.newMessage.trim() && !this.selectedMessageImage) || !this.selectedChatId) return;

			try {
				this.sendingMessage = true;
				this.pendingMessage = this.newMessage.trim();

				// Set flag to scroll to bottom after sending
				this.shouldScrollAfterSend = true;

				let requestData;
				let requestConfig = {};
				const messageText = this.newMessage.trim();
				const messageTime = new Date().toISOString();
				let messageType = 'text';
				let previewText = messageText;

				if (this.selectedMessageImage) {
					const formData = new FormData();
					formData.append('image', this.selectedMessageImage);
					messageType = this.getFileType(this.selectedMessageImage);
					formData.append('type', messageType);

					if (messageText) {
						formData.append('text', messageText);
					}

					formData.append('isForward', 'false');

					// Add reply information
					if (this.replyingToMessage) {
						formData.append('replyTo', this.replyingToMessage.id.toString());
					} else {
						formData.append('replyTo', '0');
					}

					requestData = formData;

					if (messageText) {
						previewText = messageText;
					} else {
						previewText = messageType === 'gif' ? '🎞️ GIF' : '📷 Photo';
					}

					requestConfig = {};
				} else {
					requestData = {
						type: 'text',
						text: messageText
					};

					// Add reply information
					if (this.replyingToMessage) {
						requestData.replyTo = this.replyingToMessage.id;
					}

					requestConfig = {
						headers: {
							'Content-Type': 'application/json'
						}
					};
				}

				console.log('Sending message - will scroll to bottom after');

				const response = await this.$axios.post(
					`/chats/${this.selectedChatId}/messages`,
					requestData,
					requestConfig
				);

				// Clear input and reply
				this.newMessage = '';
				this.pendingMessage = '';
				this.clearReply();

				if (this.selectedMessageImage) {
					this.clearMessageImageSelection();
				}

				// Update chat preview
				this.updateChatPreview(this.selectedChatId, {
					lastMsgText: previewText,
					lastMsgTime: messageTime,
					lastMsgType: messageType,
					lastMsgUsername: this.currentUsername
				});

				// Refresh messages - shouldScrollAfterSend flag will make it scroll to bottom
				await this.getConversation(this.selectedChatId, false);

			} catch (err) {
				console.error('Failed to send message', err);
				this.shouldScrollAfterSend = false; // Reset flag on error

				if (err.response?.status === 400) {
					alert('Invalid message format. Please try again.');
				} else if (err.response?.status === 413) {
					alert('File is too large. Please choose a smaller image.');
				} else if (err.response?.status === 415) {
					alert('Unsupported file type. Please choose a valid image.');
				} else {
					alert('Failed to send message. Please try again.');
				}
			} finally {
				this.sendingMessage = false;
			}
		},

		// New Chat Image Modal Methods
		openNewChatImageModal() {
			this.showNewChatImageModal = true;
			this.tempNewChatImage = null;
			this.tempNewChatImagePreviewUrl = null;
			this.newChatImageModalError = null;
		},

		closeNewChatImageModal() {
			this.showNewChatImageModal = false;
			this.clearTempNewChatImageSelection();
		},

		handleNewChatImageSelect(event) {
			const file = event.target.files[0];

			if (!file) {
				this.clearTempNewChatImageSelection();
				return;
			}

			// Validate file type
			const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
			if (!validTypes.includes(file.type)) {
				this.newChatImageModalError = 'Please select a valid image file (JPG, PNG, GIF, WebP)';
				this.clearTempNewChatImageSelection();
				return;
			}

			// Validate file size (10MB limit)
			const maxSize = 10 * 1024 * 1024; // 10MB in bytes
			if (file.size > maxSize) {
				this.newChatImageModalError = 'File size must be less than 10MB';
				this.clearTempNewChatImageSelection();
				return;
			}

			this.tempNewChatImage = file;
			this.newChatImageModalError = null;

			// Create preview URL
			if (this.tempNewChatImagePreviewUrl) {
				URL.revokeObjectURL(this.tempNewChatImagePreviewUrl);
			}
			this.tempNewChatImagePreviewUrl = URL.createObjectURL(file);
		},

		clearTempNewChatImageSelection() {
			this.tempNewChatImage = null;

			if (this.tempNewChatImagePreviewUrl) {
				URL.revokeObjectURL(this.tempNewChatImagePreviewUrl);
				this.tempNewChatImagePreviewUrl = null;
			}

			// Clear the file input
			if (this.$refs.newChatImageInput) {
				this.$refs.newChatImageInput.value = '';
			}
		},

		selectNewChatImage() {
			if (!this.tempNewChatImage) return;

			this.selectedNewChatImage = this.tempNewChatImage;

			// Create preview URL for new chat
			if (this.newChatImagePreviewUrl) {
				URL.revokeObjectURL(this.newChatImagePreviewUrl);
			}
			this.newChatImagePreviewUrl = URL.createObjectURL(this.selectedNewChatImage);

			// Close modal
			this.closeNewChatImageModal();
		},

		clearNewChatImageSelection() {
			this.selectedNewChatImage = null;

			if (this.newChatImagePreviewUrl) {
				URL.revokeObjectURL(this.newChatImagePreviewUrl);
				this.newChatImagePreviewUrl = null;
			}
		},

		getFileType(file) {
			const extension = file.name.split('.').pop().toLowerCase();
			if (extension === 'gif') {
				return 'gif';
			}
			return 'image';
		},

		updateChatPreview(chatId, updates) {
			const chatIndex = this.chats.findIndex(chat => chat.id === chatId);
			if (chatIndex !== -1) {
				// Update the chat with new message info
				const chat = this.chats[chatIndex];

				// Ensure we're updating the correct fields
				if (updates.lastMsgText !== undefined) {
					chat.lastMsgText = updates.lastMsgText;
				}
				if (updates.lastMsgTime !== undefined) {
					chat.lastMsgTime = updates.lastMsgTime;
				}
				if (updates.lastMsgType !== undefined) {
					chat.lastMsgType = updates.lastMsgType;
				}
				if (updates.lastMsgUsername !== undefined) {
					chat.lastMsgUsername = updates.lastMsgUsername;
				}

				console.log('Updated chat preview:', {
					chatId,
					chatName: this.getChatName(chat),
					lastMsgText: chat.lastMsgText,
					lastMsgType: chat.lastMsgType
				});

				// Move the chat to the top of the list
				if (chatIndex > 0) {
					const updatedChat = this.chats.splice(chatIndex, 1)[0];
					this.chats.unshift(updatedChat);
				}

				// Force reactivity update for Vue 3
				this.$forceUpdate();
			}
		},

		// Toggle delete dropdown
		toggleDeleteDropdown(messageId) {
			this.showDeleteDropdown = this.showDeleteDropdown === messageId ? null : messageId;
			// Close other dropdowns when opening delete dropdown
			if (this.showDeleteDropdown === messageId) {
				this.showLikesDropdown = null;
				this.showReplyDropdown = null;
			}
		},

		handleDropdownMouseLeave(dropdownType) {
			// Close the dropdown when mouse leaves
			if (dropdownType === 'likes') {
				this.showLikesDropdown = null;
			} else if (dropdownType === 'reply') {
				this.showReplyDropdown = null;
			} else if (dropdownType === 'delete') {
				this.showDeleteDropdown = null;
			}
		},

		// Keep dropdown open when mouse is over it
		handleDropdownMouseEnter(dropdownType) {
			// This prevents the dropdown from closing when hovering over it
			// The dropdown stays open as long as mouse is over it
		},

		// Confirm message deletion
		confirmDeleteMessage(message) {
			// Close dropdown first
			this.showDeleteDropdown = null;

			// Show confirmation dialog
			if (confirm('Are you sure you want to delete this message? This action cannot be undone.')) {
				this.deleteMessage(message);
			}
		},

		async leaveGroup() {
			if (!confirm('Are you sure you want to leave this group?')) return;

			try {
				console.log('Leaving group:', this.selectedChatId);

				const response = await this.$axios.delete(`/chats/${this.selectedChatId}/members`);

				console.log('Leave group response:', response.data);

				// Check if chat was deleted (user was the last member)
				if (response.data && response.data.chatDeleted) {
					console.log('Chat was deleted, removing from chat list');

					// Remove chat from the chats list
					this.chats = this.chats.filter(chat => chat.id !== this.selectedChatId);

					// Clear selected chat and reset state
					this.selectedChatId = null;
					this.messages = [];
					this.lastReadMessageId = null;

					// Reset indicators
					this.hasNewMessages = false;
					this.newMessageCount = 0;

					// Show success message
					alert('You have left the group. The group has been deleted since you were the last member.');
				} else {
					// Group still exists, just refresh chats and clear selection
					await this.getMyConversations();
					this.selectedChatId = null;

					// Show success message
					alert('You have successfully left the group.');
				}

			} catch (err) {
				console.error('Failed to leave group', err);

				// Handle specific error cases
				if (err.response?.status === 403) {
					alert('You are not authorized to leave this group.');
				} else if (err.response?.status === 404) {
					alert('Group not found or you are not a member.');
				} else if (err.response?.status === 401) {
					console.log('Authentication error during group leave');
					this.$emit('logout');
				} else {
					alert('Failed to leave group. Please try again.');
				}
			}
		},

		// Delete message function
		async deleteMessage(message) {
			if (!message || !message.id) {
				console.error('Invalid message for deletion');
				return;
			}

			try {
				// Set loading state
				this.deletingMessage = message.id;

				console.log('Deleting message:', message.id);

				// Make DELETE request to the backend
				const response = await this.$axios.delete(`/messages/${message.id}`);

				console.log('Delete response:', response.data);

				// Check if chat was deleted
				if (response.data && response.data.chatDeleted) {
					console.log('Chat was deleted, removing from chat list');

					// Remove chat from the chats list
					this.chats = this.chats.filter(chat => chat.id !== this.selectedChatId);

					// Clear selected chat
					this.selectedChatId = null;
					this.messages = [];
					this.lastReadMessageId = null;

					// Reset indicators
					this.hasNewMessages = false;
					this.newMessageCount = 0;

					// Show success message
					alert('Message deleted. The conversation has been removed since it was the last message.');
				} else {
					// Just remove the message from the current messages list
					const messageIndex = this.messages.findIndex(msg => msg.id === message.id);
					if (messageIndex !== -1) {
						this.messages.splice(messageIndex, 1);
					}

					// Update chat preview if this was the last message
					const chat = this.chats.find(c => c.id === this.selectedChatId);
					if (chat && this.messages.length > 0) {
						// Find the new last message
						const lastMessage = this.messages[this.messages.length - 1];
						this.updateChatPreview(this.selectedChatId, {
							lastMsgText: lastMessage.text || (lastMessage.type === 'image' ? '📷 Photo' : '🎞️ GIF'),
							lastMsgTime: lastMessage.createdAt,
							lastMsgType: lastMessage.type,
							lastMsgUsername: lastMessage.username
						});
					}

					// Force reactivity update
					this.$forceUpdate();
				}

			} catch (error) {
				console.error('Failed to delete message:', error);

				// Handle specific error cases
				if (error.response?.status === 403) {
					alert('You can only delete your own messages.');
				} else if (error.response?.status === 404) {
					alert('Message not found or already deleted.');
				} else if (error.response?.status === 401) {
					console.log('Authentication error during message deletion');
					this.$emit('logout');
				} else {
					alert('Failed to delete message. Please try again.');
				}
			} finally {
				// Clear loading state
				this.deletingMessage = null;
			}
		},

		// Enhanced scrollToBottom method
		scrollToBottom() {
			this.$nextTick(() => {
				const messagesContainer = this.$refs.messagesContainer;
				if (messagesContainer) {
					messagesContainer.scrollTop = messagesContainer.scrollHeight;
					this.isUserAtBottom = true;
					this.hasNewMessages = false;
					this.newMessageCount = 0;
				}
			});
		},

		// Method to handle clicking the new message indicator
		scrollToNewMessages() {
			this.scrollToBottom();
			this.hasNewMessages = false;
			this.newMessageCount = 0;
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

				let requestData;
				let requestConfig = {};

				if (this.selectedNewChatImage) {
					// Create FormData for file upload
					const formData = new FormData();
					formData.append('image', this.selectedNewChatImage);
					formData.append('type', this.getFileType(this.selectedNewChatImage));
					formData.append('receivers', JSON.stringify(this.selectedUsers.map(user => user.username)));

					if (this.initialMessage.trim()) {
						formData.append('text', this.initialMessage.trim());
					}

					if (this.selectedUsers.length > 1 && this.newChatName.trim()) {
						formData.append('chatName', this.newChatName.trim());
					}

					formData.append('isForward', 'false');

					requestData = formData;
					requestConfig = {};
				} else {
					// JSON request for text message
					requestData = {
						type: 'text',
						receivers: this.selectedUsers.map(user => user.username)
					};

					if (this.initialMessage.trim()) {
						requestData.text = this.initialMessage.trim();
					}

					if (this.selectedUsers.length > 1) {
						requestData.chatName = this.newChatName.trim();
					}

					requestConfig = {
						headers: {
							'Content-Type': 'application/json'
						}
					};
				}

				console.log('Creating chat with media support...');

				// Create the chat
				const response = await this.$axios.post('/chats', requestData, requestConfig);

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

				if (err.response?.status === 413) {
					this.newChatError = 'File is too large. Please choose a smaller image.';
				} else if (err.response?.status === 415) {
					this.newChatError = 'Unsupported file type. Please choose a valid image.';
				} else {
					this.newChatError = err.response?.data?.message || 'Failed to create conversation. Please try again.';
				}
			} finally {
				this.newChatLoading = false;
			}
		},

		// Updated closeNewChatModal method
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

			// Clear media selection
			this.clearNewChatImageSelection();
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

		// Enhanced selectChat method
		selectChat(chatId) {
			console.log('Selecting chat:', chatId);

			// Clear unread count
			const chatIndex = this.chats.findIndex(chat => chat.id === chatId);
			if (chatIndex !== -1 && this.chats[chatIndex].unread > 0) {
				this.chats[chatIndex].unread = 0;
				this.$forceUpdate();
			}

			this.selectedChatId = chatId;
			this.lastReadMessageId = null;

			// Reset indicators
			this.hasNewMessages = false;
			this.newMessageCount = 0;
			this.isUserAtBottom = true;
			this.shouldScrollAfterSend = false;

			// Load chat for first time - will scroll to bottom
			this.getConversation(chatId, true);
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
				const baseURL = this.$axios.defaults.baseURL;

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

		// Format file size for display
		formatFileSize(bytes) {
			if (bytes === 0) return '0 Bytes';
			const k = 1024;
			const sizes = ['Bytes', 'KB', 'MB', 'GB'];
			const i = Math.floor(Math.log(bytes) / Math.log(k));
			return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
		},

		// Open image selection modal
		openImageModal() {
			this.showImageModal = true;
			this.tempSelectedImage = null;
			this.tempImagePreviewUrl = null;
			this.imageModalError = null;
		},

		// Close image selection modal
		closeImageModal() {
			this.showImageModal = false;
			this.clearTempImageSelection();
		},

		// Handle image file selection in modal
		handleMessageImageSelect(event) {
			const file = event.target.files[0];

			if (!file) {
				this.clearTempImageSelection();
				return;
			}

			// Validate file type
			const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
			if (!validTypes.includes(file.type)) {
				this.imageModalError = 'Please select a valid image file (JPG, PNG, GIF, WebP)';
				this.clearTempImageSelection();
				return;
			}

			// Validate file size (10MB limit)
			const maxSize = 10 * 1024 * 1024; // 10MB in bytes
			if (file.size > maxSize) {
				this.imageModalError = 'File size must be less than 10MB';
				this.clearTempImageSelection();
				return;
			}

			this.tempSelectedImage = file;
			this.imageModalError = null;

			// Create preview URL
			if (this.tempImagePreviewUrl) {
				URL.revokeObjectURL(this.tempImagePreviewUrl);
			}
			this.tempImagePreviewUrl = URL.createObjectURL(file);
		},

		// Clear temporary image selection in modal
		clearTempImageSelection() {
			this.tempSelectedImage = null;

			if (this.tempImagePreviewUrl) {
				URL.revokeObjectURL(this.tempImagePreviewUrl);
				this.tempImagePreviewUrl = null;
			}

			// Clear the file input
			if (this.$refs.messageImageInput) {
				this.$refs.messageImageInput.value = '';
			}
		},

		// Select image for message
		selectMessageImage() {
			if (!this.tempSelectedImage) return;

			this.selectedMessageImage = this.tempSelectedImage;

			// Create preview URL for message input
			if (this.messageImagePreviewUrl) {
				URL.revokeObjectURL(this.messageImagePreviewUrl);
			}
			this.messageImagePreviewUrl = URL.createObjectURL(this.selectedMessageImage);

			// Close modal
			this.closeImageModal();
		},

		// Clear selected message image
		clearMessageImageSelection() {
			this.selectedMessageImage = null;

			if (this.messageImagePreviewUrl) {
				URL.revokeObjectURL(this.messageImagePreviewUrl);
				this.messageImagePreviewUrl = null;
			}
		},

		// Updated getMessageImageUrl method
		async getMessageImageUrl(mediaUrl) {
			if (!mediaUrl) return null;

			// Check cache first using message ID or mediaUrl as key
			const cacheKey = mediaUrl;
			if (this.messageImageUrls[cacheKey]) {
				return this.messageImageUrls[cacheKey];
			}

			try {
				const baseURL = this.$axios.defaults.baseURL;
				const imageUrl = `${baseURL}/uploads/messages/${mediaUrl}`;
				const token = localStorage.getItem('token');

				console.log('Fetching message image from:', imageUrl);

				const response = await this.$axios.get(imageUrl, {
					responseType: 'blob',
					headers: {
						'Authorization': token
					}
				});

				const blobUrl = URL.createObjectURL(response.data);

				// Vue 3: Direct assignment to reactive object
				this.messageImageUrls[cacheKey] = blobUrl;

				return blobUrl;

			} catch (error) {
				console.error('Failed to fetch message image:', error);
				return null;
			}
		},

		// Fixed loadMessageImages method for Vue 3
		async loadMessageImages() {
			const imageMessages = this.messages.filter(msg =>
				(msg.type === 'image' || msg.type === 'gif') && msg.mediaUrl
			);

			for (const message of imageMessages) {
				if (!this.messageImageUrls[message.id] && !message.imageLoading) {
					// Set loading state - Vue 3 way
					message.imageLoading = true;

					try {
						const imageUrl = await this.getMessageImageUrl(message.mediaUrl);
						if (imageUrl) {
							// Vue 3: Direct assignment works due to Proxy reactivity
							this.messageImageUrls[message.id] = imageUrl;
						} else {
							message.imageError = true;
						}
					} catch (error) {
						console.error(`Failed to load image for message ${message.id}:`, error);
						message.imageError = true;
					} finally {
						message.imageLoading = false;
					}
				}
			}
		},

		// Fixed handleMessageImageError method for Vue 3
		handleMessageImageError(message) {
			console.error('Image load error for message:', message.id);

			// Vue 3: Direct assignment
			message.imageError = true;

			// Remove from cache
			if (this.messageImageUrls[message.id]) {
				URL.revokeObjectURL(this.messageImageUrls[message.id]);
				delete this.messageImageUrls[message.id];
			}
		},


		// Open image viewer
		openImageViewer(imageUrl, title) {
			this.imageViewerUrl = imageUrl;
			this.imageViewerTitle = title;
			this.showImageViewer = true;
		},

		// Close image viewer
		closeImageViewer() {
			this.showImageViewer = false;
			this.imageViewerUrl = null;
			this.imageViewerTitle = null;
		},

		// Updated getLastMessagePreview method to handle images better
		getLastMessagePreview(chat) {
			if (!chat.lastMsgText) {
				return 'No messages yet';
			}

			let preview = chat.lastMsgText;

			// Handle different message types
			if (chat.lastMsgType === 'image') {
				// If lastMsgText is a caption, show it; otherwise show image indicator
				if (preview && !preview.includes('📷')) {
					// It's a caption, use as is
				} else {
					preview = '📷 Photo';
				}
			} else if (chat.lastMsgType === 'gif') {
				// If lastMsgText is a caption, show it; otherwise show GIF indicator
				if (preview && !preview.includes('🎞️')) {
					// It's a caption, use as is
				} else {
					preview = '🎞️ GIF';
				}
			}

			// For group chats, prepend the username (but not for media indicators)
			if (chat.isGroup && chat.lastMsgUsername) {
				if (preview === '📷 Photo' || preview === '🎞️ GIF') {
					preview = `${chat.lastMsgUsername}: ${preview}`;
				} else if (preview && preview.length > 0) {
					preview = `${chat.lastMsgUsername}: ${preview}`;
				}
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

		// Update your polling method to not force scroll
		schedulePoll() {
			if (!this.isPolling) return;

			this.pollingInterval = setTimeout(async () => {
				if (!this.isPolling) return;

				try {
					if (this.activeRequests.size < 2) {
						if (this.selectedChatId) {
							console.log('Polling - preserving scroll position');

							// Polling should NEVER be treated as first load
							await Promise.all([
								this.getConversation(this.selectedChatId, false), // false = preserve position
								this.getLastReadMessageId(this.selectedChatId)
							]);
						}

						if (Math.random() < 0.1) {
							await this.getMyConversations();
						}
					}
				} catch (error) {
					console.error('Polling error:', error);
				} finally {
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

			Object.values(this.messageImageUrls).forEach(blobUrl => {
				URL.revokeObjectURL(blobUrl);
			});

			if (this.photoPreviewUrl) {
				URL.revokeObjectURL(this.photoPreviewUrl);
			}

			if (this.profileImagePreviewUrl) {
				URL.revokeObjectURL(this.profileImagePreviewUrl);
			}

			if (this.currentUserImageUrl) {
				URL.revokeObjectURL(this.currentUserImageUrl);
			}

			if (this.newChatImagePreviewUrl) {
				URL.revokeObjectURL(this.newChatImagePreviewUrl);
			}

			if (this.tempNewChatImagePreviewUrl) {
				URL.revokeObjectURL(this.tempNewChatImagePreviewUrl);
			}

			this.imageCache = {};
			this.chatImageUrls = {};
			this.messageImageUrls = {};
			this.photoPreviewUrl = null;
			this.profileImagePreviewUrl = null;
			this.currentUserImageUrl = null;
			this.newChatImagePreviewUrl = null;
			this.tempNewChatImagePreviewUrl = null;
		}
	}
}
</script>

<style scoped>
@import url('../assets/main.css');
@import "ChatsView.css";
</style>
