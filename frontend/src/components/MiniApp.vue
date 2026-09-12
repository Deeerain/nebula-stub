<script setup>
import { useUserStore } from '../store/user';


const userStore = useUserStore()
</script>

<template>
    <div class="profile">
        <small v-if="userStore.username">{{ userStore.username }}</small>
        <div class="profile__bio">
            <img :src="userStore.photo_url" alt="avatar">
            <h2>{{ userStore.first_name }} {{ userStore.last_name }}</h2>
        </div>
    </div>
</template>