<template>
    <div v-if="!res.is_active && isDashboard" class="flex flex-col items-center text-center my-14 pt-6">
        <h3>Verify Your Email</h3>
        <p>
            Please <router-link to="/account">verify</router-link> your account email address.
        </p>
    </div>
    <div v-if="!res.is_active && !isDashboard" class="flex flex-col p-5 pb-4 my-8 bg-white dark:bg-neutral-800">
        <h1>Verify Your Email</h1>
        <p>
            We have sent a 6-digit OTP code to your email address. Please enter the code below to verify your account email. Accounts with unconfirmed email address may be deleted after 7 days.
        </p>
        <div v-if="!confirmSuccess" class="mb-4 max-w-xs">
            <div class="mb-4">
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
            <div class="flex flex-row gap-4">
                <button @click="confirmEmail" class="cta">
                    Verify
                </button>
                <button @click="sendOtp" class="cta cancel">
                    Resend OTP
                </button>
            </div>
        </div>
        <p v-if="error" class="error">Error: {{ error }}</p>
        <p v-if="resendSuccess && !error && !confirmSuccess" class="success">{{ resendSuccess }}</p>
        <p v-if="confirmSuccess" class="text-sm my-3">
            <span class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-emerald-100 text-emerald-800 dark:bg-emerald-800 dark:text-emerald-100">{{ confirmSuccess }}</span>
        </p>
    </div>
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
                error.value = 'Too many requests, please try again later'
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
                error.value = 'Too many requests, please try again later'
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