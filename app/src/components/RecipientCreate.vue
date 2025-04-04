<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-create-recipient'" class="cta">
            <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                stroke-linejoin="round">
                <path d="M5 12h14"></path>
                <path d="M12 5v14"></path>
            </svg>
            New Recipient
        </button>
        <div v-bind:id="'modal-create-recipient'" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>NEW RECIPIENT</h4>
                    </header>
                    <article>
                        <div class="mb-5">
                            <p>
                                Add a email address to receive forwarded emails. A 6 digit verification code will be sent to this email address. If expired, you can resend the verification code. Unverified recipients will not receive any forwarded emails and will be deleted after 7 days. You can add up to 10 recipients.
                            </p>
                        </div>
                        <div class="mb-3">
                            <label for="recipient_email">
                                Email:
                            </label>
                            <input
                                v-model="recipient.email"
                                v-bind:class="{ 'error': emailError }"
                                id="recipient_email"
                                placeholder="name@example.net"
                                type="text"
                            >
                            <p v-if="emailError" class="error">Required</p>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click="postRecipient" class="cta">
                                Add Recipient
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
import overlay from '@preline/overlay'
import axios from 'axios'
import { recipientApi } from '../api/recipient.ts'
import events from '../events.ts'

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
        events.emit('recipient.create', {})
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
    recipient.value = {} as any
    error.value = ''
    emailError.value = false
    const modal = document.querySelector('#modal-create-recipient') as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-create-recipient' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>