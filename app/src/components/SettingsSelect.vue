<template>
    <form @submit.prevent="submit" class="mb-8">
        <div class="max-w-xs mb-4">
            <label :for="id">
                {{ label }}
            </label>
            <select
                :id="id"
                v-model="value"
                :disabled="!options.length"
                :aria-invalid="error ? 'true' : undefined"
                :aria-describedby="error ? `${id}-error` : undefined"
                required
            >
                <option v-for="option in options" :key="option.value" :value="option.value">
                    {{ option.label }}
                </option>
            </select>
        </div>
        <button type="submit" class="cta sm" :disabled="saving || !options.length || !changed">
            {{ buttonLabel }}
        </button>
        <p v-if="error" :id="`${id}-error`" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </form>
</template>

<script setup lang="ts">
import { computed, ref, watch, type PropType } from 'vue'
import { useSettingSave } from '../lib/useSettingSave.ts'

export interface SelectOption {
    value: string
    label: string
}

const props = defineProps({
    id: { type: String, required: true },
    label: { type: String, required: true },
    buttonLabel: { type: String, required: true },
    options: { type: Array as PropType<SelectOption[]>, required: true },
    initialValue: { type: String, default: '' },
    update: { type: Function as PropType<(value: string) => Promise<{ data: { message: string } }>>, required: true },
})

const value = ref(props.initialValue)
const saved = ref(props.initialValue)
// Update is enabled only when there is something to save, including a first option
// shown in place of a value that isn't stored
const changed = computed(() => value.value !== saved.value)
const { saving, error, success, save } = useSettingSave()

// Options can load after the setting, and a value that isn't stored yet shows the first option
watch(() => props.options, (options) => {
    if (options.length && !options.some((option) => option.value === value.value)) {
        value.value = options[0].value
    }
}, { immediate: true })

const submit = async () => {
    if (!changed.value) return

    const next = value.value
    if (await save(() => props.update(next))) {
        saved.value = next
    }
}
</script>
