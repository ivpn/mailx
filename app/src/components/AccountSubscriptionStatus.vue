<template>
    <div v-if="!isActive && isDashboard" class="card-secondary m-8 mb-0">
        <p class="m-0">
            Your <router-link to="/account">subscription</router-link> is inactive
        </p>
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
const props = defineProps(['dashboard'])
const isDashboard = props.dashboard

const getSubscription = async () => {
    try {
        const response = await subscriptionApi.get()
        res.value = response.data
        isActive.value = res.value.active_until > new Date().toISOString()
    } catch (err) {
    }
}

onMounted(() => {
    getSubscription()
})

watch(currentRoute, (newRoute) => {
    route.value = newRoute.path
}, { immediate: true })
</script>