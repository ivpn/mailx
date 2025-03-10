<template>
    <div>
        <button
            v-bind:disabled="!alias.recipients.length"
            v-bind:data-hs-overlay="'#modal-send-alias' + alias.id">
            Send
        </button>
        <div v-bind:id="'modal-send-alias' + alias.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
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
                    </header>
                    <article>
                        <div class="mb-5">
                            <p>
                                Generate the proper email address to send a message from this alias. Note that to send emails using an alias, you need to do so from a verified recipient.
                            </p>
                        </div>
                        <div class="mb-5">
                            <label v-bind:for="'from_alias_' + alias.id">
                                From alias:
                            </label>
                            <input
                                v-bind:id="'from_alias_' + alias.id"
                                v-bind:value="alias.name" disabled
                                type="text"
                            >
                        </div>
                        <div class="mb-5">
                            <label v-bind:for="'to_email_' + alias.id">
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
                                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                        {{ copyText }}
                                    </span>
                                </span>
                            </div>

                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click="showAddress" class="cta">
                                Show address
                            </button>
                            <button @click="close" class="cta cancel">
                                Close
                            </button>
                        </nav>
                    </footer>
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