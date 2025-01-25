<template>
    <tr>
        <td class="pr-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <div class="mt-1 flex items-center gap-2">
                <div v-if="alias.enabled" class="hs-tooltip flex-none rounded-full bg-emerald-500/20 p-1">
                    <div class="h-1.5 w-1.5 rounded-full bg-emerald-400"></div>
                    <span
                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
                        role="tooltip">
                        Active
                    </span>
                </div>
                <div v-if="!alias.enabled" class="hs-tooltip flex-none rounded-full bg-gray-500/20 p-1">
                    <div class="h-1.5 w-1.5 rounded-full bg-gray-400 dark:bg-neutral-500"></div>
                    <span
                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
                        role="tooltip">
                        Inactive
                    </span>
                </div>
                <p class="dark:text-gray-100">{{ new Date(alias.created_at).toDateString() }}</p>
            </div>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm text-gray-800">
            <div class="hs-tooltip inline-block">
                <span class="hs-tooltip-toggle">
                    <button class="dark:text-gray-100 truncate max-w-[320px]" @click="copyAlias(alias.name)">
                        {{ alias.name }}
                    </button>
                    <span
                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
                        role="tooltip">
                        {{ copyText }}: {{ alias.name }}
                    </span>
                </span>
            </div>
            <p class="text-gray-500 dark:text-gray-400">{{ alias.description }}</p>
        </td>
        <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
            <div class="hs-tooltip inline-block">
                <span class="hs-tooltip-toggle dark:text-gray-100">
                    {{ alias.stats.forwards }}/{{ alias.stats.blocks }}/{{ alias.stats.replies }}/{{ alias.stats.sends }}
                    <span
                        class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
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
            <div class="flex items-center hs-tooltip">
                <input
                    v-bind:disabled="!alias.recipients.length"
                    type="checkbox" v-bind:checked="alias.enabled" @change="updateAlias"
                    class="form-checkbox relative w-11 h-6 p-px bg-gray-100 dark:bg-neutral-600 border-transparent text-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:ring-white dark:focus:ring-neutral-800 disabled:opacity-50 disabled:pointer-events-none checked:bg-none checked:text-bluish-500 checked:border-bluish-500 focus:ring-offset-transparent

                    before:inline-block before:size-5 before:bg-white dark:before:bg-neutral-400 checked:before:bg-bluish-200 before:translate-x-0 checked:before:translate-x-full before:rounded-full before:shadow before:transform before:transition before:ease-in-out before:duration-200">
                <span
                    v-if="!alias.recipients.length"
                    class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0 transition-opacity inline-block absolute invisible z-10 py-1 px-2 bg-gray-900 dark:bg-neutral-900 text-xs font-medium text-white rounded shadow-sm"
                    role="tooltip">
                    Disabled
                </span>
            </div>
        </td>
        <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
            <div class="flex gap-5 justify-end">
                <AliasSend :alias="alias" />
                <AliasEdit :alias="alias" :recipients="recipients" :key="rowKey" />
                <button
                    @click="deleteAlias"
                    class="text-red-600 hover:text-red-700 dark:text-red-500 dark:hover:text-red-600 font-medium text-sm py-2 focus:outline-none focus:shadow-outline"
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
import AliasEdit from './AliasEdit.vue'
import AliasSend from './AliasSend.vue'
import { aliasApi } from '../api/alias.ts'
import events from '../events.ts'

const props = defineProps(['alias', 'recipients'])
const alias = ref(props.alias)
const recipients = ref(props.recipients)
const copyText = ref('Click to copy')
const rowKey = ref(0)

const updateAlias = async () => {
    alias.value.enabled = !alias.value.enabled
    try {
        await aliasApi.update(alias.value.id, alias.value)
        renderRow()
    } catch {}
}

const deleteAlias = () => {
    events.emit('alias.delete', { id: alias.value.id })
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias)
    copyText.value = 'Copied'
    setTimeout(() => {
        copyText.value = 'Click to copy'
    }, 2000)
}

const renderRow = () => {
    rowKey.value++
    tooltip.autoInit()
}

onMounted(() => {
    tooltip.autoInit()
})
</script>