<template>
    <h4>Email Trackers</h4>
    <p>
        Remove tracking pixels from forwarded emails. This stops senders from using them to see when you open an email, your IP address, your device and your approximate location.
        When trackers are removed, the Mailx header at the top of the email shows how many, with a link to a full report.
    </p>
    <div class="mb-8">
        <label for="block-trackers">
            Block email trackers:
        </label>
        <input
            @change="save"
            v-model="enabled"
            :disabled="saving"
            :aria-invalid="error ? 'true' : undefined"
            :aria-describedby="error ? 'block-trackers-error' : undefined"
            id="block-trackers"
            type="checkbox"
            class="checkbox-switch"
        >
        <p v-if="error" id="block-trackers-error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { settingsApi } from '../api/settings.ts'

const props = defineProps({
    initialEnabled: { type: Boolean, default: true },
})

const enabled = ref(props.initialEnabled)
const saving = ref(false)
const error = ref('')
const success = ref('')

const save = async () => {
    saving.value = true
    error.value = ''
    success.value = ''

    try {
        const res = await settingsApi.updateBlockTrackers(enabled.value)
        success.value = res.data.message
    } catch (err) {
        // Show the setting as it's still stored
        enabled.value = !enabled.value
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        } else {
            error.value = 'Something went wrong'
        }
    } finally {
        saving.value = false
    }
}
</script>
