import { me } from '@/api/auth'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore({
    id: 'auth',
    state: () => ({
        user: null as {Id: string, Username: string} | null
    }),
    getters: {
        getUser(state) {
            return state.user
        }
    },
    actions: {
        setUser(user: {Id: string, Username: string} | null) {
            this.user = user
        },
        async fetchUser() {
            this.user = await me()
            return this.user
        }
    }
})
