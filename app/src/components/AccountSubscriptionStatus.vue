<template>
    <div v-if="!isActive() && !isAccountRoute()" class="flex flex-col p-5 my-8 bg-white dark:bg-neutral-800">
        <p class="text-gray-500 dark:text-gray-400 mb-2">Account subscription is inactive</p>
        <router-link class="text-bluish-500 hover:text-bluish-600" to="/account">View Details</router-link>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { subscriptionApi } from '../api/subscription.ts'

const res = ref({
    id: '',
    active_until: ''
})

const route = ref('/')
const currentRoute = useRoute()

const getSubscription = async () => {
    try {
        const response = await subscriptionApi.get()
        res.value = response.data
    } catch (err) {
    }
}

const isActive = () => {
    return res.value.active_until > new Date().toISOString()
}

const isAccountRoute = () => {
    return route.value === '/account'
}

onMounted(() => {
    getSubscription()
})

watch(currentRoute, (newRoute) => {
    route.value = newRoute.path
}, { immediate: true })
</script>