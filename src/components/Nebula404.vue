<template>
  <div class="nebula-404">
    <!-- Фон с космической туманностью -->
    <div class="space-bg"></div>
    <div class="glitch-overlay"></div>
    
    <!-- Контент -->
    <div class="content">
      <!-- Логотип Nebula -->
      <div class="logo-section">
        <div class="nebula-logo" style="background-color: white; border-radius: 50%;">
          <img src="../assets/nabula.svg" alt="nebula" height="100" width="100">
        </div>
        <h1 class="nebula-title">NEBULA</h1>
        <div class="status-badge">CONNECTION</div>
      </div>

      <!-- Основная ошибка с глитчем -->
      <div class="error-section">
        <div class="glitch-wrapper">
          <div class="glitch" data-text="404">404</div>
        </div>
        <div class="glitch-sub" data-text="CONNECTION LOST">CONNECTION LOST</div>
      </div>

      <!-- Сообщение для маскировки -->
      <div class="message-box">
        <div class="terminal">
          <p class="terminal-line">
            <span class="prompt">$</span> ping nebula.secure
          </p>
          <p class="terminal-line error-line">
            <span class="prompt">></span> Request timeout for icmp_seq 0
          </p>
          <p class="terminal-line">
            <span class="prompt">$</span> tracepath nebula.gateway
          </p>
          <p class="terminal-line error-line">
            <span class="prompt">></span> Host unreachable in 42 hops
          </p>
          <p class="terminal-line">
            <span class="prompt">$</span> nebula status
          </p>
          <p class="terminal-line blink">
            <span class="prompt">></span> TUNNEL OFFLINE [ERR_NEBULA_0042]
          </p>
        </div>
      </div>

      <!-- Кнопки действий -->
      <div class="action-buttons">
        <button @click="simulateRetry" class="btn-primary" :disabled="isRetrying">
          <span v-if="!isRetrying">⟳ ПЕРЕЗАПУСТИТЬ ТУННЕЛЬ</span>
          <span v-else>🔮 УСТАНОВЛЕНИЕ ЗАЩИЩЕННОГО СОЕДИНЕНИЯ...</span>
        </button>
        
        <button @click="showDiagnostics = !showDiagnostics" class="btn-secondary">
          🔍 ТЕЛЕМЕТРИЯ
        </button>
      </div>

      <!-- Диагностика (маскировка под тех информацию) -->
      <div v-if="showDiagnostics" class="diagnostics-panel">
        <div class="panel-header">
          <span>⚡ СИСТЕМНАЯ ТЕЛЕМЕТРИЯ ⚡</span>
          <span class="close" @click="showDiagnostics = false">✕</span>
        </div>
        <div class="panel-content">
          <div class="metric">
            <span>NODE STATUS:</span>
            <span class="value error">OFFLINE</span>
          </div>
          <div class="metric">
            <span>ENCRYPTION:</span>
            <span class="value">AES-256-GCM</span>
          </div>
          <div class="metric">
            <span>PROTOCOL:</span>
            <span class="value">WireGuard®</span>
          </div>
          <div class="metric">
            <span>PACKET LOSS:</span>
            <span class="value error">100%</span>
          </div>
          <div class="metric">
            <span>LAST HANDSHAKE:</span>
            <span class="value">42:42:42 UTC</span>
          </div>
          <div class="progress-bar">
            <div class="progress-label">RECONNECT ATTEMPTS</div>
            <div class="progress-track">
              <div class="progress-fill" :style="{ width: `${retryCount * 20}%` }"></div>
            </div>
            <div class="progress-value">{{ retryCount }}/5</div>
          </div>
        </div>
      </div>

      <!-- Футер -->
      <div class="footer">
        <div class="footer-content">
          <span>© NEBULA • ЗАЩИЩЕННОЕ ОБЛАКО</span>
          <div class="status-leds">
            <span class="led red"></span>
            <span class="led yellow"></span>
            <span class="led green"></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// Состояния
const isRetrying = ref(false)
const showDiagnostics = ref(false)
const retryCount = ref(0)


