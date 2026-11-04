import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import Aliases from './Aliases.vue'
import { aliasApi } from '../api/alias.ts'
import { settingsApi } from '../api/settings.ts'
import { subscriptionApi } from '../api/subscription.ts'

vi.mock('../api/alias.ts', () => ({
    aliasApi: { getList: vi.fn() },
}))
vi.mock('../api/settings.ts', () => ({
    settingsApi: { getDefaults: vi.fn() },
}))
vi.mock('../api/subscription.ts', () => ({
    subscriptionApi: { get: vi.fn() },
}))
vi.mock('../lib/preline.ts', () => ({
    initDropdowns: vi.fn(),
    initOverlays: vi.fn(),
    initTooltips: vi.fn(),
    closeDropdowns: vi.fn(),
}))
vi.mock('../lib/useBreakpoint.ts', async () => {
    const { ref } = await import('vue')
    return { useBreakpoint: () => ({ isDesktop: ref(true) }) }
})
// The modals aren't under test here and pull in Preline widgets jsdom can't run.
vi.mock('./AliasCreate.vue', () => ({ default: { render: () => null } }))
vi.mock('./AliasEdit.vue', () => ({ default: { render: () => null } }))
vi.mock('./AliasSend.vue', () => ({ default: { render: () => null } }))
vi.mock('./AliasPause.vue', () => ({ default: { render: () => null } }))

const api = vi.mocked(aliasApi)

let app: App | undefined
let root: HTMLElement

const flush = async () => {
    for (let i = 0; i < 3; i++) {
        await new Promise((resolve) => setTimeout(resolve))
        await nextTick()
    }
}

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
    is_custom_domain: false,
    is_domain_verified: null,
    is_domain_enabled: false,
    stats: { forwards: 0, blocks: 0, replies: 0, sends: 0 },
    ...overrides,
})

const mount = async () => {
    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp(Aliases)
    app.mount(root)
    await flush()
}

const emptyCard = () => root.querySelector('.card-empty')
const tableBlock = () => root.querySelector('.card-primary')!.parentElement!
const statusLabel = () => root.querySelector('#hs-dropdown-alias-status .text-accent')?.textContent?.trim()

const search = async (query: string) => {
    const input = root.querySelector<HTMLInputElement>('#input_search')!
    input.value = query
    input.dispatchEvent(new Event('input'))
    input.form!.dispatchEvent(new Event('submit'))
    await flush()
}

const pickStatus = async (label: string) => {
    const button = [...root.querySelectorAll<HTMLButtonElement>('[aria-labelledby="hs-dropdown-alias-status"] button')]
        .find((b) => b.textContent?.trim() === label)!
    button.click()
    await flush()
}

beforeEach(() => {
    api.getList.mockReset()
    vi.mocked(settingsApi.getDefaults).mockResolvedValue({
        data: { id: 's1', domain: 'example.net', domains: [], custom_domains: [], recipients: ['me@example.com'], recipient: 'me@example.com', from_name: '' },
    } as any)
    vi.mocked(subscriptionApi.get).mockResolvedValue({ data: { status: 'active' } } as any)
})

afterEach(() => {
    app?.unmount()
    root.remove()
})

describe('Aliases empty states', () => {
    it('shows the full-page empty card when the user has no aliases and no search', async () => {
        api.getList.mockResolvedValue({ data: { aliases: [], total: 0 } } as any)
        await mount()

        expect(emptyCard()?.textContent).toContain('You have no aliases yet')
        expect(tableBlock().classList.contains('hidden')).toBe(true)
    })

    it('keeps the Show filter when a search has no results, so the search can be switched to Deleted', async () => {
        api.getList.mockResolvedValue({ data: { aliases: [testAlias()], total: 1 } } as any)
        await mount()

        api.getList.mockResolvedValue({ data: { aliases: [], total: 0 } } as any)
        await search('promo@example.net')

        expect(emptyCard()).toBeNull()
        expect(tableBlock().classList.contains('hidden')).toBe(false)
        expect(root.querySelector('tbody')?.textContent).toContain('No aliases found')
        expect(statusLabel()).toBe('All')

        const deleted = testAlias({ deleted_at: '2026-10-28T10:00:00Z' })
        api.getList.mockResolvedValue({ data: { aliases: [deleted], total: 1 } } as any)
        await pickStatus('Deleted')

        expect(api.getList).toHaveBeenLastCalledWith(expect.objectContaining({ search: 'promo@example.net', status: 'deleted' }))
        expect(statusLabel()).toBe('Deleted')
        expect(root.querySelector('tbody')?.textContent).toContain('promo@example.net')
    })

    it('brings the empty card back once a search with no results is cleared on an empty account', async () => {
        api.getList.mockResolvedValue({ data: { aliases: [], total: 0 } } as any)
        await mount()
        await search('nothing')

        expect(emptyCard()).toBeNull()

        root.querySelector<HTMLButtonElement>('.icon.close')!.closest('button')!.click()
        await flush()

        expect(emptyCard()?.textContent).toContain('You have no aliases yet')
    })
})
