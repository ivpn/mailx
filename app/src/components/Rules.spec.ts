import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import overlay from '@preline/overlay'
import Rules from './Rules.vue'
import { ruleApi } from '../api/rule.ts'
import { recipientApi } from '../api/recipient.ts'
import { subscriptionApi } from '../api/subscription.ts'
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
vi.mock('../api/recipient.ts', () => ({
    recipientApi: { getList: vi.fn() },
}))
vi.mock('../api/subscription.ts', () => ({
    subscriptionApi: { get: vi.fn() },
}))
vi.mock('@preline/overlay', () => ({
    default: { autoInit: vi.fn(), open: vi.fn(), close: vi.fn() },
}))
vi.mock('../lib/preline.ts', () => ({
    initDropdowns: vi.fn(),
    initOverlays: vi.fn(),
    closeDropdowns: vi.fn(),
}))
vi.mock('../lib/useBreakpoint.ts', async () => {
    const { ref } = await import('vue')
    return { useBreakpoint: () => ({ isDesktop: ref(true) }) }
})

// Sortable needs real pointer events, so the drag-and-drop wrapper is replaced by one that renders
// its rows and lets a test emit the events a drop would.
const draggable = vi.hoisted(() => ({ props: {} as Record<string, any>, emit: (() => {}) as (event: string, ...args: any[]) => void }))
vi.mock('vue-draggable-plus', async () => {
    const { defineComponent, h } = await import('vue')
    return {
        VueDraggable: defineComponent({
            props: ['modelValue', 'disabled', 'tag', 'handle', 'animation'],
            emits: ['update:modelValue', 'start', 'end'],
            setup(props, { slots, emit }) {
                draggable.props = props
                draggable.emit = emit as typeof draggable.emit
                return () => h(props.tag || 'div', slots.default?.())
            },
        }),
    }
})

const api = vi.mocked(ruleApi)

let app: App | undefined
let root: HTMLElement

const flush = async () => {
    for (let i = 0; i < 3; i++) {
        await new Promise((resolve) => setTimeout(resolve))
        await nextTick()
    }
}

const rule = (id: string, overrides: Partial<Rule> = {}): Rule => ({
    id,
    created_at: '2026-09-01T10:00:00Z',
    name: 'Rule ' + id,
    order: 0,
    operator: 'AND',
    conditions: [{ type: 'subject', match: 'contains', values: ['x'] }],
    actions: [{ type: 'block' }],
    forwards: true,
    replies: false,
    sends: false,
    active: true,
    applied: 0,
    last_applied_at: null,
    ...overrides,
})

const mount = async (rules: Rule[], status = 'active') => {
    api.getList.mockResolvedValue({ data: rules } as any)
    vi.mocked(subscriptionApi.get).mockResolvedValue({ data: { status } } as any)

    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp(Rules)
    app.mount(root)
    await flush()
}

const all = <T extends Element = HTMLElement>(selector: string) => [...root.querySelectorAll<T>(selector)]
const one = <T extends Element = HTMLElement>(selector: string) => root.querySelector<T>(selector) as T
const rowNames = () => all('[data-test="rule-row"]').map((row) => row.querySelectorAll('td')[2].querySelector('p')?.textContent?.trim())
const toggles = () => all<HTMLInputElement>('[data-test="rule-active"]')
const handles = () => all('.rule-handle')

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

const press = async (el: HTMLElement, key: string) => {
    el.dispatchEvent(new KeyboardEvent('keydown', { key, bubbles: true }))
    await flush()
}

const drop = async (order: Rule[]) => {
    draggable.emit('start')
    draggable.emit('update:modelValue', order)
    draggable.emit('end')
    await flush()
}

beforeEach(() => {
    Object.values(api).forEach((fn) => fn.mockReset())
    vi.mocked(recipientApi.getList).mockResolvedValue({
        data: [
            { id: 'rcp-1', email: 'me@example.com', is_active: true },
            { id: 'rcp-2', email: 'pending@example.com', is_active: false },
        ],
    } as any)
    vi.mocked(overlay.open).mockReset()
})

afterEach(() => {
    app?.unmount()
    root.remove()
})

