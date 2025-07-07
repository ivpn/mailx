<template>
    <tr class="desktop">
        <td>
            <p>{{ new Date(recipient.created_at).toDateString() }}</p>
        </td>
        <td>
            <div class="hs-tooltip inline-block">
                <p class="hs-tooltip-toggle">
                    <button class="plain truncate max-w-[320px]" @click="copyAlias(recipient.email)">
                        {{ recipient.email }}
                    </button>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                        {{ copyText }}: {{ recipient.email }}
                    </span>
                </p>
            </div>
        </td>
        <td>
            <p>
                <span v-if="recipient.is_active" class="badge success">Verified</span>
                <span v-if="!recipient.is_active" class="badge">Unverified</span>
            </p>
        </td>
        <td>
            <div class="hs-tooltip inline-block">
                <input
                    @change="updateRecipient"
                    v-bind:checked="recipient.pgp_enabled"
                    v-bind:disabled="!recipient.pgp_key"
                    type="checkbox"
                    class="mr-4"
                >
                <span v-if="!recipient.pgp_key" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                    Add PGP key to enable encryption
                </span>
            </div>
        </td>
        <td>
            <div class="hs-dropdown [--offset:0]">
                <button v-bind:id="'hs-dropdown-recipient-edit-' + recipient.id">
                    <i class="icon icon-secondary more text-lg"></i>
                </button>
                <div
                    class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                    v-bind:aria-labelledby="'hs-dropdown-recipient-edit-' + recipient.id"
                    >
                    <button v-bind:class="{ 'hide': recipient.is_active }"v-bind:data-hs-overlay="'#modal-verify-recipient' + recipient.id">
                        <i class="icon icon-primary check text-xs"></i>
                        Verify
                    </button>
                    <button v-bind:data-hs-overlay="'#modal-edit-recipient' + recipient.id">
                        <i class="icon icon-primary edit text-xs"></i>
                        Edit
                    </button>
                    <button v-bind:class="{ 'hide': recipient.pgp_key }" v-bind:data-hs-overlay="'#modal-add-key-recipient' + recipient.id">
                        <i class="icon icon-primary key text-xs"></i>
                        Add PGP Key
                    </button>
                    <button v-if="recipient.pgp_key" @click.stop="deletePgpKey" class="delete">
                        <i class="icon icon-error trash text-xs"></i>
                        Remove PGP Key
                    </button>
                    <button @click.stop="deleteRecipient" class="delete">
                        <i class="icon icon-error trash text-xs"></i>
                        Delete
                    </button>
                </div>
            </div>
        </td>
    </tr>
    <tr class="tablet">
        <td>
            <div class="flex gap-2 justify-between">
                <div class="text-start">
                    <div>
                        <p class="mb-3">{{ new Date(recipient.created_at).toDateString() }}</p>
                    </div>
                    <div class="hs-tooltip inline-block">
                        <p class="hs-tooltip-toggle">
                            <button class="plain text-base truncate p-0 max-w-[320px]" @click="copyAlias(recipient.email)">
                                {{ recipient.email }}
                            </button>
                            <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                {{ copyText }}: {{ recipient.email }}
                            </span>
                        </p>
                    </div>
                    <div class="hs-tooltip">
                        <p class="mb-3">Encryption:</p>
                        <input
                            @change="updateRecipient"
                            v-bind:checked="recipient.pgp_enabled"
                            v-bind:disabled="!recipient.pgp_key"
                            type="checkbox"
                            class="mr-4"
                        >
                        <span v-if="!recipient.pgp_key" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                            Add PGP key to enable encryption
                        </span>
                    </div>
                </div>
                <div class="text-end">
                    <div class="hs-dropdown [--offset:0]">
                        <button class="py-0" v-bind:id="'hs-dropdown-recipient-edit-' + recipient.id">
                            <i class="icon icon-secondary more text-lg"></i>
                        </button>
                        <div
                            class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                            v-bind:aria-labelledby="'hs-dropdown-recipient-edit-' + recipient.id"
                            >
                            <button v-bind:class="{ 'hide': recipient.is_active }"v-bind:data-hs-overlay="'#modal-verify-recipient' + recipient.id">
                                <i class="icon icon-primary check text-xs"></i>
                                Verify
                            </button>
                            <button v-bind:data-hs-overlay="'#modal-edit-recipient' + recipient.id">
                                <i class="icon icon-primary edit text-xs"></i>
                                Edit
                            </button>
                            <button v-bind:class="{ 'hide': recipient.pgp_key }" v-bind:data-hs-overlay="'#modal-add-key-recipient' + recipient.id">
                                <i class="icon icon-primary key text-xs"></i>
                                Add PGP Key
                            </button>
                            <button v-if="recipient.pgp_key" @click.stop="deletePgpKey" class="delete">
                                <i class="icon icon-error trash text-xs"></i>
                                Remove PGP Key
                            </button>
                            <button @click.stop="deleteRecipient" class="delete">
                                <i class="icon icon-error trash text-xs"></i>
                                Delete
                            </button>
                        </div>
                    </div>
                    <div>
                        <p class="my-3">
                            <span v-if="recipient.is_active" class="badge success">Verified</span>
                            <span v-if="!recipient.is_active" class="badge">Unverified</span>
                        </p>
                    </div>
                </div>
            </div>
            <hr>
        </td>
    </tr>

    <RecipientAddPGPKey :recipient="recipient" />
    <RecipientVerify :recipient="recipient" />
    <RecipientEdit :recipient="recipient" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import tooltip from '@preline/tooltip'
