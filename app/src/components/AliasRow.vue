<template>
    <tr v-if="isDesktop" class="desktop-lg">
        <td v-if="selectable" class="w-10">
            <div class="flex items-center">
                <input
                    type="checkbox"
                    class="checkbox-plain"
                    v-bind:checked="selected"
                    @change="$emit('onToggleSelect', alias.id)"
                >
            </div>
        </td>
        <td>
            <div class="flex items-center hs-tooltip">
                <input
                    @change="updateAlias"
                    v-bind:checked="alias.enabled && alias.recipients.length > 0 && !isDomainUnverified && !isAliasDeleted"
                    v-bind:disabled="!alias.recipients.length || isDomainUnverified || isAliasDeleted"
                    type="checkbox"
                    class="checkbox-switch"
                >
                <span v-if="isAliasDeleted" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                    Alias deleted. Address is not forwarding mail.
                </span>
                <span v-else-if="isDomainUnverified" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                    Domain not verified or disabled. Address is not forwarding mail.
                </span>
                <span v-else-if="!alias.recipients.length" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                    Disabled
                </span>
            </div>
        </td>
        <td class="whitespace-normal">
            <div class="block break-all hs-tooltip">
                <p class="hs-tooltip-toggle m-0">{{ truncatedDescription }}</p>
                <span v-if="alias.description && alias.description.length > 45" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.description }}</span>
            </div>
        </td>
        <td>
            <div class="hs-tooltip inline-block">
                <p class="hs-tooltip-toggle m-0 break-all">
                    <button class="plain text-wrap text-start text-sm p-0 flex items-center" @click="copyAlias(alias.name)">
                        <i v-if="alias.pinned" class="icon pin icon-accent text-xs mr-1 shrink-0"></i>{{ alias.name }}
                    </button>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                        {{ copyText }}: {{ alias.name }}
                    </span>
                </p>
                <p v-if="isCreatedByWildcard" class="text-xs text-tertiary mt-1">Created by Catch-All</p>
            </div>
        </td>
        <td>
            <div class="flex items-center gap-3 mb-1">
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.forwards }}
                    <i class="icon forward text-xs icon-tertiary"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.forwards }} Forwards</span>
                </p>
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.blocks }}
                    <i class="icon block text-xs icon-tertiary"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.blocks }} Blocks</span>
                </p>
            </div>
            <div class="flex items-center gap-3 mt-1">
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.replies }}
                    <i class="icon reply text-xs icon-tertiary"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.replies }} Replies</span>
                </p>
                <p class="flex items-center gap-1 hs-tooltip">
                    {{ alias.stats.sends }}
                    <i class="icon send text-xs icon-tertiary"></i>
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
            <div class="hs-dropdown">
                <button v-bind:id="'hs-dropdown-alias-edit-' + alias.id">
                    <i class="icon icon-secondary more text-lg"></i>
                </button>
                <div
                    class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                    v-bind:aria-labelledby="'hs-dropdown-alias-edit-' + alias.id"
                >
                    <button
                        v-if="!alias.deleted_at"
                        v-bind:disabled="!alias.recipients.length"
                        @click="$emit('onSend', alias)"
                        v-bind:class="{ 'hide': alias.wildcard }"
                        >
                        <i class="icon icon-primary send text-xs"></i>
                        Send
                    </button>
                    <button
                        v-if="!alias.deleted_at"
                        @click="$emit('onEdit', alias)">
                        <i class="icon icon-primary edit text-xs"></i>
                        Edit
                    </button>
                    <button
                        v-if="!alias.deleted_at"
                        @click.stop="togglePin">
                        <i class="icon icon-primary pin text-xs"></i>
                        {{ alias.pinned ? 'Unpin' : 'Pin' }}
                    </button>
                    <button 
                        v-if="alias.deleted_at"
                        @click.stop="restoreAlias">
                        <i class="icon icon-primary reply text-xs"></i>
                        Restore
                    </button>
                    <button
                        v-if="!alias.deleted_at"
                        @click.stop="deleteAlias" class="delete">
                        <i class="icon icon-error trash text-xs"></i>
                        Delete
                    </button>
                    <button
                        v-if="alias.is_custom_domain"
                        @click.stop="forgetAlias" class="delete">
                        <i class="icon icon-error trash text-xs"></i>
                        Forget
                    </button>
                </div>
            </div>
        </td>
    </tr>
    <tr v-else class="tablet-lg">
        <td>
            <div class="flex gap-2 justify-between">
                <div class="text-start">
                    <div>
                        <p class="mb-3">{{ formatDistanceToNow(new Date(alias.created_at)) }}</p>
                    </div>
                    <div>
                        <div class="hs-tooltip inline-block mb-5 break-all">
                            <p class="hs-tooltip-toggle mb-0">
                                <button class="plain truncate text-sm p-0 text-wrap text-start" @click="copyAlias(alias.name)">
                                    <span v-if="alias.description" class="block break-words">{{ truncatedDescription }}</span>
                                    <span class="text-sm break-all flex items-center"><i v-if="alias.pinned" class="icon pin icon-accent text-xs mr-1 shrink-0"></i>{{ alias.name }}</span>
                                </button>
                                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                    {{ copyText }}: {{ alias.name }}
                                </span>
                            </p>
                            <p v-if="isCreatedByWildcard" class="text-xs text-tertiary mt-1">Created by Wildcard</p>
                        </div>
                    </div>
                    <div class="flex items-center hs-tooltip">
                        <input
                            @change="updateAlias"
                            v-bind:checked="alias.enabled && alias.recipients.length > 0 && !isDomainUnverified && !isAliasDeleted"
                            v-bind:disabled="!alias.recipients.length || isDomainUnverified || isAliasDeleted"
                            type="checkbox"
                            class="checkbox-switch"
                        >
                        <span v-if="isAliasDeleted" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                            Alias deleted. Address is not forwarding mail.
                        </span>
                        <span v-else-if="isDomainUnverified" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                            Domain not verified or disabled. Address is not forwarding mail.
                        </span>
                        <span v-else-if="!alias.recipients.length" class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                            Disabled
                        </span>
                    </div>
                </div>
                <div>
                    <div class="hs-dropdown mb-3">
                        <button class="py-0" v-bind:id="'hs-dropdown-alias-edit-' + alias.id">
                            <i class="icon icon-secondary more text-lg"></i>
                        </button>
                        <div
                            class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                            v-bind:aria-labelledby="'hs-dropdown-alias-edit-' + alias.id"
                        >
                            <button
                                v-if="!alias.deleted_at"
                                v-bind:disabled="!alias.recipients.length"
                                @click="$emit('onSend', alias)"
                                v-bind:class="{ 'hide': alias.wildcard }"
                                >
                                <i class="icon icon-primary send text-xs"></i>
                                Send
                            </button>
                            <button
                                v-if="!alias.deleted_at"
                                @click="$emit('onEdit', alias)">
                                <i class="icon icon-primary edit text-xs"></i>
                                Edit
                            </button>
                            <button
                                v-if="!alias.deleted_at"
                                @click.stop="togglePin">
                                <i class="icon icon-primary pin text-xs"></i>
                                {{ alias.pinned ? 'Unpin' : 'Pin' }}
                            </button>
                            <button
                                v-if="alias.deleted_at"
                                @click.stop="restoreAlias">
                                <i class="icon icon-primary reply text-xs"></i>
                                Restore
                            </button>
                            <button
                                v-if="!alias.deleted_at"
                                @click.stop="deleteAlias" class="delete">
                                <i class="icon icon-error trash text-xs"></i>
                                Delete
                            </button>
                            <button
                                v-if="alias.is_custom_domain"
                                @click.stop="forgetAlias" class="delete">
                                <i class="icon icon-error trash text-xs"></i>
                                Forget
                            </button>
                        </div>
                    </div>
                    <div>
                        <div class="flex items-center gap-3 mb-1">
                            <p class="flex items-center gap-1 hs-tooltip mb-2">
                                {{ alias.stats.forwards }}
                                <i class="icon forward text-xs icon-tertiary"></i>
                                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.forwards }} Forwards</span>
                            </p>
                            <p class="flex items-center gap-1 hs-tooltip mb-2">
                                {{ alias.stats.blocks }}
                                <i class="icon block text-xs icon-tertiary"></i>
                                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.blocks }} Blocks</span>
                            </p>
                        </div>
                        <div class="flex items-center gap-3 mt-1">
                            <p class="flex items-center gap-1 hs-tooltip mb-2">
                                {{ alias.stats.replies }}
                                <i class="icon reply text-xs icon-tertiary"></i>
                                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.replies }} Replies</span>
                            </p>
                            <p class="flex items-center gap-1 hs-tooltip mb-2">
                                {{ alias.stats.sends }}
                                <i class="icon send text-xs icon-tertiary"></i>
                                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ alias.stats.sends }} Sends</span>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
            <hr>
        </td>
    </tr>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { aliasApi } from '../api/alias.ts'