// Симуляция переподключения
const simulateRetry = () => {
  if (isRetrying.value) return
  
  isRetrying.value = true
  
  // Имитация попытки подключения
  setTimeout(() => {
    retryCount.value++
    
    if (retryCount.value >= 5) {
      // Эффект "перезагрузки" после 5 попыток
      alert('⚠️ НЕ УДАЛОСЬ УСТАНОВИТЬ ТУННЕЛЬ\n\nПроверьте интернет-соединение и повторите попытку.')
      retryCount.value = 0
      isRetrying.value = false
    } else {
      // Имитация успеха (для маскировки)
      setTimeout(() => {
        isRetrying.value = false
        // Здесь можно ничего не делать - это просто заглушка
        if (retryCount.value === 3) {
          // Эффект "частичного подключения" для реализма
          const msg = document.querySelector('.terminal')
          if (msg) {
            const newLine = document.createElement('p')
            newLine.className = 'terminal-line success-line'
            newLine.innerHTML = '<span class="prompt">></span> TUNNEL RE-ESTABLISHED (PARTIAL)'
            msg.appendChild(newLine)
            setTimeout(() => {
              newLine.remove()
            }, 3000)
          }
        }
      }, 1500)
    }
  }, 2000)
}

// Эффект случайного глитча на фоне
let glitchInterval = null

const randomGlitch = () => {
  const overlay = document.querySelector('.glitch-overlay')
  if (overlay && Math.random() > 0.95) {
    overlay.style.opacity = Math.random() * 0.3 + 0.1
    overlay.style.transform = `skew(${Math.random() * 5 - 2.5}deg)`
    setTimeout(() => {
      overlay.style.opacity = '0.05'
      overlay.style.transform = 'skew(0deg)'
    }, 100)
  }
}

onMounted(() => {
  glitchInterval = setInterval(randomGlitch, 200)
})

onUnmounted(() => {
  if (glitchInterval) clearInterval(glitchInterval)
})
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.nebula-404 {
  min-height: 100vh;
  background: radial-gradient(ellipse at center, #0a0a1a 0%, #05050a 100%);
  position: relative;
  overflow-x: hidden;
  font-family: 'Courier New', 'Fira Code', 'Monaco', monospace;
}

/* Космический фон */
.space-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    radial-gradient(2px 2px at 20px 30px, #fff, rgba(0,0,0,0)),
    radial-gradient(1px 1px at 40px 70px, #00ffff, rgba(0,0,0,0)),
    radial-gradient(3px 3px at 80px 120px, #ff00ff, rgba(0,0,0,0));
  background-size: 200px 200px;
  background-repeat: no-repeat;
  opacity: 0.3;
  pointer-events: none;
  animation: stars 20s linear infinite;
}

@keyframes stars {
  0% { transform: translateY(0); }
  100% { transform: translateY(-200px); }
}

/* Глитч оверлей */
.glitch-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 255, 255, 0.1) 2px,
    rgba(0, 255, 255, 0.1) 4px
  );
  pointer-events: none;
  opacity: 0.05;
  transition: all 0.05s ease;
  z-index: 1;
}

