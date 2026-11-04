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
        <button type="submit" class="cta sm" :disabled="saving || !changed">
            Update
        </button>
        <p v-if="error" id="list-unsubscribe-behaviour-error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </form>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { settingsApi } from '../api/settings.ts'
import { useSettingSave } from '../lib/useSettingSave.ts'

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
const saved = ref(props.initialBehaviour)
// Update is enabled only when the selection differs from the stored behaviour
const changed = computed(() => behaviour.value !== saved.value)
const { saving, error, success, save: saveSetting } = useSettingSave()

const save = async () => {
    if (!changed.value) return

    const next = behaviour.value
    if (await saveSetting(() => settingsApi.updateListUnsubscribeBehaviour(next))) {
        saved.value = next
    }
}
</script>
