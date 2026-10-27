<template>
    <h4>Unsubscribe Behaviour</h4>
    <p>
        Choose what happens when you click "Unsubscribe" in your email client (where supported) on emails forwarded by Mailx.
        This sets the List-Unsubscribe and List-Unsubscribe-Post headers on those emails.
    </p>
    <form @submit.prevent="save" class="mb-8">
        <div class="max-w-full mb-4">
            <label for="list-unsubscribe-behaviour">
                Set "Unsubscribe" behaviour:
            </label>
            <select
                id="list-unsubscribe-behaviour"
                v-model.number="behaviour"
                :aria-invalid="error ? 'true' : undefined"
                :aria-describedby="error ? 'list-unsubscribe-behaviour-error' : undefined"
                required
            >
                <option v-for="option in options" :key="option.value" :value="option.value">
                    {{ option.label }}
                </option>
            </select>
        </div>
        <button type="submit" class="cta sm" :disabled="saving">
            Update
        </button>
        <p v-if="error" id="list-unsubscribe-behaviour-error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { settingsApi } from '../api/settings.ts'

const props = defineProps({
    initialBehaviour: { type: Number, default: 0 },
})

const options = [
    { value: 0, label: 'Unsubscribe from sender, fallback: deactivate alias' },
    { value: 1, label: 'Deactivate alias' },
    { value: 2, label: 'Delete alias' },
    { value: 3, label: 'Block sender email address' },
    { value: 4, label: 'Block sender domain' },
]

const behaviour = ref(props.initialBehaviour)
const saving = ref(false)
const error = ref('')
const success = ref('')

const save = async () => {
    saving.value = true
    error.value = ''
    success.value = ''

    try {
        const res = await settingsApi.updateListUnsubscribeBehaviour(behaviour.value)
        success.value = res.data.message
    } catch (err) {
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
