import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Dashboard from './components/Dashboard.vue'
import QuickActions from './components/QuickActions.vue'
import Recipients from './components/Recipients.vue'
import Wildcards from './components/Wildcards.vue'
import Stats from './components/Stats.vue'
import Diagnostics from './components/Diagnostics.vue'
import Settings from './components/Settings.vue'
import Account from './components/Account.vue'
import Signup from './components/Signup.vue'
import Login from './components/Login.vue'
import InitiateResetPassword from './components/InitiateResetPassword.vue'
import ResetPassword from './components/ResetPassword.vue'
import Terms from './components/Terms.vue'
import Privacy from './components/Privacy.vue'
import Faq from './components/Faq.vue'
import NotFound from './components/NotFound.vue'
import { type IStaticMethods } from 'preline/preline'

declare global {
    interface Window {
        HSStaticMethods: IStaticMethods;
    }
}

const AppName = import.meta.env.VITE_APP_NAME

// Protected routes that require authentication
const PROTECTED_ROUTES = ['/', '/recipients', '/stats', '/settings', '/account']

// Dashboard child routes
const dashboardChildren: RouteRecordRaw[] = [
    {
        path: '',
        name: `${AppName} - Aliases`,
        component: QuickActions,
    },
    {
        path: 'wildcard',
        name: `${AppName} - Wildcard`,
        component: Wildcards,
    },
    {
        path: 'recipients',
        name: `${AppName} - Recipients`,
        component: Recipients,
    },
    {
        path: 'stats',
        name: `${AppName} - Stats`,
        component: Stats,
    },
    {
        path: 'diagnostics',
        name: `${AppName} - Diagnostics`,
        component: Diagnostics,
    },
    {
        path: 'settings',
        name: `${AppName} - Settings`,
        component: Settings,
    },
    {
        path: 'account',
        name: `${AppName} - Account`,
        component: Account,
    },
]

// Public routes
const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: AppName,
        component: Dashboard,
        children: dashboardChildren
    },
    {
        path: '/signup/:subid',
        name: `${AppName} - Sign Up`,
        component: Signup
    },
    {
        path: '/login',
        name: `${AppName} - Log In`,
        component: Login
    },
    {
        path: '/signup-complete',
        name: `${AppName} - Signup Complete`,
        component: Login
    },
    {
        path: '/reset/password/initiate',
        name: `${AppName} - Reset Password`,
        component: InitiateResetPassword
    },
    {
        path: '/reset/password/:otp',
        name: `${AppName} - Set New Password`,
        component: ResetPassword,
        props: true // Better approach than accessing params directly
    },
    {
        path: '/tos',
        name: `${AppName} - Terms`,
        component: Terms
    },
    {
        path: '/privacy',
        name: `${AppName} - Privacy Policy`,
        component: Privacy
    },
    {
        path: '/faq',
        name: `${AppName} - FAQ`,
        component: Faq
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

// Check if user is logged in
const isLoggedIn = (): boolean => {
    const email = localStorage.getItem('email')
    return email !== null && email.trim() !== ''
}

// Authentication guard
router.beforeEach((to, _) => {
    document.title = to.name as string

    if (PROTECTED_ROUTES.includes(to.path) && !isLoggedIn()) {
        return { name: `${AppName} - Log In` }
    }
})

// Reinitialize Preline plugins after route changes
router.afterEach((failure) => {
    if (!failure) {
        setTimeout(() => {
            window.HSStaticMethods.autoInit();
        }, 100)
    }
})

export default router
