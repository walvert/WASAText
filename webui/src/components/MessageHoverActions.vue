<template>
	<div class="message-hover-actions">
		<!-- Like Button -->
		<button
			class="like-btn"
			:class="{ liked: isLiked }"
			@click="$emit('toggle-like')"
		>
			<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
				<path d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.011L8 2.748zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143c.06.055.119.112.176.171a3.12 3.12 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15z"/>
			</svg>
		</button>

		<!-- Reply Button -->
		<button
			class="message-reply-btn"
			@click.stop="$emit('toggle-reply-dropdown')"
		>
			<!-- Backward arrow icon for reply -->
			<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
				<path d="M9.5 13a.5.5 0 0 1-.5-.5V10H2.5a.5.5 0 0 1 0-1H9V6.5a.5.5 0 0 1 .854-.354l3 3a.5.5 0 0 1 0 .708l-3 3A.5.5 0 0 1 9.5 13z" transform="scale(-1,1) translate(-16,0)"/>
			</svg>

			<!-- Dropdown -->
			<div
				v-if="showReplyDropdown"
				class="reply-dropdown"
				@click.stop
				@mouseleave="$emit('reply-mouse-leave')"
				@mouseenter="$emit('reply-mouse-enter')"
			>
				<div
					class="reply-option"
					@click="$emit('start-reply')"
				>
					<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
						<path d="M9.5 13a.5.5 0 0 1-.5-.5V10H2.5a.5.5 0 0 1 0-1H9V6.5a.5.5 0 0 1 .854-.354l3 3a.5.5 0 0 1 0 .708l-3 3A.5.5 0 0 1 9.5 13z" transform="scale(-1,1) translate(-16,0)"/>
					</svg>
					Reply
				</div>
			</div>
		</button>

		<!-- Forward Button -->
		<button
			class="message-forward-btn"
			@click.stop="$emit('toggle-forward-dropdown')"
		>
			<!-- Forward arrow icon -->
			<svg width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
				<path d="M9.5 13a.5.5 0 0 1-.5-.5V10H2.5a.5.5 0 0 1 0-1H9V6.5a.5.5 0 0 1 .854-.354l3 3a.5.5 0 0 1 0 .708l-3 3A.5.5 0 0 1 9.5 13z"/>
			</svg>

			<!-- Dropdown -->
			<div
				v-if="showForwardDropdown"
				class="forward-dropdown"
				@click.stop
				@mouseleave="$emit('forward-mouse-leave')"
				@mouseenter="$emit('forward-mouse-enter')"
			>
				<div
					class="forward-option"
					@click="$emit('open-forward-modal')"
				>
					<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
						<path d="M9.5 13a.5.5 0 0 1-.5-.5V10H2.5a.5.5 0 0 1 0-1H9V6.5a.5.5 0 0 1 .854-.354l3 3a.5.5 0 0 1 0 .708l-3 3A.5.5 0 0 1 9.5 13z"/>
					</svg>
					Forward
				</div>
			</div>
		</button>

		<!-- Delete Button -->
		<button
			v-if="isCurrentUser"
			class="delete-btn"
			:class="{ loading: isDeleting }"
			@click.stop="$emit('toggle-delete-dropdown')"
			:disabled="isDeleting"
		>
			<div v-if="isDeleting" class="spinner-border spinner-border-sm" role="status">
				<span class="visually-hidden">Deleting...</span>
			</div>
			<svg v-else width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
				<path d="M2.5 1a1 1 0 0 0-1 1v1a1 1 0 0 0 1 1H3v9a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V4h.5a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H10a1 1 0 0 0-1-1H7a1 1 0 0 0-1 1H2.5zm3 4a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 .5-.5zM8 5a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7A.5.5 0 0 1 8 5zm3 .5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 1 0z"/>
			</svg>

			<!-- Dropdown -->
			<div
				v-if="showDeleteDropdown"
				class="delete-dropdown"
				@click.stop
				@mouseleave="$emit('delete-mouse-leave')"
				@mouseenter="$emit('delete-mouse-enter')"
			>
				<div
					class="delete-option"
					@click="$emit('confirm-delete')"
				>
					<svg width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
						<path d="M2.5 1a1 1 0 0 0-1 1v1a1 1 0 0 0 1 1H3v9a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V4h.5a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H10a1 1 0 0 0-1-1H7a1 1 0 0 0-1 1H2.5zm3 4a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 .5-.5zM8 5a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7A.5.5 0 0 1 8 5zm3 .5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 1 0z"/>
					</svg>
					Delete Message
				</div>
			</div>
		</button>
	</div>
</template>

<script>
export default {
	name: 'MessageHoverActions',
	props: {
		message: {
			type: Object,
			required: true
		},
		isCurrentUser: {
			type: Boolean,
			required: true
		},
		isLiked: {
			type: Boolean,
			default: false
		},
		isDeleting: {
			type: Boolean,
			default: false
		},
		showReplyDropdown: {
			type: Boolean,
			default: false
		},
		showForwardDropdown: {
			type: Boolean,
			default: false
		},
		showDeleteDropdown: {
			type: Boolean,
			default: false
		}
	},
	emits: [
		'toggle-like',
		'toggle-reply-dropdown',
		'toggle-forward-dropdown',
		'toggle-delete-dropdown',
		'start-reply',
		'open-forward-modal',
		'confirm-delete',
		'reply-mouse-leave',
		'reply-mouse-enter',
		'forward-mouse-leave',
		'forward-mouse-enter',
		'delete-mouse-leave',
		'delete-mouse-enter'
	]
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";
</style>
