<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-totp-enable'"
            class="py-2 px-3 font-medium bg-bluish-500 text-white hover:bg-bluish-600">
            Enable
        </button>
        <div v-bind:id="'modal-totp-enable'"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3>Enable 2-Factor Authentication</h3>
                        <button @click="close" type="button"
                            class="flex justify-center items-center size-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 dark:text-gray-100 hover:bg-gray-100 dark:hover:bg-neutral-700 disabled:opacity-50 disabled:pointer-events-none">
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
                        <div v-if="!isEnabled">
                            <div class="mb-5">
                                <p>
                                    To enable two-factor authentication, please scan the code with a TOTP app (for example: Google Authenticator) and enter the code in the field below.
                                </p>
                                <p>
                                    If you cannot scan QR code, you can enter the following information manually. Secret: {{ resEnable.secret }}, Account: {{ resEnable.account }}
                                </p>
                            </div>
                            <div class="mb-5 container">
                                <canvas class="mx-auto" id="totp_qr_code"></canvas>
                            </div>
                            <div class="mb-3">
                                <label for="totp_enable_code">
                                    Code from TOTP app:
                                </label>
                                <input
                                    v-model="req.otp"
                                    v-bind:class="{ 'error': codeError }"
                                    id="totp_enable_code"
                                    placeholder="6-digit code"
                                    type="text"
                                    pattern="[0-9]*"
                                >
                                <p v-if="codeError" class="error">Required</p>
                            </div>
                            <div class="flex justify-start items-center gap-x-3 pt-4 border-t dark:border-neutral-600">
                                <button
                                    @click="totpEnableConfirm"
                                    class="py-2 px-3 inline-flex items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600 disabled:opacity-50 disabled:pointer-events-none">
                                    Enable
                                </button>
                                <button
                                    @click="close"
                                    class="text-gray-500 bg-gray-100 hover:bg-gray-200 dark:text-gray-300 dark:bg-neutral-600 dark:hover:bg-neutral-700 font-medium text-base py-2 px-3 focus:outline-none focus:shadow-outline">
                                    Cancel
                                </button>
                            </div>
                        </div>
                        <div v-if="isEnabled">
                            <p>
                                Two-factor authentication was set up successfully.
                            </p>
                            <p>
                                Please record the following backup codes which you will be able to use instead of TOTP in case you lost access to your device.
                            </p>
                            <p class="py-4 px-5 bg-gray-100 dark:bg-neutral-800">
                                Backup codes:
                                <span class="text-gray-800 dark:text-gray-100">
                                    {{ resConfirm.backup }}
                                </span>
                            </p>
                            <p>
                                Each of these codes can be used only once.
                            </p>
                            <div class="flex justify-start items-center gap-x-3 pt-4 border-t dark:border-neutral-600">
                                <button
                                    @click="complete"
                                    class="text-gray-500 bg-gray-100 hover:bg-gray-200 dark:text-gray-300 dark:bg-neutral-600 dark:hover:bg-neutral-700 font-medium text-base py-2 px-3 focus:outline-none focus:shadow-outline">
                                    Close
                                </button>
                            </div>
                        </div>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="px-5 error">Error: {{ error }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import overlay from '@preline/overlay'
import QRious from 'qrious'
import events from '../events.ts'

const req = ref({ otp: '' })
const resEnable = ref({ uri: '', secret: '', account: '' })
const resConfirm = ref({ backup: '' })
const error = ref('')
const codeError = ref(false)
const isEnabled = ref(false)

const close = () => {
    req.value = {} as any
    error.value = ''
    const modal = document.querySelector('#modal-totp-enable') as any
    overlay.close(modal)
}

const complete = () => {
    events.emit('totp.enable', {})
    close()
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-totp-enable' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        totpEnable()
    })
}

const totpEnable = async () => {
    try {
        const res = await userApi.totpEnable()
        resEnable.value = res.data
        generateQRCode(resEnable.value.uri)
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const totpEnableConfirm = async () => {
    if (!req.value.otp) {
        codeError.value = true
        return
    }

    req.value.otp = req.value.otp + ''

    try {
        const res = await userApi.totpEnableConfirm(req.value)
        resConfirm.value = res.data
        isEnabled.value = true
        codeError.value = false
    } catch (err) {
        if (axios.isAxiosError(err)) {
            resConfirm.value = { backup: '' }
            isEnabled.value = false
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}

const generateQRCode = (uri: string) => {
    new QRious({
        element: document.getElementById('totp_qr_code'),
        value: uri,
        size: 150,
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>