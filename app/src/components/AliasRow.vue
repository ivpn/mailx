<template>
    <tr>
        <td>
            <div class="flex items-center hs-tooltip">
                <input
                    @change="updateAlias"
                    v-bind:checked="alias.enabled"
                    v-bind:disabled="!alias.recipients.length"
                    type="checkbox"
                >
                <span v-if="!alias.recipients.length" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                    Disabled
                </span>
            </div>
        </td>
        <td>
            <div class="hs-tooltip inline-block">
                <p class="hs-tooltip-toggle">
                    <button class="plain truncate max-w-[320px] text-base p-0" @click="copyAlias(alias.name)">
                        {{ alias.name.split('@')[0] }}
                    </button>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                        {{ copyText }}: {{ alias.name }}
                    </span>
                </p>
            </div>
            <p class="foreground-secondary">{{ alias.description }}</p>
        </td>
        <td>
            <p class="py-3">@{{ alias.name.split('@')[1] }}</p>
        </td>
        <td>
            <p class="py-3">
                <span v-if="!alias.catch_all">Alias</span>
                <span v-if="alias.catch_all">Catch-all</span>
            </p>
        </td>
        <td>
            <div class="flex items-center gap-3 mb-1">
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.forwards }}
                    <i class="icon forward text-xs bg-neutrallight-10 dark:bg-neutraldark-10"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.forwards }} Forwards</span>
                </p>
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.blocks }}
                    <i class="icon block text-xs bg-neutrallight-10 dark:bg-neutraldark-10"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.blocks }} Blocks</span>
                </p>
            </div>
            <div class="flex items-center gap-3 mt-1">
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.replies }}
                    <i class="icon reply text-xs bg-neutrallight-10 dark:bg-neutraldark-10"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.replies }} Replies</span>
                </p>
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.sends }}
                    <i class="icon send text-xs bg-neutrallight-10 dark:bg-neutraldark-10"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.sends }} Sends</span>
                </p>
            </div>
        </td>
        <td>
            <div class="mt-1 flex items-center gap-2">
                <p>{{ formatDistanceToNow(new Date(alias.created_at)) }}</p>
            </div>
        </td>
        <td>
            <div class="flex gap-5 justify-end">
                <AliasSend :alias="alias" />
                <AliasEdit :alias="alias" :recipients="recipients" :key="rowKey" />
                <button @click.stop="deleteAlias" class="delete">
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
import { formatDistanceToNow } from 'date-fns'

const props = defineProps(['alias', 'recipients', 'catchAll'])
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
    const errMessage = props.catchAll ? 'WARNING: You will not be able to create the same catch-all alias in the next 90 days. Are you sure you want to delete alias? ' : 'Are you sure you want to delete alias?'
    if (!confirm(errMessage)) return

    events.emit('alias.delete', { id: alias.value.id, catchAll: props.catchAll })
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