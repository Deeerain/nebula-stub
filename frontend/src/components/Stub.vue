<script setup>
import { onMounted, useTemplateRef } from 'vue';

const movementStrenght = 40
const gyroStrenght = 2
const logo = useTemplateRef("bg")

const paralaxMouse = (e) => {
    const mouseX = (e.clientX / window.innerWidth) - 0.5
    const mouseY = (e.clientY / window.innerHeight) - 0.5 

    const moveX = -(mouseX * movementStrenght)
    const moveY = -(mouseY * movementStrenght)

    updatePosition(moveX, moveY)
}

const paralaxGyro = (e) => {
    let tiltX = e.gamma
    let tiltY = e.beta

    if (tiltX === null || tiltY == null) return

    if (tiltX > 30) tiltX = 30
    if (tiltY <- 30) tiltX = -30
    if (tiltX > 60) tiltX = 60
    if (tiltY <- 60) tiltX = -60

    const moveX = -(tiltX * gyroStrenght);
    const moveY = -((tiltY - 35) * gyroStrenght); 

    updatePosition(moveX, moveY);
}


const updatePosition = (x, y) => {
    logo.value.style.transform = `translate(calc(-50% + ${x}px), calc(-50% + ${y}px))`;
}

onMounted(() => {    
    document.title = "NEBULA"

    if (!matchMedia('(pointer: coarse)').matches) {
        document.addEventListener('mousemove', paralaxMouse)
    }
    document.addEventListener('deviceorientation', paralaxGyro)
})
</script>

<template>
    <div ref="bg" class="bg-logo-container">
        <img src="../assets/nebula.svg" alt="" srcset="">
    </div>
    <div class="content">
        <h1>NEBULA</h1>
        <div>
            <a class="btn" href="https://mtg.nebulasrv.cc/connect">Connect to TG web proxy</a>
        </div>
    </div>
</template>

<style>
.bg-logo-container {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: -1;
    width: 600px;
    height: 600px;
    opacity: 0.08;
    pointer-events: none;
}

.bg-logo-container img {
    width: 100%;
    height: 100%;
    object-fit: contain;
    animation: spinBackground 20s linear infinite;
}

@keyframes spinBackground {
    from {
        transform: rotateY(0deg)
    }
    to {
        transform: rotateY(360deg);
    }
}

.content {
    flex-grow: 1;
    max-width: 800px;
    margin: 100px auto;
    text-align: center;
    padding: 20px;
}
</style>