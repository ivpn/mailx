<template>
    <div>
        <!-- The modal is opened from code, but Preline only binds a modal that has a toggle -->
        <button v-bind:data-hs-overlay="'#' + modalId" type="button" class="hidden" tabindex="-1" aria-hidden="true"></button>
        <div v-bind:id="modalId" class="hs-overlay hidden">
            <div class="sm:!max-w-3xl">
                <div>
                    <header>
                        <button @click="close" type="button" class="close" aria-label="Close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4 class="uppercase">
                            {{ editingId ? 'Edit rule' : 'New rule' }}
                        </h4>
                    </header>
                    <form @submit.prevent="save" autocomplete="off" novalidate>
                        <article>
                            <div class="pb-5">
                                <label for="rule_name">
                                    Name
                                </label>
                                <input
                                    v-model="form.name"
                                    v-bind:class="{ 'error': nameError }"
                                    v-bind:aria-invalid="nameError"
                                    v-bind:maxlength="ruleLimits.name"
                                    id="rule_name"
                                    placeholder="Enter rule name"
                                    type="text"
                                >
                                <p v-if="nameError" class="error">Required</p>
                            </div>

                            <div class="pb-5">
                                <p class="text-secondary text-sm font-bold mb-3 mt-0">Apply to</p>
                                <div class="flex flex-wrap gap-5">
                                    <label v-for="type in relayTypes" :key="type.value" class="font-normal mb-0 gap-2">
                                        <input
                                            v-model="form[type.value]"
                                            v-bind:id="'rule_apply_' + type.value"
                                            type="checkbox"
                                            class="checkbox-plain"
                                        >
                                        {{ type.label }}
                                    </label>
                                </div>
                                <p v-if="typesError" class="error">Required</p>
                            </div>

                            <div class="pb-5">
                                <p class="text-secondary text-sm font-bold mb-3 mt-0">Conditions</p>
                                <template v-for="(row, index) in form.conditions" :key="row.key">
                                    <div v-if="index > 0" class="flex items-center gap-2 my-2">
                                        <select
                                            v-model="form.operator"
                                            v-bind:name="'rule_operator_' + index"
                                            class="w-24 mb-0 py-1.5"
                                            aria-label="Join conditions with"
                                            data-test="operator"
                                        >
                                            <option value="AND">AND</option>
                                            <option value="OR">OR</option>
                                        </select>
                                    </div>
                                    <div
                                        class="border p-3"
                                        v-bind:class="conditionError(row) ? 'border-error' : 'border-secondary'"
                                        data-test="condition"
                                    >
                                        <div class="flex flex-col sm:flex-row sm:items-center gap-2">
                                            <span class="text-sm text-secondary shrink-0">If the</span>
                                            <select
                                                v-model="row.type"
                                                @change="onConditionTypeChange(row)"
                                                v-bind:name="'rule_condition_type_' + index"
                                                class="mb-0"
                                                aria-label="Condition"
                                                data-test="condition-type"
                                            >
                                                <option
                                                    v-for="spec in conditionSpecs"
                                                    :key="spec.type"
                                                    :value="spec.type"
                                                    :disabled="isUnsupported(spec.types, selectedTypes)"
                                                >
                                                    {{ optionLabel(spec, selectedTypes) }}
                                                </option>
                                            </select>
                                            <select
                                                v-if="!isBoolean(row)"
                                                v-model="row.match"
                                                @change="row.inputError = ''"
                                                v-bind:name="'rule_condition_match_' + index"
                                                class="mb-0 sm:w-56 sm:shrink-0"
                                                aria-label="Match"
                                                data-test="condition-match"
                                            >
                                                <option v-for="match in matchOptions" :key="match.value" :value="match.value">
                                                    {{ match.label }}
                                                </option>
                                            </select>
                                            <button
                                                v-if="form.conditions.length > 1"
                                                @click="removeCondition(index)"
                                                type="button"
                                                class="delete p-1 shrink-0 self-end sm:self-auto"
                                                aria-label="Remove condition"
                                                data-test="remove-condition"
                                            >
                                                <i class="icon trash icon-error text-sm"></i>
                                            </button>
                                        </div>
                                        <div v-if="!isBoolean(row)" class="mt-2">
                                            <div class="flex gap-2">
                                                <input
                                                    v-model="row.input"
                                                    @keydown.enter.prevent="addValue(row)"
                                                    v-bind:name="'rule_condition_value_' + index"
                                                    v-bind:maxlength="ruleLimits.valueLength"
                                                    v-bind:placeholder="valuePlaceholder(row)"
                                                    aria-label="Value"
                                                    type="text"
                                                    data-test="condition-input"
                                                >
                                                <button @click="addValue(row)" type="button" class="cta plain sm shrink-0" data-test="add-value">
                                                    Add
                                                </button>
                                            </div>
                                            <div v-if="row.values.length" class="flex flex-wrap items-center gap-1 mt-2">
                                                <template v-for="(value, valueIndex) in row.values" :key="value">
                                                    <span v-if="valueIndex > 0" class="text-xs text-tertiary">or</span>
                                                    <span class="badge gap-1 break-all" data-test="chip">
                                                        {{ value }}
                                                        <button
                                                            @click="removeValue(row, valueIndex)"
                                                            type="button"
                                                            class="p-0 flex"
                                                            v-bind:aria-label="'Remove ' + value"
                                                        >
                                                            <i class="icon close icon-secondary text-xs"></i>
                                                        </button>
                                                    </span>
                                                </template>
                                            </div>
                                        </div>
                                        <p v-if="conditionError(row)" class="error mt-2 mb-0" data-test="condition-error">{{ conditionError(row) }}</p>
                                    </div>
                                </template>
                                <p v-if="form.conditions.length > 1" class="text-xs text-tertiary mt-2 mb-0">
                                    A rule uses one operator for all of its conditions, so changing one changes them all.
                                </p>
                                <button
                                    @click="addCondition"
                                    v-bind:disabled="form.conditions.length >= ruleLimits.conditions"
                                    type="button"
                                    class="plain mt-2"
                                    data-test="add-condition"
                                >
                                    + Add condition
                                </button>
                            </div>

                            <div class="pb-5">
                                <p class="text-secondary text-sm font-bold mb-3 mt-0">Actions</p>
                                <template v-for="(row, index) in form.actions" :key="row.key">
                                    <p v-if="index > 0" class="text-xs font-bold text-tertiary my-2">AND</p>
                                    <div
                                        class="border p-3"
                                        v-bind:class="actionError(row) ? 'border-error' : 'border-secondary'"
                                        data-test="action"
                                    >
                                        <div class="flex flex-col sm:flex-row sm:items-center gap-2">
                                            <span class="text-sm text-secondary shrink-0">Then</span>
                                            <select
                                                v-model="row.type"
                                                @change="onActionTypeChange(row)"
                                                v-bind:name="'rule_action_type_' + index"
                                                class="mb-0"
                                                aria-label="Action"
                                                data-test="action-type"
                                            >
                                                <option
                                                    v-for="spec in actionSpecs"
                                                    :key="spec.type"
                                                    :value="spec.type"
                                                    :disabled="actionOptionDisabled(spec, row)"
                                                >
                                                    {{ optionLabel(spec, selectedTypes) }}
                                                </option>
                                            </select>
                                            <button
                                                v-if="form.actions.length > 1"
                                                @click="removeAction(index)"
                                                type="button"
                                                class="delete p-1 shrink-0 self-end sm:self-auto"
                                                aria-label="Remove action"
                                                data-test="remove-action"
                                            >
                                                <i class="icon trash icon-error text-sm"></i>
                                            </button>
                                        </div>
                                        <div v-if="actionSpec(row.type)?.value === 'text'" class="mt-2">
                                            <input
                                                v-model="row.value"
                                                v-bind:name="'rule_action_value_' + index"
                                                v-bind:maxlength="ruleLimits.actionValue"
                                                v-bind:placeholder="actionSpec(row.type)?.placeholder"
                                                aria-label="Action value"
                                                type="text"
                                                data-test="action-value"
                                            >
                                        </div>
                                        <div v-if="actionSpec(row.type)?.value === 'recipients'" class="mt-2 flex flex-col gap-2" role="group" aria-label="Recipients">
                                            <label v-for="recipient in recipients" :key="recipient.id" class="font-normal mb-0 gap-2 break-all">
                                                <input
                                                    v-model="row.recipients"
                                                    :value="recipient.id"
                                                    :name="'rule_action_recipients_' + index"
                                                    type="checkbox"
                                                    class="checkbox-plain shrink-0"
                                                    data-test="action-recipient"
                                                >
                                                {{ recipient.email }}
                                            </label>
                                            <p v-if="!recipients.length" class="text-sm text-tertiary m-0">
                                                You don't have any verified recipients.
                                            </p>
                                        </div>
                                        <p v-if="actionSpec(row.type)?.hint" class="text-xs text-tertiary mt-2 mb-0" data-test="action-hint">{{ actionSpec(row.type)?.hint }}</p>
                                        <p v-if="actionError(row)" class="error mt-2 mb-0" data-test="action-error">{{ actionError(row) }}</p>
                                    </div>
                                </template>
                                <button
                                    @click="addAction"
                                    v-bind:disabled="form.actions.length >= ruleLimits.actions"
                                    type="button"
                                    class="plain mt-2"
                                    data-test="add-action"
                                >
                                    + Add action
                                </button>
                            </div>
                        </article>
                        <footer>
                            <nav>
                                <button v-bind:disabled="saving" type="submit" class="cta" data-test="save">
                                    {{ editingId ? 'Save changes' : 'Create rule' }}
                                </button>
                                <button @click="close" type="button" class="cancel">
                                    Cancel
                                </button>
                            </nav>
                            <p v-if="error" class="error px-5">Error: {{ error }}</p>
                        </footer>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import axios from 'axios'
