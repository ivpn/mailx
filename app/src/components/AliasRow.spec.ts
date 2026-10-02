import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, h, nextTick, type App } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import AliasRow from './AliasRow.vue'
import { aliasApi } from '../api/alias.ts'
import events from '../events.ts'
import { formatPausedUntil } from '../lib/aliasPause.ts'

vi.mock('../api/alias.ts', () => ({
    aliasApi: { update: vi.fn(), pin: vi.fn(), pause: vi.fn(), unpause: vi.fn(), restore: vi.fn() },
}))
vi.mock('../lib/preline.ts', () => ({
    closeDropdowns: vi.fn(),
}))
const breakpoint = vi.hoisted(() => ({ desktop: true }))
vi.mock('../lib/useBreakpoint.ts', async () => {
    const { ref } = await import('vue')
    return { useBreakpoint: () => ({ isDesktop: ref(breakpoint.desktop) }) }
})

let app: App | undefined
let root: HTMLElement

const flush = async () => {
    for (let i = 0; i < 3; i++) {
        await Promise.resolve()
        await nextTick()
    }
}

const inHours = (hours: number) => new Date(Date.now() + hours * 60 * 60 * 1000).toISOString()

const testAlias = (overrides: Record<string, any> = {}) => ({
    id: 'alias-1',
    created_at: '2026-09-01T10:00:00Z',
    deleted_at: null,
    name: 'promo@example.net',
    enabled: true,
    description: '',
    recipients: 'me@example.com',
    from_name: '',
    pinned: false,
    paused_until: null,
    origin: 0,
    is_custom_domain: false,
    is_domain_verified: null,
    is_domain_enabled: false,
    stats: { forwards: 0, blocks: 0, replies: 0, sends: 0 },
    ...overrides,
})

const mount = async (alias: Record<string, any>, props: Record<string, any> = {}) => {
    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp({ render: () => h('table', [h('tbody', [h(AliasRow, { alias, wildcard: false, ...props })])]) })
    app.mount(root)
    await flush()
}

const menuLabels = () => [...root.querySelectorAll('.hs-dropdown-menu button')].map((b) => b.textContent?.trim())
const pauseButton = () => root.querySelector<HTMLButtonElement>('[data-test="alias-pause"]')!
const pauseTooltip = () => pauseButton().parentElement?.querySelector('[role="tooltip"]')?.textContent?.trim()
const text = (selector: string) => root.querySelector(selector)?.textContent?.trim()

beforeEach(() => {
    Object.values(vi.mocked(aliasApi)).forEach((fn) => fn.mockReset())
    breakpoint.desktop = true
})

afterEach(() => {
    app?.unmount()
    root.remove()
    vi.useRealTimers()
    vi.restoreAllMocks()
})

describe.each([['desktop', true], ['mobile', false]])('AliasRow pause (%s)', (_, desktop) => {
    beforeEach(() => {
        breakpoint.desktop = desktop
    })

    it('offers Pause directly after Pin, and Unpause while paused', async () => {
        await mount(testAlias())
        expect(menuLabels()).toEqual(['Send', 'Edit', 'Pin', 'Pause', 'Delete'])
        app?.unmount()
        root.remove()

        await mount(testAlias({ pinned: true, paused_until: inHours(3) }))
        expect(menuLabels()).toEqual(['Send', 'Edit', 'Unpin', 'Unpause', 'Delete'])
    })

    it('hides the action for a deleted alias', async () => {
        await mount(testAlias({ deleted_at: '2026-09-10T10:00:00Z' }))
        expect(root.querySelector('[data-test="alias-pause"]')).toBeNull()
    })

    it('replaces the Catch-All label while paused', async () => {
        const pausedUntil = inHours(3)
        await mount(testAlias({ origin: 1, paused_until: pausedUntil }))

        expect(text('[data-test="alias-paused-until"]')).toBe('Paused until ' + formatPausedUntil(pausedUntil))
        expect(root.querySelector('[data-test="alias-catch-all"]')).toBeNull()
    })

    it('shows the Catch-All label for a pause that has ended', async () => {
        await mount(testAlias({ origin: 1, paused_until: inHours(-1) }))

        expect(root.querySelector('[data-test="alias-paused-until"]')).toBeNull()
        expect(text('[data-test="alias-catch-all"]')).toBe('Created by Catch-All')
        expect(menuLabels()).toContain('Pause')
    })
})

