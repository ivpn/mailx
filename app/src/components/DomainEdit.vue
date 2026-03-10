<template>
    <div>
        <div v-bind:id="'modal-edit-domain' + domain.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>EDIT DOMAIN</h4>
                    </header>
                    <article>
                        <div class="mb-5">
                            <h4>Default Recipient</h4>
                            <p>
                                Set the default recipient for this domain.
                            </p>
                            <div class="mb-6">
                                <label for="recipient">
                                    Select default recipient:
                                </label>
                                <select id="recipient" :disabled="!recipients.length">
                                    <option
                                        v-for="recipient in recipients"
                                        v-bind:value=recipient
                                        :selected="recipient == domain.recipient"
                                        :key="recipient">
                                        {{ recipient }}
                                    </option>
                                </select>
                            </div>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button
                                @click="updateDomain"
                                class="cta">
                                Save
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
import { domainApi } from '../api/domain.ts'

const props = defineProps(['domain', 'recipients'])
const domain = ref(props.domain)
const recipients = ref(props.recipients)
const error = ref('')

const updateDomain = async () => {
    const payload = {
        id: domain.value.id,
        recipient: domain.value.recipient
    }

    try {
        await domainApi.update(domain.value.id, payload)
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            const errorMsg = err.response?.data.error || err.message
            error.value = err.response?.status === 429
                ? 'Too many requests, please try again later.'
                : errorMsg
        }
    }
}

const close = () => {
    error.value = ''
    const modal = document.querySelector('#modal-edit-domain' + domain.value.id) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-edit-domain' + domain.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
    // Prepend an empty option to allow unsetting default recipient
    recipients.value.unshift('')
})
</script>