<template>
    <div class="mt-3 mb-5">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="">
                <h2>Subscription</h2>
                <p class="mb-0">Status:</p>
                <p v-if="sub.id && !syncing">
                    <span v-if="isActive() && sub.id" class="badge success">Active</span>
                    <span v-if="!isActive() && sub.id" class="badge">Inactive</span>
                </p>
                <p v-if="syncing" class="text-sm">
                    <span v-if="isActive()" class="badge progress">Syncing...</span>
                </p>
                <p class="mb-0">Subscription active until:</p>
                <div v-if="isActive()" class="mb-3">
                    <p class="mb-3 text-primary">
                        {{ activeUntilDate() }}
                    </p>
                </div>
            </div>
            <div class="border-r border-transparent">
                <h2>Account Info</h2>
                <div class="mb-3">
                    <p class="mb-0">Mailx ID:</p>
                    <p class="mb-3 text-primary">
                        {{ email }}
                    </p>
                </div>
                <div class="mb-3">
                    <p class="mb-0">Email status:</p>
                    <p class="mb-3">
                        <span v-if="user.is_active && user.id" class="badge success">Verified</span>
                        <span v-if="!user.is_active && user.id" class="badge">Not verified</span>
                    </p>
                </div>
            </div>
        </div>
        <div v-if="isManaged()" class="card-tertiary">
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
        <div v-if="isLimited()" class="card-tertiary">
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
        <div v-if="isPendingDelete()" class="card-tertiary">
            <footer>
                <div>
                    <i class="icon info icon-primary"></i>
                </div>
                <div>
                    <h4>This account has been replaced and is scheduled for deletion.</h4>
                    <p>
                        A new Mailx signup was completed for your IVPN account, so this account will be deleted in 48 hours. Export any data you need before then.
                    </p>
                </div>
            </footer>
        </div>
        <div v-if="isOutage()" class="card-tertiary">
            <footer>
                <div>
                    <i class="icon info icon-primary"></i>
                </div>
                <div>
                    <h4>Out of sync</h4>
                    <p>
                        Your last account status update was {{ updatedAtDate() }}. <a target="_blank" :href="resyncUrl">Sync with IVPN</a>
                    </p>
                </div>
            </footer>
        </div>
        <p v-if="error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import tooltip from '@preline/tooltip'
import axios from 'axios'
import { subscriptionApi } from '../api/subscription.ts'
import { userApi } from '../api/user.ts'
import events from '../events.ts'

const sub = ref({
    id: '',
    updated_at: '',
    active_until: '',
    status: '',
    outage: false,
    type: '',
})
const user = ref({
    id: '',
    is_active: false
})
const error = ref('')
const success = ref('')
const email = ref(localStorage.getItem('email'))
const subid = ref('')
const sessionid = ref('')
const currentRoute = useRoute()
const syncing = ref(false)
const activateUrl = import.meta.env.VITE_RESYNC_URL
const resyncUrl = import.meta.env.VITE_RESYNC_URL + '?action=sync&service=mail'

const getSubscription = async () => {
    try {
        const res = await subscriptionApi.get()
        sub.value = res.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const getUser = async () => {
    try {
        const response = await userApi.get()
        user.value = response.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const updateSubscription = async () => {
    syncing.value = true
    try {
        const res = await subscriptionApi.update({
            id: sub.value.id,
            subid: subid.value,
        })
        success.value = res.data.message
        error.value = ''
        await getSubscription()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.response?.data.error || err.message
        }
    } finally {
        syncing.value = false
    }
}

const rotateSessionId = async () => {
    if (!sessionid.value) {
        return
    }

    syncing.value = true
    try {
        await subscriptionApi.rotateSessionId({
            sessionid: sessionid.value,
        })
        await getSubscription()
        await updateSubscription()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    } finally {
        syncing.value = false
    }
}

const isActive = () => {
    return sub.value.status === 'active' || sub.value.status === 'grace_period'
}

const isLimited = () => {
    return sub.value.status === 'limited_access'
}

const isPendingDelete = () => {
    return sub.value.status === 'pending_delete'
}

const isManaged = () => {
    return sub.value.type === 'Managed'
}

const activeUntilDate = () => {
    return new Date(sub.value.active_until).toDateString()
}

const updatedAtDate = () => {
    return new Date(sub.value.updated_at).toLocaleString()
}

const onUpdateEmail = (event: any) => {
    email.value = event.email
    getUser()
}

const isOutage = () => {
    return sub.value.outage
}

const parseParams = () => {
    const route = useRoute()
    const q = route.query
    const first = (v: unknown) => typeof v === 'string' ? v : Array.isArray(v) ? v[0] : ''
    subid.value = first(q.subid) || (route.params.subid as string) || ''
    sessionid.value = first(q.sessionid) || (route.params.sessionid as string) || ''

    if (!sessionid.value || !sessionid.value.match(/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/)) {
        return
    }

    rotateSessionId()
}

onMounted(() => {
    getSubscription()
    getUser()
    tooltip.autoInit()
    events.on('user.update', onUpdateEmail)
})

watch(currentRoute, () => {
    parseParams()
}, { immediate: true })
</script>