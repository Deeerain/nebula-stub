<script setup>
import { onMounted, ref } from 'vue';
import ServiceStatus from './components/ServiceStatus.vue';
import MiniApp from './components/MiniApp.vue';
import { useUserStore } from './store/user.js';

const userStore = useUserStore()
onMounted(async () => {
  await userStore.init()
})
</script>

<template>
  <header>
    <h1>nebula</h1>
  </header>
  <main>
    <MiniApp v-if="userStore.tgInited"/>
    <ServiceStatus v-else/>
  </main>
</template>

<style lang="css">
header {
  text-transform: uppercase;
}

main {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  flex-grow: 1;
  padding: 5px;
}
</style>