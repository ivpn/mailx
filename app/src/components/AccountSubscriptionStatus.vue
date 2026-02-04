<template>
    <div v-if="isLimited() && isDashboard && sub.id" class="card-tertiary m-8 mb-0">
        <footer>
            <div>
                <i class="icon info icon-primary"></i>
            </div>
            <div>
                <h4>Limited Access Mode</h4>
                <p>
                    Your MailX account is in limited access mode. To regain full access add time to your <a target="_blank" :href="activateUrl">IVPN account</a>.
                </p>
            </div>
        </footer>
    </div>
    <div v-if="isPendingDelete() && isDashboard && sub.id" class="card-tertiary m-8 mb-0">
        <footer>
            <div>
                <i class="icon info icon-primary"></i>
            </div>
            <div>
                <h4>Pending Deletion</h4>
                <p>
                    Your account is pending deletion. To reinstate access add time to your <a target="_blank" :href="activateUrl">IVPN account</a>.
                </p>
            </div>
        </footer>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { subscriptionApi } from '../api/subscription.ts'

const sub = ref({
    id: '',
    updated_at: '',
    active_until: '',
    status: '',
    outage: false,
})

const route = ref('/')
const currentRoute = useRoute()
const props = defineProps(['dashboard'])
const isDashboard = props.dashboard
const activateUrl = import.meta.env.VITE_RESYNC_URL

const getSubscription = async () => {
    try {
        const res = await subscriptionApi.get()
        sub.value = res.data
    } catch (err) {
    }
}

const isLimited = () => {
    return sub.value.status === 'limited_access'
}

const isPendingDelete = () => {
    return sub.value.status === 'pending_delete'
}

onMounted(() => {
    getSubscription()
})

watch(currentRoute, (newRoute) => {
    route.value = newRoute.path
}, { immediate: true })
</script>