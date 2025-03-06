<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-verify-recipient' + recipient.id">
            Verify
        </button>
        <div v-bind:id="'modal-verify-recipient' + recipient.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <h3>Verify recipient</h3>
                        <button @click="close" class="close">
                            <span class="sr-only">Close</span>
                            <svg class="flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                stroke-linecap="round" stroke-linejoin="round">
                                <path d="M18 6 6 18"></path>
                                <path d="m6 6 12 12"></path>
                            </svg>
                        </button>
                    </header>
                    <article>
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
                    </article>
                    <footer>
                        <nav>
                            <button @click="verifyRecipient" class="pcta">
                                Verify
                            </button>
                            <button @click="sendOtp" class="cta cancel">
                                Resend OTP
                            </button>
                            <button @click="close" class="cta cancel">
                                Cancel
                            </button>
                        </nav>
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                        <p v-if="resendSuccess && !error" class="success px-5">{{ resendSuccess }}</p>
                    </footer>
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