import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import RuleEdit from './RuleEdit.vue'
import { ruleApi } from '../api/rule.ts'
import type { Rule } from '../lib/rules.ts'

vi.mock('../api/rule.ts', () => ({
    ruleApi: {
        getList: vi.fn(),
        create: vi.fn(),
        update: vi.fn(),
        updateActive: vi.fn(),
        reorder: vi.fn(),
        delete: vi.fn(),
    },
}))
vi.mock('@preline/overlay', () => ({
    default: { autoInit: vi.fn(), open: vi.fn(), close: vi.fn() },
}))

const api = vi.mocked(ruleApi)
const recipients = [
    { id: 'rcp-1', email: 'me@example.com' },
    { id: 'rcp-2', email: 'work@example.com' },
]

let app: App | undefined
let root: HTMLElement
let editor: { open: (rule?: Rule | null) => void }
let onSaved: ReturnType<typeof vi.fn>
let onInactiveSub: ReturnType<typeof vi.fn>

const flush = async () => {
    await new Promise((resolve) => setTimeout(resolve))
    await nextTick()
}

const mount = async (rule: Rule | null = null) => {
    root = document.createElement('div')
    document.body.appendChild(root)
    onSaved = vi.fn()
    onInactiveSub = vi.fn()
    app = createApp(RuleEdit, { recipients, onSaved, onInactiveSub })
    editor = app.mount(root) as unknown as typeof editor
    editor.open(rule)
    await flush()
}

const all = <T extends Element = HTMLElement>(selector: string) => [...root.querySelectorAll<T>(selector)]
const one = <T extends Element = HTMLElement>(selector: string) => root.querySelector<T>(selector) as T
const conditions = () => all('[data-test="condition"]')
const actions = () => all('[data-test="action"]')

const choose = async (select: HTMLSelectElement, value: string) => {
    select.value = value
    select.dispatchEvent(new Event('change'))
    await flush()
}

const type = async (input: HTMLInputElement, value: string) => {
    input.value = value
    input.dispatchEvent(new Event('input'))
    await flush()
}

const enter = async (input: HTMLInputElement) => {
    input.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }))
    await flush()
}

const addChip = async (row: HTMLElement, value: string) => {
    const input = row.querySelector<HTMLInputElement>('[data-test="condition-input"]')!
    await type(input, value)
    await enter(input)
}

const click = async (el: Element) => {
    (el as HTMLElement).click()
    await flush()
}

const checkApply = async (value: string, checked: boolean) => {
    const box = one<HTMLInputElement>('#rule_apply_' + value)
    if (box.checked !== checked) await click(box)
}

const save = async () => {
    one('form').dispatchEvent(new Event('submit'))
    await flush()
}

const option = (select: HTMLSelectElement, value: string) =>
    [...select.options].find((o) => o.value === value) as HTMLOptionElement

const apiError = (message: string) => {
    const config = { headers: new AxiosHeaders() }
    return new AxiosError('Request failed with status code 400', 'ERR_BAD_REQUEST', config, undefined, {
        status: 400,
        statusText: 'Bad Request',
        headers: {},
        config,
        data: { error: message },
    })
}

const savedRule = (overrides: Partial<Rule> = {}): Rule => ({
    id: 'rule-1',
    created_at: '2026-09-01T10:00:00Z',
    name: 'Newsletters',
    order: 0,
    operator: 'OR',
    conditions: [
        { type: 'is_mailing_list', match: '', values: [] },
        { type: 'sender', match: 'ends_with', values: ['@substack.com'] },
    ],
    actions: [
        { type: 'replace_subject', value: '[News] {subject}' },
        { type: 'forward_to', value: 'rcp-2,deleted-rcp' },
    ],
    forwards: true,
    replies: false,
    sends: false,
    active: true,
    applied: 4,
    last_applied_at: null,
    ...overrides,
})

