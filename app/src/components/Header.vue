<template>
    <header class="flex flex-col justify-between w-full text-sm bg-white shadow-sm">
        <div class="container mx-auto max-w-screen-lg px-5">
            <div class="flex flex-row justify-between max-w-screen-lg">
                <Menu />
                <div class="hs-dropdown relative inline-flex [--placement:bottom-right] my-3">
                    <button id="hs-dropdown-default" type="button"
                        class="flex items-center hs-dropdown-toggle text-gray-500 bg-gray-100 pl-4 pr-3 rounded-md hover:bg-gray-200 font-medium">
                        {{ jwt().email }}
                        <svg class="ms-1 flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                            stroke-linecap="round" stroke-linejoin="round">
                            <path d="m6 9 6 6 6-6" />
                        </svg>
                    </button>
                    <div class="hs-dropdown-menu transition-[opacity,margin] duration hs-dropdown-open:opacity-100 opacity-0 hidden min-w-60 bg-white border border-gray-200 shadow-sm rounded-lg p-2 mt-2 after:h-4 after:absolute after:-bottom-4 after:start-0 after:w-full before:h-4 before:absolute before:-top-4 before:start-0 before:w-full"
                        aria-labelledby="hs-dropdown-default">
                        <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-500 hover:bg-gray-100 focus:outline-none"
                            href="/account">
                            Account
                        </a>
                        <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-500 hover:bg-gray-100 focus:outline-none"
                            href="/subscription">
                            Subscription
                        </a>
                        <a @click.prevent="logout"
                            class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-red-600 hover:bg-gray-100 focus:outline-none"
                            href="#">
                            Log out
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
import Menu from './Menu.vue'
import { jwt } from '../utils/jwt'
import dropdown from '@preline/dropdown'
import { userApi } from '../api/user.ts'

dropdown.autoInit

const logout = async () => {
    try {
        await userApi.logout()
        window.location.href = '/login'
    } catch { }
}
</script>