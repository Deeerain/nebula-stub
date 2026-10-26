<script setup>
import { computed } from 'vue';
import { useUserStore } from '../store/user';

const userStore = useUserStore()

const expiryDate = computed(() => {
    const timestamp = userStore.client?.expiry_time
    if (!timestamp) {
        return 'Никогда'
    }

    // 3x-ui returns milliseconds, but seconds are accepted as well.
    const milliseconds = timestamp < 1e12 ? timestamp * 1000 : timestamp

    return new Date(milliseconds).toLocaleDateString()
})

const usedTraffic = computed(() => {
    const bytes = userStore.client?.used_traffic ?? 0
    if (!bytes) {
        return '0 ГБ'
    }

    return `${(bytes / 1024 ** 3).toFixed(2)} ГБ`
})
</script>

<template>
    <div class="subscription">
        <h2>Информация подписки</h2>

        <div v-if="userStore.isLoading">
            <h3>Загрузка...</h3>
        </div>

        <table v-else-if="userStore.client">
            <tbody>
                <tr>
                    <td>Статус</td>
                    <td>{{ userStore.client.enable ? "Активна" : "Не активна" }}</td>
                </tr>
                <tr>
                    <td>Ограничение IP</td>
                    <td>{{ userStore.client.limit_ip > 0 ? userStore.client.limit_ip : "Нет" }}</td>
                </tr>
                <tr>
                    <td>Ограничение устройств</td>
                    <td>{{ userStore.client.limit_hwid > 0 ? userStore.client.limit_hwid : "Нет" }}</td>
                </tr>
                <tr>
                    <td>Окончание</td>
                    <td>{{ expiryDate }}</td>
                </tr>
                <tr>
                    <td>Использовано</td>
                    <td>{{ usedTraffic }}</td>
                </tr>
            </tbody>
        </table>

        <span v-else>Нет активной подписки</span>
    </div>
</template>
