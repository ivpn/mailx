import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import Settings from './Settings.vue'
import { settingsApi } from '../api/settings.ts'

vi.mock('../api/settings.ts', () => ({
    settingsApi: {
        getDefaults: vi.fn(),
        update: vi.fn(),
        updateRemoveHeader: vi.fn(),
        updateLogIssues: vi.fn(),
        updateBlockTrackers: vi.fn(),
        updateListUnsubscribeBehaviour: vi.fn(),
    },
}))
vi.mock('../api/recipient.ts', () => ({
    recipientApi: { getList: vi.fn().mockResolvedValue({ data: [] }) },
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
    vi.stubEnv('VITE_DOMAINS', 'mailx.net')
    Object.values(api).forEach((fn) => fn.mockReset())
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
        expect(api.update).not.toHaveBeenCalled()
        expect(section('remove-header').querySelector('.success')?.textContent).toBe('Mailx header setting updated successfully.')
        expect(section('log-issues').querySelector('.success')).toBeNull()
    })

    it('saves the diagnostics toggle through its own endpoint', async () => {
        api.updateLogIssues.mockResolvedValue({ data: { message: 'Diagnostics setting updated successfully.' } } as any)
        await mount({ log_issues: false })

        input('log-issues').click()
        await flush()

        expect(api.updateLogIssues).toHaveBeenCalledWith(true)
        expect(api.update).not.toHaveBeenCalled()
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

    it('keeps the Save Settings message at the bottom', async () => {
        api.update.mockResolvedValue({ data: { message: 'Settings updated successfully.' } } as any)
        await mount()

        const save = [...root.querySelectorAll('button')].find((b) => b.textContent?.includes('Save Settings'))!
        save.click()
        await flush()

        expect(section('remove-header').querySelector('.success')).toBeNull()
        expect(section('log-issues').querySelector('.success')).toBeNull()
        const messages = root.querySelectorAll('.success')
        expect(messages[messages.length - 1].textContent).toBe('Settings updated successfully.')
    })
})
