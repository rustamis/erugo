<script setup>
import { ref } from 'vue'
import { CircleCheckBig } from 'lucide-vue-next'
const showCopiedMessage = ref(false)

const props = defineProps({
  variable: {
    type: String,
    required: true
  }
})

const copyToClipboard = () => {
  navigator.clipboard.writeText(props.variable)
  showCopiedMessage.value = true
  setTimeout(() => {
    showCopiedMessage.value = false
  }, 5000)
}
</script>
<template>
  <div class="env-variable-name" @click="copyToClipboard">
    <div class="label">env</div>
    <div class="value">{{ variable }}</div>
    <div class="copied-message" :class="{ 'show': showCopiedMessage }">Copied to clipboard <CircleCheckBig /></div>
  </div>
</template>
<style scoped lang="scss">
.env-variable-name {
  position: relative;
  --height: 40px;
  display: flex;
  flex-direction: row;
  align-items: center;
  background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 95%);
  height: var(--height);
  border-radius: 5px;
  cursor: pointer;
  .label {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    padding-left: 10px;
    padding-right: 10px;
    border-radius: 3px;
    border: none;
    color: var(--secondary-color);
    outline: none;
    height: var(--height);
    background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 95%);
    font-weight: bold;
  }
  .value {
    position: relative !important;
    background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 100%);
    height: var(--height);
    border: none;
    border-radius: 0 3px 3px 0;
    text-align: center;
    margin: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-grow: 1;
    font-style: italic;
    font-size: 0.8rem;
  }
  .copied-message {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 0.8rem;
    border-radius: 5px;
    background-color: color-mix(in srgb, var(--accent-color-light), rgba(255, 255, 255, 0) 5%);
    opacity: 0;
    transition: opacity 0.3s ease-in-out;
    pointer-events: none;
    svg {
      width: 16px;
      height: 16px;
      margin-left: 5px;
    }
    &.show {
      opacity: 1;
      pointer-events: auto;
    }
  }
}
</style>
