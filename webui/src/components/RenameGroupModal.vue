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

				<ErrorMsg v-if="displayError" :msg="displayError"/>
			</div>
			<div class="vue-modal-footer">
				<button class="btn btn-secondary" @click="closeModal">Cancel</button>
				<button
					class="btn btn-primary"
					@click="handleRename"
					:disabled="loading || !groupName.trim()"
				>
					<span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status">
						<span class="visually-hidden">Loading...</span>
					</span>
					Rename
				</button>
			</div>
		</div>
	</div>
</template>

<script>
import ErrorMsg from "./ErrorMsg.vue";

export default {
	name: 'RenameGroupModal',
	components: {ErrorMsg},
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
			groupName: '',
			validationError: ''
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

	computed: {
		displayError () {
			return this.validationError || this.error
		}
	},

	methods: {
		closeModal() {
			this.$emit('close');
		},

		handleRename() {
			if (!this.groupName.trim()) {
				this.validationError = "Group name must be at least 3 characters"
				return;
			}

			if (this.groupName.trim() === this.currentName) {
				this.validationError = "Select a different group name"
				return
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
