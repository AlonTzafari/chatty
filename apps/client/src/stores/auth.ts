import { me } from '@/api/auth'
import type authSchema from '@/schemas/auth-schema'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { z } from 'zod'

export const useAuthStore = defineStore('auth',() => {
    const user = ref<z.infer<typeof authSchema> | null>(null)
    const fetchUser = async () => {
        const auth = await me()
        user.value = auth 
        return auth
    }
    return {user, fetchUser}
})