describe('Rules list', () => {
    it('shows the empty state', async () => {
        await mount([])

        expect(one('[data-test="empty"]').textContent).toContain("You don't have any rules yet")
        expect(all('[data-test="rule-row"]')).toHaveLength(0)
    })

    it('lists the rules in order with their email types, counts and warnings', async () => {
        await mount([
            rule('a', { name: 'Newsletters', forwards: true, sends: true, applied: 3, last_applied_at: '2026-09-30T08:00:00Z' }),
            rule('b', { name: '<b>Bills</b>', invalid_recipients: true }),
        ])

        expect(rowNames()).toEqual(['Newsletters', '<b>Bills</b>'])
        expect(root.querySelector('b')).toBeNull()

        const [first, second] = all('[data-test="rule-row"]')
        expect(first.textContent).toContain('Forwards · Sends')
        expect(first.textContent).toContain('Last applied:')
        expect(first.querySelectorAll('td')[4].textContent).toContain('3')
        expect(first.querySelector('[data-test="invalid-recipients"]')).toBeNull()
        expect(second.textContent).toContain('Never')
        expect(second.querySelector('[data-test="invalid-recipients"]')).not.toBeNull()
        expect(handles().every((h) => h.getAttribute('aria-label') === 'Reorder' && h.getAttribute('role') === 'button')).toBe(true)
        expect(handles().every((h) => h.tabIndex === 0 && h.classList.contains('cursor-grab'))).toBe(true)
    })

    it('opens the editor for a new rule and for an existing one', async () => {
        await mount([rule('a')])

        one<HTMLButtonElement>('[data-test="new-rule"]').click()
        await flush()
        expect(one<HTMLInputElement>('#rule_name').value).toBe('')

        one<HTMLButtonElement>('[data-test="rule-edit"]').click()
        await flush()
        expect(one<HTMLInputElement>('#rule_name').value).toBe('Rule a')
        expect(overlay.open).toHaveBeenCalledTimes(2)
    })

    it('deletes a rule after confirmation', async () => {
        api.delete.mockResolvedValue({ data: {} } as any)
        await mount([rule('a'), rule('b')])

        all<HTMLButtonElement>('[data-test="rule-delete"]')[1].click()
        await flush()
        expect(overlay.open).toHaveBeenCalled()
        expect(api.delete).not.toHaveBeenCalled()
        expect(root.textContent).toContain('delete the rule Rule b')

        api.getList.mockResolvedValue({ data: [rule('a')] } as any)
        one<HTMLButtonElement>('[data-test="confirm-delete"]').click()
        await flush()
        expect(api.delete).toHaveBeenCalledWith('b')
        expect(rowNames()).toEqual(['Rule a'])
    })
})

describe('Rules active toggle', () => {
    it('updates optimistically', async () => {
        api.updateActive.mockResolvedValue({ data: {} } as any)
        await mount([rule('a')])

        toggles()[0].click()
        expect(toggles()[0].checked).toBe(false)
        await flush()

        expect(api.updateActive).toHaveBeenCalledWith('a', false)
        expect(toggles()[0].checked).toBe(false)
    })

    it('reverts and shows the error when the update fails', async () => {
        api.updateActive.mockRejectedValue(apiError('Unable to update rule. Please try again.'))
        await mount([rule('a', { active: false })])

        toggles()[0].click()
        await flush()

        expect(api.updateActive).toHaveBeenCalledWith('a', true)
        expect(toggles()[0].checked).toBe(false)
        expect(one('[data-test="rules-error"]').textContent).toContain('Unable to update rule.')
    })
})

