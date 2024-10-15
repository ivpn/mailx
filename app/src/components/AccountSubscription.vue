<template>
    <div class="mb-5">
        <h2 class="text-2xl font-semibold dark:text-gray-100 mb-5">Subscription</h2>
        <p v-if="res.id" class="text-sm text-gray-500 mb-5">
            <span v-if="isActive()"
                class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-emerald-100 text-emerald-800 dark:bg-emerald-800 dark:text-emerald-100">Active</span>
            <span v-if="!isActive()"
                class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-gray-100 text-gray-500 dark:bg-gray-500 dark:text-gray-100">Inactive</span>
        </p>
        <div v-if="isActive()" class="mb-3">
            <h2 class="font-semibold text-gray-800 dark:text-gray-100 mb-3">
                Active until:
            </h2>
            <p class="text-gray-500 dark:text-gray-400 mb-3">
                {{ activeUntilDate() }}
            </p>
        </div>
        <div class="mb-3">
            <h2 class="font-semibold text-gray-800 dark:text-gray-100 mb-3">
                Subscription ID:
            </h2>
            <div class="hs-tooltip text-gray-500 dark:text-gray-400 mb-3">
                <span class="hs-tooltip-toggle">
                    <button @click="copyAlias(res.id)">
                        {{ res.id }}
                    </button>
                    <span
                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
                        role="tooltip">
                        {{ copyText }}
                    </span>
                </span>
            </div>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import tooltip from '@preline/tooltip'
import axios from 'axios'
import { subscriptionApi } from '../api/subscription.ts'

const res = ref({
    id: '',
    active_until: ''
})
const error = ref('')
const copyText = ref('Click to copy')

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

const activeUntilDate = () => {
    return new Date(res.value.active_until).toDateString()
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias)
    copyText.value = 'Copied!'
    setTimeout(() => {
        copyText.value = 'Click to copy'
    }, 2000)
}

onMounted(() => {
    getSubscription()
    tooltip.autoInit()
})
</script>