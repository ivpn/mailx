<template>
    <div class="mb-5">
        <h2>Delete Recipient</h2>
        <button
            v-bind:data-hs-overlay="'#modal-delete-recipient'"
            class="cta delete mb-4">
            Delete Recipient
        </button>
        <div v-bind:id="'modal-delete-recipient'" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>DELETE RECIPIENT</h4>
                    </header>
                    <article>
                        <div>
                            <div class="mb-5">
                                <p>
                                    <strong>WARNING:</strong> this operation cannot be undone. Deleting your recipient will permanently remove it from all associated aliases.
                                </p>
                            </div>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click.stop="deleteRecipient" class="cta delete">
                                Delete Recipient
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
import { recipientApi } from '../api/recipient.ts'
import axios from 'axios'
import overlay from '@preline/overlay'
import events from '../events.ts'

const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const recipients = ref([] as string[])
const req = ref({ recipients: '' })
const error = ref('')

const deleteRecipient = async () => {
    if (!confirm('Are you sure you want to delete recipient? Note that aliases without recipient(s) will be disabled.')) return

    const data = { recipients: recipients.value }

    try {
        await recipientApi.delete(recipient.value.id, data)
        events.emit('recipient.reload', {})
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const getRecipients = async () => {
    try {
        const res = await recipientApi.getList()
        const list = res.data.filter((item: { is_active: boolean }) => item.is_active)
        recipients.value = list.map((recipient: { email: string }) => recipient.email)

        // Remove current recipient from the list
        const index = recipients.value.indexOf(recipient.value.email)
        if (index > -1) {
            recipients.value.splice(index, 1)
        }

        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const close = () => {
    req.value = {} as any
    error.value = ''
    const modal = document.querySelector('#modal-delete-recipient') as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-delete-recipient' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        getRecipients()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>