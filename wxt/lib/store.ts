export interface StoredData {
    apiToken?: string
}

const STORAGE_KEYS = {
    apiToken: 'apiToken',
}

type Listener<T> = (value: T) => void

const apiTokenListeners = new Set<Listener<string | undefined>>()

// ---- Storage reference (session-only) ----
const sessionStorage = browser.storage.session

if (!sessionStorage) {
    throw new Error(
        "browser.storage.session is not available."
    )
}

// ---- Propagate changes to listeners ----
browser.storage.onChanged.addListener((changes, area) => {
    if (area !== 'session') return

    if (STORAGE_KEYS.apiToken in changes) {
        const newValue = changes[STORAGE_KEYS.apiToken].newValue
        for (const listener of apiTokenListeners) {
            listener(newValue)
        }
    }
})

export const store = {
    async getApiToken(): Promise<string | undefined> {
        const result = await sessionStorage.get(STORAGE_KEYS.apiToken)
        return result[STORAGE_KEYS.apiToken]
    },

    async setApiToken(apiToken: string): Promise<void> {
        await sessionStorage.set({ [STORAGE_KEYS.apiToken]: apiToken })
    },

    async removeApiToken(): Promise<void> {
        await sessionStorage.remove(STORAGE_KEYS.apiToken)
    },

    onApiTokenChange(listener: Listener<string | undefined>): () => void {
        apiTokenListeners.add(listener)
        return () => apiTokenListeners.delete(listener)
    },

    async clearAll(): Promise<void> {
        await sessionStorage.clear()
    },
}
