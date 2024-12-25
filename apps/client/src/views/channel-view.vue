<script setup lang="ts">
import { computed, onMounted, onUnmounted, onUpdated, ref, useTemplateRef } from 'vue';
import { useRoute } from 'vue-router';
import Textarea from 'primevue/textarea'
import Button from 'primevue/button';
import { wsClient } from '@/ws';
import { z } from 'zod';
import type messageSchema from '@/schemas/message-schema';
import { getChannel } from '@/api/channel';
import { getMessagesByChannel, sendMessageToChannel } from '@/api/message';
import type channelSchema from '@/schemas/channel-schema';
import Message from '@/components/message-item.vue'
import { useAuthStore } from '@/stores/auth';
const auth = useAuthStore()
const route = useRoute('channel')
const channelInfo = ref<z.infer<typeof channelSchema> | null>(null)
const message = ref("")
const messages = ref<z.infer<typeof messageSchema>[]>([])
const disabled = computed(() => message.value === "")
const main = useTemplateRef<HTMLElement>('main')
const qScroll = ref(true)
let unsub: () => void
function scrollHandler(e: Event) {
    if(e.target) {
        const el = e.target as HTMLElement
        const isBottom = el.scrollTop === (el.scrollHeight - el.clientHeight) 
        qScroll.value = isBottom
    }
}
onMounted(async () => {
    try {
        const [channelRes, messagesRes] = await Promise.all([
            getChannel(route.params.id),
            getMessagesByChannel(route.params.id),
        ])
        channelInfo.value = channelRes
        messages.value = messagesRes
        unsub = await wsClient.subscribe('message-updates', (data) => {
            const i = messages.value.findIndex(msg => msg.Id === data.Id)
            if(i === -1) {
                messages.value.push(data)
            } else {
                messages.value.splice(i, 1, data)
            }
        })
        if(main.value) {
            main.value.addEventListener('scroll', scrollHandler) 
        } 
    } catch(e) {
        console.error(e)
    }
})
onUpdated(() => {
    if(main.value && qScroll.value) {
        main.value.scroll({behavior: 'instant', top: main.value.scrollHeight}) 
    } 
})
onUnmounted(() => {
    unsub?.()
    main.value?.removeEventListener('scroll', scrollHandler) 
})
async function sendMessage() {
    try {
        await sendMessageToChannel(route.params.id, message.value)
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
    <main ref="main">
        <Message v-for="message of messages" :key="message.Id" :message="message" :stick-end="message.UserId === auth.user?.Id"/>
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
        height: calc(100% - 10rem);
        display: flex;
        flex-direction: column-reverse;
        justify-content: flex-start;
        align-items: stretch;
        overflow-y: scroll;
        gap: 1rem;
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
    .clickable:hover {
        cursor: pointer;
    }
    .clickable:active {
        cursor: auto;
    }
</style>
