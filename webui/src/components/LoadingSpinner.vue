<script>
export default {
	props: ["loading"]
}
</script>

<template>
	<div>
		<div v-if="loading" style="text-align: center">
			<div class="spinner-border" role="status">
				<span class="visually-hidden">Loading...</span>
			</div>
		</div>
		<div v-if="!loading"><slot /></div>
	</div>
</template>
