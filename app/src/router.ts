import { createRouter, createWebHistory } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import Register from './components/Register.vue'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: HelloWorld,
        props: { msg: 'Hello' }
    },
    {
        path: '/register',
        name: 'Register',
        component: Register
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router
