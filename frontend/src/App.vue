<script setup>
import { onMounted } from 'vue';
import HelloWorld from './components/HelloWorld.vue'

const user = ref(null)

onMounted(async () => {
  const tg = window.Telegram?.WebApp

  user.value = tg.initDataUnsafe?.user

  if (tg) {
    tg.ready()
    tg.expand()
  }
})
</script>

<template>
  <h1>Hello {{ user }}</h1>
</template>
