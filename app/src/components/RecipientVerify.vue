<template>
    <div>
        <div v-bind:id="'modal-verify-recipient' + recipient.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>VERIFY RECIPIENT</h4>
                    </header>
                    <article>
                        <div class="mb-5">
                            <p>
                                We have sent a 6-digit OTP code to your new recipient email address. Enter the code below to complete the verification.
                            </p>
                        </div>
                        <div class="mb-5">
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
                            <button @click="verifyRecipient" class="cta">
                                Verify Recipient
                            </button>
                            <button @click="sendOtp" class="cancel">
                                Resend OTP
                            </button>
                            <button @click="close" class="cancel">
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
                error.value = 'Too many requests, please try again later.'
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
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const close = () => {
    req.value.otp = ''
    resendSuccess.value = ''
    error.value = ''
    otpError.value = false
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