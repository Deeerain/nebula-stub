<script setup>
import { onMounted } from 'vue';
import { useUserStore } from '../store/user';
import SubscriptionInfo from './SubscriptionInfo.vue';


const userStore = useUserStore()

onMounted(async () => {
    try {
        await userStore.loadUserData()
    } finally {
        try {
            await userStore.loadSubscriptions()
        } catch(err) {
            console.error(err)
        }
    }
})

</script>

<template>
    <div class="profile">
        <small v-if="userStore.username">@{{ userStore.username }}</small>
        <div class="profile__bio">
            <img :src="userStore.photo_url" alt="avatar">
            <h1>{{ userStore.first_name }} {{ userStore.last_name }}</h1>
            <button>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.73-5.73"/>
                </svg>
            </button>
        </div>
    </div>
    <SubscriptionInfo/>
    <div class="subscriptions">
        <div class="subscription__item subscription" v-for="sub, subIndex in userStore.get_links" :key="subIndex">
            {{ subIndex}} {{ sub }}
        </div>
    </div>
</template>

<style lang="css">
.profile {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.profile button {
    background: none;
    border: none;
}

.profile button:hover {
    cursor: pointer;
}

.profile .profile__bio {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    padding: 5px;
    flex-grow: 1;
    gap: 1rem;
    background-color: #3f3f3f;
    border-radius: 5px;
    overflow: hidden;
}

.profile__bio img {
    max-width: 100px;
    max-height: 100px;
    border-radius: 5px;
}

.subscriptions {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-top: 5px;
}

.subscriptions button {
    border: none;
    background: none;
}


.subscription {
    background: #3f3f3f;
    border-radius: 5px;
    padding: 5px;
}

.subscription__connections {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.connection {
    border-radius: 5px;
    background: #545454;
    padding: 0 20px;
}
</style>