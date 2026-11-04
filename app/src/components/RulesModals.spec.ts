import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import Rules from './Rules.vue'
import { ruleApi } from '../api/rule.ts'
import { recipientApi } from '../api/recipient.ts'
import { subscriptionApi } from '../api/subscription.ts'
import type { Rule } from '../lib/rules.ts'

// Unlike Rules.spec.ts, Preline's overlay is real here: it only opens a modal it has bound to a
// toggle, which a mocked overlay can't catch.
vi.mock('../api/rule.ts', () => ({
    ruleApi: { getList: vi.fn(), create: vi.fn(), update: vi.fn(), updateActive: vi.fn(), reorder: vi.fn(), delete: vi.fn() },
}))
vi.mock('../api/recipient.ts', () => ({
    recipientApi: { getList: vi.fn() },
}))
vi.mock('../api/subscription.ts', () => ({
    subscriptionApi: { get: vi.fn() },
}))
vi.mock('../lib/useBreakpoint.ts', async () => {
    const { ref } = await import('vue')
    return { useBreakpoint: () => ({ isDesktop: ref(true) }) }
})

let app: App | undefined
let root: HTMLElement

const wait = async (ms = 0) => {
    await new Promise((resolve) => setTimeout(resolve, ms))
    await nextTick()
}

const rule: Rule = {
    id: 'rule-1',
    created_at: '2026-09-01T10:00:00Z',
    name: 'Bills',
    order: 0,
    operator: 'AND',
    conditions: [{ type: 'subject', match: 'contains', values: ['invoice'] }],
    actions: [{ type: 'block' }],
    forwards: true,
    replies: false,
    sends: false,
    active: true,
    applied: 0,
    last_applied_at: null,
}

const isOpen = (id: string) => {
    const modal = document.getElementById(id)!
    return !modal.classList.contains('hidden') && modal.classList.contains('open')
}

beforeEach(async () => {
    vi.mocked(ruleApi.getList).mockResolvedValue({ data: [rule] } as any)
    vi.mocked(recipientApi.getList).mockResolvedValue({ data: [] } as any)
    vi.mocked(subscriptionApi.get).mockResolvedValue({ data: { status: 'active' } } as any)

    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp(Rules)
    app.mount(root)
    await wait()
})

afterEach(() => {
    app?.unmount()
    root.remove()
    document.querySelectorAll('[data-hs-overlay-backdrop-template], .hs-overlay-backdrop').forEach((el) => el.remove())
    document.body.style.overflow = ''
})

describe('Rules modals', () => {
    it('opens the New Rule modal', async () => {
        expect(isOpen('modal-rule-edit')).toBe(false)

        root.querySelector<HTMLButtonElement>('[data-test="new-rule"]')!.click()
        await wait(100)

        expect(isOpen('modal-rule-edit')).toBe(true)
        expect(root.querySelector('#modal-rule-edit h4')?.textContent).toContain('New rule')
    })

    it('opens the edit modal with the rule loaded', async () => {
        root.querySelector<HTMLButtonElement>('[data-test="rule-edit"]')!.click()
        await wait(100)

        expect(isOpen('modal-rule-edit')).toBe(true)
        expect(root.querySelector<HTMLInputElement>('#rule_name')!.value).toBe('Bills')
    })

    it('opens the delete confirmation', async () => {
        root.querySelector<HTMLButtonElement>('[data-test="rule-delete"]')!.click()
        await wait(100)

        expect(isOpen('modal-delete-rule')).toBe(true)
        expect(root.querySelector('#modal-delete-rule')?.textContent).toContain('Bills')
    })

    it('gives every field in the editor an id or a name', async () => {
        const fields = [...root.querySelectorAll('#modal-rule-edit input, #modal-rule-edit select')]
        expect(fields.length).toBeGreaterThan(0)
        expect(fields.filter((f) => !f.id && !f.getAttribute('name'))).toEqual([])
    })
})
