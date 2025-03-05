<template>
    <tr>
        <td class="pr-5 py-4 whitespace-nowrap text-start text-sm text-gray-800 dark:text-gray-100">
            <p>{{ new Date(recipient.created_at).toDateString() }}</p>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800 dark:text-gray-100">
            <div class="hs-tooltip inline-block">
                <span class="hs-tooltip-toggle">
                    <button class=" dark:text-gray-100 truncate max-w-[320px]" @click="copyAlias(recipient.email)">
                        {{ recipient.email }}
                    </button>
                    <span
                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
                        role="tooltip">
                        {{ copyText }}: {{ recipient.email }}
                    </span>
                </span>
            </div>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <p>
                <span v-if="recipient.is_active" class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-emerald-100 text-emerald-800 dark:bg-emerald-800 dark:text-emerald-100">Verified</span>
                <span v-if="!recipient.is_active" class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-gray-100 text-gray-500 dark:bg-gray-500 dark:text-gray-100">Unverified</span>
            </p>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <RecipientAddPGPKey v-if="!recipient.pgp_key" :recipient="recipient" />
            <div v-if="recipient.pgp_key">
                <input
                    @change="updateRecipient"
                    v-bind:checked="recipient.pgp_enabled"
                    v-bind:disabled="!recipient.pgp_key"
                    type="checkbox"
                    class="mr-5"
                >
                <button
                    @click="deletePgpKey"
                    class="text-red-600 hover:text-red-700 dark:text-red-500 dark:hover:text-red-600 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline">
                    Delete Key
                </button>
            </div>
        </td>
        <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
            <div class="flex gap-5 justify-end">
                <RecipientVerify v-if="!recipient.is_active" :recipient="recipient" />
                <RecipientEdit :recipient="recipient" />
                <button
                    @click.stop="deleteRecipient"
                    class="text-red-600 hover:text-red-700 dark:text-red-500 dark:hover:text-red-600 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline">
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