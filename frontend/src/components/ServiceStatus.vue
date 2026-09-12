<script setup>
import { onMounted } from 'vue';
import { useServerStatus } from '../store/server_status';


const serverStore = useServerStatus()
onMounted(async () => {
    await serverStore.loadStatuses()
})

</script>
<template>
    <table>
        <caption>Service status</caption>
        <thead>
            <tr>
                <th>Service</th>
                <th>Status</th>
            </tr>
        </thead> 
        <tbody>
            <tr v-for="service, index in serverStore.serverStatuses" :key="index">
                <td>{{ service.name }}</td>
                <td>{{ service.status }}</td>
            </tr>
        </tbody>
    </table>
</template>

<style lang="css">
</style>