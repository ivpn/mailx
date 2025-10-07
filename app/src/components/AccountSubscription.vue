<template>
    <div class="mb-5">
        <h2>Account</h2>
        <p v-if="sub.id" class="text-sm">
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
        <div v-if="isLimited()" class="card-tertiary">
            <footer>
                <div>
                    <i class="icon info icon-primary"></i>
                </div>
                <div>
                    <h4>Limited Access Mode</h4>
                    <p>
                        Your MailX account is in limited access mode. To regain full access add time to your <a href="https://www.ivpn.net/account/">IVPN account</a>.
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
                        Your account is pending deletion. To reinstate access add time to your <a href="https://www.ivpn.net/account/">IVPN account</a>.
                    </p>
                </div>
            </footer>
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

const sub = ref({
    id: '',
    active_until: '',
    is_active: false,
    is_grace_period: false,
})
const error = ref('')
const email = ref(localStorage.getItem('email'))

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

const onUpdateEmail = (event: any) => {
    email.value = event.email
}

onMounted(() => {
    getSubscription()
    tooltip.autoInit()
    events.on('user.update', onUpdateEmail)
})
</script>