import events from '../events.ts'
import { formatDistanceToNow } from 'date-fns'
import { closeDropdowns } from '../lib/preline.ts'
import { useBreakpoint } from '../lib/useBreakpoint.ts'

const props = defineProps(['alias', 'recipients', 'wildcard', 'selectable', 'selected'])
defineEmits(['onToggleSelect', 'onEdit', 'onSend'])
const { isDesktop } = useBreakpoint()
// Computed, not ref(props.alias): rows are keyed by alias id and survive list refreshes, so the
// row has to track the replacement object rather than the one captured at mount.
const alias = computed(() => props.alias)
const isDomainUnverified = computed(() => alias.value.is_custom_domain === true && (alias.value.is_domain_verified === false || alias.value.is_domain_enabled === false))
const isAliasDeleted = computed(() => alias.value.deleted_at !== null)
const isCreatedByWildcard = computed(() => alias.value.origin === 1)
const truncatedDescription = computed(() => {
    const desc = alias.value.description
    if (!desc) return ''
    return desc.length > 45 ? desc.slice(0, 45) + '...' : desc
})
const copyText = ref('Click to copy')

const updateAlias = async () => {
    alias.value.enabled = !alias.value.enabled
    try {
        await aliasApi.update(alias.value.id, alias.value)
        // The row may no longer belong in the list the parent is showing.
        events.emit('alias.enabled', { id: alias.value.id, enabled: alias.value.enabled })
    } catch {}
}

