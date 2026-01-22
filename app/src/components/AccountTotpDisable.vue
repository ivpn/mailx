<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-totp-disable'" class="cta">
            Disable
        </button>
        <div v-bind:id="'modal-totp-disable'" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>DISABLE 2-FACTOR AUTHENTICATION</h4>
                    </header>
                    <article>
                        <div class="mb-5">
                            <p>
                                To disable two-factor authentication, please enter code from TOTP app or a backup code.
                            </p>
                        </div>
                        <div class="mb-5">
                            <label for="totp_disable_code">
                                Code from TOTP app:
                            </label>
                            <input
                                v-model="req.otp"
                                v-bind:class="{ 'error': codeError }"
                                id="totp_disable_code"
                                placeholder="6-digit code"
                                type="text"
                                pattern="[0-9]*"
                            >
                            <p v-if="codeError" class="error">Required</p>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click="disableTotp" class="cta">
                                Disable 2-Factor Authentication
                            </button>
                            <button @click="close" class="cancel">
                                Cancel
                            </button>
                        </nav>
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                    </footer>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import overlay from '@preline/overlay'
import events from '../events.ts'

const req = ref({ otp: '' })
const error = ref('')
const codeError = ref(false)

const close = () => {
    req.value = { otp: '' }
    error.value = ''
    const modal = document.querySelector('#modal-totp-disable') as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-totp-disable' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

const disableTotp = async () => {
    if (!req.value.otp) {
        codeError.value = true
        return
    }

    req.value.otp = req.value.otp + ''

    try {
        await userApi.totpDisable(req.value)
        codeError.value = false
        error.value = ''
        events.emit('totp.disable', {})
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>