import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import Settings from './Settings.vue'
import { settingsApi } from '../api/settings.ts'
import { recipientApi } from '../api/recipient.ts'

vi.mock('../api/settings.ts', () => ({
    settingsApi: {
        getDefaults: vi.fn(),
        updateDefaultDomain: vi.fn(),
        updateDefaultRecipient: vi.fn(),
        updateAliasFormat: vi.fn(),
        updateFromName: vi.fn(),
        updateRemoveHeader: vi.fn(),
        updateLogIssues: vi.fn(),
        updateBlockTrackers: vi.fn(),
        updateListUnsubscribeBehaviour: vi.fn(),
    },
}))
vi.mock('../api/recipient.ts', () => ({
    recipientApi: { getList: vi.fn() },
}))
vi.mock('../api/log.ts', () => ({
    logApi: { deleteAll: vi.fn() },
}))

const api = vi.mocked(settingsApi)

let app: App | undefined
let root: HTMLElement

const flush = async () => {
    await new Promise((resolve) => setTimeout(resolve))
    await nextTick()
}

const mount = async (settings: Record<string, unknown> = {}) => {
    api.getDefaults.mockResolvedValue({
        data: {
            id: 'settings-1',
            domain: 'mailx.net',
            recipient: '',
            from_name: '',
            alias_format: 'words',
            log_issues: false,
            remove_header: false,
            list_unsubscribe_behaviour: 0,
            block_trackers: true,
            custom_domains: [],
            ...settings,
        },
    } as any)

    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp(Settings)
    app.component('router-link', { template: '<a><slot /></a>' })
    app.mount(root)
    await flush()
}

const input = (id: string) => root.querySelector(`#${id}`) as HTMLInputElement
const section = (id: string) => input(id).parentElement as HTMLElement
const form = (id: string) => input(id).closest('form') as HTMLFormElement

const submit = async (id: string) => {
    form(id).dispatchEvent(new Event('submit'))
    await flush()
}

const choose = (id: string, value: string) => {
    const select = input(id) as unknown as HTMLSelectElement
    select.value = value
    select.dispatchEvent(new Event('change'))
}

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

beforeEach(() => {
    vi.stubEnv('VITE_DOMAINS', 'mailx.net,mailx.io')
    Object.values(api).forEach((fn) => fn.mockReset())
    vi.mocked(recipientApi.getList).mockResolvedValue({
        data: [
            { email: 'me@example.com', is_active: true },
            { email: 'pending@example.com', is_active: false },
            { email: 'work@example.com', is_active: true },
        ],
    } as any)
})

afterEach(() => {
    app?.unmount()
    root.remove()
    vi.unstubAllEnvs()
})

describe('Settings toggles', () => {
    it('saves the Mailx header toggle through its own endpoint', async () => {
        api.updateRemoveHeader.mockResolvedValue({ data: { message: 'Mailx header setting updated successfully.' } } as any)
        await mount({ remove_header: false })

        input('remove-header').click()
        await flush()

        expect(api.updateRemoveHeader).toHaveBeenCalledWith(true)
        expect(api.updateLogIssues).not.toHaveBeenCalled()
        expect(section('remove-header').querySelector('.success')?.textContent).toBe('Mailx header setting updated successfully.')
        expect(section('log-issues').querySelector('.success')).toBeNull()
    })

    it('saves the diagnostics toggle through its own endpoint', async () => {
        api.updateLogIssues.mockResolvedValue({ data: { message: 'Diagnostics setting updated successfully.' } } as any)
        await mount({ log_issues: false })

        input('log-issues').click()
        await flush()

        expect(api.updateLogIssues).toHaveBeenCalledWith(true)
        expect(api.updateRemoveHeader).not.toHaveBeenCalled()
        expect(section('log-issues').querySelector('.success')?.textContent).toBe('Diagnostics setting updated successfully.')
        expect(section('remove-header').querySelector('.success')).toBeNull()
    })

    it('reverts the toggle and shows the API error below it', async () => {
        api.updateLogIssues.mockRejectedValue(apiError('Unable to update diagnostics setting.'))
        await mount({ log_issues: true })

        input('log-issues').click()
        await flush()

        expect(input('log-issues').checked).toBe(true)
        expect(section('log-issues').querySelector('.error')?.textContent).toContain('Unable to update diagnostics setting.')
    })

    it('has no Save Settings button', async () => {
        await mount()

        const labels = [...root.querySelectorAll('button')].map((b) => b.textContent?.trim())
        expect(labels).not.toContain('Save Settings')
    })
})

