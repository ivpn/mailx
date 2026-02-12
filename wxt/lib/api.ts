const BASE_URL = import.meta.env.VITE_API_URL

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

    return res.json()
}

async function fetchAliases(apiToken: string, search = '') {
    const res = await fetch(`${BASE_URL}/v1/api/aliases?search=${encodeURIComponent(search)}`, {
        headers: { Authorization: `Bearer ${apiToken}` },
    })

    return res.json()
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

    return res.json()
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

    return res.json()
}

async function deleteAlias(apiToken: string, aliasId: string) {
    const res = await fetch(`${BASE_URL}/v1/api/alias/${aliasId}`, {
        method: 'DELETE',
        headers: {
            Authorization: `Bearer ${apiToken}`,
        },
    })

    return res.json()
}

async function fetchDefaults(apiToken: string) {
    const res = await fetch(`${BASE_URL}/v1/api/defaults`, {
        headers: { Authorization: `Bearer ${apiToken}` },
    })

    return res.json()
}

async function logout(apiToken: string) {
    const res = await fetch(`${BASE_URL}/v1/api/logout`, {
        method: 'POST',
        headers: {
            Authorization: `Bearer ${apiToken}`,
        },
    })

    return res.json()
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
}
