<template>
    <tr>
        <td class="pr-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <div class="mt-1 flex items-center gap-2">
                <div v-if="alias.enabled" class="flex-none rounded-full bg-emerald-500/20 p-1">
                    <div class="h-1.5 w-1.5 rounded-full bg-emerald-500"></div>
                </div>
                <div v-if="!alias.enabled" class="flex-none rounded-full bg-gray-500/20 p-1">
                    <div class="h-1.5 w-1.5 rounded-full bg-gray-400"></div>
                </div>
                <p>{{ new Date(alias.created_at).toDateString() }}</p>
            </div>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            {{ alias.name }}
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            {{ alias.recipients.split(',').length }}
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
            <div class="hs-tooltip inline-block">
                <span class="hs-tooltip-toggle">
                    {{ alias.stats.forwards }}/{{ alias.stats.blocks }}/{{ alias.stats.replies }}/{{ alias.stats.sends }}
                    <span
                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 text-xs font-medium text-white rounded shadow-sm"
                        role="tooltip">
                        {{ alias.stats.forwards }} Forwards<br>
                        {{ alias.stats.blocks }} Blocks<br>
                        {{ alias.stats.replies }} Replies<br>
                        {{ alias.stats.sends }} Sends
                    </span>
                </span>
            </div>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
            <div class="flex items-center">
                <input type="checkbox" v-bind:checked="alias.enabled"
                    class="form-checkbox relative w-11 h-6 p-px bg-gray-100 border-transparent text-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:ring-white disabled:opacity-50 disabled:pointer-events-none checked:bg-none checked:text-violet-600 checked:border-violet-600 

                    before:inline-block before:size-5 before:bg-white checked:before:bg-violet-200 before:translate-x-0 checked:before:translate-x-full before:rounded-full before:shadow before:transform before:ring-0 before:transition before:ease-in-out before:duration-200">
            </div>
        </td>
        <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
            <div class="flex gap-5 justify-end">
                <AliasEdit :alias="alias" :recipients="recipients" />
                <button
                    @click="deleteAlias"
                    class="text-red-600 hover:text-red-700 font-semibold text-sm py-2 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Delete
                </button>
            </div>
        </td>
    </tr>
</template>

<script setup lang="ts">
import { ref, onMounted, defineEmits } from 'vue'
import tooltip from '@preline/tooltip'
import AliasEdit from './AliasEdit.vue'

const props = defineProps(['alias', 'recipients'])
let alias = props.alias
const recipients = ref(props.recipients)
const emit = defineEmits(['onDeleteAlias'])

const deleteAlias = () => {
    emit('onDeleteAlias', alias.id)
}

onMounted(() => {
    tooltip.autoInit()
})
</script>