beforeEach(() => {
    Object.values(api).forEach((fn) => fn.mockReset())
    api.create.mockResolvedValue({ data: {} } as any)
    api.update.mockResolvedValue({ data: {} } as any)
})

afterEach(() => {
    app?.unmount()
    root.remove()
})

describe('RuleEdit rows', () => {
    it('adds and removes condition rows, hiding delete on the last one', async () => {
        await mount()

        expect(conditions()).toHaveLength(1)
        expect(all('[data-test="remove-condition"]')).toHaveLength(0)

        for (let i = 0; i < 4; i++) await click(one('[data-test="add-condition"]'))
        expect(conditions()).toHaveLength(5)
        expect(all('[data-test="remove-condition"]')).toHaveLength(5)
        expect(one<HTMLButtonElement>('[data-test="add-condition"]').disabled).toBe(true)

        await click(all('[data-test="remove-condition"]')[0])
        expect(conditions()).toHaveLength(4)
        expect(one<HTMLButtonElement>('[data-test="add-condition"]').disabled).toBe(false)

        for (let i = 0; i < 3; i++) await click(all('[data-test="remove-condition"]')[0])
        expect(conditions()).toHaveLength(1)
        expect(all('[data-test="remove-condition"]')).toHaveLength(0)
    })

    it('adds and removes action rows, hiding delete on the last one', async () => {
        await mount()

        expect(actions()).toHaveLength(1)
        expect(all('[data-test="remove-action"]')).toHaveLength(0)

        for (let i = 0; i < 4; i++) await click(one('[data-test="add-action"]'))
        expect(actions()).toHaveLength(5)
        expect(one<HTMLButtonElement>('[data-test="add-action"]').disabled).toBe(true)
        // Rows are joined by a static AND
        expect(all('p').filter((p) => p.textContent?.trim() === 'AND')).toHaveLength(4)

        await click(all('[data-test="remove-action"]')[4])
        expect(actions()).toHaveLength(4)
    })

    it('offers each action once, except forward to', async () => {
        await mount()

        await click(one('[data-test="add-action"]'))
        const [first, second] = all<HTMLSelectElement>('[data-test="action-type"]')
        expect(first.value).toBe('block')
        expect(option(second, 'block').disabled).toBe(true)
        expect(second.value).not.toBe('block')

        await choose(first, 'forward_to')
        await choose(second, 'forward_to')
        expect(option(first, 'forward_to').disabled).toBe(false)
        expect(second.value).toBe('forward_to')
    })

    it('keeps every operator select in sync', async () => {
        await mount()
        await click(one('[data-test="add-condition"]'))
        await click(one('[data-test="add-condition"]'))

        const operators = () => all<HTMLSelectElement>('[data-test="operator"]')
        expect(operators()).toHaveLength(2)
        expect(operators().map((s) => s.value)).toEqual(['AND', 'AND'])

        await choose(operators()[1], 'OR')
        expect(operators().map((s) => s.value)).toEqual(['OR', 'OR'])
        expect(root.textContent).toContain('changing one changes them all')
    })
})

