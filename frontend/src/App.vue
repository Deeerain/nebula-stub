<script setup>
import { onMounted, ref } from 'vue';
import ServiceStatus from './components/ServiceStatus.vue';
import MiniApp from './components/MiniApp.vue';
import { useUserStore } from './store/user.js';

const userStore = useUserStore()

function LoadScript(){
  return new Promise((resolve, reject) => {
      const url = new URL(window.location.href)
      console.log(url)

      if (url.hash.startsWith("#tgWebAppData")) {
          console.log("Telegram mode")

          const tgImport = document.createElement("script")
          tgImport.src = "https://telegram.org/js/telegram-web-app.js?63"
          tgImport.async = true;
          tgImport.onload = () => resolve()
          tgImport.onerror = (err) => reject(err);
          document.head.appendChild(tgImport)   
      }
  })
}

onMounted(async () => {
  try {

    await LoadScript()
    await userStore.init()

    console.log(userStore.tg)
  } catch (e) {
  }
})
</script>

<template>
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