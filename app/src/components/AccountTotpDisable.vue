<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-totp-disable'"
            class="py-2 px-3 font-medium bg-bluish-500 text-white hover:bg-bluish-600">
            Disable
        </button>
        <div v-bind:id="'modal-totp-disable'"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white border shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b">
                        <h3 class="text-xl text-gray-800 font-semibold">
                            Disable 2-Factor Authentication
                        </h3>
                        <button @click="close" type="button"
                            class="flex justify-center items-center size-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none">
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
                            <p class="text-gray-500 mb-3">
                                To disable two-factor authentication, please enter code from TOTP app or a backup code.
                            </p>
                        </div>
                        <div class="mb-3">
                            <label for="totp_disable_code" class="block text-gray-500 mb-3">
                                Code from TOTP app:
                            </label>
                            <input
                                v-model="req.otp"
                                v-bind:class="{ 'border-gray-500': !codeError, 'border-red-600': codeError }"
                                id="totp_disable_code"
                                placeholder="6-digit code"
                                class="appearance-none outline-none border w-full py-3 px-4 text-gray-500 leading-tight focus:border-bluish-500 mb-2"
                                type="number">
                                <p v-if="codeError" class="text-red-600 text-sm">Required</p>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-3 py-4 px-4 border-t">
                        <button @click="disableTotp"
                            class="py-2 px-3 inline-flex items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600 disabled:opacity-50 disabled:pointer-events-none">
                            Disable
                        </button>
                        <button @click="close"
                            class="text-gray-500 bg-gray-100 hover:bg-gray-200 font-medium text-base py-2 px-3 focus:outline-none focus:shadow-outline">
                            Cancel
                        </button>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="px-5 text-red-600 text-sm mb-3">Error: {{ error }}</p>
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

const emit = defineEmits(['onTotpDisable'])

const req = ref({ otp: '' })
const error = ref('')
const codeError = ref(false)

const close = () => {
    req.value = { otp: '' }
    error.value = ''
    const modal = document.querySelector('#modal-totp-disable') as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-totp-disable' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

const disableTotp = async () => {
    if (!req.value.otp) {
        codeError.value = true
        return
    }

    req.value.otp = req.value.otp + ''

    try {
        await userApi.totpDisable(req.value)
        codeError.value = false
        error.value = ''
        emit('onTotpDisable')
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>