const togglePin = async () => {
    closeDropdowns()
    const pinned = !alias.value.pinned
    try {
        await aliasApi.pin(alias.value.id, pinned)
        events.emit('alias.update', {})
    } catch {}
}

const deleteAlias = () => {
    closeDropdowns()
    const errMessage = props.wildcard ? 'WARNING: You will not be able to create the same wildcard alias in the next 90 days. Are you sure you want to delete alias? ' : 'Are you sure you want to delete alias? A deleted email alias can be restored within 90 days.'
    if (!confirm(errMessage)) return

    events.emit('alias.delete', { id: alias.value.id, wildcard: props.wildcard })
}

const restoreAlias = async () => {
    closeDropdowns()
    try {
        await aliasApi.restore(alias.value.id)
        events.emit('alias.update', {})
    } catch {}
}

const forgetAlias = () => {
    const errMessage = 'WARNING: This operation cannot be undone. You will not be able to restore this alias. Are you sure you want to permanently delete alias?'
    if (!confirm(errMessage)) return

    events.emit('alias.forget', { id: alias.value.id })
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias).then(() => {
        copyText.value = 'Copied'
        setTimeout(() => {
            copyText.value = 'Click to copy'
        }, 2000)
    }).catch(() => {
        copyText.value = 'Failed to copy'
        setTimeout(() => {
            copyText.value = 'Click to copy'
        }, 2000)
    })
}
</script>
