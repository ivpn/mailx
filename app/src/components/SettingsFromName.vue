<template>
    <form @submit.prevent="submit" class="mb-8">
        <div class="max-w-xs mb-4">
            <label for="from-name">
                From name:
            </label>
            <input
                v-model="fromName"
                :maxlength="maxLength"
                :aria-invalid="error ? 'true' : undefined"
                :aria-describedby="error ? 'from-name-error' : undefined"
                id="from-name"
                type="text"
            >
        </div>
        <button type="submit" class="cta sm" :disabled="saving">
            Update From Name
        </button>
        <p v-if="error" id="from-name-error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { settingsApi } from '../api/settings.ts'
import { useSettingSave } from '../lib/useSettingSave.ts'

const props = defineProps({
    initialFromName: { type: String, default: '' },
})

// Matches the API limit
const maxLength = 64

const fromName = ref(props.initialFromName)
const { saving, error, success, save } = useSettingSave()

const submit = async () => {
    if (await save(() => settingsApi.updateFromName(fromName.value))) {
        // Show the name as it's stored
        fromName.value = fromName.value.trim()
    }
}
</script>
