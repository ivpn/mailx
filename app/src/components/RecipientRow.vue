<template>
    <tr>
        <td>
            <p>{{ new Date(recipient.created_at).toDateString() }}</p>
        </td>
        <td>
            <div class="hs-tooltip inline-block">
                <span class="hs-tooltip-toggle">
                    <button class="plain truncate max-w-[320px]" @click="copyAlias(recipient.email)">
                        {{ recipient.email }}
                    </button>
                    <span class="hs-tooltip-content" role="tooltip">
                        {{ copyText }}: {{ recipient.email }}
                    </span>
                </span>
            </div>
        </td>
        <td>
            <p>
                <span v-if="recipient.is_active" class="badge success">Verified</span>
                <span v-if="!recipient.is_active" class="badge">Unverified</span>
            </p>
        </td>
        <td>
            <RecipientAddPGPKey v-if="!recipient.pgp_key" :recipient="recipient" />
            <div v-if="recipient.pgp_key">
                <input
                    @change="updateRecipient"
                    v-bind:checked="recipient.pgp_enabled"
                    v-bind:disabled="!recipient.pgp_key"
                    type="checkbox"
                    class="mr-5"
                >
                <button @click.stop="deletePgpKey" class="delete">
                    Delete Key
                </button>
            </div>
        </td>
        <td>
            <div class="flex gap-5 justify-end">
                <RecipientVerify v-if="!recipient.is_active" :recipient="recipient" />
                <RecipientEdit :recipient="recipient" />
                <button @click.stop="deleteRecipient" class="delete">
                    Delete
                </button>
            </div>
        </td>
    </tr>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import tooltip from '@preline/tooltip'
import RecipientVerify from './RecipientVerify.vue'
import RecipientEdit from './RecipientEdit.vue'
import RecipientAddPGPKey from './RecipientAddPGPKey.vue'
import { recipientApi } from '../api/recipient.ts'
import events from '../events.ts'

const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const copyText = ref('Click to copy')

const deleteRecipient = () => {
    if (!confirm('Are you sure you want to delete recipient? Note that aliases with this recipient will be disabled.')) return
    
    events.emit('recipient.delete', { id: recipient.value.id })
}

const updateRecipient = async () => {
    // Toggle pgp_enabled option
    recipient.value.pgp_enabled = !recipient.value.pgp_enabled

    const payload = {
        id: recipient.value.id,
        pgp_key: recipient.value.pgp_key,
        pgp_enabled: recipient.value.pgp_enabled
    }

    try {
        await recipientApi.update(payload)
    } catch {}
}

const deletePgpKey = async () => {
    if (!confirm('Are you sure you want to delete PGP public key?')) return

    const payload = {
        id: recipient.value.id,
        pgp_key: '',
        pgp_enabled: false,
    }

    try {
        await recipientApi.update(payload)
        events.emit('recipient.update', {})
    } catch {}
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
})
</script>