describe('AliasRow pause', () => {
    it('drops the label when the pause ends while the page is open', async () => {
        vi.useFakeTimers({ toFake: ['setInterval', 'clearInterval', 'Date'] })
        vi.setSystemTime(new Date('2026-10-02T10:00:00Z'))
        await mount(testAlias({ origin: 1, paused_until: '2026-10-02T10:00:30Z' }))
        expect(root.querySelector('[data-test="alias-paused-until"]')).not.toBeNull()

        vi.advanceTimersByTime(60 * 1000)
        await flush()

        expect(root.querySelector('[data-test="alias-paused-until"]')).toBeNull()
        expect(text('[data-test="alias-catch-all"]')).toBe('Created by Catch-All')
        expect(pauseButton().textContent?.trim()).toBe('Pause')
    })

    it('Pause asks the parent to open the modal', async () => {
        const onPause = vi.fn()
        const alias = testAlias()
        await mount(alias, { onOnPause: onPause })

        pauseButton().click()
        await flush()

        expect(onPause).toHaveBeenCalledWith(alias)
        expect(aliasApi.unpause).not.toHaveBeenCalled()
    })

    it('Unpause calls DELETE and reloads the list', async () => {
        vi.mocked(aliasApi.unpause).mockResolvedValue({ data: {} } as any)
        const emit = vi.spyOn(events, 'emit')
        await mount(testAlias({ paused_until: inHours(3) }))

        pauseButton().click()
        await flush()

        expect(aliasApi.unpause).toHaveBeenCalledWith('alias-1')
        expect(emit).toHaveBeenCalledWith('alias.update', {})
    })

    it('passes an Unpause error to the parent', async () => {
        const config = { headers: new AxiosHeaders() }
        const message = 'Your subscription is not active. Please renew to unpause aliases.'
        vi.mocked(aliasApi.unpause).mockRejectedValue(new AxiosError('Request failed', 'ERR_BAD_REQUEST', config, undefined, {
            status: 400, statusText: 'Bad Request', headers: {}, config, data: { error: message },
        }))
        const onActionError = vi.fn()
        await mount(testAlias({ paused_until: inHours(3) }), { onOnActionError: onActionError })

        pauseButton().click()
        await flush()

        expect(onActionError).toHaveBeenCalledWith(message)
    })

    it('is disabled with a tooltip while the subscription is limited', async () => {
        await mount(testAlias(), { readOnlyPause: true })

        expect(pauseButton().disabled).toBe(true)
        expect(pauseTooltip()).toBe('Renew your subscription to pause aliases.')
    })

    it('is enabled without a tooltip otherwise', async () => {
        await mount(testAlias(), { readOnlyPause: false })

        expect(pauseButton().disabled).toBe(false)
        expect(pauseTooltip()).toBeUndefined()
    })

    it('tells a paused alias when the pause ends while the subscription is limited', async () => {
        const onPause = vi.fn()
        const pausedUntil = inHours(5)
        await mount(testAlias({ paused_until: pausedUntil }), { readOnlyPause: true, onOnPause: onPause })

        expect(pauseButton().disabled).toBe(true)
        expect(pauseButton().textContent?.trim()).toBe('Unpause')
        expect(pauseTooltip()).toBe(`Renew your subscription to unpause. The pause ends automatically on ${formatPausedUntil(pausedUntil)}.`)

        pauseButton().dispatchEvent(new MouseEvent('click'))
        await flush()
        expect(aliasApi.unpause).not.toHaveBeenCalled()
        expect(onPause).not.toHaveBeenCalled()
    })
})
