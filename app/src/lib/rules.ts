// Mirrors the condition and action tables in api/internal/model/rule.go. The server validates
// every rule again, so this only drives what the editor offers.

export type RelayType = 'forwards' | 'replies' | 'sends'

export interface RuleCondition {
    type: string
    match: string
    values: string[]
}

export interface RuleAction {
    type: string
    value?: string
}

export interface Rule {
    id: string
    created_at: string
    name: string
    order: number
    operator: 'AND' | 'OR'
    conditions: RuleCondition[]
    actions: RuleAction[]
    forwards: boolean
    replies: boolean
    sends: boolean
    active: boolean
    applied: number
    last_applied_at: string | null
    invalid_recipients?: boolean
}

export const ruleLimits = {
    conditions: 5,
    actions: 5,
    values: 20,
    valueLength: 255,
    name: 50,
    actionValue: 100,
}

export const relayTypes: { value: RelayType, label: string }[] = [
    { value: 'forwards', label: 'Forwards' },
    { value: 'replies', label: 'Replies' },
    { value: 'sends', label: 'Sends' },
]

const allTypes: RelayType[] = ['forwards', 'replies', 'sends']
const forwardsOnly: RelayType[] = ['forwards']

export interface ConditionSpec {
    type: string
    // Reads after "If the"
    label: string
    boolean: boolean
    types: RelayType[]
}

export const conditionSpecs: ConditionSpec[] = [
    { type: 'sender', label: 'sender email', boolean: false, types: allTypes },
    { type: 'sender_domain', label: 'sender domain', boolean: false, types: forwardsOnly },
    { type: 'sender_name', label: 'sender name', boolean: false, types: forwardsOnly },
    { type: 'subject', label: 'subject', boolean: false, types: allTypes },
    { type: 'alias', label: 'alias', boolean: false, types: allTypes },
    { type: 'alias_description', label: 'alias description', boolean: false, types: allTypes },
    { type: 'destination', label: 'destination email', boolean: false, types: ['replies', 'sends'] },
    { type: 'is_mailing_list', label: 'email is from a mailing list / newsletter', boolean: true, types: forwardsOnly },
    { type: 'has_attachments', label: 'email has attachments', boolean: true, types: allTypes },
    { type: 'alias_is_wildcard', label: 'alias is a wildcard alias', boolean: true, types: allTypes },
]

export const matchOptions = [
    { value: 'is', label: 'is exactly' },
    { value: 'is_not', label: 'is not' },
    { value: 'contains', label: 'contains' },
    { value: 'not_contains', label: 'does not contain' },
    { value: 'starts_with', label: 'starts with' },
    { value: 'not_starts_with', label: 'does not start with' },
    { value: 'ends_with', label: 'ends with' },
    { value: 'not_ends_with', label: 'does not end with' },
    { value: 'regex', label: 'matches regex' },
    { value: 'not_regex', label: 'does not match regex' },
]

export const defaultMatch = 'contains'

export const isRegexMatch = (match: string) => match === 'regex' || match === 'not_regex'

export interface ActionSpec {
    type: string
    // Reads after "Then"
    label: string
    value: 'none' | 'text' | 'recipients'
    required: boolean
    types: RelayType[]
    placeholder?: string
    hint?: string
}

export const actionSpecs: ActionSpec[] = [
    { type: 'block', label: 'block the email', value: 'none', required: false, types: allTypes },
    { type: 'blocklist_sender', label: 'add sender email to the blocklist', value: 'none', required: false, types: forwardsOnly },
    { type: 'blocklist_domain', label: 'add sender domain to the blocklist', value: 'none', required: false, types: forwardsOnly },
    {
        type: 'replace_subject', label: 'replace the subject with', value: 'text', required: true, types: allTypes,
        placeholder: '[News] {subject}', hint: '{subject} is replaced with the current subject.',
    },
    { type: 'set_display_name', label: 'set the sender name to', value: 'text', required: true, types: allTypes, placeholder: 'Sender name' },
    { type: 'forward_to', label: 'forward to recipient(s)', value: 'recipients', required: true, types: forwardsOnly },
    { type: 'disable_encryption', label: 'turn off PGP encryption', value: 'none', required: false, types: forwardsOnly },
    { type: 'hide_banner', label: 'hide the Mailx banner', value: 'none', required: false, types: forwardsOnly },
    { type: 'show_banner', label: 'show the Mailx banner', value: 'none', required: false, types: forwardsOnly },
    { type: 'deactivate_alias', label: 'deactivate the alias', value: 'none', required: false, types: forwardsOnly },
    {
        type: 'set_alias_description', label: 'set the alias description to', value: 'text', required: false, types: forwardsOnly,
        placeholder: 'Leave empty to clear it',
    },
    { type: 'remove_attachments', label: 'remove attachments', value: 'none', required: false, types: forwardsOnly },
    {
        type: 'stop', label: 'stop processing more rules', value: 'none', required: false, types: allTypes,
        hint: 'Note: other rules below this one won\'t run on emails this rule matches.',
    },
]

// forward_to is the only action that can be used more than once in a rule
export const repeatableActions = ['forward_to']

export const conditionSpec = (type: string) => conditionSpecs.find(spec => spec.type === type)

export const actionSpec = (type: string) => actionSpecs.find(spec => spec.type === type)

export const joinLabels = (labels: string[]) => {
    if (labels.length <= 1) return labels.join('')
    return labels.slice(0, -1).join(', ') + ' and ' + labels[labels.length - 1]
}

// Whether a condition or action that supports `types` can't be used with every selected type
export const isUnsupported = (types: RelayType[], selected: RelayType[]) => selected.some(type => !types.includes(type))

const typeLabels = (types: RelayType[]) => relayTypes.filter(t => types.includes(t.value)).map(t => t.label)

export const onlyAvailableFor = (types: RelayType[]) => 'Only available for ' + joinLabels(typeLabels(types))

// Disabled options don't show a tooltip in most browsers, so the label says why it can't be picked
export const optionLabel = (spec: { label: string, types: RelayType[] }, selected: RelayType[]) =>
    isUnsupported(spec.types, selected)
        ? `${spec.label} (${joinLabels(typeLabels(spec.types)).toLowerCase()} only)`
        : spec.label

export const ruleRelayLabels = (rule: Pick<Rule, RelayType>) =>
    relayTypes.filter(t => rule[t.value]).map(t => t.label)

export const isInactiveSubError = (message: string) => message.startsWith('Your subscription is not active')

// The server rejects a reorder sent from a list that's out of date, e.g. after a change in another tab
export const isStaleRulesError = (message: string) => message.startsWith('Your rules have changed')
