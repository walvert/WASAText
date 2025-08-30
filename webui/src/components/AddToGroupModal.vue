<template>
	<div v-if="show" class="vue-modal" @click.self="closeModal">
		<div class="vue-modal-dialog">
			<div class="vue-modal-header">
				<h5 class="vue-modal-title">Add Member to Group</h5>
				<button class="vue-modal-close" @click="closeModal">×</button>
			</div>
			<div class="vue-modal-body">
				<div class="mb-3">
					<label class="form-label">Username</label>
					<input
						type="text"
						class="form-control"
						v-model="localUsername"
						placeholder="Enter username to add"
						@keyup.enter="handleAddMember"
						ref="memberUsernameInput"
					>
					<div class="form-text">Enter the username of the person you want to add to this group.</div>
				</div>
				<div v-if="error" class="error-msg">{{ error }}</div>
			</div>
			<div class="vue-modal-footer">
				<button class="btn btn-secondary" @click="closeModal">Cancel</button>
				<button
					class="btn btn-primary"
					@click="handleAddMember"
					:disabled="loading || !localUsername.trim()"
				>
					  <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status">
						<span class="visually-hidden">Loading...</span>
					  </span>
					Add Member
				</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'AddToGroupModal',

	props: {
		show: {
			type: Boolean,
			required: true
		},
		loading: {
			type: Boolean,
			default: false
		},
		error: {
			type: String,
			default: null
		}
	},

	data() {
		return {
			localUsername: ''
		}
	},

	watch: {
		show(newVal) {
			if (newVal) {
				// Reset form when modal opens
				this.localUsername = ''

				// Focus input after modal is shown
				this.$nextTick(() => {
					if (this.$refs.memberUsernameInput) {
						this.$refs.memberUsernameInput.focus()
					}
				})
			}
		}
	},

	methods: {
		closeModal() {
			this.$emit('close')
		},

		handleAddMember() {
			if (!this.localUsername.trim()) {
				return
			}

			this.$emit('add-member', this.localUsername.trim())
		}
	}
}
</script>

<style scoped>
@import "../views/ChatsView.css";
@import url('../assets/modals.css');
</style>

