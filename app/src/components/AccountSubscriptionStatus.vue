<template>
    <div v-if="!isActive && !isAccountRoute()" class="flex flex-col p-5 my-8 bg-white dark:bg-neutral-800">
        <p class="mb-2">Account subscription is inactive</p>
        <router-link to="/account">View Details</router-link>
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
const isActive = ref(true)

const getSubscription = async () => {
    try {
        const response = await subscriptionApi.get()
        res.value = response.data
        isActive.value = res.value.active_until > new Date().toISOString()
    } catch (err) {
    }
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