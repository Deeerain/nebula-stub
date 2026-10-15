<script setup>
import { useUserStore } from '../store/user';

const userStore = useUserStore()

</script>

<template>
    <div class="subscription">
        <h2>Информация подписки</h2>


        <div v-if="userStore.isLoading">
            <h3>Loading....</h3>
        </div>

        <div v-if="userStore.client || !userStore.isLoading">
            <table>
                <body>
                    <tr>
                        <td>Статус</td>
                        <td>{{ userStore.client.enable ? "Активна" : "Не автивна" }}</td>
                    </tr>
                    <tr>
                        <td>Ограничение IP</td>
                        <td>{{ userStore.client.limitIP > 0 ? userStore.client.limitIP : "Много" }}</td>
                    </tr>
                    <tr>
                        <td>Ограничение Устройств</td>
                        <td>{{ userStore.client.limitHwid > 0 ? userStore.client.limitHwid : "Много" }}</td>
                    </tr>
                    <tr>
                        <td>Окончание</td>
                        <td>{{ userStore.client.limitHwid > 0 ? new Date(userStore.client.expiryTime).toLocaleDateString() : "Никогда" }}</td>
                    </tr>
                </body>
            </table>
        </div>
        <div v-else>
            <span>Нет активной подписки</span>
        </div>
    </div>
</template>