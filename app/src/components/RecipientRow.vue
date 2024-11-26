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
        <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
            <div class="flex gap-5 justify-end">
                <RecipientVerify v-if="!recipient.is_active" :recipient="recipient" @onVerifyRecipient="onVerifyRecipient" />
                <button
                    @click="deleteRecipient"
                    class="text-red-600 hover:text-red-700 dark:text-red-500 dark:hover:text-red-600 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
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
// import { recipientApi } from '../api/recipient.ts'

const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const emit = defineEmits(['onDeleteRecipient', 'onVerifyRecipient'])
const copyText = ref('Click to copy')

const deleteRecipient = () => {
    emit('onDeleteRecipient', recipient.value.id)
}

const onVerifyRecipient = () => {
    emit('onVerifyRecipient')
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