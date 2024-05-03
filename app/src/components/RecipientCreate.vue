<template>
    <div>
        <button v-bind:data-hs-overlay="'#hs-modal-create-recipient'"
            class="mt-3 py-2 pl-2 pr-3 inline-flex justify-center items-center gap-x-2 text-sm font-medium rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700">
            <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                stroke-linejoin="round">
                <path d="M5 12h14"></path>
                <path d="M12 5v14"></path>
            </svg>
            Add Recipient
        </button>
        <div v-bind:id="'hs-modal-create-recipient'"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white border shadow-sm rounded-xl pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b">
                        <h3 class="font-bold text-gray-800">
                            Add Recipient
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
                            <p class="text-sm text-gray-500 mb-3">
                                Add a email address to receive forwarded emails. A 6 digit verification code will be sent to this email address. If expired, you can resend the verification code.
                            </p>
                        </div>
                        <div class="mb-3">
                            <label for="recipient_email" class="block text-gray-500 text-sm font-semibold mb-3">
                                Email:
                            </label>
                            <input
                                v-model="recipient.email"
                                v-bind:class="{ 'border-red-600': emailError }"
                                id="recipient_email"
                                placeholder="name@example.net"
                                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-violet-600 mb-2"
                                type="text">
                            <p v-if="emailError" class="text-red-600 text-sm">Required field</p>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-2 py-3 px-4 border-t">
                        <button @click="postRecipient"
                            class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-medium rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700 disabled:opacity-50 disabled:pointer-events-none">
                            Add Recipient
                        </button>
                        <button @click="close"
                            class="text-gray-500 bg-gray-100 hover:bg-gray-200 font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline">
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
import overlay from '@preline/overlay'
import axios from 'axios'
import { recipientApi } from '../api/recipient.ts'

const emit = defineEmits(['onCreateRecipient'])
const recipient = ref({
    email: '',
})
const error = ref('')
const emailError = ref(false)

const validateEmail = () => {
    emailError.value = !recipient.value.email
    return !emailError.value
}

const postRecipient = async () => {
    if (!validateEmail()) {
        return
    }

    try {
        await recipientApi.create(recipient.value)
        error.value = ''
        emit('onCreateRecipient')
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const close = () => {
    recipient.value = {} as any
    error.value = ''
    const modal = document.querySelector('#hs-modal-create-recipient')
    if (modal instanceof HTMLElement) {
        overlay.close(modal)
    }
}

onMounted(() => {
    overlay.autoInit()
})
</script>