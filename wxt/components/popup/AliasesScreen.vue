<template>
    <div class="page center">
        <h1>Aliases</h1>
        <p v-if="isLoading">Loading...</p>
        <p v-else-if="error" class="error">{{ error }}</p>
        <div v-else>
            <ul>
                <li v-for="alias in list" :key="alias.id">
                    {{ alias.name }}
                </li>
            </ul>
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
