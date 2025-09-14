<template src="./ChatsView.html"></template>

<script>
import Sidebar from '../components/Sidebar.vue'
import ConversationSection from '../components/ConversationSection.vue'

import ProfileModal from '../components/ProfileModal.vue'
import CreateChatModal from '../components/CreateChatModal.vue'
import NewChatImageModal from '../components/NewChatImageModal.vue'
import ForwardMessageModal from '../components/ForwardMessageModal.vue'
import RenameGroupModal from '../components/RenameGroupModal.vue'
import ImageSelectionModal from '../components/ImageSelectionModal.vue'
import AddToGroupModal from '../components/AddToGroupModal.vue'
import SetGroupPhotoModal from '../components/SetGroupPhotoModal.vue'



export default {
	name: 'ChatsView',
	components: {
		Sidebar,
		ConversationSection,

		ProfileModal,
		CreateChatModal,
		NewChatImageModal,
		ForwardMessageModal,
		RenameGroupModal,
		ImageSelectionModal,
		AddToGroupModal,
		SetGroupPhotoModal,
	},

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
			pendingMessage: '',

			// Read status tracking
			lastReadMessageId: null,

			// User info
			currentUserId: null,
			currentUsername: null,

			// Rename Group Modal
			showRenameGroupModal: false,
			setGroupNameLoading: false,
			setGroupNameError: null,

			// Add member modal
			showAddToGroupModal: false,
			addToGroupLoading: false,
			addToGroupError: null,

			// Set Group Photo Modal
			showSetGroupPhotoModal: false,
			setGroupPhotoLoading: false,
			setGroupPhotoError: null,

			// Profile modal
			showProfileModal: false,
			profileLoading: false,
			profileError: null,

			chatImageUrls: {},
			imageCache: {},

			// Polling improvements
			pollingInterval: null,
			isPolling: false,
			pollingRetryCount: 0,
			maxRetries: 3,
			basePollingInterval: 5000,
			currentPollingInterval: 5000,

			// Request tracking
			activeRequests: new Set(),
			requestController: null,

			// New Chat Modal
			showNewChatModal: false,
			newChatLoading: false,
			newChatError: null,
			newChatName: '',
			initialMessage: '',

			// User Selection
			users: [],
			filteredUsers: [],
			selectedUsers: [],
			loadingUsers: false,
			userSearchQuery: '',

			// User profile image
			currentUserImageUrl: null,
			selectedProfileImage: null,
			profileImagePreviewUrl: null,

			// Image message support
			showImageModal: false,
			selectedMessageImage: null,
			messageImagePreviewUrl: null,
			messageImageUrls: {},


			// Scroll management properties
			lastMessageCount: 0,
			shouldScrollToBottom: false,
			lastScrollTop: 0,
			isUserAtBottom: true,
			hasNewMessages: false,
			newMessageCount: 0,
			scrollThreshold: 100,
			preserveScrollPosition: false,
			savedScrollPosition: 0,

			// Flags for different scenarios
			isFirstChatLoad: false,
			shouldScrollAfterSend: false,

			// New Chat Image Support
			showNewChatImageModal: false,
			selectedNewChatImage: null,
			newChatImagePreviewUrl: null,

			// Like dropdown
			showLikesDropdown: null,

			// Delete dropdown
			showDeleteDropdown: null,
			deletingMessage: null,

			// Reply dropdown
			showReplyDropdown: null,
			replyingToMessage: null,

			// Forward modal states
			showForwardModal: false,
			forwardingMessage: null,

			// Forward users
			forwardUserSearchQuery: '',

			// Forward dropdown
			showForwardDropdown: null,

			// Chat header info dropdown
			showChatInfoDropdown: false,
			chatMembers: [],
			loadingChatMembers: false,
			memberImageUrls: {},
			chatInfoDropdownTimeout: null,

			// User profile images for modal
			userImageUrls: {},
		}
	},

	computed: {
		selectedChat() {
			return this.chats.find(chat => chat.id === this.selectedChatId);
		},
		forwardGroupChats() {
			return this.chats.filter(chat => chat.isGroup);
		},
		canCreateChat() {
			const hasRecipients = this.selectedUsers.length > 0;
			const hasGroupNameIfNeeded = this.selectedUsers.length <= 1 || (this.selectedUsers.length > 1 && !!this.newChatName.trim());
			const hasInitialContent = !!(this.initialMessage.trim() || this.selectedNewChatImage);

			return hasRecipients && hasGroupNameIfNeeded && hasInitialContent && !this.newChatLoading;
		}
	},

	async created() {
		document.addEventListener('visibilitychange', this.handleVisibilityChange);
		const token = localStorage.getItem('token');
		if (!token) {
			console.log('No token found, emitting logout');
			this.$emit('logout');
			return;
		}

		console.log('Token found, setting up ChatsView');

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

		this.$axios.defaults.headers.common['Authorization'] = `${token}`;
		console.log('Set Authorization header');

		await this.getMyPhoto();

		await this.getMyConversations();

		this.startPolling();
	},
	beforeUnmount() {
		const messagesContainer = this.$refs.messagesContainer;
		if (messagesContainer) {
			messagesContainer.removeEventListener('scroll', this.handleScroll);
		}
		// Cleanup
		this.cleanup();

		// Remove visibility change listener
		document.removeEventListener('visibilitychange', this.handleVisibilityChange);

		// Clean up blob URLs to prevent memory leaks
		Object.values(this.imageCache).forEach(blobUrl => {
			URL.revokeObjectURL(blobUrl);
		});
	},

	mounted() {
		this.$nextTick(() => {
			const messagesContainer = this.$refs.messagesContainer;
			if (messagesContainer) {
				messagesContainer.addEventListener('scroll', this.handleScroll);
			}

			// Focus input if there's already a selected chat (page refresh case)
			if (this.selectedChatId) {
				setTimeout(() => {
					this.focusMessageInput();
				}, 500);
			}
		});

		// Close dropdowns when clicking outside
		document.addEventListener('click', () => {
			this.showLikesDropdown = null;
			this.showDeleteDropdown = null;
			this.showReplyDropdown = null;
			this.showForwardDropdown = null;
			this.showChatInfoDropdown = false;
		});

		document.addEventListener('click', (event) => {
			const chatContainer = document.querySelector('.chat-container');
			const isClickInChatArea = chatContainer && chatContainer.contains(event.target);
			const isInteractiveElement = event.target.closest('button, input, textarea, select, a, [role="button"]');

			if (isClickInChatArea && !isInteractiveElement && this.selectedChatId) {
				this.focusMessageInput();
			}
		});
	},


	methods: {
		async updateProfile(formData) {
			const { username, profileImage } = formData

			try {
				this.profileLoading = true
				this.profileError = null

				// Update username if provided
				if (username) {
					await this.setMyUserName(username)
				}

				// Update profile image if provided
				if (profileImage) {
					await this.setMyPhoto(profileImage)
				}

				// Close modal and show success
				this.showProfileModal = false
				alert('Profile updated successfully!')

			} catch (err) {
				console.error('Failed to update profile', err)
				console.error('Error response:', err.response)

				// Check for authentication errors
				if (err.response?.status === 401) {
					console.log('Authentication error during profile update, logging out')
					this.$emit('logout')
					return
				}

				if (err.response?.status === 413) {
					this.profileError = 'File is too large. Please choose a smaller image.'
				} else if (err.response?.status === 400) {
					this.profileError = 'Invalid file format. Please choose a valid image.'
				} else if (err.response?.status === 409) {
					this.profileError = 'Username already exists. Please choose a different username.'
				} else {
					this.profileError = err.response?.data?.message || 'Failed to update profile. Please try again.'
				}
			} finally {
				this.profileLoading = false
			}
		},

		async setMyUserName(username) {
			console.log('Updating username to:', username)

			const response = await this.$axios.put('/users/username', {
				username: username,
			}, {
				headers: {
					'Content-Type': 'application/json'
				}
			})

			console.log('Username update response:', response.data)

			// Update current user info from server response
			if (response.data && response.data.username) {
				this.currentUsername = response.data.username

				const userData = {
					id: response.data.id || this.currentUserId,
					username: this.currentUsername,
				}
				localStorage.setItem('user', JSON.stringify(userData))
			}
		},

		async setMyPhoto(imageFile) {
			console.log('Updating profile image...')

			const formData = new FormData()
			formData.append('image', imageFile)

			const response = await this.$axios.put('/users/image', formData, {
				headers: {
					'Content-Type': 'multipart/form-data'
				}
			})

			console.log('Image update response:', response.data)

			if (response.data && response.data.imageUrl) {
				console.log('Loading new image from:', response.data.imageUrl)
				this.currentUserImageUrl = await this.getImage(response.data.imageUrl)
			} else {
				await this.getMyPhoto()
			}
		},

		closeProfileModal() {
			this.showProfileModal = false
			this.profileError = ''
		},

		async getMyConversations() {
			const controller = new AbortController();
			const requestId = 'getMyConversations-' + Date.now();

			try {
				this.loading = true;
				this.error = null;

				this.activeRequests.add(requestId);

				console.log('Fetching chats with token:', localStorage.getItem('token') ? 'present' : 'missing');

				const response = await this.$axios.get('/chats', {
					signal: controller.signal,
					timeout: 10000
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

		async getConversation(chatId, isFirstLoad = false) {
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

				// Sort messages by timestamp (oldest first)
				const newMessages = response.data.sort((a, b) => {
					return new Date(a.createdAt) - new Date(b.createdAt);
				});

				// Update messages
				this.messages = newMessages;

				// Load images for image messages
				await this.loadMessageImages();

				// Load likes for all messages
				await this.getComments();

				// Fetch last read message ID
				await this.getLastRead(chatId);

				// Focus message input after loading (only for first load)
				if (isFirstLoad) {
					setTimeout(() => {
						this.focusMessageInput();
					}, 100);
				}

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

		async createChat() {
			if (!this.canCreateChat) {
				console.warn('Attempted to create chat when conditions not met');
				return;
			}

			try {
				this.newChatLoading = true;
				this.newChatError = null;

				let requestData;
				let requestConfig = {};

				if (this.selectedNewChatImage) {
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
					requestData = {
						type: 'text',
						receivers: this.selectedUsers.map(user => user.username),
						text: this.initialMessage.trim()
					};

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

				await this.getMyConversations();

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

		async sendMessage() {
			if ((!this.newMessage.trim() && !this.selectedMessageImage) || !this.selectedChatId) return;

			try {
				this.sendingMessage = true;
				this.pendingMessage = this.newMessage.trim();

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

					if (this.replyingToMessage) {
						requestData.replyTo = this.replyingToMessage.id;
					}

					requestConfig = {
						headers: {
							'Content-Type': 'application/json'
						}
					};
				}

				console.log('Sending message...');

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

				this.updateChatPreview(this.selectedChatId, {
					lastMsgText: previewText,
					lastMsgTime: messageTime,
					lastMsgType: messageType,
					lastMsgUsername: this.currentUsername
				});

				await this.getConversation(this.selectedChatId, false);

				setTimeout(() => {
					this.focusMessageInput();
				}, 100);

			} catch (err) {
				console.error('Failed to send message', err);

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

		async setGroupName(newName) {
			if (!newName.trim()) {
				return;
			}

			if (!this.selectedChatId) {
				return;
			}

			try {
				this.setGroupNameLoading = true;
				this.setGroupNameError = null;

				const response = await this.$axios.put(`/chats/${this.selectedChatId}/chat-name`, {
					chatName: newName.trim()
				});

				console.log('Group renamed successfully:', response.data);

				const chatIndex = this.chats.findIndex(chat => chat.id === this.selectedChatId);
				if (chatIndex !== -1) {
					this.chats[chatIndex].name = newName.trim();
				}

				// Close modal and reset form
				this.showRenameGroupModal = false;

				// Refresh chats to get latest data from server
				await this.getMyConversations();

				// Show success message
				alert('Group renamed successfully!');

			} catch (err) {
				console.error('Failed to rename group', err);
				// Error will be handled by the modal's error prop binding
			} finally {
				this.setGroupNameLoading = false;
			}
		},

		async setGroupPhoto(photoFile) {
			if (!photoFile) {
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
				formData.append('image', photoFile);

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
				this.closeSetGroupPhotoModal();

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

		async addToGroup(username) {
			// Add proper parameter validation
			if (!username || typeof username !== 'string' || !username.trim()) {
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
					username: username.trim()
				});

				console.log('Member added successfully:', response.data);

				// Close modal and reset form
				this.closeAddToGroupModal();

				// Refresh chats to get latest data
				await this.getMyConversations();

				// Show success message
				alert(`${username.trim()} has been added to the group!`);

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

		async leaveGroup() {
			if (!confirm('Are you sure you want to leave this group?')) return;

			// Store the current chat ID before we start the async operation
			const leavingChatId = this.selectedChatId;

			if (!leavingChatId) {
				console.error('No chat selected to leave');
				return;
			}

			try {
				console.log('Leaving group:', leavingChatId);

				const response = await this.$axios.delete(`/chats/${leavingChatId}/members`);

				console.log('Leave group response:', response.data);

				this.selectedChatId = null;
				this.messages = [];
				this.lastReadMessageId = null;
				this.hasNewMessages = false;
				this.newMessageCount = 0;
				this.chatMembers = [];

				this.chats = this.chats.filter(chat => chat.id !== leavingChatId);

				await this.getMyConversations();

				// Show success message
				alert('You have successfully left the group.');
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

		async deleteMessage(message) {
			if (!message || !message.id) {
				console.error('Invalid message for deletion');
				return;
			}

			try {
				this.deletingMessage = message.id;

				console.log('Deleting message:', message.id);

				const response = await this.$axios.delete(`/messages/${message.id}`);

				console.log('Delete response:', response.data);

				if (response.data && response.data.chatDeleted) {
					console.log('Chat was deleted, removing from chat list');

					this.chats = this.chats.filter(chat => chat.id !== this.selectedChatId);

					this.selectedChatId = null;
					this.messages = [];
					this.lastReadMessageId = null;

					this.hasNewMessages = false;
					this.newMessageCount = 0;

					alert('Message deleted. The conversation has been removed since it was the last message.');
				} else {
					const messageIndex = this.messages.findIndex(msg => msg.id === message.id);
					if (messageIndex !== -1) {
						this.messages.splice(messageIndex, 1);
					}

					const chat = this.chats.find(c => c.id === this.selectedChatId);
					if (chat && this.messages.length > 0) {
						// find the new most recent message
						const lastMessage = this.messages[this.messages.length - 1];
						this.updateChatPreview(this.selectedChatId, {
							lastMsgText: lastMessage.text || (lastMessage.type === 'image' ? '📷 Photo' : '🎞️ GIF'),
							lastMsgTime: lastMessage.createdAt,
							lastMsgType: lastMessage.type,
							lastMsgUsername: lastMessage.username
						});
					}
				}

			} catch (error) {
				console.error('Failed to delete message:', error);

				// Handle specific error cases
				if (error.response?.status === 404) {
					alert('Message not found.');
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

		async getUsers() {
			try {
				this.loadingUsers = true;
				this.newChatError = null;

				const response = await this.$axios.get('/users');
				this.users = response.data;
				this.filteredUsers = response.data;

				console.log('Users fetched successfully:', this.users);

				// Load user profile images
				await this.loadUserImages();

			} catch (err) {
				console.error('Failed to fetch users', err);
				this.newChatError = 'Failed to load users. Please try again.';
			} finally {
				this.loadingUsers = false;
			}
		},

		async getGroupMembers() {
			if (!this.selectedChatId) return;

			try {
				this.loadingChatMembers = true;

				const response = await this.$axios.get(`/chats/${this.selectedChatId}/members`);
				this.chatMembers = response.data || [];

				console.log('Chat members loaded:', this.chatMembers);

				await this.loadMemberImages();

			} catch (error) {
				console.error('Failed to load chat members:', error);
				this.chatMembers = [];

				if (error.response?.status === 401) {
					this.$emit('logout');
				}
			} finally {
				this.loadingChatMembers = false;
			}
		},

		async loadMemberImages() {
			for (const member of this.chatMembers) {
				if (member.imageUrl && !this.memberImageUrls[member.id]) {
					try {
						const imageUrl = await this.getImage(member.imageUrl);
						if (imageUrl) {
							this.memberImageUrls[member.id] = imageUrl;
						}
					} catch (error) {
						console.error(`Failed to load image for member ${member.id}:`, error);
					}
				}
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
			} catch (error) {
				console.error('Failed to add like:', error);
			}
		},

		async uncommentMessage(message) {
			try {
				const response = await this.$axios.delete(`/messages/${message.id}/comments`);
				message.likes = response.data || [];
			} catch (error) {
				console.error('Failed to remove like:', error);
			}
		},

		async getImage(imagePath) {
			if (!imagePath) return null;

			if (this.imageCache[imagePath]) {
				return this.imageCache[imagePath];
			}

			try {
				const pathParts = imagePath.split('/');

				if (pathParts.length < 4 || pathParts[0] !== 'uploads' || pathParts[2] !== 'images') {
					console.error('Invalid image path format:', imagePath);
					return null;
				}

				const folder = pathParts[1];
				const filename = pathParts[3];

				const validFolders = ['user', 'chats', 'messages'];
				if (!validFolders.includes(folder)) {
					console.error('Invalid folder in image path:', folder);
					return null;
				}

				const imageUrl = `/uploads/${folder}/images/${filename}`;

				console.log('Fetching image from:', imageUrl);

				const token = localStorage.getItem('token');

				const response = await this.$axios.get(imageUrl, {
					responseType: 'blob',
					headers: {
						'Authorization': token
					}
				});

				const blobUrl = URL.createObjectURL(response.data);

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
						const imageUrl = await this.getImage(chat.image);
						if (imageUrl) {
							this.chatImageUrls[chat.id] = imageUrl;
						}
					} catch (error) {
						console.error(`Failed to load image for chat ${chat.id}:`, error);
					}
				}
			}
		},

		async loadUserImages() {
			for (const user of this.users) {
				if (user.imageUrl && !this.userImageUrls[user.id]) {
					try {
						const imageUrl = await this.getImage(user.imageUrl);
						if (imageUrl) {
							this.userImageUrls[user.id] = imageUrl;
						}
					} catch (error) {
						console.error(`Failed to load image for user ${user.id}:`, error);
					}
				}
			}
		},

		async getMyPhoto() {
			try {
				console.log('Loading current user image...');

				const response = await this.$axios.get('/users/image');

				if (response.data && response.data.imageUrl) {
					this.currentUserImageUrl = await this.getImage(response.data.imageUrl);
				}
			} catch (error) {
				console.error('Failed to load user image:', error);

				if (error.response?.status === 404) {
					console.log('User has no profile image set');
				} else if (error.response?.status === 401) {
					console.log('Authentication error loading user image');
				}

				this.currentUserImageUrl = null;
			}
		},

		handleUsersImageError(userId) {
			if (this.userImageUrls[userId]) {
				URL.revokeObjectURL(this.userImageUrls[userId]);
				delete this.userImageUrls[userId];
			}
		},

		handleMemberImageError(memberId) {
			if (this.memberImageUrls[memberId]) {
				URL.revokeObjectURL(this.memberImageUrls[memberId]);
				delete this.memberImageUrls[memberId];
			}
		},

		async toggleMessageLike(message) {
			const isLiked = this.isMessageLikedByUser(message);

			if (isLiked) {
				await this.uncommentMessage(message);
			} else {
				await this.commentMessage(message);
			}
		},

		isMessageLikedByUser(message) {
			return message.likes && message.likes.includes(this.currentUsername);
		},

		handleScroll() {
			const messagesContainer = this.$refs.messagesContainer;
			if (!messagesContainer) return;

			const { scrollTop, scrollHeight, clientHeight } = messagesContainer;

			// Calculate distance from bottom
			const distanceFromBottom = scrollHeight - (scrollTop + clientHeight);

			// Consider user "at bottom" if within threshold
			this.isUserAtBottom = distanceFromBottom <= this.scrollThreshold;

			// Store the current scroll position for debugging
			this.lastScrollTop = scrollTop;

			// Hide indicator if user scrolls to bottom
			if (this.isUserAtBottom && this.hasNewMessages) {
				this.hasNewMessages = false;
				this.newMessageCount = 0;
			}
		},


		async getLastRead(chatId) {
			const controller = new AbortController();
			const requestId = 'getLastRead-' + chatId + '-' + Date.now();

			try {
				this.activeRequests.add(requestId);

				const response = await this.$axios.get(`/chats/${chatId}/last-read`, {
					signal: controller.signal,
					timeout: 5000
				});

				this.lastReadMessageId = response.data.lastReadId;

			} catch (err) {
				if (err.name === 'AbortError' || err.code === 'ECONNABORTED') {
					console.log('Last read request aborted or timed out');
					return;
				}

				console.error('Failed to fetch last read message ID', err);
			} finally {
				this.activeRequests.delete(requestId);
			}
		},

		startReply(message) {
			this.replyingToMessage = message;
			this.showReplyDropdown = null;
			this.focusMessageInput();
		},

		clearReply() {
			this.replyingToMessage = null;
			this.focusMessageInput();
		},

		openForwardModal(message) {
			this.forwardingMessage = message;
			this.showForwardModal = true;
			this.showForwardDropdown = null;
		},

		closeForwardModal() {
			this.showForwardModal = false;
			this.forwardingMessage = null;
		},

		async forwardMessage(recipients) {
			if (!this.forwardingMessage || recipients.length === 0) {
				return { success: false, error: 'Please select at least one recipient' };
			}

			try {
				console.log('Forwarding message:', this.forwardingMessage.id, 'to recipients:', recipients);

				const response = await this.$axios.post(`/messages/${this.forwardingMessage.id}/forwards`, {
					recipients: recipients
				});

				console.log('Forward response:', response.data);

				this.closeForwardModal();

				const recipientCount = recipients.length;
				alert(`Message forwarded successfully to ${recipientCount} recipient${recipientCount !== 1 ? 's' : ''}!`);

				await this.getMyConversations();

				return { success: true, data: response.data };

			} catch (err) {
				console.error('Failed to forward message', err);

				let errorMessage;
				if (err.response?.status === 403) {
					errorMessage = 'You are not authorized to forward this message.';
				} else if (err.response?.status === 404) {
					errorMessage = 'Message not found or some recipients are invalid.';
				} else if (err.response?.status === 401) {
					console.log('Authentication error during message forward');
					this.$emit('logout');
					return { success: false, error: 'Authentication error' };
				} else {
					errorMessage = err.response?.data?.message || 'Failed to forward message. Please try again.';
				}

				return { success: false, error: errorMessage };
			}
		},

		async loadForwardUsers() {
			try {
				const response = await this.$axios.get('/users');
				return { success: true, data: response.data };
			} catch (err) {
				console.error('Failed to fetch forward users', err);

				if (err.response?.status === 401) {
					this.$emit('logout');
					return { success: false, error: 'Authentication error' };
				}

				return { success: false, error: 'Failed to load users. Please try again.' };
			}
		},

		handleNewChatImageSelect(data) {
			this.selectedNewChatImage = data.file

			if (this.newChatImagePreviewUrl) {
				URL.revokeObjectURL(this.newChatImagePreviewUrl)
			}
			this.newChatImagePreviewUrl = data.previewUrl

			this.closeNewChatImageModal()
		},

		openNewChatImageModal() {
			this.showNewChatImageModal = true;
		},

		closeNewChatImageModal() {
			this.showNewChatImageModal = false;
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
				const chat = this.chats[chatIndex];

				const updatedChat = {
					...chat,
					...updates
				};

				console.log('Updated chat preview:', {
					chatId,
					chatName: this.getChatName(updatedChat),
					lastMsgText: updatedChat.lastMsgText,
					lastMsgType: updatedChat.lastMsgType
				});

				this.chats.splice(chatIndex, 1);
				this.chats.unshift(updatedChat);
			}
		},

		confirmDeleteMessage(message) {
			this.showDeleteDropdown = null;

			if (confirm('Are you sure you want to delete this message? This action cannot be undone.')) {
				this.deleteMessage(message);
			}
		},

		handleReplyImageError(replyMessageId) {
			if (this.messageImageUrls[replyMessageId]) {
				URL.revokeObjectURL(this.messageImageUrls[replyMessageId]);
				delete this.messageImageUrls[replyMessageId];
			}
		},

		retryLoadMessages() {
			if (this.selectedChatId) {
				this.getConversation(this.selectedChatId, false);
			}
		},

		// Helper method to handle scroll position changes
		handleScrollPositionChanged(scrollData) {
			this.isUserAtBottom = scrollData.isUserAtBottom;
			this.lastScrollTop = scrollData.scrollTop;

			// Hide indicator if user scrolls to bottom
			if (scrollData.isUserAtBottom && this.hasNewMessages) {
				this.hasNewMessages = false;
				this.newMessageCount = 0;
			}
		},

		// Helper method to handle scroll to new messages
		scrollToNewMessages() {
			this.hasNewMessages = false;
			this.newMessageCount = 0;
		},

		focusMessageInput() {
			if (this.$refs.conversationSection) {
				this.$refs.conversationSection.focusMessageInput();
			}
		},

		restoreScrollPosition(savedScrollTop, heightDifference = 0) {
			if (this.$refs.conversationSection) {
				this.$refs.conversationSection.restoreScrollPosition(savedScrollTop, heightDifference);
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
			console.log('toggleUserSelection called with:', user);
			console.log('selectedUsers before:', this.selectedUsers);

			const index = this.selectedUsers.findIndex(u => u.id === user.id);

			if (index > -1) {
				// User is already selected, remove them
				this.selectedUsers.splice(index, 1);
			} else {
				// User is not selected, add them
				this.selectedUsers.push(user);
			}

			console.log('selectedUsers after:', this.selectedUsers);

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

		closeNewChatModal() {
			this.showNewChatModal = false;

			this.newChatName = '';
			this.initialMessage = '';
			this.selectedUsers = [];
			this.users = [];
			this.filteredUsers = [];
			this.userSearchQuery = '';
			this.newChatError = null;
			this.newChatLoading = false;
			this.loadingUsers = false;
			this.userImageUrls = {};

			this.clearNewChatImageSelection();

			if (this.selectedChatId) {
				this.focusMessageInput();
			}
		},

		openNewChatModal() {
			console.log('Opening new chat modal...'); // Debug log
			this.showNewChatModal = true;

			this.getUsers();
		},

		openRenameGroupModal() {
			this.showRenameGroupModal = true;
		},

		closeRenameGroupModal() {
			this.showRenameGroupModal = false;
			this.setGroupNameError = null;
		},

		openAddToGroupModal() {
			this.addToGroupError = null;
			this.showAddToGroupModal = true;
		},

		closeAddToGroupModal() {
			this.showAddToGroupModal = false;
			this.addToGroupError = null;
		},

		openSetGroupPhotoModal() {
			this.setGroupPhotoError = null;
			this.showSetGroupPhotoModal = true;
		},

		closeSetGroupPhotoModal() {
			this.showSetGroupPhotoModal = false;
			this.setGroupPhotoError = null;
		},

		// Add method to handle visibility changes (pause polling when tab is hidden)
		handleVisibilityChange() {
			if (document.hidden) {
				this.stopPolling();
			} else {
				this.startPolling();
			}
		},

		selectChat(chatId) {
			console.log('Selecting chat:', chatId);

			// Close chat info dropdown when switching chats
			this.showChatInfoDropdown = false;
			this.chatMembers = [];

			// Clear unread count
			const chatIndex = this.chats.findIndex(chat => chat.id === chatId);
			if (chatIndex !== -1 && this.chats[chatIndex].unread > 0) {
				this.chats[chatIndex].unread = 0;
			}

			this.selectedChatId = chatId;
			this.lastReadMessageId = null;
			this.hasNewMessages = false;
			this.newMessageCount = 0;
			this.isUserAtBottom = true;
			this.shouldScrollAfterSend = false;

			// Load chat for first time and scroll to bottom
			this.getConversation(chatId, true);

			// Focus input on message input bar
			this.$nextTick(() => {
				this.focusMessageInput();
			});
		},

		getChatName(chat) {
			if (chat.isGroup) {
				return chat.name || 'Unnamed Group';
			} else {
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

		// Handle user image load error
		handleUserImageError() {
			this.currentUserImageUrl = null;
		},

		handleFilterUsers(searchQuery) {
			this.userSearchQuery = searchQuery;
			this.filterUsers();
		},

		async loadUserImagesForForward(users) {
			for (const user of users) {
				if (user.imageUrl && !this.userImageUrls[user.id]) {
					try {
						const imageUrl = await this.getImage(user.imageUrl);
						if (imageUrl) {
							this.userImageUrls[user.id] = imageUrl;
						}
					} catch (error) {
						console.error(`Failed to load image for user ${user.id}:`, error);
					}
				}
			}
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

		// Handle image selection from modal
		handleImageSelect(data) {
			// data contains { file, previewUrl }
			this.selectedMessageImage = data.file

			// Create preview URL for message input
			if (this.messageImagePreviewUrl) {
				URL.revokeObjectURL(this.messageImagePreviewUrl)
			}
			this.messageImagePreviewUrl = data.previewUrl

			// Close modal
			this.closeImageModal()
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
				const imageUrl = `${baseURL}/uploads/messages/images/${mediaUrl}`;
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

		handleMessageImageError(message) {
			console.error('Image load error for message:', message.id);

			message.imageError = true;

			if (this.messageImageUrls[message.id]) {
				URL.revokeObjectURL(this.messageImageUrls[message.id]);
				delete this.messageImageUrls[message.id];
			}
		},

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
				return date.toLocaleDateString([], {weekday: 'short'}) + ' ' +
					date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
			}

			// If this year, show month and day with time
			if (date.getFullYear() === now.getFullYear()) {
				return date.toLocaleDateString([], {month: 'short', day: 'numeric'}) + ' ' +
					date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
			}

			// Otherwise show full date with time
			return date.toLocaleDateString() + ' ' +
				date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
		},

		handleImageError(event) {
			// Get the chat ID from the image element or find it another way
			const imgElement = event.target;
			const chatId = imgElement.getAttribute('data-chat-id');

			if (chatId) {
				// Remove the failed image URL from cache
				delete this.chatImageUrls[chatId];

				// Also remove from image cache
				const chat = this.chats.find(c => c.id === chatId);
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
					if (this.activeRequests.size < 2) {
						if (this.selectedChatId) {

							await Promise.all([
								this.getConversation(this.selectedChatId, false),
								this.getLastRead(this.selectedChatId),
								this.getMyConversations()
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
			console.log('Cancelling active requests:', this.activeRequests.size);
			this.activeRequests.clear();
		},

		async logout() {
			if (confirm('Are you sure you want to logout?')) {
				this.cleanup();

				localStorage.removeItem('token');
				localStorage.removeItem('user');

				delete this.$axios.defaults.headers.common['Authorization'];

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

			Object.values(this.memberImageUrls).forEach(blobUrl => {
				URL.revokeObjectURL(blobUrl);
			});

			Object.values(this.userImageUrls).forEach(blobUrl => {
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
			this.memberImageUrls = {};
			this.profileImagePreviewUrl = null;
			this.currentUserImageUrl = null;
			this.newChatImagePreviewUrl = null;
			this.tempNewChatImagePreviewUrl = null;
			this.userImageUrls = {};
		}
	}
}
</script>

<style scoped>
@import url('../assets/main.css');
@import url('../assets/modals.css');
@import url('../assets/message.css');
@import "ChatsView.css";
</style>
