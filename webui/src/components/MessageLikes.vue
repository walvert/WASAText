<template>
	<div
		class="message-likes"
		@click="$emit('toggle-dropdown')"
		@click.stop
	>
		<svg width="12" height="12" fill="currentColor" viewBox="0 0 16 16">
			<path d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.011L8 2.748zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143c.06.055.119.112.176.171a3.12 3.12 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15z"/>
		</svg>
		{{ likes.length }}

		<!-- Likes Dropdown -->
		<div
			v-if="showDropdown"
			class="likes-dropdown"
			@click.stop
			@mouseleave="$emit('mouse-leave')"
			@mouseenter="$emit('mouse-enter')"
		>
			<div class="likes-list">
				<div
					v-for="username in likes"
					:key="username"
					class="like-user"
				>
					{{ username }}
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'MessageLikes',
	props: {
		likes: {
			type: Array,
			required: true
		},
		showDropdown: {
			type: Boolean,
			default: false
		}
	},
	emits: [
		'toggle-dropdown',
		'mouse-leave',
		'mouse-enter'
	]
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>
