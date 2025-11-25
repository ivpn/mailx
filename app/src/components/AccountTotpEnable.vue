<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-totp-enable'" class="cta">
            Enable
        </button>
        <div v-bind:id="'modal-totp-enable'" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>ENABLE 2-FACTOR AUTHENTICATION</h4>
                    </header>
                    <article>
                        <div v-if="!isEnabled">
                            <div class="mb-5">
                                <p>
                                    To enable 2FA, please scan the QR code in a TOTP app (such as Ente or Aegis), and enter the output code in the field below.
                                </p>
                                <p>
                                    You can also add the following information for manual setup:<br>
                                    - Secret: {{ resEnable.secret }}<br>
                                    - Account: {{ resEnable.account }}
                                </p>
                            </div>
                            <div class="mb-5 container">
                                <canvas class="mx-auto" id="totp_qr_code"></canvas>
                            </div>
                            <div class="mb-3">
                                <label for="totp_enable_code">
                                    Code from TOTP app:
                                </label>
                                <input
                                    v-model="req.otp"
                                    v-bind:class="{ 'error': codeError }"
                                    id="totp_enable_code"
                                    placeholder="6-digit code"
                                    type="text"
                                    pattern="[0-9]*"
                                >
                                <p v-if="codeError" class="error">Required</p>
                            </div>
                        </div>
                        <div v-if="isEnabled">
                            <p>
                                Two-factor authentication is active.
                            </p>
                            <p>
                                Save the following backup codes. You will need them if you lose access to your TOTP authenticator.
                            </p>
                            <p class="py-4 px-5 bg-secondary">
                                Backup codes:
                                <span class="text-primary">
                                    {{ resConfirm.backup }}
                                </span>
                            </p>
                            <p>
                                Each code can be used once.
                            </p>
                        </div>
                    </article>
                    <footer>
                        <nav v-if="!isEnabled">
                            <button @click="totpEnableConfirm" class="cta">
                                Enable 2-Factor Authentication
                            </button>
                            <button @click="close" class="cancel">
                                Cancel
                            </button>
                        </nav>
                        <nav v-if="isEnabled">
                            <button @click="complete" class="cancel">
                                Close
                            </button>
                        </nav>
                        <p v-if="error" class="px-5 error">Error: {{ error }}</p>
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
import QRious from 'qrious'
import events from '../events.ts'

const req = ref({ otp: '' })
const resEnable = ref({ uri: '', secret: '', account: '' })
const resConfirm = ref({ backup: '' })
const error = ref('')
const codeError = ref(false)
const isEnabled = ref(false)

const close = () => {
    req.value = {} as any
    error.value = ''
    document.removeEventListener('keydown', handleKeydown)
    const modal = document.querySelector('#modal-totp-enable') as any
    overlay.close(modal)
}

const complete = () => {
    events.emit('totp.enable', {})
    close()
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-totp-enable' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        totpEnable()
        focusFirstInput()
        document.addEventListener('keydown', handleKeydown)
    })
}

const focusFirstInput = () => {
    const input = document.getElementById('totp_enable_code')
    input?.focus()
}

const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
        event.preventDefault()
        totpEnableConfirm()
    }
}

const totpEnable = async () => {
    try {
        const res = await userApi.totpEnable()
        resEnable.value = res.data
        generateQRCode(resEnable.value.uri)
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const totpEnableConfirm = async () => {
    if (!req.value.otp) {
        codeError.value = true
        return
    }

    req.value.otp = req.value.otp + ''

    try {
        const res = await userApi.totpEnableConfirm(req.value)
        resConfirm.value = res.data
        isEnabled.value = true
        codeError.value = false
    } catch (err) {
        if (axios.isAxiosError(err)) {
            resConfirm.value = { backup: '' }
            isEnabled.value = false
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const generateQRCode = (uri: string) => {
    new QRious({
        element: document.getElementById('totp_qr_code'),
        value: uri,
        size: 150,
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>