<script>
export default {
	props: ['msg']
}
</script>

<template>
	<div class="alert alert-danger" role="alert">
		{{ msg || "error" }}
	</div>
</template>

<style>
</style>
