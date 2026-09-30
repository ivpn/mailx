import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import SettingsBlockTrackers from './SettingsBlockTrackers.vue'
import { settingsApi } from '../api/settings.ts'

vi.mock('../api/settings.ts', () => ({
    settingsApi: {
        updateBlockTrackers: vi.fn(),
    },
}))

const update = vi.mocked(settingsApi.updateBlockTrackers)

let app: App | undefined
let root: HTMLElement

const mount = (props: Record<string, unknown> = {}) => {
    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp(SettingsBlockTrackers, props)
    app.mount(root)
    return root
}

const flush = async () => {
    await new Promise((resolve) => setTimeout(resolve))
    await nextTick()
}

const toggle = () => root.querySelector('#block-trackers') as HTMLInputElement

const click = async () => {
    toggle().click()
    await flush()
}

beforeEach(() => {
    update.mockReset()
})

afterEach(() => {
    app?.unmount()
    root.remove()
})

describe('SettingsBlockTrackers', () => {
    it('is on by default', () => {
        mount()
        expect(toggle().checked).toBe(true)
        expect(root.textContent).toContain('Remove tracking pixels from forwarded emails')
        expect(root.querySelector('label[for="block-trackers"]')?.textContent).toContain('Block email trackers')
    })

    it('reflects initialEnabled', () => {
        mount({ initialEnabled: false })
        expect(toggle().checked).toBe(false)
    })

    it('saves the new value when toggled', async () => {
        update.mockResolvedValue({ data: { message: 'Block email trackers setting updated successfully.' } } as any)
        mount({ initialEnabled: true })

        await click()
        expect(update).toHaveBeenCalledWith(false)
        expect(root.querySelector('.success')?.textContent).toContain('updated successfully')

        await click()
        expect(update).toHaveBeenLastCalledWith(true)
    })

    it('reverts the toggle and shows the API error when saving fails', async () => {
        const config = { headers: new AxiosHeaders() }
        update.mockRejectedValue(new AxiosError('Request failed with status code 422', 'ERR_BAD_REQUEST', config, undefined, {
            status: 422,
            statusText: 'Unprocessable Entity',
            headers: {},
            config,
            data: { error: 'The block_trackers field must be true or false.' },
        }))
        mount({ initialEnabled: true })

        await click()

        expect(toggle().checked).toBe(true)
        expect(root.querySelector('#block-trackers-error')?.textContent).toContain('must be true or false')
        expect(toggle().getAttribute('aria-invalid')).toBe('true')
    })
})
