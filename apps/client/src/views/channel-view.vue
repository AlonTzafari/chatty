<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import Textarea from 'primevue/textarea'
import Button from 'primevue/button';
import { wsClient } from '@/ws';
type Message = {
    Id: string, UserId: string, Content: string, CreatedAt: string, Username: string
}
const route = useRoute('channel')
const channelInfo = ref<{Id: string, Name: string} | null>(null)
const message = ref("")
const messages = ref<Message[]>([])
const disabled = computed(() => message.value === "")
let unsub: () => void
onMounted(async () => {
    try {
        console.log('route.params.id', route.params.id);
        const [channelInfoRes, messagesRes] = await Promise.all([
            fetch(`/api/channels/${route.params.id}`),
            fetch(`/api/messages?channel_id=${route.params.id}`),
        ])
        channelInfo.value = await channelInfoRes.json()
        messages.value = await messagesRes.json()
       unsub = wsClient.subscribe<Message>('message-updates', (data) => {
            messages.value.push(data)
        })
    } catch(e) {
        console.error(e)
    }
})
onUnmounted(() => {
    unsub?.()
})
watch(
    () => route.params.id,
    async (newId) => {
        console.log('route.params.id', route.params.id, "newId", newId);
    }
)
async function sendMessage() {
    try {
        const data = {
            channelId: route.params.id,
            content: message.value,
        }
        await fetch('/api/messages', {method: 'post', body: JSON.stringify(data), headers: {
            "Content-Type": "application/json"
        }})
        message.value = ""
    } catch(e) {
        console.error(e)
    }
}

</script>
<template>
    <header>
        <RouterLink class="unset clickable" to="/"><i class="pi pi-home" style="font-size: 2rem;"></i></RouterLink>
        <h1>
            {{ channelInfo?.Name }}
        </h1>
    </header>
    <main>
        <div v-for="message of messages" :key="message.Id">{{ message.Content }}</div>
    </main>
    <footer>
        <Textarea v-model="message" rows="1" cols="100" />
        <Button :disabled @click="sendMessage"><i class="pi pi-send"></i></Button>
    </footer>
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
    main {
        width: 100%;
        height: calc(100% - 10em);
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }
    footer {
        height: 5rem;
        width: 100%;
        border-top: 1px solid rgb(37, 35, 35);
        display: flex;
        gap: 2rem;
        align-items: center;
        padding: 1rem 1rem;

    }
    .unset {
        all: unset;
    }
    .clickable:hover {
        cursor: pointer;
    }
    .clickable:active {
        cursor: auto;
    }
</style>
