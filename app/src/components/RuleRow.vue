<template>
    <tr v-if="isDesktop" class="desktop-lg" data-test="rule-row">
        <td class="w-8">
            <!-- Not a <button>: Firefox won't start a drag from one -->
            <span
                @keydown.up.prevent="move(-1)"
                @keydown.down.prevent="move(1)"
                v-bind:tabindex="readOnly ? -1 : 0"
                v-bind:aria-disabled="readOnly"
                v-bind:class="readOnly ? 'opacity-50 cursor-not-allowed' : 'cursor-grab'"
                v-bind:data-rule-id="rule.id"
                role="button"
                class="rule-handle p-1 inline-flex"
                aria-label="Reorder"
                title="Drag, or use the arrow keys, to reorder"
            >
                <i class="icon reorder icon-tertiary"></i>
            </span>
        </td>
        <td>
            <p>{{ new Date(rule.created_at).toDateString() }}</p>
        </td>
        <td class="whitespace-normal">
            <div class="flex items-center gap-1">
                <p class="m-0 break-all">{{ rule.name }}</p>
                <span v-if="rule.invalid_recipients" class="hs-tooltip inline-flex shrink-0" data-test="invalid-recipients">
                    <i class="icon alert icon-error text-xs hs-tooltip-toggle" aria-label="Missing recipient"></i>
                    <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ invalidRecipientsText }}</span>
                </span>
            </div>
            <p class="text-xs text-tertiary mt-1">{{ relayLabels }}</p>
        </td>
        <td>
            <div class="flex items-center">
                <input
                    @change="emit('onToggle', rule)"
                    v-bind:name="'rule_active_' + rule.id"
                    v-bind:checked="rule.active"
                    v-bind:disabled="readOnly"
                    v-bind:aria-label="'Active: ' + rule.name"
                    type="checkbox"
                    class="checkbox-switch"
                    data-test="rule-active"
                >
            </div>
        </td>
        <td>
            <div class="hs-tooltip inline-block">
                <p class="hs-tooltip-toggle m-0">{{ rule.applied }}</p>
                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">{{ lastAppliedText }}</span>
            </div>
        </td>
        <td>
            <div class="hs-dropdown [--offset:0]">
                <button v-bind:id="'hs-dropdown-rule-' + rule.id" aria-label="Rule actions">
                    <i class="icon icon-secondary more text-lg"></i>
                </button>
                <div
                    class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                    v-bind:aria-labelledby="'hs-dropdown-rule-' + rule.id"
                    >
                    <button @click.stop="edit" v-bind:disabled="readOnly" data-test="rule-edit">
                        <i class="icon icon-primary edit text-xs"></i>
                        Edit
                    </button>
                    <button @click.stop="remove" v-bind:disabled="readOnly" class="delete" data-test="rule-delete">
                        <i class="icon icon-error trash text-xs"></i>
                        Delete
                    </button>
                </div>
            </div>
        </td>
    </tr>
    <tr v-else class="tablet-lg" data-test="rule-row">
        <td>
            <div class="flex gap-3 justify-between">
                <div class="flex gap-3 text-start min-w-0">
                    <div>
                        <span
                            @keydown.up.prevent="move(-1)"
                            @keydown.down.prevent="move(1)"
                            v-bind:tabindex="readOnly ? -1 : 0"
                            v-bind:aria-disabled="readOnly"
                            v-bind:class="readOnly ? 'opacity-50 cursor-not-allowed' : 'cursor-grab'"
                            v-bind:data-rule-id="rule.id"
                            role="button"
                            class="rule-handle p-1 inline-flex"
                            aria-label="Reorder"
                        >
                            <i class="icon reorder icon-tertiary"></i>
                        </span>
                    </div>
                    <div class="min-w-0">
                        <p class="mb-3">{{ new Date(rule.created_at).toDateString() }}</p>
                        <div class="break-all">
                            <p class="plain text-base text-primary text-wrap">
                                <span class="text-sm text-secondary">Name:</span><br>
                                {{ rule.name }}
                                <i v-if="rule.invalid_recipients" class="icon alert icon-error text-xs" v-bind:title="invalidRecipientsText"></i>
                                <br>
                                <span class="text-xs text-tertiary">{{ relayLabels }}</span>
                            </p>
                        </div>
                        <div class="mb-5">
                            <span class="text-sm text-secondary font-primary">Active:</span>
                            <p>
                                <input
                                    @change="emit('onToggle', rule)"
                                    v-bind:name="'rule_active_' + rule.id"
                                    v-bind:checked="rule.active"
                                    v-bind:disabled="readOnly"
                                    v-bind:aria-label="'Active: ' + rule.name"
                                    type="checkbox"
                                    class="checkbox-switch"
                                    data-test="rule-active"
                                >
                            </p>
                        </div>
                        <div>
                            <p class="plain text-base text-primary">
                                <span class="text-sm text-secondary">Applied:</span><br>
                                {{ rule.applied }}
                                <span class="text-xs text-tertiary">· {{ lastAppliedText }}</span>
                            </p>
                        </div>
                    </div>
                </div>
                <div class="text-end">
                    <div class="hs-dropdown">
                        <button class="py-0" v-bind:id="'hs-dropdown-rule-' + rule.id" aria-label="Rule actions">
                            <i class="icon icon-secondary more text-lg"></i>
                        </button>
                        <div
                            class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                            v-bind:aria-labelledby="'hs-dropdown-rule-' + rule.id"
                            >
                            <button @click.stop="edit" v-bind:disabled="readOnly" data-test="rule-edit">
                                <i class="icon icon-primary edit text-xs"></i>
                                Edit
                            </button>
                            <button @click.stop="remove" v-bind:disabled="readOnly" class="delete" data-test="rule-delete">
                                <i class="icon icon-error trash text-xs"></i>
                                Delete
                            </button>
                        </div>
                    </div>
                </div>
            </div>
            <hr>
        </td>
    </tr>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useBreakpoint } from '../lib/useBreakpoint.ts'
import { closeDropdowns } from '../lib/preline.ts'
import { type Rule, ruleRelayLabels } from '../lib/rules.ts'

const props = defineProps<{ rule: Rule, readOnly: boolean }>()
const emit = defineEmits<{
    onToggle: [rule: Rule],
    onEdit: [rule: Rule],
    onDelete: [rule: Rule],
    onMove: [rule: Rule, delta: number],
}>()
const { isDesktop } = useBreakpoint()

const invalidRecipientsText = "A recipient this rule forwards to was deleted or deactivated. The alias's own recipients are used instead."

const relayLabels = computed(() => ruleRelayLabels(props.rule).join(' · '))

const lastAppliedText = computed(() => props.rule.last_applied_at
    ? 'Last applied: ' + new Date(props.rule.last_applied_at).toLocaleString()
    : 'Never')

const move = (delta: number) => {
    if (props.readOnly) return
    emit('onMove', props.rule, delta)
}

const edit = () => {
    closeDropdowns()
    emit('onEdit', props.rule)
}

const remove = () => {
    closeDropdowns()
    emit('onDelete', props.rule)
}
</script>
