export interface StoredData {
    apiToken?: string
}

const STORAGE_KEYS = {
    apiToken: 'apiToken',
}

type Listener<T> = (value: T) => void

const apiTokenListeners = new Set<Listener<string | undefined>>()

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
}
