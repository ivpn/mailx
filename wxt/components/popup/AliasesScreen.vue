<template>
    <div class="page p-5">
        <p v-if="isLoading">Loading...</p>
        <p v-else-if="error" class="error">{{ error }}</p>
        <div v-else class="w-full">
            <hr class="m-0">
            <div v-for="alias in list" :key="alias.id" class="py-3 border-b border-secondary flex items-center gap-x-4">
                <div class="flex items-center hs-tooltip">
                    <input
                        @change=""
                        v-bind:checked="alias.enabled"
                        type="checkbox"
                        class="xs"
                    >
                </div>
                <div class="grow font-medium">{{ alias.name }}</div>
                <button class="">
                    <i class="icon icon-secondary trash text-xs"></i>
                </button>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'
import { Alias } from '@/lib/types'

const props = defineProps<{ apiToken: string }>()
const list = ref([] as Alias[])
const isLoading = ref(false)
const error = ref<string | null>(null)

const fetchAliases = async () => {
    error.value = null

    try {
        isLoading.value = true
        const res = await api.fetchAliases(props.apiToken)
        list.value = res.aliases
        console.log('Fetched aliases:', res.aliases)
    } catch (err) {
        error.value = 'An unexpected error occurred'
        console.error('Fetch aliases error:', err)
    } finally {
        isLoading.value = false
    }
}

onMounted(() => {
    fetchAliases()
})
</script>
