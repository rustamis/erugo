<script setup>
  import { ref, computed, getCurrentInstance, watch } from 'vue'

  const props = defineProps({
    modelValue: {
      type: [String, Object],
      required: true
    },
    label: {
      type: String,
      required: true
    }
  })

  const emit = defineEmits(['update:modelValue'])
  const colourPickerInput = ref(null)
  const colourRaw = ref(props.modelValue)
  const opacity = ref(1)

  const componentId = computed(() => {
    const instance = getCurrentInstance()
    return instance.uid
  })

  const colour = computed(() => {
    let hex = colourRaw.value + ''
    hex = hex.toUpperCase()
    return {
      hex: hex,
      rgb: hexToRgb(hex)
    }
  })

  watch(colour, (newVal) => {
    emit('update:modelValue', newVal.rgb)
  })

  function hexToRgb(hex) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
    return result
      ? {
          r: parseInt(result[1], 16),
          g: parseInt(result[2], 16),
          b: parseInt(result[3], 16)
        }
      : null
  }

  function openColourPicker() {
    colourPickerInput.value.click()
  }
</script>

<template>
  <div class="colour-picker">
    <div class="input-container" @click="openColourPicker">
      <label :for="componentId">{{ label }}</label>
      <div class="colour-display">
        <div class="colour-display-swatch" :style="{ backgroundColor: colour.hex }"></div>
        <div class="colour-display-label">
          <div class="colour-display-label-hex">{{ colour.hex }}</div>
          <div class="colour-display-label-rgb">rgb({{ colour.rgb.r }}, {{ colour.rgb.g }}, {{ colour.rgb.b }})</div>
        </div>
      </div>
      <div class="colour-display-opacity">
        <div class="colour-display-opacity-label">Opacity</div>
        <div class="colour-display-opacity-slider">
          <input type="range" min="0" max="1" step="0.01" v-model="opacity" />
        </div>
      </div>
      <input type="color" :id="componentId" v-model="colourRaw" ref="colourPickerInput" />
    </div>
  </div>
</template>

<style scoped>
  .colour-picker {
    position: relative;

    input[type="color"] {
      opacity: 0;
      width: 0;
      height: 0;
      position: absolute;
      top: 0;
      left: 0;
    }

    .colour-display {
      border-radius: 50%;
      border: 1px solid var(--accent-colour-light-transparent);
      cursor: pointer;
      border: none;
      outline: none;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: flex-start;
      .colour-display-label {
        margin-left: 10px;
        font-size: 0.8rem;
        font-weight: 300;
        color: var(--secondary-colour);
      }
      .colour-display-swatch {
        width: 50px;
        height: 50px;
        border-radius: 50%;
      }
    }
  }
</style>
