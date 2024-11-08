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
            const res = await fetch('/api/.me')
            const user = await res.json()
            this.user = user
            return user
        }
    }
})
