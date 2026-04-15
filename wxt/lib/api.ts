import { store } from './store'

const BASE_URL = import.meta.env.VITE_API_URL

function clearSession() {
    store.clearAll()
}

async function livez() {
    const res = await fetch(`${BASE_URL}/livez`)
    return res.text()
}

async function authenticate(apiKey: string) {
    const res = await fetch(`${BASE_URL}/v1/api/authenticate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ access_key: apiKey }),
    })

    const body = await res.json()
    if (!res.ok) throw new Error(body.error ?? 'Authentication failed')
    return body
}

async function fetchAliases(apiToken: string, search = '') {
    const res = await fetch(`${BASE_URL}/v1/api/aliases?search=${encodeURIComponent(search)}`, {
        headers: { Authorization: `Bearer ${apiToken}` },
    })

    const body = await res.json()
    if (res.status === 401) { clearSession(); throw new Error(body.error ?? 'Unauthorized') }
    if (!res.ok) throw new Error(body.error ?? 'Failed to fetch aliases')
    return body
}

async function createAlias(apiToken: string, data: any) {
    const res = await fetch(`${BASE_URL}/v1/api/alias`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${apiToken}`,
        },
        body: JSON.stringify(data),
    })

    const body = await res.json()
    if (res.status === 401) { clearSession(); throw new Error(body.error ?? 'Unauthorized') }
    if (!res.ok) throw new Error(body.error ?? 'Failed to create alias')
    return body
}

async function updateAlias(apiToken: string, aliasId: string, data: any) {
    const res = await fetch(`${BASE_URL}/v1/api/alias/${aliasId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${apiToken}`,
        },
        body: JSON.stringify(data),
    })

    const body = await res.json()
    if (res.status === 401) { clearSession(); throw new Error(body.error ?? 'Unauthorized') }
    if (!res.ok) throw new Error(body.error ?? 'Failed to update alias')
    return body
}

async function deleteAlias(apiToken: string, aliasId: string) {
    const res = await fetch(`${BASE_URL}/v1/api/alias/${aliasId}`, {
        method: 'DELETE',
        headers: {
            Authorization: `Bearer ${apiToken}`,
        },
    })

    const body = await res.json()
    if (res.status === 401) { clearSession(); throw new Error(body.error ?? 'Unauthorized') }
    if (!res.ok) throw new Error(body.error ?? 'Failed to delete alias')
    return body
}

async function fetchDefaults(apiToken: string) {
    const res = await fetch(`${BASE_URL}/v1/api/defaults`, {
        headers: { Authorization: `Bearer ${apiToken}` },
    })

    const body = await res.json()
    if (res.status === 401) { clearSession(); throw new Error(body.error ?? 'Unauthorized') }
    if (!res.ok) throw new Error(body.error ?? 'Failed to fetch defaults')
    return body
}

async function logout(apiToken: string) {
    const res = await fetch(`${BASE_URL}/v1/api/logout`, {
        method: 'POST',
        headers: {
            Authorization: `Bearer ${apiToken}`,
        },
    })

    const body = await res.json()
    if (res.status === 401) { clearSession(); throw new Error(body.error ?? 'Unauthorized') }
    if (!res.ok) throw new Error(body.error ?? 'Failed to log out')
    return body
}

export const api = {
    livez,
    authenticate,
    fetchAliases,
    createAlias,
    updateAlias,
    deleteAlias,
    fetchDefaults,
    logout,
    clearSession,
}
