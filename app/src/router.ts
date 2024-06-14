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

declare global {
    interface Window {
        HSStaticMethods: IStaticMethods;
    }
}

const routes = [
    {
        path: '/',
        name: 'App',
        component: Dashboard,
        children: [
            {
                path: '',
                name: 'App',
                component: QuickActions,
            },
            {
                path: 'aliases',
                name: 'App - Aliases',
                component: Aliases,
            },
            {
                path: 'recipients',
                name: 'App - Recipients',
                component: Recipients,
            },
            {
                path: 'stats',
                name: 'App - Stats',
                component: Stats,
            },
            {
                path: 'settings',
                name: 'App - Settings',
                component: Settings,
            },
            {
                path: 'account',
                name: 'App - Account',
                component: Account,
            },
        ]
    },
    {
        path: '/signup',
        name: 'App - Sign Up',
        component: Signup
    },
    {
        path: '/login',
        name: 'App - Log In',
        component: Login
    },
    {
        path: '/initiate-password-reset',
        name: 'App - Reset Password',
        component: InitiateResetPassword
    },
    {
        path: '/reset/password/:otp',
        name: 'App - Set New Password',
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
