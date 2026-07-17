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
                            <h4>Catch-All Recipient</h4>
                            <p>
                                Set the default recipient for this domain. This overrides the default recipient selected in the Settings.
                            </p>
                            <div class="mb-6">
                                <label v-bind:for="'recipient_' + domain.id">
                                    Select default recipient:
                                </label>
                                <select v-bind:id="'recipient_' + domain.id" v-model="selectedRecipient">
                                    <option value="">—</option>
                                    <option
                                        v-for="recipient in recipients"
                                        v-bind:value="recipient"
                                        :key="recipient">
                                        {{ recipient }}
                                    </option>
                                </select>
                            </div>
                        </div>
                        <div class="mb-5">
                            <h4>Catch-All From Name</h4>
                            <p>
                                Set the default "From" name for this domain. This overrides the default "From" name selected in the Settings.
                            </p>
                            <div class="mb-6">
                                <label v-bind:for="'from_name_' + domain.id">
                                    From name:
                                </label>
                                <input
                                    type="text"
                                    v-bind:id="'from_name_' + domain.id"
                                    v-model="fromName"
                                    placeholder="e.g. John Doe"
                                />
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
import events from '../events.ts'

const props = defineProps(['domain', 'recipients'])
const domain = ref(props.domain)
const recipients = ref(props.recipients)
const selectedRecipient = ref(props.domain.recipient ?? '')
const fromName = ref(props.domain.from_name ?? '')
const error = ref('')

const updateDomain = async () => {
    const payload = {
        id: domain.value.id,
        recipient: selectedRecipient.value,
        from_name: fromName.value,
        enabled: domain.value.enabled,
        catch_all: domain.value.catch_all,
    }

    try {
        await domainApi.update(domain.value.id, payload)
        error.value = ''
        events.emit('domain.update', {})
        close()
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
    selectedRecipient.value = props.domain.recipient ?? ''
    fromName.value = props.domain.from_name ?? ''
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
})
</script>