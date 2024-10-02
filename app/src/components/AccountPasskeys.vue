<template>
    <div class="mb-5">
        <h2 class="text-2xl font-semibold dark:text-gray-100 mb-5">Passkeys</h2>
        <p class="text-gray-500 dark:text-gray-400 mb-5">
            Add or remove Passkeys.<br>
        </p>
        <div class="flex justify-start items-center gap-x-3 mb-3">
            <button @click=""
                class="py-2 px-3 inline-flex items-center gap-x-2 font-medium text-base bg-bluish-500 text-white hover:bg-bluish-600 disabled:opacity-50 disabled:pointer-events-none">
                <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                    viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="M5 12h14"></path>
                    <path d="M12 5v14"></path>
                </svg>
                Add Passkey
            </button>
        </div>

        <div class="flex flex-col">
            <div class="-m-1.5 overflow-x-auto">
                <div class="p-1.5 min-w-full inline-block align-middle">
                    <div class="overflow-hidden">
                        <table class="min-w-full divide-y divide-gray-200 dark:divide-neutral-600">
                            <thead>
                                <tr>
                                    <th scope="col"
                                        class="pr-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                        CREATED
                                    </th>
                                    <th scope="col"
                                        class="px-5 py-3 text-start text-xs font-medium text-gray-500 dark:text-gray-400">
                                        ID
                                    </th>
                                    <th scope="col"
                                        class="pl-5 py-3 text-end text-xs font-medium text-gray-500 dark:text-gray-400">
                                        ACTIONS
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-200 dark:divide-neutral-600">
                                <tr>
                                    <td class="pr-5 py-4 whitespace-nowrap text-start text-sm">
                                        2021-09-01 12:34:56
                                    </td>
                                    <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
                                        1234567890
                                    </td>
                                    <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
                                        <button @click=""
                                            class="text-red-600 hover:text-red-700 dark:text-red-500 dark:hover:text-red-600 font-medium text-sm py-2 focus:outline-none focus:shadow-outline"
                                            type="submit">
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                                <tr>
                                    <td class="pr-5 py-4 whitespace-nowrap text-start text-sm">
                                        2021-09-01 12:34:56
                                    </td>
                                    <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
                                        1234567890
                                    </td>
                                    <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
                                        <button @click=""
                                            class="text-red-600 hover:text-red-700 dark:text-red-500 dark:hover:text-red-600 font-medium text-sm py-2 focus:outline-none focus:shadow-outline"
                                            type="submit">
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                                <tr>
                                    <td class="pr-5 py-4 whitespace-nowrap text-start text-sm">
                                        2021-09-01 12:34:56
                                    </td>
                                    <td class="px-5 py-4 whitespace-nowrap text-start text-sm">
                                        1234567890
                                    </td>
                                    <td class="pl-5 py-4 whitespace-nowrap text-end text-sm">
                                        <button @click=""
                                            class="text-red-600 hover:text-red-700 dark:text-red-500 dark:hover:text-red-600 font-medium text-sm py-2 focus:outline-none focus:shadow-outline"
                                            type="submit">
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>

    </div>
</template>

<script setup lang="ts">

</script>