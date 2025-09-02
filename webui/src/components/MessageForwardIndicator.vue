<template>
	<div class="message-forwarded-indicator">
		<svg class="forwarded-icon" viewBox="0 0 24 24" fill="currentColor">
			<path d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z"/>
			<path d="M2.59 16.59L7.17 12 2.59 7.41 4 6l6 6-6 6-1.41-1.41z"/>
		</svg>
		<span class="forwarded-text">Forwarded</span>
	</div>
</template>

<script>
export default {
	name: 'MessageForwardIndicator'
}
</script>

<style scoped>
@import "../views/ChatsView.css";
</style>
