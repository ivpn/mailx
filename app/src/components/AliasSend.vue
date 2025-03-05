<template>
    <div>
        <button
            v-bind:disabled="!alias.recipients.length"
            v-bind:data-hs-overlay="'#modal-send-alias' + alias.id">
            Send
        </button>
        <div v-bind:id="'modal-send-alias' + alias.id"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white dark:bg-neutral-800 border dark:border-neutral-600 shadow-sm rounded pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b dark:border-neutral-600">
                        <h3>Send from alias</h3>
                        <button @click="close" class="close">
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
                            <p>
                                Generate the proper email address to send a message from this alias. Note that to send emails using an alias, you need to do so from a verified recipient.
                            </p>
                        </div>
                        <div class="mb-5">
                            <label v-bind:for="'from_alias_' + alias.id"
                                class="block text-gray-500 dark:text-gray-400 mb-3">
                                From alias:
                            </label>
                            <input
                                v-bind:id="'from_alias_' + alias.id"
                                v-bind:value="alias.name" disabled
                                type="text"
                            >
                        </div>
                        <div class="mb-5">
                            <label v-bind:for="'to_email_' + alias.id"
                                class="block text-gray-500 dark:text-gray-400 mb-3">
                                To email:
                            </label>
                            <input
                                v-bind:id="'to_email_' + alias.id"
                                v-bind:class="{ 'error': emailError }"
                                v-model="toEmail"
                                type="text"
                            >
                            <p v-if="emailError" class="error">Valid email required</p>
                        </div>
                        <div v-bind:class="{ 'hidden': generatedEmail == '' }" class="mb-5">
                            <p>Send message to this email:</p>
                            <div class="hs-tooltip mb-3">
                                <span class="hs-tooltip-toggle">
                                    <button @click="copy(generatedEmail)" class="plain">
                                        {{ generatedEmail }}
                                    </button>
                                    <span
                                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
                                        role="tooltip">
                                        {{ copyText }}
                                    </span>
                                </span>
                            </div>

                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-3 py-4 px-4 border-t dark:border-neutral-600">
                        <button @click="showAddress" class="cta">
                            Show address
                        </button>
                        <button @click="close" class="cta cancel">
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

const isValidEmail = (email: string) => {
    const re = /\S+@\S+\.\S+/
    return re.test(email)
}

const validateEmail = () => {
    emailError.value = !toEmail.value || !isValidEmail(toEmail.value)
    return !emailError.value
}

const showAddress = () => {
    if (!validateEmail()) {
        generatedEmail.value = ''
        return
    }

    generatedEmail.value = alias.value.name.replace('@', `+${toEmail.value.replace('@', '=')}@`)
}

const close = () => {
    toEmail.value = ''
    generatedEmail.value = ''
    emailError.value = false
    const modal = document.querySelector('#modal-send-alias' + alias.value.id) as any
    overlay.close(modal)
}

const copy = (text: string) => {
    navigator.clipboard.writeText(text)
    copyText.value = 'Copied!'
    setTimeout(() => {
        copyText.value = 'Click to copy'
    }, 2000)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-send-alias' + alias.value.id as any, true) as any
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