describe('Settings defaults', () => {
    it('preselects the stored values and lists custom domains and active recipients', async () => {
        await mount({
            domain: 'shop.example',
            recipient: 'work@example.com',
            alias_format: 'uuid',
            from_name: 'Jane',
            custom_domains: [{ name: 'shop.example' }],
        })

        const options = (id: string) => [...input(id).querySelectorAll('option')].map((o) => o.value)
        expect(options('domain')).toEqual(['mailx.net', 'mailx.io', 'shop.example'])
        expect(options('recipient')).toEqual(['me@example.com', 'work@example.com'])
        expect(options('format')).toEqual(['words', 'random', 'uuid'])
        expect(input('domain').value).toBe('shop.example')
        expect(input('recipient').value).toBe('work@example.com')
        expect(input('format').value).toBe('uuid')
        expect(input('from-name').value).toBe('Jane')
    })

    it.each([
        ['domain', 'mailx.io', 'updateDefaultDomain', 'Default domain updated successfully.'],
        ['recipient', 'work@example.com', 'updateDefaultRecipient', 'Default recipient updated successfully.'],
        ['format', 'random', 'updateAliasFormat', 'Default alias format updated successfully.'],
    ] as const)('saves %s through its own endpoint', async (id, value, method, message) => {
        api[method].mockResolvedValue({ data: { message } } as any)
        await mount()

        choose(id, value)
        await submit(id)

        expect(api[method]).toHaveBeenCalledWith(value)
        expect(form(id).querySelector('.success')?.textContent).toBe(message)
        expect(root.querySelectorAll('.success')).toHaveLength(1)
    })

    it('saves the first option when no default is stored yet', async () => {
        api.updateDefaultRecipient.mockResolvedValue({ data: { message: 'Default recipient updated successfully.' } } as any)
        await mount({ recipient: '' })

        expect(input('recipient').value).toBe('me@example.com')
        await submit('recipient')
        expect(api.updateDefaultRecipient).toHaveBeenCalledWith('me@example.com')
    })

    it('disables the recipient form when there are no active recipients', async () => {
        vi.mocked(recipientApi.getList).mockResolvedValue({ data: [] } as any)
        await mount()

        expect(input('recipient').disabled).toBe(true)
        expect((form('recipient').querySelector('button') as HTMLButtonElement).disabled).toBe(true)
    })

    it('saves the from name, including an empty one', async () => {
        api.updateFromName.mockResolvedValue({ data: { message: 'From name updated successfully.' } } as any)
        await mount({ from_name: 'Jane' })

        input('from-name').value = '  Jane Doe '
        input('from-name').dispatchEvent(new Event('input'))
        await submit('from-name')
        expect(api.updateFromName).toHaveBeenCalledWith('  Jane Doe ')
        expect(input('from-name').value).toBe('Jane Doe')
        expect(form('from-name').querySelector('.success')?.textContent).toBe('From name updated successfully.')

        input('from-name').value = ''
        input('from-name').dispatchEvent(new Event('input'))
        await submit('from-name')
        expect(api.updateFromName).toHaveBeenLastCalledWith('')
    })

    it('shows the API error below the form that failed', async () => {
        api.updateDefaultDomain.mockRejectedValue(apiError('The selected domain is invalid.'))
        await mount()

        await submit('domain')

        expect(root.querySelector('#domain-error')?.textContent).toContain('The selected domain is invalid.')
        expect(input('domain').getAttribute('aria-invalid')).toBe('true')
        expect(form('recipient').querySelector('.error')).toBeNull()
    })
})
