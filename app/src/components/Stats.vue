<template>
    <div class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 class="text-2xl font-bold text-gray-800 mb-5">Stats</h1>
        <div class="mb-5">
            <p>Forwards: {{ stats.forwards }}</p>
            <p>Blocks: {{ stats.blocks }}</p>
            <p>Replies: {{ stats.replies }}</p>
            <p>Sends: {{ stats.sends }}</p>
            <p>Bandwidth: {{ getBandwidth() }}kb</p>
            <p>Aliases: {{ stats.aliases }}</p>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'

const stats = ref({
    forwards: 0,
    blocks: 0,
    replies: 0,
    sends: 0,
    bandwidth: 0,
    aliases: 0,
    messages: {
        created_at: '',
        type: 0,
    }
})
const error = ref('')

const getStats = async () => {
    try {
        const response = await userApi.stats()
        stats.value = response.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const getBandwidth = () => {
    return (stats.value.bandwidth / 1024).toFixed(2)
}

onMounted(() => {
    getStats()
})
</script>