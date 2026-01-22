<template>
    <div>
        <div v-bind:id="'modal-send-alias' + alias.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>SEND FROM ALIAS</h4>
                    </header>
                    <article>
                        <div class="mb-5">
                            <p>
                                To send an email using an alias, you need a specific forwarding address generated after providing the target email address. Only emails originating from your verified recipient email addresses are forwarded.
                            </p>
                        </div>
                        <div class="mb-7">
                            <label v-bind:for="'from_alias_' + alias.id">
                                From alias:
                            </label>
                            <input
                                v-bind:id="'from_alias_' + alias.id"
                                v-bind:value="alias.name" disabled
                                type="text"
                            >
                        </div>
                        <div class="mb-7">
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
                            <div class="hs-tooltip [--strategy:absolute] mb-3">
                                <span class="hs-tooltip-toggle">
                                    <button @click="copy(generatedEmail)" class="plain break-all">
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
                            <button @click="close" class="cancel">
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
    modal.element.on('open', () => {
        focusFirstInput()
    })
}

const focusFirstInput = () => {
    const input = document.getElementById('to_email_' + alias.value.id) as HTMLInputElement
    input?.focus()
}

onMounted(() => {
    overlay.autoInit()
    tooltip.autoInit()
    addEvents()
})
</script>