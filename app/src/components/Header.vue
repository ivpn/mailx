<template>
    <header class="flex flex-col justify-between w-full text-sm border-b border-gray-200">
        <div class="flex flex-row justify-between pt-6 pb-2">
            <a class="inline-flex items-center gap-x-2 text-xl font-semibold text-black" href="/">
                <span class="bg-black w-10 h-10 text-white flex items-center justify-center rounded-md">B</span>
                Brand
            </a>
            <div class="hs-dropdown relative inline-flex [--placement:bottom-right]">
                <button id="hs-dropdown-default" type="button"
                    class="hs-dropdown-toggle text-gray-500 bg-gray-100 px-3 rounded-md hover:text-blue-5000 font-medium">
                    {{ jwt().email }}
                    <svg class="ms-1 flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                        stroke-linejoin="round">
                        <path d="m6 9 6 6 6-6" />
                    </svg>
                </button>
                <div class="hs-dropdown-menu transition-[opacity,margin] duration hs-dropdown-open:opacity-100 opacity-0 hidden min-w-60 bg-gray-100 border border-gray-200 shadow- rounded-lg p-2 mt-2 after:h-4 after:absolute after:-bottom-4 after:start-0 after:w-full before:h-4 before:absolute before:-top-4 before:start-0 before:w-full"
                    aria-labelledby="hs-dropdown-default">
                    <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-500 hover:bg-gray-200 focus:outline-none"
                        href="/account">
                        Account
                    </a>
                    <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-500 hover:bg-gray-200 focus:outline-none"
                        href="/subscription">
                        Subscription
                    </a>
                    <a class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-500 hover:bg-gray-200 focus:outline-none"
                        href="#">
                        Log Out
                    </a>
                </div>
            </div>
        </div>
        <Menu />
    </header>
</template>

<script setup lang="ts">
import Menu from './Menu.vue'
import { jwt } from '../utils/jwt'
import dropdown from '@preline/dropdown'

dropdown.autoInit
</script>