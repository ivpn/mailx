<template>
    <div class="flex flex-col bg-white border border-gray-200 rounded-xl p-5 pb-4 my-8">
        <h1 class="text-lg font-bold text-gray-800 mb-4">Subscription</h1>
        <p v-if="res.id" class="text-sm text-gray-500 mb-3">
            <span v-if="isActive()" class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-teal-100 text-teal-800">Active</span>
            <span v-if="!isActive()" class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-gray-100 text-gray-500">Inactive</span>
        </p>
        <p v-if="res.id" class="text-sm text-gray-500 mb-3">Subscription ID:
            <span class="text-black font-semibold">{{ res.id }}</span>
        </p>
        <p v-if="error" class="text-red-500 text-sm mb-3">{{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { subscriptionApi } from '../api/subscription.ts'

const res = ref({
    id: '',
    active_until: ''
})
const error = ref('')

const getSubscription = async () => {
    try {
        const response = await subscriptionApi.get()
        res.value = response.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const isActive = () => {
    return res.value.active_until > new Date().toISOString()
}

onMounted(() => {
    getSubscription()
})
</script>