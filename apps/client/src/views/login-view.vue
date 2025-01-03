<script setup lang="ts">
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import FloatLabel from 'primevue/floatlabel';
import Message from 'primevue/message';
import Card from 'primevue/card';
import {useAuthStore} from '@/stores/auth';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const authStore = useAuthStore()
const router = useRouter()
const username = ref('')
const password = ref('')
const err = ref('')
const loading = ref(false)
async function submit() {
    err.value = ""
    try {
        loading.value = true
        const loginRes = await fetch('/api/login', {
            body: JSON.stringify({
                username: username.value, 
                password: password.value
            }),
            method: 'post',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        if(loginRes.status === 401) {
            err.value = "wrong username or password"
            throw new Error(loginRes.statusText)
        }
        if(!loginRes.ok) {
            err.value = "login failed"
            throw new Error(loginRes.statusText)
        }
        const meRes = await fetch('/api/.me')
        const me: {Id: string, Username: string} | null = await meRes.json()
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
                    Login
                </h1>
            </template>
            <template #content>
                <form @submit.prevent="submit">
                    <FloatLabel variant="on">
                        <InputText id="username" name="username" v-model="username" required :disabled="loading"/>
                        <label for="username">username</label>
                    </FloatLabel>
                    <FloatLabel variant="on">
                        <InputText id="password" name="password" type="password" v-model="password" required :disabled="loading"/>
                    <label for="password">password</label>
                </FloatLabel>
                <Message v-if="err.length > 0" severity="error">{{ err }}</Message>
                <Button class="button" type="submit" :loading="loading">login</Button>
            </form>
        </template>
        <template #footer>
            <RouterLink to="/register">register</RouterLink>
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
        justify-content: flex-start;
        align-items: center;
        gap: 1.5rem;
        height: 17rem;
    }
</style>
