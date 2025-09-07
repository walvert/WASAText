<template>
	<div class="chat-header-bar p-3 border-bottom bg-white shadow-sm">
		<div class="d-flex justify-content-between align-items-center">
			<div class="d-flex align-items-center">
				<!-- Chat Image -->
				<div class="chat-header-avatar me-3">
					<div class="avatar-circle" v-if="!chatImageUrl" style="width: 40px; height: 40px;">
						<span class="avatar-text" style="font-size: 16px;">{{ getChatInitials() }}</span>
					</div>
					<img
						v-else
						:src="chatImageUrl"
						:alt="getChatName()"
						class="avatar-image"
						style="width: 40px; height: 40px;"
						@error="$emit('image-error')"
					>
				</div>

				<!-- Chat Name -->
				<div class="chat-header-info">
					<h5 class="chat-header-name mb-0">{{ getChatName() }}</h5>
					<small v-if="selectedChat && selectedChat.isGroup" class="text-muted">
						{{ chatMembers.length > 0 ? `${chatMembers.length} members` : 'Group chat' }}
					</small>
				</div>

				<!-- Info Icon -->
				<div class="chat-info-dropdown ms-2" style="position: relative;">
					<button
						class="btn btn-sm btn-ghost chat-info-btn"
						@click="toggleChatInfoDropdown"
						@click.stop
						title="Chat info"
					>
						<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
							<path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
							<path d="m8.93 6.588-2.29.287-.082.38.45.083c.294.07.352.176.288.469l-.738 3.468c-.194.897.105 1.319.808 1.319.545 0 1.178-.252 1.465-.598l.088-.416c-.2.176-.492.246-.686.246-.275 0-.375-.193-.304-.533L8.93 6.588zM9 4.5a1 1 0 1 1-2 0 1 1 0 0 1 2 0z"/>
						</svg>
					</button>

					<!-- Chat Info Dropdown -->
					<div
						v-if="showChatInfoDropdown"
						class="chat-info-dropdown-menu"
						@click.stop
						@mouseleave="handleChatInfoDropdownMouseLeave"
						@mouseenter="handleChatInfoDropdownMouseEnter"
					>
						<div class="dropdown-header">
							<strong>{{ selectedChat.isGroup ? 'Group Members' : 'Chat Info' }}</strong>
						</div>

						<!-- Loading state -->
						<div v-if="loadingChatMembers" class="dropdown-item-loading">
							<div class="spinner-border spinner-border-sm me-2" role="status">
								<span class="visually-hidden">Loading...</span>
							</div>
							Loading members...
						</div>

						<!-- Members list -->
						<div v-else-if="chatMembers.length > 0">
							<div
								v-for="member in chatMembers"
								:key="member.id"
								class="chat-member-item"
							>
								<div class="member-avatar me-2">
									<div class="avatar-circle" v-if="!memberImageUrls[member.id]" style="width: 28px; height: 28px;">
										<span class="avatar-text" style="font-size: 11px;">{{ getUserInitials(member.username) }}</span>
									</div>
									<img
										v-else
										:src="memberImageUrls[member.id]"
										:alt="member.username"
										class="avatar-image"
										style="width: 28px; height: 28px;"
										@error="() => $emit('member-image-error', member.id)"
									>
								</div>
								<div class="member-info">
									<span class="member-name">{{ member.username }}</span>
									<span v-if="member.username === currentUsername" class="member-badge">You</span>
								</div>
							</div>
						</div>

						<!-- Private chat info -->
						<div v-else-if="!selectedChat.isGroup" class="dropdown-item-text">
							<small class="text-muted">Private conversation</small>
						</div>

						<!-- Error state -->
						<div v-else class="dropdown-item-text">
							<small class="text-muted">Unable to load members</small>
						</div>
					</div>
				</div>
			</div>

			<!-- Actions Dropdown for groups -->
			<div v-if="selectedChat && selectedChat.isGroup" class="dropdown">
				<button class="btn btn-sm btn-ghost chat-actions-btn" data-bs-toggle="dropdown">
					<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
						<path d="M3 9.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3z"/>
					</svg>
				</button>
				<ul class="dropdown-menu dropdown-menu-end">
					<li>
						<button class="dropdown-item" @click="$emit('open-rename-group-modal')">
							<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16" class="me-2">
								<path d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5v11z"/>
								<path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293zm-1.75 2.456-2-2L4.939 9.21a.5.5 0 0 0-.121.196l-.805 2.414a.25.25 0 0 0 .316.316l2.414-.805a.5.5 0 0 0 .196-.12l6.813-6.814z"/>
								<path d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5v11z"/>
							</svg>
							Rename Group
						</button>
					</li>
					<li>
						<button class="dropdown-item" @click="$emit('open-add-to-group-modal')">
							<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16" class="me-2">
								<path d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm2-3a2 2 0 1 1-4 0 2 2 0 0 1 4 0zm4 8c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4zm-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10z"/>
								<path fill-rule="evenodd" d="M13.5 5a.5.5 0 0 1 .5.5V7h1.5a.5.5 0 0 1 0 1H14v1.5a.5.5 0 0 1-1 0V8h-1.5a.5.5 0 0 1 0-1H13V5.5a.5.5 0 0 1 .5-.5z"/>
							</svg>
							Add Member
						</button>
					</li>
					<li>
						<button class="dropdown-item" @click="$emit('open-set-group-photo-modal')">
							<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16" class="me-2">
								<path d="M10.5 8.5a2.5 2.5 0 1 1-5 0 2.5 2.5 0 0 1 5 0z"/>
								<path d="M2 4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2h-1.172a2 2 0 0 1-1.414-.586L9.828 1.828A2 2 0 0 0 8.414 1H7.586a2 2 0 0 0-1.414.586L4.586 3.414A2 2 0 0 1 3.172 4H2zm.5 2a.5.5 0 1 1 0-1 .5.5 0 0 1 0 1zm9 2.5a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0z"/>
								<path d="M13.5 1a.5.5 0 0 1 .5.5v2a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1 0-1h1.5V1.5a.5.5 0 0 1 .5-.5z"/>
								<path d="M13 0a1 1 0 0 1 1 1v2a1 1 0 0 1-1 1h-2a1 1 0 0 1-1-1V1a1 1 0 0 1 1-1h2zm0 1h-2v2h2V1z"/>
								<circle cx="15" cy="2" r="1" fill="currentColor"/>
							</svg>
							Set Group Photo
						</button>
					</li>
					<li>
						<hr class="dropdown-divider">
					</li>
					<li>
						<button class="dropdown-item text-danger" @click="$emit('leave-group')">
							<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16" class="me-2">
								<path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0v2z"/>
								<path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3z"/>
							</svg>
							Leave Group
						</button>
					</li>
				</ul>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'ChatHeader',
	props: {
		selectedChat: {
			type: Object,
			required: true
		},
		chatImageUrl: {
			type: String,
			default: null
		},
		currentUsername: {
			type: String,
			required: true
		},
		chatMembers: {
			type: Array,
			default: () => []
		},
		loadingChatMembers: {
			type: Boolean,
			default: false
		},
		memberImageUrls: {
			type: Object,
			default: () => ({})
		}
	},

	emits: [
		'image-error',
		'member-image-error',
		'open-rename-group-modal',
		'open-add-to-group-modal',
		'open-set-group-photo-modal',
		'leave-group',
		'toggle-chat-info-dropdown',
		'get-group-members'
	],

	data() {
		return {
			showChatInfoDropdown: false,
			chatInfoDropdownTimeout: null
		}
	},

	methods: {
		getChatName() {
			if (this.selectedChat.isGroup) {
				return this.selectedChat.name || 'Unnamed Group';
			} else {
				return this.selectedChat.name || 'Private Chat';
			}
		},

		getChatInitials() {
			const name = this.getChatName();
			if (this.selectedChat.isGroup) {
				return name.charAt(0).toUpperCase();
			} else {
				const words = name.split(' ');
				if (words.length >= 2) {
					return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase();
				}
				return name.charAt(0).toUpperCase();
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

		toggleChatInfoDropdown() {
			if (this.showChatInfoDropdown) {
				this.showChatInfoDropdown = false;
			} else {
				this.showChatInfoDropdown = true;
				this.$emit('get-group-members');
			}

			if (this.chatInfoDropdownTimeout) {
				clearTimeout(this.chatInfoDropdownTimeout);
				this.chatInfoDropdownTimeout = null;
			}
		},

		handleChatInfoDropdownMouseLeave() {
			this.chatInfoDropdownTimeout = setTimeout(() => {
				this.showChatInfoDropdown = false;
			}, 300);
		},

		handleChatInfoDropdownMouseEnter() {
			if (this.chatInfoDropdownTimeout) {
				clearTimeout(this.chatInfoDropdownTimeout);
				this.chatInfoDropdownTimeout = null;
			}
		}
	},

	beforeUnmount() {
		if (this.chatInfoDropdownTimeout) {
			clearTimeout(this.chatInfoDropdownTimeout);
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";

.chat-header-bar {
	flex-shrink: 0;
	border-bottom: 1px solid #dee2e6;
	z-index: 10;
	background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
	border-bottom: 1px solid #e0e6ed;
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.chat-header-avatar {
	flex-shrink: 0;
	position: relative;
}

.chat-header-avatar .avatar-circle,
.chat-header-avatar .avatar-image {
	border: 2px solid #e9ecef;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	transition: transform 0.2s ease;
}

.chat-header-avatar:hover .avatar-circle,
.chat-header-avatar:hover .avatar-image {
	transform: scale(1.05);
	border-color: #007bff;
}

.chat-header-info {
	flex: 1;
	min-width: 0;
}

.chat-header-name {
	font-weight: 700;
	font-size: 1.1rem;
	color: #2c3e50;
	margin: 0;
	line-height: 1.2;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

</style>
