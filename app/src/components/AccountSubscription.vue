<template>
    <div class="mb-5">
        <h2>Account</h2>
        <p v-if="sub.id && !syncing" class="text-sm">
            <span v-if="isActive()" class="badge success">Active</span>
            <span v-if="!isActive()" class="badge">Inactive</span>
        </p>
        <p v-if="syncing" class="text-sm">
            <span v-if="isActive()" class="badge progress">Syncing...</span>
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
        <div v-if="isLimited()" class="card-tertiary">
            <footer>
                <div>
                    <i class="icon info icon-primary"></i>
                </div>
                <div>
                    <h4>Limited Access Mode</h4>
                    <p>
                        Your MailX account is in limited access mode. To regain full access add time to your <a target="_blank" href="https://www.ivpn.net/account/">IVPN account</a>.
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
                    <h4>Pending Deletion</h4>
                    <p>
                        Your account is pending deletion. To reinstate access add time to your <a target="_blank" href="https://www.ivpn.net/account/">IVPN account</a>.
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
                        Your account status last update was at {{ updatedAtDate() }}. <a target="_blank" href="https://www.ivpn.net/account/">Sync with IVPN</a>
                    </p>
                </div>
            </footer>
        </div>
        <p v-if="error" class="error">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import tooltip from '@preline/tooltip'
import axios from 'axios'
import { subscriptionApi } from '../api/subscription.ts'
import events from '../events.ts'

const sub = ref({
    id: '',
    active_until: '',
    is_active: false,
    updated_at: '',
    is_grace_period: false,
    outage: false,
})
const error = ref('')
const email = ref(localStorage.getItem('email'))
const subid = ref('')
const preauthid = ref('')
const preauthtokenhash = ref('')
const currentRoute = useRoute()
const syncing = ref(false)

const getSubscription = async () => {
    try {
        const res = await subscriptionApi.get()
        sub.value = res.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const updateSubscription = async () => {
    if (!subid.value || !preauthid.value || !preauthtokenhash.value) {
        return
    }

    syncing.value = true
    try {
        await subscriptionApi.update({
            id: sub.value.id,
            subid: subid.value,
            preauthid: preauthid.value,
            preauthtokenhash: preauthtokenhash.value,
        })
        await getSubscription()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    } finally {
        syncing.value = false
    }
}

const isActive = () => {
    return sub.value.active_until > new Date().toISOString()
}

const isLimited = () => {
    return sub.value.is_grace_period && !isActive()
}

const isPendingDelete = () => {
    return !sub.value.is_grace_period && !isActive()
}

const activeUntilDate = () => {
    return new Date(sub.value.active_until).toDateString()
}

const updatedAtDate = () => {
    return new Date(sub.value.updated_at).toLocaleString()
}

const onUpdateEmail = (event: any) => {
    email.value = event.email
}

const isOutage = () => {
    return sub.value.outage
}

const parseParams = () => {
    const route = useRoute()
    const q = route.query
    const first = (v: unknown) => typeof v === 'string' ? v : Array.isArray(v) ? v[0] : ''
    subid.value = first(q.subid) || (route.params.subid as string) || ''
    preauthid.value = first(q.preauthid) || (route.params.preauthid as string) || ''
    preauthtokenhash.value = first(q.preauthtokenhash) || (route.params.preauthtokenhash as string) || ''
    preauthtokenhash.value = preauthtokenhash.value.replace(/ /g, '+')

    if (!subid.value || !subid.value.match(/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/)) {
        return
    }

    if (!preauthid.value || !preauthid.value.match(/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/)) {
        return
    }

    if (!preauthtokenhash.value) {
        return
    }

    updateSubscription()
}

onMounted(() => {
    getSubscription()
    tooltip.autoInit()
    events.on('user.update', onUpdateEmail)
})

watch(currentRoute, () => {
    parseParams()
}, { immediate: true })
</script>