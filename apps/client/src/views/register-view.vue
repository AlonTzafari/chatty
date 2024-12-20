<script setup lang="ts">
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import FloatLabel from 'primevue/floatlabel';
import Message from 'primevue/message';
import Card from 'primevue/card';
import { computed, ref } from 'vue';
import { useRouter, RouterLink } from 'vue-router';
const router = useRouter()
const username = ref('')
const usernameErr = ref('')
const password = ref('')
const password2 = ref('')
const formErr = ref('')
const password2Valid = computed(() => {
    return password.value === password2.value && password2.value.length > 0
})
const loading = ref(false)
async function submit(e: Event) {
    usernameErr.value = ''
    usernameErr.value = ''
    try {
        loading.value = true
        const form = e.target as HTMLFormElement
        const formData = new FormData(form)
        const res = await fetch('/api/register', {body: formData, method: 'post'})
        if (res.ok) {
            router.push("/login")
            return
        }
        if (res.status === 400) {
            const err = await res.json() as {Username: string, Password: string}
            usernameErr.value = err.Username
            return
        }
        formErr.value = 'register failed'
    } catch(e) {
        formErr.value = 'register failed'
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
                    Register
                </h1>
            </template> 
            <template #content>
                <form @submit.prevent="submit">
                    <FloatLabel variant="on">
                        <InputText id="username" name="username" :invalid="usernameErr.length > 0" v-model="username" required :disabled="loading"/>
                            <Message v-if="usernameErr.length > 0"  severity="error" >{{ usernameErr }}</Message>
                        <label for="username">username</label>
                    </FloatLabel>
                    <FloatLabel variant="on">
                        <InputText id="password" name="password" type="password" v-model="password" required :disabled="loading"/>
                        <label for="password">password</label>
                    </FloatLabel>
                    <FloatLabel variant="on">
                        <InputText id="password2" name="password2" type="password" v-model="password2" required :disabled="loading"/>
                        <label for="password2">confirm password</label>
                    </FloatLabel>
                <Button class="button" type="submit" :disabled="loading || !password2Valid">register</Button>
            </form>
        </template>
        <template #footer>
            <div class="link-container">
                <RouterLink to="/login">login</RouterLink>
            </div>
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
