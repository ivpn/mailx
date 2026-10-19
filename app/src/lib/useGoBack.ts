import { useRouter } from 'vue-router'

// Back button for pages that can be the first page in the tab, e.g. a footer page opened
// in a new tab or from a link outside the app. There router.back() has no app page to return
// to and does nothing, so it goes to the account (signed in) or the landing page instead.
export const useGoBack = () => {
    const router = useRouter()

    return () => {
        if (router.options.history.state.back) {
            router.back()
            return
        }

        const email = localStorage.getItem('email')
        const signedIn = email !== null && email.trim() !== ''
        router.replace(signedIn ? '/account' : '/')
    }
}
