import { afterEach, describe, expect, it } from 'vitest'
import { createApp, defineComponent, h, nextTick, type App } from 'vue'
import { createRouter, createWebHistory, RouterView, type Router } from 'vue-router'
import { useGoBack } from './useGoBack.ts'

const Page = defineComponent({
    setup() {
        const goBack = useGoBack()
        return () => h('button', { onClick: goBack }, 'Back')
    },
})

const Blank = defineComponent({ render: () => h('div') })

describe('useGoBack', () => {
    let app: App | undefined
    let root: HTMLElement
    let router: Router

    // Starts the app on `path` as the tab's first page
    const mount = async (path: string) => {
        window.history.replaceState(null, '', path)
        router = createRouter({
            history: createWebHistory(),
            routes: [
                { path: '/', component: Blank },
                { path: '/account', component: Blank },
                { path: '/trackers-report', component: Blank },
                { path: '/privacy', component: Page },
            ],
        })
        root = document.createElement('div')
        document.body.appendChild(root)
        app = createApp({ render: () => h(RouterView) })
        app.use(router)
        app.mount(root)
        await router.isReady()
    }

    const clickBack = async () => {
        root.querySelector('button')!.click()
        // router.back() completes on popstate, which jsdom fires asynchronously
        await new Promise((resolve) => setTimeout(resolve, 20))
        await nextTick()
    }

    afterEach(() => {
        app?.unmount()
        root.remove()
        localStorage.clear()
    })

    it('returns to the previous app page', async () => {
        await mount('/trackers-report#eyJuIjoxfQ')
        await router.push('/privacy')

        await clickBack()

        expect(router.currentRoute.value.fullPath).toBe('/trackers-report#eyJuIjoxfQ')
    })

    it('goes to the landing page when it is the first page and signed out', async () => {
        await mount('/privacy')

        await clickBack()

        expect(router.currentRoute.value.path).toBe('/')
    })

    it('goes to the account when it is the first page and signed in', async () => {
        localStorage.setItem('email', 'user@example.com')
        await mount('/privacy')

        await clickBack()

        expect(router.currentRoute.value.path).toBe('/account')
    })
})
