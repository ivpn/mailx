import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from './components/Dashboard.vue'
import QuickActions from './components/QuickActions.vue'
import Recipients from './components/Recipients.vue'
import Stats from './components/Stats.vue'
import Settings from './components/Settings.vue'
import Account from './components/Account.vue'
import Signup from './components/Signup.vue'
import Login from './components/Login.vue'
import InitiateResetPassword from './components/InitiateResetPassword.vue'
import ResetPassword from './components/ResetPassword.vue'
import Terms from './components/Terms.vue'
import NotFound from './components/NotFound.vue'
import { type IStaticMethods } from 'preline/preline'

declare global {
    interface Window {
        HSStaticMethods: IStaticMethods;
    }
}

const AppName = import.meta.env.VITE_APP_NAME

const routes = [
    {
        path: '/',
        name: AppName,
        component: Dashboard,
        children: [
            {
                path: '',
                name: AppName,
                component: QuickActions,
            },
            {
                path: 'recipients',
                name: AppName + ' - Recipients',
                component: Recipients,
            },
            {
                path: 'stats',
                name: AppName + ' - Stats',
                component: Stats,
            },
            {
                path: 'settings',
                name: AppName + ' - Settings',
                component: Settings,
            },
            {
                path: 'account',
                name: AppName + ' - Account',
                component: Account,
            },
        ]
    },
    {
        path: '/signup/:subid',
        name: AppName + ' - Sign Up',
        component: Signup
    },
    {
        path: '/login',
        name: AppName + ' - Log In',
        component: Login
    },
    {
        path: '/signup-complete',
        name: AppName + ' - Signup Complete',
        component: Login
    },
    {
        path: '/reset/password/initiate',
        name: AppName + ' - Reset Password',
        component: InitiateResetPassword
    },
    {
        path: '/reset/password/:otp',
        name: AppName + ' - Set New Password',
        component: ResetPassword,
        params: true
    },
    {
        path: '/tos',
        name: AppName + ' - Terms',
        component: Terms
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

    const protectedRoutes = ['/', '/recipients', '/stats', '/settings', '/account']
    
    if (protectedRoutes.includes(to.path) && !isLoggedIn()) {
        return { name: AppName + ' - Log In' }
    }
})

router.afterEach((failure) => {
    // Reinitialize Preline plugins
    // https://preline.co/docs/preline-javascript.html
    if (!failure) {
        setTimeout(() => {
            window.HSStaticMethods.autoInit();
        }, 100)
    }
})

const isLoggedIn = (): boolean => {
    const email = localStorage.getItem('email')
    return email !== null && email.trim() !== ''
}

export default router
