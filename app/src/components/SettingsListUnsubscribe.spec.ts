import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import SettingsListUnsubscribe from './SettingsListUnsubscribe.vue'
import { settingsApi } from '../api/settings.ts'

vi.mock('../api/settings.ts', () => ({
    settingsApi: {
        updateListUnsubscribeBehaviour: vi.fn(),
    },
}))

const update = vi.mocked(settingsApi.updateListUnsubscribeBehaviour)

let app: App | undefined
let root: HTMLElement

const mount = (props: Record<string, unknown> = {}) => {
    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp(SettingsListUnsubscribe, props)
    app.mount(root)
    return root
}

const flush = async () => {
    await new Promise((resolve) => setTimeout(resolve))
    await nextTick()
}

const select = () => root.querySelector('select') as HTMLSelectElement
const button = () => root.querySelector('button') as HTMLButtonElement

const choose = (value: string) => {
    select().value = value
    select().dispatchEvent(new Event('change'))
}

const submit = async () => {
    root.querySelector('form')!.dispatchEvent(new Event('submit'))
    await flush()
}

beforeEach(() => {
    update.mockReset()
})

afterEach(() => {
    app?.unmount()
    root.remove()
})

describe('SettingsListUnsubscribe', () => {
    it('renders five options and preselects initialBehaviour', () => {
        mount({ initialBehaviour: 2 })

        const options = select().querySelectorAll('option')
        expect(options).toHaveLength(5)
        expect([...options].map((o) => o.value)).toEqual(['0', '1', '2', '3', '4'])
        expect(select().value).toBe('2')
    })

    it('defaults to the original List-Unsubscribe behaviour', () => {
        mount()
        expect(select().value).toBe('0')
    })

    it('sends the selected behaviour as a number', async () => {
        update.mockResolvedValue({ data: { message: 'Unsubscribe behaviour updated successfully.' } } as any)
        mount({ initialBehaviour: 0 })

        choose('3')
        await submit()

        expect(update).toHaveBeenCalledTimes(1)
        expect(update).toHaveBeenCalledWith(3)
        expect(root.querySelector('.success')?.textContent).toContain('updated successfully')
        expect(select().getAttribute('aria-invalid')).toBeNull()
    })

    it('enables Update only while the selection differs from the saved behaviour', async () => {
        update.mockResolvedValue({ data: { message: 'Unsubscribe behaviour updated successfully.' } } as any)
        mount({ initialBehaviour: 1 })

        expect(button().disabled).toBe(true)

        choose('4')
        await nextTick()
        expect(button().disabled).toBe(false)

        choose('1')
        await nextTick()
        expect(button().disabled).toBe(true)

        choose('4')
        await submit()
        expect(button().disabled).toBe(true)
    })

    it('does not save when the selection is unchanged', async () => {
        mount({ initialBehaviour: 1 })

        await submit()

        expect(update).not.toHaveBeenCalled()
    })

    it('shows the API error and marks the select invalid on 422', async () => {
        const config = { headers: new AxiosHeaders() }
        update.mockRejectedValue(new AxiosError('Request failed with status code 422', 'ERR_BAD_REQUEST', config, undefined, {
            status: 422,
            statusText: 'Unprocessable Entity',
            headers: {},
            config,
            data: { error: 'The selected behaviour is invalid.' },
        }))
        mount()

        choose('2')
        await submit()

        const error = root.querySelector('#list-unsubscribe-behaviour-error')
        expect(error?.textContent).toContain('The selected behaviour is invalid.')
        expect(select().getAttribute('aria-invalid')).toBe('true')
        expect(select().getAttribute('aria-describedby')).toBe('list-unsubscribe-behaviour-error')
        expect(button().disabled).toBe(false)
    })
})
