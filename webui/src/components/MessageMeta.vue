<template>
	<div class="message-meta">
		<span class="message-time">{{ formattedTime }}</span>
		<span v-if="isCurrentUser" class="message-status">
      <!-- Double checkmark for read messages -->
      <svg
		  v-if="isRead"
		  class="message-status-read"
		  viewBox="0 0 18 12"
		  fill="none"
		  xmlns="http://www.w3.org/2000/svg"
	  >
        <path
			d="M1 6L4 9L10 3"
			stroke="currentColor"
			stroke-width="1.5"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>
        <path
			d="M7 6L10 9L16 3"
			stroke="currentColor"
			stroke-width="1.5"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>
      </svg>

			<!-- Single checkmark for sent but not read messages -->
      <svg
		  v-else
		  class="message-status-sent"
		  viewBox="0 0 12 12"
		  fill="none"
		  xmlns="http://www.w3.org/2000/svg"
	  >
        <path
			d="M1 6L4 9L10 3"
			stroke="currentColor"
			stroke-width="1.5"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>
      </svg>
    </span>
	</div>
</template>

<script>
export default {
	name: 'MessageMeta',
	props: {
		message: {
			type: Object,
			required: true
		},
		isCurrentUser: {
			type: Boolean,
			required: true
		},
		isRead: {
			type: Boolean,
			default: false
		}
	},
	computed: {
		formattedTime() {
			return this.formatMessageTime(this.message.createdAt)
		}
	},
	methods: {
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
		}
	}
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";
</style>
