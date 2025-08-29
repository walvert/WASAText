<script>
export default {
	props: ['msg', 'message']
}
</script>

<template>
	<div class="alert alert-danger" role="alert">
		{{ msg || message || "error" }}
	</div>
</template>

<style>
</style>
