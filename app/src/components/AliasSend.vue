<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-send-alias' + alias.id"
            class="text-blue-600 hover:text-blue-700 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
            type="submit">
            Send
        </button>
        <div v-bind:id="'modal-send-alias' + alias.id"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white border shadow-sm rounded-xl pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b">
                        <h3 class="font-bold text-gray-800">
                            Send from alias
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
                            <p class="text-gray-800 mb-3">
                                Generate the proper email address to send a message from this alias. Note that to send emails using an alias, you need to do so from an account-verified recipient.
                            </p>
                        </div>
                        <div class="mb-5">
                            <label v-bind:for="'from_alias_' + alias.id"
                                class="block text-gray-500 text-sm font-semibold mb-3">
                                From alias:
                            </label>
                            <input v-bind:id="'from_alias_' + alias.id" v-bind:value="alias.name" disabled
                                class="disabled appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-600 mb-2"
                                type="text">
                        </div>
                        <div class="mb-5">
                            <label v-bind:for="'to_email_' + alias.id"
                                class="block text-gray-500 text-sm font-semibold mb-3">
                                To email:
                            </label>
                            <input
                                v-bind:id="'to_email_' + alias.id"
                                v-bind:class="{ 'border-red-600': emailError }"
                                v-model="toEmail"
                                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-600 mb-2"
                                type="text">
                            <p v-if="emailError" class="text-red-600 text-sm">Required field</p>
                        </div>
                        <div v-bind:class="{ 'hidden': generatedEmail == '' }" class="mb-5">
                            <p class="text-gray-500 text-sm font-semibold mb-3">
                                Send message to this email:
                            </p>
                            <div class="hs-tooltip text-gray-800 mb-3">
                                <span class="hs-tooltip-toggle">
                                    <button @click="copy(generatedEmail)">
                                        {{ generatedEmail }}
                                    </button>
                                    <span
                                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 text-xs font-medium text-white rounded shadow-sm"
                                        role="tooltip">
                                        {{ copyText }}
                                    </span>
                                </span>
                            </div>

                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-2 py-3 px-4 border-t">
                        <button @click="showAddress"
                            class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-medium rounded-md bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none">
                            Show address
                        </button>
                        <button @click="close"
                        class="text-gray-500 bg-gray-100 hover:bg-gray-200 font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline">
                            Close
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import overlay from '@preline/overlay'
import tooltip from '@preline/tooltip'

const props = defineProps(['alias'])
const alias = ref(props.alias)
const toEmail = ref('')
const generatedEmail = ref('')
const emailError = ref(false)
const copyText = ref('Click to copy')

const validateEmail = () => {
    emailError.value = !toEmail.value
    return !emailError.value
}

const showAddress = () => {
    if (!validateEmail()) return
    generatedEmail.value = alias.value.name.replace('@', `+${toEmail.value.replace('@', '=')}@`)
}

const close = () => {
    toEmail.value = ''
    generatedEmail.value = ''
    emailError.value = false
    const modal = document.querySelector('#modal-send-alias' + alias.value.id)
    if (modal instanceof HTMLElement) {
        overlay.close(modal)
    }
}

const copy = (text: string) => {
    navigator.clipboard.writeText(text)
    copyText.value = 'Copied!'
    setTimeout(() => {
        copyText.value = 'Click to copy'
    }, 2000)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-send-alias' + alias.value.id, true)
    modal.element.on('close', () => {
        close()
    })
}

onMounted(() => {
    overlay.autoInit()
    tooltip.autoInit()
    addEvents()
})
</script>