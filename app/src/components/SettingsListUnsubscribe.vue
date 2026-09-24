<template>
    <h4>List-Unsubscribe</h4>
    <p>
        Control how List-Unsubscribe and List-Unsubscribe-Post headers are set on forwarded emails.
        On email clients that support it, clicking "Unsubscribe" will perform one of these actions.
    </p>
    <form @submit.prevent="save" class="mb-8">
        <div class="max-w-full mb-4">
            <label for="list-unsubscribe-behaviour">
                Select behaviour:
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
            Update List-Unsubscribe Behaviour
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
    { value: 0, label: 'Use original List-Unsubscribe, fallback to one-click deactivate' },
    { value: 1, label: 'Use one-click deactivate' },
    { value: 2, label: 'Use one-click delete' },
    { value: 3, label: 'Use one-click block sender email' },
    { value: 4, label: 'Use one-click block sender domain' },
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
