<template>
    <div class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 class="text-xl font-bold text-gray-800 mb-5">Settings</h1>
        <h2 class="font-semibold text-gray-800 mb-3">
            Default Domain
        </h2>
        <p class="text-sm text-gray-500 mb-3">
            The default alias domain is the domain to be selected by default in the drop down options when generating a
            new alias.
        </p>
        <div class="max-w-xs mb-6">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="domain">
                Select default domain:
            </label>
            <select id="domain" 
                class="form-select py-3 px-4 pe-9 block w-full border-2 border-gray-200 rounded-lg text-gray-700 focus:border-violet-600 disabled:opacity-50 disabled:pointer-events-none outline-none">
                <option>Domain 1</option>
                <option>Domain 2</option>
                <option>Domain 3</option>
            </select>
        </div>
        <h2 class="font-semibold text-gray-800 mb-3">
            Default Recipient
        </h2>
        <p class="text-sm text-gray-500 mb-3">
            The default recipient to be selected by default in the drop down options when creating a new recipient.
        </p>
        <div class="max-w-xs mb-6">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="recipient">
                Select default recipient:
            </label>
            <select id="recipient" 
                class="form-select py-3 px-4 pe-9 block w-full border-2 border-gray-200 rounded-lg text-gray-700 focus:border-violet-600 disabled:opacity-50 disabled:pointer-events-none outline-none">
                <option>Recipient 1</option>
                <option>Recipient 2</option>
                <option>Recipient 3</option>
            </select>
        </div>
        <h2 class="font-semibold text-gray-800 mb-3">
            From Name
        </h2>
        <p class="text-sm text-gray-500 mb-3">
            The 'From Name' is shown when you send an email from an alias or reply anonymously to a forwarded email. If
            left blank, then the email alias will be used as the 'From Name'.
        </p>
        <div class="max-w-xs mb-5">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="from-name">
                From name:
            </label>
            <input
                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-violet-600 mb-2"
                id="from-name" type="text">
        </div>
        <div class="mb-3">
            <button
                class="bg-violet-600 hover:bg-violet-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                type="submit">
                Save Settings
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
</script>