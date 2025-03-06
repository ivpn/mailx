<template>
    <div>
        <h2>Delete Account</h2>
        <button
            v-bind:data-hs-overlay="'#modal-delete-account'"
            class="delete">
            Delete Account
        </button>

        <div v-bind:id="'modal-delete-account'"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3>Delete Account</h3>
                        <button
                            @click="close"
                            class="close">
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
                        <div>
                            <div class="mb-5">
                                <p>
                                    <strong>Warning:</strong> this operation cannot be undone. Deleting your account will permanently remove data from our systems.
                                </p>
                            </div>
                            <div class="mb-5">
                                <p>
                                    To confirm permanent deletion of your account, please enter the following symbols in the field below:
                                    <span class="text-black dark:text-white">{{ otp }}</span>
                                </p>
                            </div>
                            <div class="mb-5">
                                <input
                                    v-model="req.otp"
                                    v-bind:class="{ 'error': otpError }"
                                    id="totp_enable_code"
                                    placeholder="8-symbol code"
                                    type="text"
                                    pattern="[0-9]*"
                                >
                                <p v-if="otpError" class="error">Required</p>
                            </div>
                            <div class="flex justify-start items-center gap-x-3 pt-4 border-t dark:border-neutral-600">
                                <button
                                    @click.stop="promptDeleteAccount"
                                    class="cta delete">
                                    Delete Account
                                </button>
                                <button
                                    @click="close"
                                    class="cta cancel">
                                    Cancel
                                </button>
                            </div>
                        </div>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
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

const req = ref({ otp: '' })
const otp = ref('')
const otpError = ref(false)
const error = ref('')

const validateOtp = () => {
    otpError.value = !req.value.otp
    return !otpError.value
}

const promptDeleteAccount = () => {
    if (!validateOtp()) return
    if (!confirm('Are you sure you want to delete your account? This action cannot be undone.')) return
    deleteAccount()
}

const deleteAccount = async () => {
    try {
        await userApi.delete(req.value)
        alert('Account is deleted successfully. You will be logged out.')
        userApi.clearSession()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}

const deleteAccountRequest = async () => {
    try {
        const res = await userApi.deleteRequest()
        otp.value = res.data.otp
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}

const close = () => {
    req.value = {} as any
    otp.value = ''
    error.value = ''
    const modal = document.querySelector('#modal-delete-account') as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-delete-account' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        deleteAccountRequest()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>