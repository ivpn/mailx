const BASE_URL = import.meta.env.VITE_API_URL

export async function authenticate(apiKey: string) {
    const res = await fetch(`${BASE_URL}/v1/api/authenticate`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ access_key: apiKey }),
    })

    return res.json()
}

export async function fetchAliases(apiToken: string) {
    const res = await fetch(`${BASE_URL}/v1/api/aliases`, {
        headers: { Authorization: `Bearer ${apiToken}` },
    })

    return res.json()
}