import overlay from '@preline/overlay'
import { ruleApi } from '../api/rule.ts'
import { initOverlays } from '../lib/preline.ts'
import {
    type ActionSpec,
    type RelayType,
    type Rule,
    type RuleAction,
    type RuleCondition,
    actionSpec,
    actionSpecs,
    conditionSpec,
    conditionSpecs,
    defaultMatch,
    isInactiveSubError,
    isRegexMatch,
    isUnsupported,
    matchOptions,
    onlyAvailableFor,
    optionLabel,
    relayTypes,
    repeatableActions,
    ruleLimits,
} from '../lib/rules.ts'

interface Recipient {
    id: string
    email: string
}

interface ConditionRow {
    key: number
    type: string
    match: string
    values: string[]
    input: string
    inputError: string
}

interface ActionRow {
    key: number
    type: string
    value: string
    recipients: string[]
}

const props = defineProps<{ recipients: Recipient[] }>()
const emit = defineEmits<{ saved: [], inactiveSub: [] }>()

const modalId = 'modal-rule-edit'
let nextKey = 0

const editingId = ref('')
const form = reactive({
    name: '',
    operator: 'AND' as 'AND' | 'OR',
    conditions: [] as ConditionRow[],
    actions: [] as ActionRow[],
    forwards: true,
    replies: false,
    sends: false,
})
const submitted = ref(false)
const saving = ref(false)
const error = ref('')

