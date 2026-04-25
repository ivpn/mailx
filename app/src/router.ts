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
import Logout from './components/Logout.vue'
import InitiateResetPassword from './components/InitiateResetPassword.vue'
import ResetPassword from './components/ResetPassword.vue'
import Terms from './components/Terms.vue'
import Privacy from './components/Privacy.vue'
import Faq from './components/Faq.vue'
import NotFound from './components/NotFound.vue'
import Landing from './components/Landing.vue'
import { type IStaticMethods } from 'preline/preline'

declare global {
    interface Window {
        HSStaticMethods: IStaticMethods;
    }
}

const AppName = import.meta.env.VITE_APP_NAME

// Dashboard child routes
const dashboardChildren: RouteRecordRaw[] = [
    {
        path: '',
        name: `${AppName}: Aliases`,
        component: QuickActions,
    },
    {
        path: 'wildcard',
        name: `${AppName}: Wildcard`,
        component: Wildcards,
    },
    {
        path: 'recipients',
        name: `${AppName}: Recipients`,
        component: Recipients,
    },
    {
        path: 'stats',
        name: `${AppName}: Stats`,
        component: Stats,
    },
    {
        path: 'diagnostics',
        name: `${AppName}: Diagnostics`,
        component: Diagnostics,
    },
    {
        path: 'settings',
        name: `${AppName}: Settings`,
        component: Settings,
    },
    {
        path: 'profile',
        name: `${AppName}: Account`,
        component: Account,
    },
]

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: `${AppName}: Home`,
        component: Landing,
    },
    {
        path: '/account',
        name: AppName,
        component: Dashboard,
        children: dashboardChildren
    },
    {
        path: '/login',
        name: `${AppName}: Log In`,
        component: Login
    },
    {
        path: '/logout',
        name: `${AppName}: Log Out`,
        component: Logout
    },
    {
        path: '/signup/:subid?',
        name: `${AppName}: Sign Up`,
        component: Signup
    },
    {
        path: '/signup-complete',
        name: `${AppName}: Signup Complete`,
        component: Login
    },
    {
        path: '/forgot-password',
        name: `${AppName}: Reset Password`,
        component: InitiateResetPassword
    },
    {
        path: '/reset-password/:token',
        name: `${AppName}: Set New Password`,
        component: ResetPassword,
    },
    {
        path: '/tos',
        name: `${AppName}: Terms`,
        component: Terms
    },
    {
        path: '/privacy',
        name: `${AppName}: Privacy Policy`,
        component: Privacy
    },
    {
        path: '/faq',
        name: `${AppName}: FAQ`,
        component: Faq
    },
    {
        path: '/:pathMatch(.*)*',
        name: '404: Not Found',
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

    // Protect all /account/* routes
    if (to.path.startsWith('/account') && !isLoggedIn()) {
        return { path: '/login', query: { redirect: to.fullPath } }
    }

    // Redirect logged-in users away from login/signup
    if ((to.path === '/login' || to.path === '/signup') && isLoggedIn()) {
        return { path: '/account' }
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
