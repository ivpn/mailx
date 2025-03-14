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
                    <button class="plain truncate max-w-[320px] text-base" @click="copyAlias(alias.name)">
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
            <p>@{{ alias.name.split('@')[1] }}</p>
        </td>
        <td>
            <div class="hs-tooltip inline-block">
                <p class="hs-tooltip-toggle dark:text-zinc-100">
                    {{ alias.stats.forwards }}/{{ alias.stats.blocks }}/{{ alias.stats.replies }}/{{ alias.stats.sends }}
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                        {{ alias.stats.forwards }} Forwards<br>
                        {{ alias.stats.blocks }} Blocks<br>
                        {{ alias.stats.replies }} Replies<br>
                        {{ alias.stats.sends }} Sends
                    </span>
                </p>
            </div>
        </td>
        <td>
            <div class="mt-1 flex items-center gap-2">
                <!-- <div v-if="alias.enabled" class="hs-tooltip flex-none rounded-full bg-emerald-500/20 p-1">
                    <div class="h-1.5 w-1.5 rounded-full bg-emerald-400"></div>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                        Active
                    </span>
                </div>
                <div v-if="!alias.enabled" class="hs-tooltip flex-none rounded-full bg-zinc-500/20 p-1">
                    <div class="h-1.5 w-1.5 rounded-full bg-zinc-400 dark:bg-neutral-500"></div>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                        Inactive
                    </span>
                </div> -->
                <p>{{ new Date(alias.created_at).toDateString() }}</p>
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