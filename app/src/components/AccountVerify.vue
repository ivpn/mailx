<template>
    <div v-if="!res.is_active && isDashboard" class="card-secondary m-8 mb-0">
        <p class="m-0"><router-link to="/account">Verify</router-link> your primary email address to start receiving forwarded mail</p>
    </div>
    <div v-if="!res.is_active && !isDashboard">
        <h2>Verify Your Email</h2>
        <p>
            Enter the 6-digit OTP code we have sent via email to verify your primary address. You have 7 days to complete this step after signing up to prevent account deletion.
        </p>
        <div v-if="!confirmSuccess" class="mb-9 max-w-xs">
            <div class="mb-5">
                <label for="account-otp">
                    6-digit OTP code:
                </label>
                <input
                    v-model="otp"
                    v-bind:class="{ 'error': otpError }"
                    id="account-otp"
                    type="text"
                    pattern="[0-9]*"
                >
                <p v-if="otpError" class="error">Required</p>
            </div>
            <div class="flex flex-row gap-2">
                <button @click="confirmEmail" class="cta">
                    Verify
                </button>
                <button @click="sendOtp" class="cancel">
                    Resend OTP
                </button>
            </div>
            <p v-if="error" class="error my-5">Error: {{ error }}</p>
            <p v-if="resendSuccess && !error && !confirmSuccess" class="success my-5">{{ resendSuccess }}</p>
            
        </div>
        <p v-if="confirmSuccess" class="text-sm mb-9">
            <span class="badge success">{{ confirmSuccess }}</span>
        </p>
    </div>
    <hr v-if="!res.is_active && !isDashboard">
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'
import events from '../events.ts'

const res = ref({
    id: '',
    is_active: true
})
const otp = ref('')
const otpError = ref(false)
const confirmSuccess = ref('')
const resendSuccess = ref('')
const error = ref('')
const props = defineProps(['dashboard'])
const isDashboard = props.dashboard

const getUser = async () => {
    try {
        const response = await userApi.get()
        res.value = response.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const confirmEmail = async () => {
    if (!validateOtp()) return

    const req = {
        otp: otp.value + ''
    }

    try {
        const response = await userApi.activate(req)
        confirmSuccess.value = response.data.message
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            confirmSuccess.value = ''
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const sendOtp = async () => {
    try {
        const response = await userApi.sendOtp()
        resendSuccess.value = response.data.message
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            resendSuccess.value = ''
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const validateOtp = () => {
    otpError.value = !otp.value
    return !otpError.value
}

const onUserUpdate = () => {
    getUser()
}

onMounted(() => {
    getUser()
    events.on('user.update', onUserUpdate)
})
</script>