import RecipientVerify from './RecipientVerify.vue'
import RecipientEdit from './RecipientEdit.vue'
import RecipientAddPGPKey from './RecipientAddPGPKey.vue'
import { recipientApi } from '../api/recipient.ts'
import events from '../events.ts'
import dropdown from '@preline/dropdown'
import axios from 'axios'

const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const copyText = ref('Click to copy')

const deleteRecipient = async () => {
    if (!confirm('Are you sure you want to delete recipient? Note that aliases with this recipient will be disabled.')) return

    try {
        await recipientApi.delete(recipient.value.id)
        events.emit('recipient.reload', {})
    } catch (err) {
        if (axios.isAxiosError(err)) {
            const error = err.response?.data.error || err.message
            events.emit('recipient.delete.error', { error: error })
        }
    }
}

const updateRecipient = async () => {
    // Toggle pgp_enabled option
    const temp_pgp_enabled = recipient.value.pgp_enabled
    recipient.value.pgp_enabled = !recipient.value.pgp_enabled

    const payload = {
        id: recipient.value.id,
        pgp_key: recipient.value.pgp_key,
        pgp_enabled: recipient.value.pgp_enabled
    }

    try {
        await recipientApi.update(payload)
    } catch (err) {
        if (axios.isAxiosError(err)) {
            var errMsg = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                errMsg = 'Too many requests, please try again later.'
            }

            recipient.value.pgp_enabled = temp_pgp_enabled // revert the change

            alert(errMsg)
        }
    }
}

const deletePgpKey = async () => {
    if (!confirm('Are you sure you want to remove PGP public key?')) return

    const payload = {
        id: recipient.value.id,
        pgp_key: '',
        pgp_enabled: false,
    }

    try {
        await recipientApi.update(payload)
        events.emit('recipient.update', {})
    } catch (err) {
        if (axios.isAxiosError(err)) {
            var errMsg = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                errMsg = 'Too many requests, please try again later.'
            }

            alert(errMsg)
        }
    }
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias)
    copyText.value = 'Copied'
    setTimeout(() => {
        copyText.value = 'Click to copy'
    }, 2000)
}

onMounted(() => {
    tooltip.autoInit()
    dropdown.autoInit()
})
</script>