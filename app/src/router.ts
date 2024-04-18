import { createRouter, createWebHistory } from 'vue-router'
import { getCookie } from 'typescript-cookie'
import Dashboard from './components/Dashboard.vue'
import Signup from './components/Signup.vue'
import Login from './components/Login.vue'

const routes = [
    {
        path: '/',
        name: 'App',
        component: Dashboard,
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
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, _, next) => {
    document.title = to.name as string

    const authCookie = getCookie('auth')

    if (to.path === '/' && !authCookie) {
        next('/login')
    } else {
        next()
    }
})

export default router