describe('RuleEdit conditions', () => {
    it('turns typed values into chips joined by "or"', async () => {
        await mount()
        const row = conditions()[0]

        await addChip(row, ' invoice ')
        await addChip(row, 'receipt')
        await click(row.querySelector('[data-test="add-value"]')!)

        const chips = () => all('[data-test="chip"]').map((c) => c.textContent?.trim())
        expect(chips()).toEqual(['invoice', 'receipt'])
        expect([...row.querySelectorAll('span')].filter((s) => s.textContent === 'or')).toHaveLength(1)

        await addChip(row, 'INVOICE')
        expect(chips()).toEqual(['invoice', 'receipt'])
        expect(row.querySelector('[data-test="condition-error"]')?.textContent).toBe('This value has already been added')

        await click(row.querySelector('[aria-label="Remove invoice"]')!)
        expect(chips()).toEqual(['receipt'])
        expect(row.querySelector('[data-test="condition-error"]')).toBeNull()
    })

    it('checks regex values before adding them', async () => {
        await mount()
        const row = conditions()[0]
        await choose(row.querySelector('[data-test="condition-match"]')!, 'regex')

        await addChip(row, '(unclosed')
        expect(all('[data-test="chip"]')).toHaveLength(0)
        expect(row.querySelector('[data-test="condition-error"]')?.textContent).toContain('Invalid regex')

        await addChip(row, 'invoice #\\d+')
        expect(all('[data-test="chip"]').map((c) => c.textContent?.trim())).toEqual(['invoice #\\d+'])
    })

    it('resets the match and values when the condition type changes', async () => {
        await mount()
        const row = conditions()[0]
        await choose(row.querySelector('[data-test="condition-match"]')!, 'starts_with')
        await addChip(row, 'invoice')

        await choose(row.querySelector('[data-test="condition-type"]')!, 'sender')
        expect(row.querySelector<HTMLSelectElement>('[data-test="condition-match"]')!.value).toBe('contains')
        expect(all('[data-test="chip"]')).toHaveLength(0)
    })

    it('hides the match and value inputs for boolean conditions', async () => {
        await mount()
        const row = conditions()[0]

        await choose(row.querySelector('[data-test="condition-type"]')!, 'has_attachments')
        expect(row.querySelector('[data-test="condition-match"]')).toBeNull()
        expect(row.querySelector('[data-test="condition-input"]')).toBeNull()
        expect(row.textContent).toContain('If the')
    })
})

describe('RuleEdit apply to', () => {
    it('disables options that not every checked email type supports', async () => {
        await mount()
        const conditionType = one<HTMLSelectElement>('[data-test="condition-type"]')
        const actionType = one<HTMLSelectElement>('[data-test="action-type"]')

        expect(one<HTMLInputElement>('#rule_apply_forwards').checked).toBe(true)
        expect(option(conditionType, 'sender_domain').disabled).toBe(false)
        expect(option(conditionType, 'destination').disabled).toBe(true)
        expect(option(conditionType, 'destination').title).toBe('Only available for Replies and Sends')

        await checkApply('replies', true)
        expect(option(conditionType, 'sender_domain').disabled).toBe(true)
        expect(option(conditionType, 'sender_domain').title).toBe('Only available for Forwards')
        expect(option(actionType, 'blocklist_sender').disabled).toBe(true)
        expect(option(actionType, 'blocklist_sender').title).toBe('Only available for Forwards')
        expect(option(conditionType, 'subject').disabled).toBe(false)
        expect(option(actionType, 'block').disabled).toBe(false)
    })

    it('highlights rows that a newly checked type does not support instead of removing them', async () => {
        await mount()
        await choose(one('[data-test="action-type"]'), 'blocklist_sender')
        await choose(one('[data-test="condition-type"]'), 'is_mailing_list')

        await checkApply('replies', true)

        expect(actions()).toHaveLength(1)
        expect(one<HTMLSelectElement>('[data-test="action-type"]').value).toBe('blocklist_sender')
        expect(one('[data-test="action-error"]').textContent).toBe('Only available for Forwards')
        expect(actions()[0].className).toContain('border-error')
        expect(one('[data-test="condition-error"]').textContent).toBe('Only available for Forwards')

        await save()
        expect(api.create).not.toHaveBeenCalled()

        await checkApply('replies', false)
        expect(root.querySelector('[data-test="action-error"]')).toBeNull()
    })

    it('requires at least one email type', async () => {
        await mount()
        await checkApply('forwards', false)
        await type(one('#rule_name'), 'Rule')
        await addChip(conditions()[0], 'x')

        await save()
        expect(api.create).not.toHaveBeenCalled()
        expect(root.textContent).toContain('Required')
    })
})

