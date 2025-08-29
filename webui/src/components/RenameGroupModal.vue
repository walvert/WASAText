<template>
	<div v-if="show" class="vue-modal" @click.self="closeModal">
		<div class="vue-modal-dialog">
			<div class="vue-modal-header">
				<h5 class="vue-modal-title">Rename Group</h5>
				<button class="vue-modal-close" @click="closeModal">×</button>
			</div>
			<div class="vue-modal-body">
				<div class="mb-3">
					<label class="form-label">Group Name</label>
					<input
						type="text"
						class="form-control"
						v-model="groupName"
						placeholder="Enter new group name"
						@keyup.enter="handleRename"
						ref="groupNameInput"
					>
				</div>
				<div v-if="error" class="error-msg">{{ error }}</div>
			</div>
			<div class="vue-modal-footer">
				<button class="btn btn-secondary" @click="closeModal">Cancel</button>
				<button
					class="btn btn-primary"
					@click="handleRename"
					:disabled="loading || !groupName.trim()"
				>
					<div v-if="loading" class="spinner-border"></div>
					Rename
				</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'RenameGroupModal',
	props: {
		show: {
			type: Boolean,
			default: false
		},
		currentName: {
			type: String,
			default: ''
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
			groupName: ''
		}
	},

	watch: {
		show(newVal) {
			if (newVal) {
				this.groupName = this.currentName;
				this.$nextTick(() => {
					this.focusInput();
				});
			} else {
				this.resetForm();
			}
		}
	},

	methods: {
		closeModal() {
			this.$emit('close');
		},

		handleRename() {
			if (!this.groupName.trim()) {
				return;
			}
			this.$emit('rename', this.groupName.trim());
		},

		focusInput() {
			if (this.$refs.groupNameInput) {
				this.$refs.groupNameInput.focus();
				this.$refs.groupNameInput.select();
			}
		},

		resetForm() {
			this.groupName = '';
		}
	}
}
</script>

<style scoped>
@import url('../views/ChatsView.css');
@import url('../assets/modals.css');
</style>
