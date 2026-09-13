<script setup>
import { useUserStore } from '../store/user';


const userStore = useUserStore()
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
    <div class="subscriptions">
        <div class="subscription__item subscription" v-for="sub, subIndex in userStore.subscriptions" :key="subIndex">
            <h3>{{ sub.name }}</h3>
            <div class="subscription__connections">
                <div class="connection" v-for="conn, connIndex in sub.connections">
                    <h4>{{ conn.name }}</h4>
                </div>
            </div>
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