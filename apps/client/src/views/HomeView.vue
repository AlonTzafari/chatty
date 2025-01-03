<script setup lang="ts">
import { logout } from '@/api/auth';
import { getUserChannels } from '@/api/channel';
import type channelSchema from '@/schemas/channel-schema';
import { useAuthStore } from '@/stores/auth';
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import Card from 'primevue/card'
import { onMounted, ref } from 'vue'; 
import { useRouter,RouterLink } from 'vue-router';
import type { z } from 'zod';
const authStore = useAuthStore()
const router = useRouter()
const loading = ref(false)
const channels = ref<z.infer<typeof channelSchema>[]>([])

onMounted(async () => {
    try {
        if(authStore.user) {
            channels.value = await getUserChannels(authStore.user.Id)
        }
    } catch(e) {
        console.error(e)
    }
})

async function logoutHandler() {
    loading.value = true
    try {
        await logout()
        router.push('/login')
    } catch(e) {
        console.error(e)
    }
    loading.value = false
}
</script>

<template>
    <header>
        <Avatar shape="circle" :label="authStore.user?.Username.charAt(0).toUpperCase()" size="large"/>
        <h1>channels</h1>
        <i :class="`pi pi-sign-out clickable ${loading && 'disabled'}`" @click="!loading && logoutHandler()"></i>  
    </header>
  <main>
    <RouterLink class="unset" to="/create-channel"><Button>+</Button></RouterLink>
    <RouterLink class="unset clickable" v-for="channel of channels" :key="channel.Id" :to="`/channel/${channel.Id}`">
        <Card><template #title>{{ channel.Name }}</template></Card>
    </RouterLink>
  </main>
</template>
<style scoped>
    header {
        height: 5rem;
        width: 100%;
        border-bottom: 1px solid rgb(37, 35, 35);
        display: flex;
        gap: 2rem;
        align-items: center;
        padding: 1rem 1rem;
    }
    .clickable:hover {
        cursor: pointer;
    }
    .clickable:active {
        cursor: auto;
    }
    .disabled {
        color: gray;
    }
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 1rem;
    }
</style>