const selectedTypes = computed<RelayType[]>(() => relayTypes.filter(type => form[type.value]).map(type => type.value))
const nameError = computed(() => submitted.value && !form.name.trim())
const typesError = computed(() => submitted.value && !selectedTypes.value.length)

const isBoolean = (row: ConditionRow) => conditionSpec(row.type)?.boolean ?? false

const newCondition = (condition?: RuleCondition): ConditionRow => {
    const type = condition?.type ?? 'subject'
    return {
        key: nextKey++,
        type,
        match: condition?.match || (conditionSpec(type)?.boolean ? '' : defaultMatch),
        values: [...(condition?.values ?? [])],
        input: '',
        inputError: '',
    }
}

const newAction = (action?: RuleAction): ActionRow => {
    const type = action?.type ?? firstAvailableAction()
    const recipients = actionSpec(type)?.value === 'recipients'
    return {
        key: nextKey++,
        type,
        value: recipients ? '' : action?.value ?? '',
        recipients: recipients ? (action?.value ?? '').split(',').map(id => id.trim()).filter(Boolean) : [],
    }
}

const firstAvailableAction = () => {
    const available = actionSpecs.find(spec => !actionOptionDisabled(spec))
    return (available ?? actionSpecs[0]).type
}

const actionOptionDisabled = (spec: ActionSpec, row?: ActionRow) => {
    if (isUnsupported(spec.types, selectedTypes.value)) return true
    if (repeatableActions.includes(spec.type)) return false
    return form.actions.some(other => other !== row && other.type === spec.type)
}

// Recipients that were deleted or deactivated since the rule was saved are dropped
const selectedRecipients = (row: ActionRow) => row.recipients.filter(id => props.recipients.some(rcp => rcp.id === id))

const conditionError = (row: ConditionRow) => {
    const spec = conditionSpec(row.type)
    if (spec && isUnsupported(spec.types, selectedTypes.value)) return onlyAvailableFor(spec.types)
    if (row.inputError) return row.inputError
    if (submitted.value && !isBoolean(row) && !row.values.length) return 'Add at least one value'
    return ''
}

