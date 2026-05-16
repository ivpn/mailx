<template>
    <div v-if="isManaged() && isDashboard && sub.id" class="card-tertiary md:m-8 md:mb-0 sm:mb-0 m-5">
        <footer>
            <div class="pt-1.5">
                <i class="icon info icon-primary"></i>
            </div>
            <div class="pt-1.5">
                <p>
                    Mailx beta ends May 19. To keep access, follow  <a target="_blank" :href="resyncUrl">this link</a> and sync with your IVPN account.
                </p>
            </div>
        </footer>
    </div>
    <div v-if="isLimited() && isDashboard && sub.id" class="card-tertiary md:m-8 md:mb-0 sm:mb-0 m-5">
        <footer>
            <div>
                <i class="icon info icon-primary"></i>
            </div>
            <div>
                <h4>Limited Access Mode</h4>
                <p>
                    Existing aliases forward normally. New aliases are disabled. Add time to your <a target="_blank" :href="activateUrl">IVPN account</a> to restore access.
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
    type: '',
})

const route = ref('/')
const currentRoute = useRoute()
const props = defineProps(['dashboard'])
const isDashboard = props.dashboard
const activateUrl = import.meta.env.VITE_RESYNC_URL
const resyncUrl = import.meta.env.VITE_RESYNC_URL + '?action=sync&service=mail'

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

const isManaged = () => {
    return sub.value.type === 'Managed'
}

onMounted(() => {
    getSubscription()
})

watch(currentRoute, (newRoute) => {
    route.value = newRoute.path
}, { immediate: true })
</script>