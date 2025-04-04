<template>
    <div class="mb-5">
        <h2>Account</h2>
        <p v-if="res.id" class="text-sm">
            <span v-if="isActive()" class="badge success">Active</span>
            <span v-if="!isActive()" class="badge">Inactive</span>
        </p>
        <div class="mb-3">
            <h4>Account email:</h4>
            <p class="mb-3">
                {{ email }}
            </p>
        </div>
        <div v-if="isActive()" class="mb-3">
            <h4>Subscription active until:</h4>
            <p class="mb-3">
                {{ activeUntilDate() }}
            </p>
        </div>
        <div class="mb-3">
            <h4>Subscription ID:</h4>
            <div class="hs-tooltip mb-3">
                <span class="hs-tooltip-toggle">
                    <button @click="copyAlias(res.id)" class="plain">
                        {{ res.id }}
                    </button>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                        {{ copyText }}
                    </span>
                </span>
            </div>
        </div>
        <p v-if="error" class="error">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import tooltip from '@preline/tooltip'
import axios from 'axios'
import { subscriptionApi } from '../api/subscription.ts'
import events from '../events.ts'

const res = ref({
    id: '',
    active_until: ''
})
const error = ref('')
const copyText = ref('Click to copy')
const email = ref(localStorage.getItem('email'))

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

const onUpdateEmail = (event: any) => {
    email.value = event.email
}

onMounted(() => {
    getSubscription()
    tooltip.autoInit()
    events.on('user.update', onUpdateEmail)
})
</script>