const actionError = (row: ActionRow) => {
    const spec = actionSpec(row.type)
    if (!spec) return ''
    if (isUnsupported(spec.types, selectedTypes.value)) return onlyAvailableFor(spec.types)
    if (!submitted.value) return ''
    if (spec.value === 'text' && spec.required && !row.value.trim()) return 'Required'
    if (spec.value === 'recipients' && !selectedRecipients(row).length) return 'Choose at least one recipient'
    return ''
}

const valuePlaceholder = (row: ConditionRow) => {
    if (isRegexMatch(row.match)) return 'Regular expression, then press Enter'
    return 'Type a value and press Enter'
}

const onConditionTypeChange = (row: ConditionRow) => {
    row.match = isBoolean(row) ? '' : defaultMatch
    row.values = []
    row.input = ''
    row.inputError = ''
}

const onActionTypeChange = (row: ActionRow) => {
    row.value = ''
    row.recipients = []
}

// Turns the typed value into a chip. Returns false if the value can't be added.
const addValue = (row: ConditionRow) => {
    const value = row.input.trim()
    row.inputError = ''
    if (!value) return true

    const regex = isRegexMatch(row.match)
    const key = (v: string) => regex ? v : v.toLowerCase()
    if (row.values.some(v => key(v) === key(value))) {
        row.inputError = 'This value has already been added'
        return false
    }
    if (row.values.length >= ruleLimits.values) {
        row.inputError = `A condition can have at most ${ruleLimits.values} values`
        return false
    }
    if (regex) {
        // The server is authoritative: Go's RE2 differs from JavaScript, e.g. it has no lookarounds
        try {
            new RegExp(value)
        } catch (err) {
            row.inputError = 'Invalid regex: ' + (err as Error).message
            return false
        }
    }

    row.values.push(value)
    row.input = ''
    return true
}

const removeValue = (row: ConditionRow, index: number) => {
    row.values.splice(index, 1)
    row.inputError = ''
}

const addCondition = () => {
    if (form.conditions.length >= ruleLimits.conditions) return
    form.conditions.push(newCondition())
}

const removeCondition = (index: number) => {
    if (form.conditions.length <= 1) return
    form.conditions.splice(index, 1)
}

const addAction = () => {
    if (form.actions.length >= ruleLimits.actions) return
    form.actions.push(newAction())
}

const removeAction = (index: number) => {
    if (form.actions.length <= 1) return
    form.actions.splice(index, 1)
}

const payload = () => ({
    name: form.name.trim(),
    operator: form.operator,
    conditions: form.conditions.map(row => isBoolean(row)
        ? { type: row.type, match: '', values: [] }
        : { type: row.type, match: row.match, values: [...row.values] }),
    actions: form.actions.map(row => {
        const spec = actionSpec(row.type)
        if (spec?.value === 'recipients') return { type: row.type, value: selectedRecipients(row).join(',') }
        if (spec?.value === 'text') return { type: row.type, value: row.value.trim() }
        return { type: row.type }
    }),
    forwards: form.forwards,
    replies: form.replies,
    sends: form.sends,
})

const save = async () => {
    if (saving.value) return

    submitted.value = true
    error.value = ''

    // A value that was typed but not added yet is part of the rule
    const pendingAdded = form.conditions.every(row => isBoolean(row) || addValue(row))
    const valid = pendingAdded && !nameError.value && !typesError.value &&
        form.conditions.every(row => !conditionError(row)) &&
        form.actions.every(row => !actionError(row))
    if (!valid) return

    saving.value = true
    try {
        if (editingId.value) {
            await ruleApi.update(editingId.value, payload())
        } else {
            await ruleApi.create(payload())
        }
        emit('saved')
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
            if (isInactiveSubError(error.value)) emit('inactiveSub')
        }
    } finally {
        saving.value = false
    }
}

const modal = () => document.getElementById(modalId) as HTMLElement

const open = (rule?: Rule | null) => {
    editingId.value = rule?.id ?? ''
    form.name = rule?.name ?? ''
    form.operator = rule?.operator ?? 'AND'
    form.forwards = rule ? rule.forwards : true
    form.replies = rule?.replies ?? false
    form.sends = rule?.sends ?? false
    form.conditions = rule?.conditions.length ? rule.conditions.map(newCondition) : [newCondition()]
    form.actions = []
    form.actions = rule?.actions.length ? rule.actions.map(newAction) : [newAction()]
    submitted.value = false
    saving.value = false
    error.value = ''

    overlay.open(modal())
    nextTick(() => document.getElementById('rule_name')?.focus())
}

const close = () => {
    overlay.close(modal())
}

defineExpose({ open })

onMounted(() => {
    initOverlays()
})
</script>
