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
import { formatMessageTime } from "../utils/helpers";

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
		formatMessageTime,
	}
}
</script>

<style scoped>
@import url('../assets/message.css');
@import "../views/ChatsView.css";
</style>
