<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-edit-recipient' + recipient.id"
            class="text-bluish-500 hover:text-bluish-600 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
            type="submit">
            Edit
        </button>
        <div v-bind:id="'modal-edit-recipient' + recipient.id"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3 class="text-xl text-gray-800 dark:text-gray-100 font-semibold">
                            Edit recipient
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
                        <h1 class="text-xl text-gray-800 dark:text-gray-100 font-semibold mb-7">{{ recipient.email }}</h1>
                        <div class="mb-5">
                            <h3 class="font-semibold text-gray-800 dark:text-gray-100 mb-3">
                                PGP/Inline Encryption
                            </h3>
                            <p v-if="!recipient.pgp_key" class="text-gray-800 dark:text-gray-100 mb-3">
                                To use this option, please add a PGP key first.
                            </p>
                            <p class="text-gray-500 dark:text-gray-400 mb-3">
                                Enable this option to use PGP/Inline instead of the default PGP/MIME encryption for forwarded emails.
                            </p>
                        </div>
                        <div class="mb-5">
                            <input
                                v-bind:disabled="!recipient.pgp_key"
                                type="checkbox" v-bind:checked="recipient.pgp_inline" @change="updateRecipient"
                                class="form-checkbox relative w-11 h-6 p-px bg-gray-100 dark:bg-neutral-600 border-transparent text-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:ring-white dark:focus:ring-neutral-800 disabled:opacity-50 disabled:pointer-events-none checked:bg-none checked:text-bluish-500 checked:border-bluish-500 focus:ring-offset-transparent

                                before:inline-block before:size-5 before:bg-white dark:before:bg-neutral-400 checked:before:bg-bluish-200 before:translate-x-0 checked:before:translate-x-full before:rounded-full before:shadow before:transform before:transition before:ease-in-out before:duration-200">
                        </div>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="px-5 text-red-600 text-sm mb-5">Error: {{ error }}</p>
                        <p v-if="success" class="px-5 text-emerald-600 dark:text-emerald-500 text-sm mb-5">{{ success }}</p>
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
// const req = ref({
//     id: props.recipient.id,
//     pgp_enabled: props.recipient.pgp_enabled,
//     pgp_inline: props.recipient.pgp_inline
// })
const recipient = ref(props.recipient)
const error = ref('')
const success = ref('')

const updateRecipient = async () => {
    try {

        recipient.value.pgp_inline = !recipient.value.pgp_inline

        const req = {
            id: recipient.value.id,
            pgp_enabled: recipient.value.pgp_enabled,
            pgp_inline: recipient.value.pgp_inline
        }

        console.log(req)
        const res = await recipientApi.update(req)
        error.value = ''
        events.emit('recipient.update', {})
        success.value = res.data.message
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
    const modal = document.querySelector('#modal-edit-recipient' + recipient.value.id) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-edit-recipient' + recipient.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>