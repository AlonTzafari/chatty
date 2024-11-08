<script setup lang="ts">
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import FloatLabel from 'primevue/floatlabel';
import Card from 'primevue/card';
import {useAuthStore} from '@/stores/auth';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const authStore = useAuthStore()
const router = useRouter()
const name = ref('')
const member = ref('')
const loading = ref(false)
async function submit(e: Event) {
    try {
        loading.value = true
        const form = e.target as HTMLFormElement
        const formData = new FormData(form)
        await fetch('/api/login', {body: formData, method: 'post'})
        const res = await fetch('/api/.me')
        const me: {Id: string, Username: string} | null = await res.json()
        authStore.user = me
        if(me) {
            router.push(router.currentRoute.value.query.from as string ?? '/')
        }
    } catch(e) {
        console.error(e);
    }
    loading.value = false
}
</script>

<template>
    <main>

        <Card class="card">
            <template #title>
                <h1>
                    Create Channel
                </h1>
            </template>
            <template #content>
                <form @submit.prevent="submit">
                    <FloatLabel variant="on">
                        <InputText id="name" name="name" v-model="name" required :disabled="loading"/>
                        <label for="name">name</label>
                    </FloatLabel>
                    <FloatLabel variant="on">
                        <InputText id="member" name="member" v-model="member" required :disabled="loading"/>
                        <Button>Add</Button>
                    <label for="member">member</label>
                </FloatLabel>
                <Button class="button" type="submit" :disabled="loading">create</Button>
            </form>
        </template>
        <template #footer>
            <h2>
                footer
            </h2>
        </template>
    </Card>
</main>
</template>
<style scoped>
    main {
        width: 100%;
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
    }
    .button {
        width: 5rem;
    }
    .card {
        display: flex;
        flex-direction: column;
        text-align: center;
        width: 25rem;
    }
    form {
        margin: auto;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: baseline;
        gap: 1.5rem;
        height: 15rem;
    }
</style>
