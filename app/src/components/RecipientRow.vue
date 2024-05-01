<template>
    <tr>
        <td class="pr-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <p>{{ new Date(recipient.created_at).toDateString() }}</p>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <p>{{ recipient.email }}</p>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <p>
                <span v-if="recipient.is_active" class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-teal-100 text-teal-800">verified</span>
                <span v-if="!recipient.is_active" class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-gray-100 text-gray-500">unverified</span>
            </p>
        </td>
        <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
            <div class="flex gap-5 justify-end">
                <!-- TODO: Edit recipient component -->
                <button
                    @click="deleteRecipient"
                    class="text-red-600 hover:text-red-700 font-medium text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Delete
                </button>
            </div>
        </td>
    </tr>
</template>

<script setup lang="ts">
import { ref } from 'vue'
// import { recipientApi } from '../api/recipient.ts'

const props = defineProps(['recipient'])
const recipient = ref(props.recipient)
const emit = defineEmits(['onDeleteRecipient', 'onEditRecipient'])

const deleteRecipient = () => {
    emit('onDeleteRecipient', recipient.value.id)
}

// const onEditRecipient = () => {
//     emit('onEditRecipient')
// }
</script>