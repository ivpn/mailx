import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from './components/Dashboard.vue'
import QuickActions from './components/QuickActions.vue'
import Aliases from './components/Aliases.vue'
import Recipients from './components/Recipients.vue'
import Stats from './components/Stats.vue'
import Settings from './components/Settings.vue'
import Account from './components/Account.vue'
import Signup from './components/Signup.vue'
import Login from './components/Login.vue'
import InitiateResetPassword from './components/InitiateResetPassword.vue'
import ResetPassword from './components/ResetPassword.vue'
import NotFound from './components/NotFound.vue'
import { type IStaticMethods } from 'preline/preline'
import env from "./env.json"

declare global {
    interface Window {
        HSStaticMethods: IStaticMethods;
    }
}

const routes = [
    {
        path: '/',
        name: env.APP_NAME,
        component: Dashboard,
        children: [
            {
                path: '',
                name: env.APP_NAME,
                component: QuickActions,
            },
            {
                path: 'aliases',
                name: env.APP_NAME + ' - Aliases',
                component: Aliases,
            },
            {
                path: 'recipients',
                name: env.APP_NAME + ' - Recipients',
                component: Recipients,
            },
            {
                path: 'stats',
                name: env.APP_NAME + ' - Stats',
                component: Stats,
            },
            {
                path: 'settings',
                name: env.APP_NAME + ' - Settings',
                component: Settings,
            },
            {
                path: 'account',
                name: env.APP_NAME + ' - Account',
                component: Account,
            },
        ]
    },
    {
        path: '/signup',
        name: env.APP_NAME + ' - Sign Up',
        component: Signup
    },
    {
        path: '/login',
        name: env.APP_NAME + ' - Log In',
        component: Login
    },
    {
        path: '/reset/password/initiate',
        name: env.APP_NAME + ' - Reset Password',
        component: InitiateResetPassword
    },
    {
        path: '/reset/password/:otp',
        name: env.APP_NAME + ' - Set New Password',
        component: ResetPassword,
        params: true
    },
    {
        path: '/:pathMatch(.*)*',
        name: '404 - Not Found',
        component: NotFound
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, _) => {
    document.title = to.name as string
})

router.afterEach((failure) => {
    if (!failure) {
        setTimeout(() => {
            window.HSStaticMethods.autoInit();
        }, 100)
    }
})

export default router
