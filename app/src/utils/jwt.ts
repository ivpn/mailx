import { jwtDecode } from 'jwt-decode'
import { getCookie } from 'typescript-cookie'

interface Jwt {
    exp: number
    user_id: string
    email: string
}

export const jwt = () => {
    const auth = getCookie('auth')

    try {
        const payload = jwtDecode<Jwt>(auth as string)
        return payload
    } catch {
        return <Jwt>{}
    }
}
