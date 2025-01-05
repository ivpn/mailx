<template>
    <div class="flex flex-col items-center justify-center p-10 pb-10">
        <p class="text-gray-500 dark:text-gray-400 mt-3 mb-8 text-center">
            Here to try MailX? You need an active IVPN account.<br>
            Sign up or log in on <a class="text-bluish-500 hover:text-bluish-600" href="https://www.ivpn.net/account/">ivpn.net</a> and look for "Email Beta" in your account settings.
        </p>
        <div class="grid grid-cols-3 gap-4 text-center mb-8">
            <p>
                <router-link class="text-bluish-500 hover:text-bluish-600" to="/privacy">Privacy Policy</router-link>
            </p>
            <p>
                <router-link class="text-bluish-500 hover:text-bluish-600" to="/tos">Terms of Service</router-link>
            </p>
            <p>
                <a class="text-bluish-500 hover:text-bluish-600" href="https://www.ivpn.net/">IVPN Limited</a>
            </p>
        </div>
    </div>
</template>