describe('Rules reordering', () => {
    it('moves a rule with the arrow keys and saves the new order', async () => {
        api.reorder.mockResolvedValue({ data: {} } as any)
        await mount([rule('a'), rule('b'), rule('c')])

        await press(handles()[0], 'ArrowDown')
        expect(rowNames()).toEqual(['Rule b', 'Rule a', 'Rule c'])
        expect(api.reorder).toHaveBeenLastCalledWith(['b', 'a', 'c'])
        expect(document.activeElement?.getAttribute('data-rule-id')).toBe('a')

        await press(handles()[2], 'ArrowUp')
        expect(api.reorder).toHaveBeenLastCalledWith(['b', 'c', 'a'])

        // Already at the top
        await press(handles()[0], 'ArrowUp')
        expect(api.reorder).toHaveBeenCalledTimes(2)
    })

    it('saves the order after a drop and reverts it on error', async () => {
        const rules = [rule('a'), rule('b'), rule('c')]
        api.reorder.mockResolvedValueOnce({ data: {} } as any)
        await mount(rules)

        await drop([rules[2], rules[0], rules[1]])
        expect(api.reorder).toHaveBeenLastCalledWith(['c', 'a', 'b'])
        expect(rowNames()).toEqual(['Rule c', 'Rule a', 'Rule b'])

        api.reorder.mockRejectedValueOnce(apiError('Unable to reorder rules. Please try again.'))
        await drop([rules[0], rules[2], rules[1]])
        expect(api.reorder).toHaveBeenLastCalledWith(['a', 'c', 'b'])
        expect(rowNames()).toEqual(['Rule c', 'Rule a', 'Rule b'])
        expect(one('[data-test="rules-error"]').textContent).toContain('Unable to reorder rules.')
    })

    it('does not save a drop that leaves the order unchanged', async () => {
        const rules = [rule('a'), rule('b')]
        await mount(rules)

        await drop([...rules])
        expect(api.reorder).not.toHaveBeenCalled()
    })

    it('reloads the list when the server says it is out of date', async () => {
        api.reorder.mockRejectedValue(apiError('Your rules have changed since this page was loaded. Please reload and try again.'))
        await mount([rule('a'), rule('b')])
        api.getList.mockResolvedValue({ data: [rule('a'), rule('b'), rule('new')] } as any)

        await press(handles()[0], 'ArrowDown')

        expect(rowNames()).toEqual(['Rule a', 'Rule b', 'Rule new'])
        expect(one('[data-test="rules-error"]').textContent).toContain('Your rules have changed')
    })
})

describe('Rules with an inactive subscription', () => {
    it('is read-only with a renew notice', async () => {
        await mount([rule('a'), rule('b')], 'limited_access')

        expect(one('[data-test="read-only-notice"]').textContent).toContain('Your existing rules are still applied')
        expect(one('[data-test="read-only-notice"] a')).not.toBeNull()
        expect(one<HTMLButtonElement>('[data-test="new-rule"]').disabled).toBe(true)
        expect(all<HTMLButtonElement>('[data-test="rule-edit"]').every((b) => b.disabled)).toBe(true)
        expect(all<HTMLButtonElement>('[data-test="rule-delete"]').every((b) => b.disabled)).toBe(true)
        expect(toggles().every((t) => t.disabled)).toBe(true)
        expect(handles().every((h) => h.getAttribute('aria-disabled') === 'true' && h.tabIndex === -1)).toBe(true)
        expect(draggable.props.disabled).toBe(true)

        await press(handles()[0], 'ArrowDown')
        expect(api.reorder).not.toHaveBeenCalled()
        expect(rowNames()).toEqual(['Rule a', 'Rule b'])
    })

    it.each(['active', 'grace_period'])('is editable while the subscription is %s', async (status) => {
        await mount([rule('a')], status)

        expect(root.querySelector('[data-test="read-only-notice"]')).toBeNull()
        expect(one<HTMLButtonElement>('[data-test="new-rule"]').disabled).toBe(false)
        expect(toggles()[0].disabled).toBe(false)
        expect(draggable.props.disabled).toBe(false)
    })

    it('shows the error and reloads when the subscription ran out while the page was open', async () => {
        api.updateActive.mockRejectedValue(apiError('Your subscription is not active. Please renew to update rules.'))
        await mount([rule('a')])
        vi.mocked(subscriptionApi.get).mockResolvedValue({ data: { status: 'limited_access' } } as any)

        toggles()[0].click()
        await flush()

        expect(one('[data-test="rules-error"]').textContent).toContain('Please renew to update rules.')
        expect(toggles()[0].checked).toBe(true)
        expect(api.getList).toHaveBeenCalledTimes(2)
        expect(one('[data-test="read-only-notice"]')).not.toBeNull()
        expect(toggles()[0].disabled).toBe(true)
    })
})
