<template>
    <div>
        <div v-bind:id="'modal-edit-recipient' + recipient.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>EDIT RECIPIENT</h4>
                    </header>
                    <article>
                        <h3>{{ recipient.email }}</h3>
                        <div class="mb-5">
                            <h4>PGP/Inline Encryption</h4>
                            <p v-if="!recipient.pgp_key">
                                To use this option, please add a PGP key first.
                            </p>
                            <p>
                                Enable this option to use PGP/Inline instead of the default PGP/MIME encryption. Do not enable if you want to recieve encrypted emails with attachments or HTML content. Only forwarded emails are encrypted. 
                            </p>
                        </div>
                        <div class="pb-5">
                            <input
                                @change="updateRecipient"
                                v-bind:checked="pgp_inline"
                                v-bind:disabled="!recipient.pgp_key"
                                type="checkbox"
                            >
                        </div>
                    </article>
                    <footer>
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                        <p v-if="success" class="success px-5">{{ success }}</p>
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

const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const pgp_inline = ref(props.recipient.pgp_inline)
const error = ref('')
const success = ref('')

const updateRecipient = async () => {
    // Toggle pgp_inline option
    pgp_inline.value = !pgp_inline.value

    const temp_pgp_inline = recipient.value.pgp_inline

    const payload = {
        id: recipient.value.id,
        pgp_key: recipient.value.pgp_key,
        pgp_enabled: recipient.value.pgp_enabled,
        pgp_inline: pgp_inline.value
    }

    try {
        const res = await recipientApi.update(payload)
        error.value = ''
        success.value = res.data.message
    } catch (err) {
        if (axios.isAxiosError(err)) {
            const errorMsg = err.response?.data.error || err.message
            error.value = err.response?.status === 429
                ? 'Too many requests, please try again later.'
                : errorMsg
        }

        pgp_inline.value = temp_pgp_inline
    }
}

const close = () => {
    error.value = ''
    success.value = ''
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