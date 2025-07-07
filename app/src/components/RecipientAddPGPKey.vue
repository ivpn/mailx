<template>
    <div>
        <div v-bind:id="'modal-add-key-recipient' + recipient.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>ADD PGP PUBLIC KEY</h4>
                    </header>
                    <article>
                        <label for="recipient_pgp">
                            Enter your public PGP key:
                        </label>
                        <textarea
                            v-model="pgp_key"
                            v-bind:class="{ 'error': pgpError }"
                            id="recipient_pgp"
                            placeholder="Starts with '-----BEGIN PGP PUBLIC KEY BLOCK-----'"
                        >
                        </textarea>
                        <p v-if="pgpError" class="error">Required</p>
                    </article>
                    <footer>
                        <nav>
                            <button @click="addKey" class="cta">
                                Add PGP Public Key
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
        pgp_key: pgp_key.value.trim()
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
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const close = () => {
    error.value = ''
    pgpError.value = false
    pgp_key.value = ''
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