<template>
    <div>
        <button
            v-bind:data-hs-overlay="'#modal-verify-recipient' + recipient.id"
            type="button">
            Verify
        </button>
        <div v-bind:id="'modal-verify-recipient' + recipient.id"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3>Verify recipient</h3>
                        <button @click="close" type="button"
                            class="flex justify-center items-center size-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 dark:text-gray-100 hover:bg-gray-100 dark:hover:bg-neutral-700  disabled:opacity-50 disabled:pointer-events-none">
                            <span class="sr-only">Close</span>
                            <svg class="flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                stroke-linecap="round" stroke-linejoin="round">
                                <path d="M18 6 6 18"></path>
                                <path d="m6 6 12 12"></path>
                            </svg>
                        </button>
                    </div>
                    <div class="p-4 whitespace-normal text-left text-base">
                        <div class="mb-5">
                            <p>
                                We have sent a 6-digit OTP code to this recipient email address. Please enter the code below to verify the recipient email. Recipients with unconfirmed email address may be deleted after 7 days.
                            </p>
                        </div>
                        <div class="mb-3">
                            <label for="otp">
                                6-digit OTP code:
                            </label>
                            <input
                                v-model="req.otp"
                                v-bind:class="{ 'error': otpError }"
                                id="otp"
                                type="text"
                                pattern="[0-9]*"
                            >
                            <p v-if="otpError" class="error">Required</p>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-3 py-4 px-4 border-t dark:border-neutral-600">
                        <button
                            @click="verifyRecipient"
                            class="py-2 px-3 inline-flex items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600 disabled:opacity-50 disabled:pointer-events-none">
                            Verify
                        </button>
                        <button
                            @click="sendOtp"
                            class="text-gray-500 bg-gray-100 hover:bg-gray-200 dark:text-gray-300 dark:bg-neutral-600 dark:hover:bg-neutral-700 font-medium text-base py-2 px-3 focus:outline-none focus:shadow-outline"
                            type="submit">
                            Resend OTP
                        </button>
                        <button
                            @click="close"
                            class="text-gray-500 bg-gray-100 hover:bg-gray-200 dark:text-gray-300 dark:bg-neutral-600 dark:hover:bg-neutral-700 font-medium text-base py-2 px-3 focus:outline-none focus:shadow-outline">
                            Cancel
                        </button>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                        <p v-if="resendSuccess && !error" class="success px-5">{{ resendSuccess }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import overlay from '@preline/overlay'
import axios from 'axios'
import { recipientApi } from '../api/recipient.ts'
import events from '../events.ts'

const req = ref({
    otp: '',
})
const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const resendSuccess = ref('')
const error = ref('')
const otpError = ref(false)

const validateOtp = () => {
    otpError.value = !req.value.otp
    return !otpError.value
}

const verifyRecipient = async () => {
    if (!validateOtp()) return

    req.value.otp = req.value.otp + ''

    try {
        await recipientApi.activate(recipient.value.id, req.value)
        error.value = ''
        events.emit('recipient.verify', {})
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}

const sendOtp = async () => {
    try {
        const response = await recipientApi.sendOtp(recipient.value.id)
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

const close = () => {
    req.value.otp = ''
    resendSuccess.value = ''
    error.value = ''
    const modal = document.querySelector('#modal-verify-recipient' + recipient.value.id) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-verify-recipient' + recipient.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>