<template>
    <div class="mb-5">
        <h2>Delete Account</h2>
        <button
            v-bind:data-hs-overlay="'#modal-delete-account'"
            class="cta delete mb-4">
            Delete Account
        </button>
        <div v-bind:id="'modal-delete-account'" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>DELETE ACCOUNT</h4>
                    </header>
                    <article>
                        <div>
                            <div class="mb-5">
                                <p>
                                    <strong>WARNING:</strong> This operation cannot be undone. Deleting your account will permanently remove data from our systems.
                                </p>
                            </div>
                            <div class="mb-5">
                                <p>
                                    To confirm permanent deletion of your account, please enter the following symbols in the field below:
                                    <span class="text-black dark:text-white">{{ otp }}</span>
                                </p>
                            </div>
                            <div class="mb-7">
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
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click.stop="promptDeleteAccount" class="cta delete">
                                Delete Account
                            </button>
                            <button @click="close" class="cancel">
                                Cancel
                            </button>
                        </nav>
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                    </footer>
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
                error.value = 'Too many requests, please try again later.'
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
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const close = () => {
    req.value = {} as any
    otp.value = ''
    error.value = ''
    otpError.value = false
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