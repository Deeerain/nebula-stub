<script setup>
import { onMounted, ref } from 'vue';
import ServiceStatus from './components/ServiceStatus.vue';
</script>

<template>
  <header>
    <h1>nebula</h1>
  </header>
  <main>
    <ServiceStatus/>
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