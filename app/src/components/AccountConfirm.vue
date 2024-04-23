<template>
    <div class="flex flex-col bg-white border border-gray-200 rounded-xl p-5 pb-4 my-8">
        <h1 class="text-lg font-bold text-gray-800 mb-4">Confirm Your Email</h1>
        <p class="text-sm text-gray-500 mb-3">
            We have sent a 6-digit OTP code to your email address. Please enter the code below to confirm your email. Accounts with unconfirmed email address may be deleted after 7 days.
        </p>
        <div v-if="!res.active" class="mb-4 max-w-xs">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="account-otp">
                6-digit OTP code:
            </label>
            <input
                v-model="otp"
                v-bind:class="{ 'border-red-500': otpError }"
                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-600 mb-4"
                id="account-otp" type="text">
            <p v-if="otpError" class="text-red-500 text-sm mb-2">Required field</p>
            <div class="flex flex-row gap-4">
                <button
                    @click="confirmEmail"
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Confirm
                </button>
                <button
                    @click="sendOtp"
                    class="text-gray-500 bg-gray-100 hover:bg-gray-200 font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Resend OTP
                </button>
            </div>
        </div>
        <p v-if="error" class="text-red-500 text-sm mb-3">{{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'

const res = ref({
    id: '',
    active: true
})
const otp = ref('')
const otpError = ref(false)
const success = ref('')
const error = ref('')

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
    if (!validateOtp()) {
        return
    }

    const req = {
        otp: otp.value
    }

    try {
        const response = await userApi.activate(req)
        success.value = response.data.message
        error.value = ''
        getUser()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.message
        }
    }
}

const sendOtp = async () => {
    try {
        const response = await userApi.sendOtp()
        success.value = response.data.message
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.message
        }
    }
}

const validateOtp = () => {
    otpError.value = !otp.value
    return !otpError.value
}

onMounted(() => {
    getUser()
})
</script>