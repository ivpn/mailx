export interface StoredData {
    apiToken?: string
}

const STORAGE_KEYS = {
    apiToken: "apiToken",
}

export const store = {
    async getApiToken(): Promise<string | undefined> {
        const { apiToken } = await browser.storage.local.get(STORAGE_KEYS.apiToken)
        return apiToken
    },

    async setApiToken(apiToken: string): Promise<void> {
        await browser.storage.local.set({ [STORAGE_KEYS.apiToken]: apiToken })
    },

    async removeApiToken(): Promise<void> {
        await browser.storage.local.remove(STORAGE_KEYS.apiToken)
    },

    async clearAll(): Promise<void> {
        await browser.storage.local.clear()
    },
}