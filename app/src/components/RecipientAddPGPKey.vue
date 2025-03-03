<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-add-key-recipient' + recipient.id"
            class="text-bluish-500 hover:text-bluish-600 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
            type="submit">
            Add PGP key
        </button>
        <div v-bind:id="'modal-add-key-recipient' + recipient.id"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div
                    class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3 class="text-xl text-gray-800 dark:text-gray-100 font-semibold">
                            Add PGP Public Key
                        </h3>
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
                        <div class="mb-3">
                            <label for="recipient_pgp" class="block text-gray-500 dark:text-gray-400 mb-3">
                                Enter your public PGP key:
                            </label>
                            <textarea
                                v-model="pgp_key"
                                v-bind:class="{ 'border-gray-500 dark:border-neutral-400': !pgpError, 'border-red-600 dark:border-red-600': pgpError }"
                                id="recipient_pgp"
                                placeholder="Starts with '-----BEGIN PGP PUBLIC KEY BLOCK-----'"
                                class="appearance-none outline-none border w-full py-3 px-4 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 leading-tight focus:border-bluish-500 mb-2 h-[350px]"
                            >
                            </textarea>
                            <p v-if="pgpError" class="text-red-600 text-sm">Required</p>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-3 py-4 px-4 border-t dark:border-neutral-600">
                        <button @click="addKey"
                            class="py-2 px-3 inline-flex items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600 disabled:opacity-50 disabled:pointer-events-none">
                            Add Key
                        </button>
                        <button @click="close"
                            class="text-gray-500 bg-gray-100 hover:bg-gray-200 dark:text-gray-300 dark:bg-neutral-600 dark:hover:bg-neutral-700 font-medium text-base py-2 px-3 focus:outline-none focus:shadow-outline">
                            Cancel
                        </button>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="px-5 text-red-600 text-sm mb-5">Error: {{ error }}</p>
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

const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const pgp_key = ref('')
const error = ref('')
const pgpError = ref(false)

const validatePgp = () => {
    pgpError.value = !pgp_key.value
    return !pgpError.value
}

const addKey = async () => {
    if (!validatePgp()) {
        return
    }

    const payload = {
        id: recipient.value.id,
        pgp_enabled: true,
        pgp_key: pgp_key.value
    }

    try {
        await recipientApi.update(payload)
        error.value = ''
        events.emit('recipient.update', {})
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

const close = () => {
    error.value = ''
    pgpError.value = false
    const modal = document.querySelector('#modal-add-key-recipient' + recipient.value.id) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-add-key-recipient' + recipient.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>