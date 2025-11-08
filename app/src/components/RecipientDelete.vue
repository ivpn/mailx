<template>
    <div>
        <div v-bind:id="'modal-delete-recipient' + recipient.id" class="hs-overlay hidden">
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
                                    <strong>WARNING:</strong> This operation cannot be undone. Deleting this recipient will remove it from all associated aliases, and any alias without a recipient will be disabled.
                                </p>
                            </div>
                            <div class="mb-5">
                                <p>
                                    You may choose a new recipient for the affected aliases:
                                </p>
                                <label v-bind:for="'recipient_' + recipient.id">
                                    Recipient(s):
                                </label>
                                <select
                                    v-model="newEmails"
                                    v-bind:id="'recipient_' + recipient.id"
                                    :disabled="!emails.length"
                                    :multiple="true"
                                    data-hs-select='{
                                    "placeholder": "Select recipient(s)",
                                    "toggleTag": "<button type=\"button\" aria-expanded=\"false\"></button>",
                                    "toggleClasses": "hs-select-disabled:pointer-events-none hs-select-disabled:opacity-50 relative py-2.5 ps-4 pe-9 flex gap-x-2 text-nowrap w-full cursor-pointer border border-primary text-secondary bg-secondary leading-tight focus:border-accent",
                                    "dropdownClasses": "mt-2 z-50 w-full max-h-72 p-1 space-y-0.5 bg-primary border border-tertiary overflow-hidden overflow-y-auto",
                                    "optionClasses": "py-2 px-4 w-full text-secondary cursor-pointer hover:bg-secondary bg-primary",
                                    "optionTemplate": "<div class=\"flex justify-between items-center w-full\"><span data-title></span><span class=\"hidden hs-selected:block\"><svg class=\"shrink-0 size-3.5 text-accent \" xmlns=\"http:.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><polyline points=\"20 6 9 17 4 12\"/></svg></span></div>",
                                    "extraMarkup": "<div class=\"absolute top-1/2 end-3 -translate-y-1/2\"><svg class=\"shrink-0 size-3.5 text-secondary \" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"m7 15 5 5 5-5\"/><path d=\"m7 9 5-5 5 5\"/></svg></div>"
                                    }' class="hidden">
                                    <option v-for="rcp in emails"
                                        v-bind:value=rcp
                                        :selected="newEmails.includes(rcp)"
                                        :key="rcp">
                                        {{ rcp }}
                                    </option>
                                </select>
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
import { ref, onMounted, onBeforeMount } from 'vue'
import { recipientApi } from '../api/recipient.ts'
import axios from 'axios'
import overlay from '@preline/overlay'
import select from '@preline/select'
import events from '../events.ts'

const props = defineProps(['recipient', 'recipients'])
const recipient = ref(props.recipient)
const recipients = ref(props.recipients)
const emails = ref([] as string[])
const newEmails = ref([] as string[])
const error = ref('')

const deleteRecipient = async () => {
    if (!confirm('Are you sure you want to delete recipient? Aliases without recipient(s) will be disabled.')) return

    const data = { recipients: newEmails.value.join(',') }

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

const updateEmails = () => {
        emails.value = recipients.value.map((r: { email: string }) => r.email)
    
        // Remove current recipient from the list
        const index = emails.value.indexOf(recipient.value.email)
        if (index > -1) {
            emails.value.splice(index, 1)
        }
}

const close = () => {
    error.value = ''
    const modal = document.querySelector('#modal-delete-recipient' + recipient.value.id) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-delete-recipient' + recipient.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

onBeforeMount(() => {
    updateEmails()
})

onMounted(() => {
    overlay.autoInit()
    select.autoInit()
    addEvents()
})
</script>