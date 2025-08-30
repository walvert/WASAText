<!-- UserAvatar.vue -->
<template>
	<div class="user-profile-avatar" @click="handleClick" :style="{ cursor: clickable ? 'pointer' : 'default' }">
		<div class="avatar-circle" v-if="!imageUrl" :style="sizeStyles">
			<span class="avatar-text" :style="textSizeStyles">{{ initials }}</span>
		</div>
		<img
			v-else
			:src="imageUrl"
			:alt="username || 'User'"
			class="avatar-image"
			:style="sizeStyles"
			@error="$emit('image-error')"
		>
	</div>
</template>

<script>
export default {
	name: 'UserAvatar',
	props: {
		username: String,
		imageUrl: String,
		size: { type: Number, default: 42 },
		clickable: { type: Boolean, default: false }
	},
	emits: ['click', 'image-error'],
	computed: {
		initials() {
			return this.getUserInitials(this.username || 'User')
		},
		sizeStyles() {
			return {
				width: `${this.size}px`,
				height: `${this.size}px`
			}
		},
		textSizeStyles() {
			return {
				fontSize: `${Math.round(this.size * 0.38)}px`
			}
		}
	},
	methods: {
		handleClick() {
			if (this.clickable) {
				this.$emit('click')
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
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>
