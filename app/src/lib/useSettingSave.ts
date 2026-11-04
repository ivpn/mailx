import { ref } from 'vue'
import axios from 'axios'

// Save state for a setting that is updated through its own endpoint, with the API's
// success or error message shown next to it.
export const useSettingSave = () => {
    const saving = ref(false)
    const error = ref('')
    const success = ref('')

    // Returns whether the save succeeded
    const save = async (request: () => Promise<{ data: { message: string } }>): Promise<boolean> => {
        saving.value = true
        error.value = ''
        success.value = ''

        try {
            const res = await request()
            success.value = res.data.message
            return true
        } catch (err) {
            if (axios.isAxiosError(err)) {
                error.value = err.response?.data.error || err.message
            } else {
                error.value = 'Something went wrong'
            }
            return false
        } finally {
            saving.value = false
        }
    }

    return { saving, error, success, save }
}
