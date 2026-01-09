import { Defaults, Preferences } from "./types"

export interface StoredData {
    apiToken?: string
}

const STORAGE_KEYS = {
    apiToken: 'apiToken',
    defaults: 'defaults',
    preferences: 'preferences',
}

type Listener<T> = (value: T) => void

const apiTokenListeners = new Set<Listener<string | undefined>>()
const defaultsListeners = new Set<Listener<Defaults | undefined>>()
const preferencesListeners = new Set<Listener<Preferences | { input_button: true }>>()

export const store = {
    async getApiToken(): Promise<string | undefined> {
        const { apiToken } = await browser.storage.local.get(STORAGE_KEYS.apiToken)
        return apiToken
    },

    async setApiToken(apiToken: string): Promise<void> {
        await browser.storage.local.set({ [STORAGE_KEYS.apiToken]: apiToken })
        for (const listener of apiTokenListeners) {
            listener(apiToken)
        }
    },

    async removeApiToken(): Promise<void> {
        await browser.storage.local.remove(STORAGE_KEYS.apiToken)
        for (const listener of apiTokenListeners) {
            listener(undefined)
        }
    },

    async getDefaults(): Promise<Defaults | undefined> {
        const { defaults } = await browser.storage.local.get(STORAGE_KEYS.defaults)
        return defaults
    },

    async setDefaults(defaults: Defaults | undefined): Promise<void> {
        await browser.storage.local.set({ [STORAGE_KEYS.defaults]: defaults })
        for (const listener of defaultsListeners) {
            listener(defaults)
        }
    },

    async removeDefaults(): Promise<void> {
        await browser.storage.local.remove(STORAGE_KEYS.defaults)
    },

    async getPreferences(): Promise<Preferences> {
        const { preferences } = await browser.storage.local.get(STORAGE_KEYS.preferences)
        return preferences || { input_button: true }
    },

    async setPreferences(preferences: Preferences): Promise<void> {
        await browser.storage.local.set({ [STORAGE_KEYS.preferences]: preferences })
        for (const listener of preferencesListeners) {
            listener(preferences)
        }
    },

    async removePreferences(): Promise<void> {
        await browser.storage.local.remove(STORAGE_KEYS.preferences)
    },

    async clearAll(): Promise<void> {
        await browser.storage.local.clear()
        for (const listener of apiTokenListeners) {
            listener(undefined)
        }
    },

    onApiTokenChange(listener: Listener<string | undefined>): () => void {
        apiTokenListeners.add(listener)
        return () => apiTokenListeners.delete(listener)
    },

    onDefaultsChange(listener: Listener<Defaults | undefined>): () => void {
        defaultsListeners.add(listener)
        return () => defaultsListeners.delete(listener)
    },
}
