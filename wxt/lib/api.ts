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

async function fetchAliases(apiToken: string) {
    const res = await fetch(`${BASE_URL}/v1/api/aliases`, {
        headers: { Authorization: `Bearer ${apiToken}` },
    })

    return res.json()
}

export const api = {
    livez,
    authenticate,
    fetchAliases,
}
