<script setup>
</script>

<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style>
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  background-color: #f4f4f4;
}
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
}
.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
.Vue-Toastification__toast {
  font-family: inherit;
  border-radius: 8px !important;
  padding: 15px 20px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.Vue-Toastification__toast--success {
  background: linear-gradient(135deg, #2ecc71 0%, #27ae60 100%) !important;
}

.Vue-Toastification__toast--error {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%) !important;
}

.Vue-Toastification__icon {
  font-size: 1.4rem !important;
}
</style>
