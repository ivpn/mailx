<template>
    <div class="container mx-auto max-w-screen-lg">
        <header class="flex flex-col justify-between w-full text-sm border-b border-gray-200">
            <div class="flex flex-row justify-between pt-6 pb-2">
                <a class="inline-flex items-center gap-x-2 text-xl font-semibold text-black" href="#">
                    <span class="bg-black w-10 h-10 text-white flex items-center justify-center rounded-md">B</span>
                    Brand
                </a>
                <button type="button"
                    class="flex items-center text-gray-500 bg-gray-100 px-3 rounded-md hover:text-blue-5000 font-medium">
                    user.name@example.net
                    <svg class="ms-1 flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                        stroke-linejoin="round">
                        <path d="m6 9 6 6 6-6" />
                    </svg>
                </button>
            </div>
            <nav>
                <div class="flex flex-row items-center gap-6">
                    <a class="py-4 font-medium text-blue-500 border-blue-500 border-b-2" href="#"
                        aria-current="page">Dashboard</a>
                    <a class="py-4 font-medium text-gray-500 hover:text-blue-500" href="#">Aliases</a>
                    <a class="py-4 font-medium text-gray-500 hover:text-blue-500" href="#">Recipients</a>
                    <a class="py-4 font-medium text-gray-500 hover:text-blue-500" href="#">Settings</a>
                    <a class="py-4 font-medium text-gray-500 hover:text-blue-500" href="#">Subscription</a>
                </div>
            </nav>
        </header>
    </div>
</template>

<script setup lang="ts">
</script>