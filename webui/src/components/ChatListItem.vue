<template>
	<div
		class="chat-item"
		:class="{ active: isSelected }"
		@click="$emit('select', chat.id)"
	>
		<div class="chat-avatar">
			<div class="avatar-circle" v-if="!chat.image">
				<span class="avatar-text">{{ getChatInitials(chat) }}</span>
			</div>
			<img
				v-else
				:src="chatImageUrl"
				:alt="chatName"
				class="avatar-image"
				@error="$emit('image-error', chat.id)"
			>
		</div>
		<div class="chat-info">
			<div class="chat-header">
				<h6 class="chat-name">{{ chatName }}</h6>
				<span class="chat-time">{{ formattedTime }}</span>
			</div>
			<div class="chat-preview">
				<p class="last-message">{{ lastMessagePreview }}</p>
				<div class="chat-preview-meta d-flex align-items-center justify-content-between">
					<div class="message-type-indicator">
						<span v-if="chat.lastMsgType === 'text'" class="message-type-icon">💬</span>
						<span v-else-if="chat.lastMsgType === 'image'" class="message-type-icon">📷</span>
						<span v-else class="message-type-icon">🎞</span>
					</div>
					<div v-if="chat.unread > 0" class="unread-badge">
						{{ chat.unread > 99 ? '99+' : chat.unread }}
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'ChatListItem',
	props: {
		chat: { type: Object, required: true },
		isSelected: { type: Boolean, default: false },
		chatImageUrl: String,
		currentUsername: String
	},
	emits: ['select', 'image-error'],
	computed: {
		chatName() {
			return this.getChatName(this.chat)
		},
		formattedTime() {
			return this.formatMessageTime(this.chat.lastMsgTime)
		},
		lastMessagePreview() {
			return this.getLastMessagePreview(this.chat)
		}
	},
	methods: {
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
				return name.charAt(0).toUpperCase();
			} else {
				const words = name.split(' ');
				if (words.length >= 2) {
					return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase();
				}
				return name.charAt(0).toUpperCase();
			}
		},

		formatMessageTime(timestamp) {
			if (!timestamp) return '';

			const date = new Date(timestamp);
			const now = new Date();

			if (date.toDateString() === now.toDateString()) {
				return date.toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit',
					hour12: false
				});
			}

			const yesterday = new Date(now);
			yesterday.setDate(now.getDate() - 1);
			if (date.toDateString() === yesterday.toDateString()) {
				return 'Yesterday ' + date.toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit',
					hour12: false
				});
			}

			const weekAgo = new Date(now);
			weekAgo.setDate(now.getDate() - 7);
			if (date > weekAgo) {
				return date.toLocaleDateString([], {weekday: 'short'}) + ' ' +
					date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
			}

			if (date.getFullYear() === now.getFullYear()) {
				return date.toLocaleDateString([], {month: 'short', day: 'numeric'}) + ' ' +
					date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
			}

			return date.toLocaleDateString() + ' ' +
				date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
		},

		getLastMessagePreview(chat) {
			if (!chat.lastMsgText) {
				return 'No messages yet';
			}

			let preview = chat.lastMsgText;

			if (chat.lastMsgType === 'image') {
				if (preview && !preview.includes('📷')) {
				} else {
					preview = '📷 Photo';
				}
			} else if (chat.lastMsgType === 'gif') {
				if (preview && !preview.includes('🎞️')) {
				} else {
					preview = '🎞️ GIF';
				}
			}

			if (chat.isGroup && chat.lastMsgUsername) {
				if (preview === '📷 Photo' || preview === '🎞️ GIF') {
					preview = `${chat.lastMsgUsername}: ${preview}`;
				} else if (preview && preview.length > 0) {
					preview = `${chat.lastMsgUsername}: ${preview}`;
				}
			}

			const maxLength = 40;
			if (preview.length > maxLength) {
				return preview.substring(0, maxLength) + '...';
			}

			return preview;
		},
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>
