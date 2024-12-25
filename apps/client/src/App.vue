<script setup lang="ts">
import { onBeforeMount, onBeforeUnmount } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth';

const router = useRouter()
const authStore = useAuthStore()
let interval: ReturnType<typeof setInterval>
onBeforeMount(async () => {
    interval = setInterval(() => {
        const user = authStore.fetchUser()
        if (!user) {
            console.log(`router.push('/login?from=${router.currentRoute}')`);
            router.push(`/login?from=${router.currentRoute}`)
        }
    }, 30*1000)
})
onBeforeUnmount(() => {
    clearInterval(interval)
})
</script>

<template>
<RouterView />
</template>
