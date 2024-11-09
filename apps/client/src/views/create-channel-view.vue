<script setup lang="ts">
import InputText from 'primevue/inputtext';
import AutoComplete from 'primevue/autocomplete';
import Button from 'primevue/button';
import FloatLabel from 'primevue/floatlabel';
import Card from 'primevue/card';
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
const authStore = useAuthStore()
const router = useRouter()
const name = ref('')
const member = ref('')
const members = ref<{Id: string, Username: string}[]>([])
const searchResults = ref<{Id: string, Username: string}[]>([])
const suggestions = computed(() => searchResults.value
.filter(result => members.value.every(member => member.Id !== result.Id))
.map(result => result.Username))
const loading = ref(false)
async function searchUsers() {
    try {
        const res = await fetch(`/api/users?username=${member.value}`)
        const users = await res.json()
        searchResults.value = users ?? []
    } catch(e) {
        console.error(e)
    }
}

function addMember() {
    const memberToAdd = searchResults.value.find(result => result.Username === member.value)
    if (memberToAdd == null) {
        console.error(new Error('Failed to find user in search results'))
        return
    } 
    members.value = [...members.value, memberToAdd]
    member.value = ""
}

async function submit() {
    try {
        loading.value = true
        const id = authStore.getUser?.Id
        if(!id) {
            throw new Error("NOT AUTHENTICATED")
        }
        const data = {
            name: name.value,
            members: members.value.map(member => member.Id).concat(id)
        }
        await fetch('/api/channels', {
            body: JSON.stringify(data), 
            method: 'post', 
            headers: {
                "Content-Type": "application/json"
            }
        })
        router.push('/')
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
                <form>
                    <FloatLabel variant="on">
                        <InputText id="name" name="name" v-model="name" required :disabled="loading"/>
                        <label for="name">name</label>
                    </FloatLabel>
                    <FloatLabel variant="on">
                        <AutoComplete id="member" name="member" v-model="member" :suggestions="suggestions" @complete="searchUsers" required :disabled="loading"/>
                        <Button @click="addMember">Add</Button>
                    <label for="member">member</label>
                </FloatLabel>
                <Button class="button" @click="submit" :disabled="loading">create</Button>
            </form>
        </template>
        <template #footer>
            <div v-for="member of members" :key="member.Id">
                {{ member.Username }}
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
        justify-content: center;
        align-items: baseline;
        gap: 1.5rem;
        height: 15rem;
    }
</style>
