<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import Textarea from 'primevue/textarea'
import Button from 'primevue/button';
const route = useRoute('channel')
const channelInfo = ref<{Id: string, Name: string} | null>(null)
const message = ref("")
onMounted(async () => {
    try {
        console.log('route.params.id', route.params.id);
        const res = await fetch(`/api/channels/${route.params.id}`)
        channelInfo.value = await res.json()
    } catch(e) {
        console.error(e)
    }
})
watch(
    () => route.params.id,
    async (newId) => {
        console.log('route.params.id', route.params.id, "newId", newId);
    }
)

</script>
<template>
    <header>
        <RouterLink class="unset clickable" to="/"><i class="pi pi-home" style="font-size: 2rem;"></i></RouterLink>
        <h1>
            {{ channelInfo?.Name }}
        </h1>
    </header>
    <main>
        
    </main>
    <footer>
        <Textarea v-model="message" rows="1" cols="100" />
        <Button><i class="pi pi-send"></i></Button>
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