describe('RuleEdit saving', () => {
    it('creates the rule, including a value that was typed but not added', async () => {
        await mount()
        await type(one('#rule_name'), '  Bills ')
        await addChip(conditions()[0], 'invoice')
        await type(conditions()[0].querySelector('[data-test="condition-input"]')!, 'receipt')

        await click(one('[data-test="add-condition"]'))
        await choose(conditions()[1].querySelector('[data-test="condition-type"]')!, 'has_attachments')
        await choose(one('[data-test="operator"]'), 'OR')

        await choose(one('[data-test="action-type"]'), 'replace_subject')
        await type(one('[data-test="action-value"]'), '[Bills] {subject}')
        await click(one('[data-test="add-action"]'))
        await choose(all<HTMLSelectElement>('[data-test="action-type"]')[1], 'forward_to')
        await click(all('[data-test="action-recipient"]')[1])
        await checkApply('sends', true)
        // forward_to doesn't support Sends, so the form isn't valid yet
        await save()
        expect(api.create).not.toHaveBeenCalled()
        await checkApply('sends', false)

        await save()

        expect(api.create).toHaveBeenCalledWith({
            name: 'Bills',
            operator: 'OR',
            conditions: [
                { type: 'subject', match: 'contains', values: ['invoice', 'receipt'] },
                { type: 'has_attachments', match: '', values: [] },
            ],
            actions: [
                { type: 'replace_subject', value: '[Bills] {subject}' },
                { type: 'forward_to', value: 'rcp-2' },
            ],
            forwards: true,
            replies: false,
            sends: false,
        })
        expect(onSaved).toHaveBeenCalled()
        expect(one('[data-test="save"]').textContent?.trim()).toBe('Create rule')
    })

    it('shows inline errors for missing fields', async () => {
        await mount()
        await choose(one('[data-test="action-type"]'), 'set_display_name')

        await save()

        expect(api.create).not.toHaveBeenCalled()
        expect(one('#rule_name').getAttribute('aria-invalid')).toBe('true')
        expect(one('[data-test="condition-error"]').textContent).toBe('Add at least one value')
        expect(one('[data-test="action-error"]').textContent).toBe('Required')
    })

    it('loads a saved rule and updates it, dropping recipients that no longer exist', async () => {
        await mount(savedRule())

        expect(one<HTMLInputElement>('#rule_name').value).toBe('Newsletters')
        expect(one<HTMLSelectElement>('[data-test="operator"]').value).toBe('OR')
        expect(all('[data-test="chip"]').map((c) => c.textContent?.trim())).toEqual(['@substack.com'])
        expect(all<HTMLInputElement>('[data-test="action-recipient"]').map((c) => c.checked)).toEqual([false, true])
        expect(one('[data-test="save"]').textContent?.trim()).toBe('Save changes')

        await save()

        expect(api.update).toHaveBeenCalledWith('rule-1', expect.objectContaining({
            name: 'Newsletters',
            operator: 'OR',
            actions: [
                { type: 'replace_subject', value: '[News] {subject}' },
                { type: 'forward_to', value: 'rcp-2' },
            ],
        }))
        expect(api.create).not.toHaveBeenCalled()
    })

    it('shows the server error and reports an inactive subscription', async () => {
        api.update.mockRejectedValue(apiError('Block the email can\'t be combined with Replace the subject with, because the email isn\'t delivered.'))
        await mount(savedRule())

        await save()
        expect(root.textContent).toContain('can\'t be combined with')
        expect(onSaved).not.toHaveBeenCalled()
        expect(onInactiveSub).not.toHaveBeenCalled()

        api.update.mockRejectedValue(apiError('Your subscription is not active. Please renew to update rules.'))
        await save()
        expect(root.textContent).toContain('Please renew to update rules.')
        expect(onInactiveSub).toHaveBeenCalled()
    })

    it('escapes markup in names and values', async () => {
        await mount(savedRule({ conditions: [{ type: 'subject', match: 'is', values: ['<img src=x onerror=alert(1)>'] }] }))

        expect(root.querySelector('img')).toBeNull()
        expect(one('[data-test="chip"]').textContent).toContain('<img src=x onerror=alert(1)>')
    })
})
