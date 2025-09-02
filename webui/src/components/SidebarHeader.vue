<template>
	<div class="sidebar-header-bar">
		<div class="d-flex justify-content-between align-items-center">
			<!-- User Profile Section -->
			<div class="d-flex align-items-center">
				<div class="user-profile-avatar me-3" @click="$emit('open-profile')" style="cursor: pointer;" title="Edit Profile">
					<div class="avatar-circle" v-if="!currentUserImageUrl" style="width: 42px; height: 42px;">
						<span class="avatar-text" style="font-size: 16px;">{{ getUserInitials(currentUsername || 'User') }}</span>
					</div>
					<img
						v-else
						:src="currentUserImageUrl"
						:alt="currentUsername || 'User'"
						class="avatar-image"
						style="width: 42px; height: 42px;"
						@error="$emit('user-image-error')"
					>
				</div>
				<div class="sidebar-user-info">
					<h5 class="sidebar-title">Chats</h5>
					<small class="sidebar-username">{{ currentUsername || 'User' }}</small>
				</div>
			</div>

			<!-- Action Icons -->
			<div class="d-flex sidebar-actions-gap">
				<button
					class="btn-sidebar-action"
					@click="$emit('new-chat')"
					title="New Chat"
				>
					<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
						<path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
					</svg>
				</button>
				<button
					class="btn-sidebar-action btn-sidebar-danger"
					@click="$emit('logout')"
					title="Logout"
				>
					<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
						<path fill-rule="evenodd" d="M6 12.5a.5.5 0 0 0 .5.5h8a.5.5 0 0 0 .5-.5v-9a.5.5 0 0 0-.5-.5h-8a.5.5 0 0 0-.5.5v2a.5.5 0 0 1-1 0v-2A1.5 1.5 0 0 1 6.5 2h8A1.5 1.5 0 0 1 16 3.5v9a1.5 1.5 0 0 1-1.5 1.5h-8A1.5 1.5 0 0 1 5 12.5v-2a.5.5 0 0 1 1 0v2z"/>
						<path fill-rule="evenodd" d="M.146 8.354a.5.5 0 0 1 0-.708l3-3a.5.5 0 1 1 .708.708L1.707 7.5H10.5a.5.5 0 0 1 0 1H1.707l2.147 2.146a.5.5 0 0 1-.708.708l-3-3z"/>
					</svg>
				</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'SidebarHeader',
	props: {
		currentUserImageUrl: String,
		currentUsername: String
	},
	emits: ['open-profile', 'new-chat', 'logout', 'user-image-error'],
	methods: {
		getUserInitials(username) {
			if (!username) return '?';

			const words = username.split(' ');
			if (words.length >= 2) {
				return (words[0].charAt(0) + words[1].charAt(0)).toUpperCase();
			}
			return username.charAt(0).toUpperCase();
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>