/* Основной контент */
.content {
  position: relative;
  z-index: 2;
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* Логотип */
.logo-section {
  text-align: center;
  margin-bottom: 3rem;
}

.nebula-logo {
  display: inline-block;
  filter: drop-shadow(0 0 20px #00ffff);
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.nebula-title {
  font-size: 3rem;
  letter-spacing: 10px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-top: 1rem;
  font-weight: 900;
}

.status-badge {
  display: inline-block;
  margin-top: 0.5rem;
  padding: 4px 12px;
  background: rgba(0, 255, 255, 0.1);
  border: 1px solid #00ffff;
  color: #00ffff;
  font-size: 0.7rem;
  letter-spacing: 2px;
}

/* 404 глитч эффект */
.error-section {
  text-align: center;
  margin: 2rem 0;
}

.glitch-wrapper {
  position: relative;
}

.glitch {
  font-size: 10rem;
  font-weight: 900;
  color: white;
  position: relative;
  text-shadow: 0.05em 0 0 rgba(255, 0, 0, 0.75),
               -0.05em -0.025em 0 rgba(0, 255, 0, 0.75),
               0.025em 0.05em 0 rgba(0, 0, 255, 0.75);
  animation: glitch-skew 3s infinite;
}

.glitch::before,
.glitch::after {
  content: attr(data-text);
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.glitch::before {
  animation: glitch-anim-1 0.3s infinite linear alternate-reverse;
  clip-path: polygon(0 0, 100% 0, 100% 45%, 0 45%);
  text-shadow: -4px 0 red;
}

.glitch::after {
  animation: glitch-anim-2 0.3s infinite linear alternate-reverse;
  clip-path: polygon(0 80%, 100% 20%, 100% 100%, 0 100%);
  text-shadow: 4px 0 blue;
}

@keyframes glitch-skew {
  0% { transform: skew(0deg); }
  10% { transform: skew(2deg); }
  20% { transform: skew(-2deg); }
  30% { transform: skew(1deg); }
  40% { transform: skew(-1deg); }
  50% { transform: skew(0deg); }
}

@keyframes glitch-anim-1 {
  0% { clip-path: polygon(0 0, 100% 0, 100% 45%, 0 45%); }
  100% { clip-path: polygon(0 30%, 100% 30%, 100% 75%, 0 75%); }
}

@keyframes glitch-anim-2 {
  0% { clip-path: polygon(0 80%, 100% 20%, 100% 100%, 0 100%); }
  100% { clip-path: polygon(0 60%, 100% 60%, 100% 100%, 0 100%); }
}

.glitch-sub {
  font-size: 1.5rem;
  margin-top: 1rem;
  color: #88ffff;
  letter-spacing: 4px;
  position: relative;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; text-shadow: 0 0 5px #00ffff; }
  50% { opacity: 0.7; text-shadow: 0 0 20px #00ffff; }
}

/* Терминал сообщений */
.message-box {
  margin: 2rem 0;
  background: rgba(0, 0, 0, 0.6);
  border-left: 3px solid #00ffff;
  padding: 1.5rem;
  backdrop-filter: blur(5px);
}

.terminal {
  font-family: monospace;
  font-size: 0.9rem;
}

.terminal-line {
  margin: 0.5rem 0;
  color: #88ffaa;
}

.prompt {
  color: #00ffff;
  margin-right: 0.5rem;
}

.error-line {
  color: #ff3366;
}

.success-line {
  color: #00ff88;
}

.blink {
  animation: blink-text 1s infinite;
}

@keyframes blink-text {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Кнопки */
.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin: 2rem 0;
  flex-wrap: wrap;
}

button {
  padding: 12px 24px;
  font-family: monospace;
  font-weight: bold;
  font-size: 0.85rem;
  letter-spacing: 2px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.btn-primary {
  background: linear-gradient(135deg, #00ffff, #0088ff);
  color: #000;
  box-shadow: 0 0 15px rgba(0, 255, 255, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 0 25px rgba(0, 255, 255, 0.6);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: wait;
}

.btn-secondary {
  background: transparent;
  border: 2px solid #00ffff;
  color: #00ffff;
}

.btn-secondary:hover {
  background: rgba(0, 255, 255, 0.1);
  transform: translateY(-2px);
}

/* Диагностическая панель */
.diagnostics-panel {
  margin-top: 2rem;
  background: rgba(0, 0, 0, 0.9);
  border: 1px solid #00ffff;
  border-radius: 4px;
  overflow: hidden;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.panel-header {
  background: rgba(0, 255, 255, 0.1);
  padding: 0.8rem 1rem;
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid #00ffff;
  color: #00ffff;
  font-size: 0.8rem;
}

.close {
  cursor: pointer;
  font-size: 1.2rem;
}

.panel-content {
  padding: 1rem;
}

.metric {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
  font-size: 0.85rem;
}

.metric span:first-child {
  color: #88ffff;
}

.value {
  color: #00ff88;
}

.value.error {
  color: #ff3366;
  animation: glitch-text 0.3s infinite;
}

@keyframes glitch-text {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; text-shadow: 0 0 2px red; }
}

.progress-bar {
  margin-top: 1rem;
  padding: 0.5rem;
}

.progress-label {
  font-size: 0.75rem;
  color: #88ffff;
  margin-bottom: 0.3rem;
}

.progress-track {
  background: rgba(0, 255, 255, 0.2);
  height: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  transition: width 0.5s ease;
}

.progress-value {
  font-size: 0.7rem;
  color: #ffaa00;
  margin-top: 0.3rem;
  text-align: center;
}

/* Футер */
.footer {
  margin-top: 3rem;
  padding-top: 1.5rem;
  border-top: 1px solid rgba(0, 255, 255, 0.2);
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.7rem;
  color: #5588aa;
  flex-wrap: wrap;
  gap: 1rem;
}

.status-leds {
  display: flex;
  gap: 0.5rem;
}

.led {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: blink 1s infinite;
}

.led.red {
  background: #ff3366;
  animation-delay: 0s;
}

.led.yellow {
  background: #ffaa00;
  animation-delay: 0.3s;
}

.led.green {
  background: #00ff88;
  animation-delay: 0.6s;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

/* Адаптивность */
@media (max-width: 768px) {
  .content {
    padding: 1rem;
  }
  
  .glitch {
    font-size: 5rem;
  }
  
  .glitch-sub {
    font-size: 1rem;
  }
  
  .nebula-title {
    font-size: 2rem;
    letter-spacing: 5px;
  }
  
  button {
    padding: 8px 16px;
    font-size: 0.75rem;
  }
  
  .terminal-line {
    font-size: 0.7rem;
  }
}
</style>