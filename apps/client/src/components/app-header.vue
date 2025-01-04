<script setup lang="ts">
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import {useAuthStore} from '@/stores/auth'
import { logout } from '@/api/auth';
import { computed, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router';
import { useChannelStore } from '@/stores/channel';


const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const channelStore = useChannelStore()
const avatarProps = computed<{src?: string, label: string}>(() => {
    if (router.currentRoute.value.name === 'channel' && channelStore.channel != null) {
        return {src: channelStore.channel?.Avatar || undefined, label: channelStore.channel?.Name.charAt(0).toUpperCase() ?? ""}
    } else {
        return {label: authStore.user?.Username.charAt(0).toUpperCase() ?? ""}
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
        <div class="title">
            <Avatar shape="circle"  :image="avatarProps.src" :label="avatarProps.label" size="large"/>
            <h1><slot></slot></h1>
        </div>
        <div class="controls">
            <RouterLink v-if="router.currentRoute.value.name !== 'home'" class="unset" to="/"><Button icon="pi pi-home" variant="link" size="large"></Button></RouterLink>
            <Button icon="pi pi-sign-out" :loading="loading" @click="logoutHandler()" variant="link"></Button>  
        </div>
    </header>
</template>
<style  scoped>
    header {
        height: var(--header-height);
        width: 100%;
        border-bottom: 1px solid rgb(37, 35, 35);
        display: flex;
        padding: 1rem 1rem;
        justify-content: space-between;
    }
    .title {
        align-items: center;
        display: flex;
        justify-content: flex-start;
        gap: 2rem;
    }
    .controls {
        display: flex;
        justify-content: flex-end;
    }
</style>