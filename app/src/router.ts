import { createRouter, createWebHistory } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import Signup from './components/Signup.vue'
import Login from './components/Login.vue'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: HelloWorld,
        props: { msg: 'Hello' }
    },
    {
        path: '/signup',
        name: 'Sign Up',
        component: Signup
    },
    {
        path: '/login',
        name: 'Log In',
        component: Login
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to) => {
    document.title = to.name as string
})

export default router
