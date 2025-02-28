<script setup>
import { ref, onMounted, watch, defineEmits } from 'vue'
import { X, Star } from 'lucide-vue-next'
import ColourPicker from '../colourPicker.vue'

import { saveTheme } from '../../api'
import injectThemeVariables from '../../lib/injectThemeVariables'
const activeTab = ref('none')

const themeUnderEdit = ref({})

const props = defineProps({
  currentTheme: {
    type: Object,
    required: true
  }
})

watch(
  themeUnderEdit,
  (newVal) => {
    injectThemeVariables('.theme-editor-preview-area', newVal)
  },
  { deep: true }
)

onMounted(() => {
  setTimeout(() => {
    activeTab.value = 'links'
  }, 10)
  injectThemeVariables('.theme-editor-preview-area', themeUnderEdit.value)
})

const handleSaveTheme = async () => {
  const theme = await saveTheme({
    name: 'EditorUserTheme',
    theme: themeUnderEdit.value
  })
  injectThemeVariables('body', themeUnderEdit.value)

}

defineEmits(['close'])
</script>

<template>
  <div class="theme-editor-overlay">
    <div class="theme-editor-content">
      <div class="theme-editor-header">
        <h3>Theme Editor</h3>
        <button class="theme-editor-save-button" @click="handleSaveTheme">Save Theme</button>
        <button class="theme-editor-close-button icon-only" @click="$emit('close')">
          <X />
        </button>
      </div>

      <div class="theme-editor-body">


        <div class="theme-editor-editor">
          {{ currentTheme }}
        </div>


      </div>
    </div>
  </div>
</template>
<style scoped lang="scss">
.theme-editor-content {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: var(--panel-background-color);
  border-radius: var(--panel-border-radius);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  z-index: 230;
  .theme-editor-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    height: 80px;
    background: var(--panel-header-background-color);
    border-radius: var(--panel-border-radius) var(--panel-border-radius) 0 0;
    width: 100%;
    h3 {
      margin: 0;
      font-size: 1.2rem;
      font-weight: 600;
      color: var(--panel-header-text-color);
    }
  }
  .theme-editor-body {
    display: flex;
    flex-direction: row;
    width: 100%;
    flex-grow: 1;
    height: calc(100% - 80px);
  }
}





.theme-editor-editor {
  background: var(--panel-background-color);
  height: 100%;
  width: 100%;
  padding: 20px;
  overflow-y: auto;
}

.theme-editor-preview-area {
  width: 500px;
  min-width: 500px;
  padding: 20px;
  height: 100%;
  overflow-y: auto;
  background: var(--panel-section-background-color);
}

.preview-card-header {
  margin-bottom: 10px;
  background: var(--panel-header-background-color);
  padding: 20px;
  border-radius: var(--panel-border-radius);
  color: var(--panel-header-text-color);
  h4 {
    margin: 0;
    margin-bottom: 10px;
  }
  p {
    margin: 0;
  }
}

.preview-card-preview-area {
  padding: 20px;
  background: var(--panel-background-color);
  border-radius: var(--panel-border-radius);
  margin-top: 20px;

  .preview-card-preview-area-item {
    margin-bottom: 10px;
    .help-text {
      margin-top: 5px;
      font-size: 0.8rem;
      color: var(--panel-text-color-alt);
    }
  }
}
.preview-link {
  color: var(--link-color);
}

.preview-link-hover {
  color: var(--link-color-hover);
}

.preview-link-active {
  color: var(--link-color-active);
}

.preview-link-disabled {
  color: var(--link-color-disabled);
  pointer-events: none;
}

.preview-link-interactive {
  color: var(--link-color);
  text-decoration: none;
  &:hover {
    color: var(--link-color-hover);
  }
}

.preview-button {
  background: var(--primary-button-background-color);
  color: var(--primary-button-text-color);
  display: block;
  width: 100%;
  border-radius: var(--button-border-radius);
}

.preview-button-hover {
  background: var(--primary-button-background-color-hover);
  color: var(--primary-button-text-color-hover);
}

.preview-button-active {
  background: var(--primary-button-background-color-active);
  color: var(--primary-button-text-color-active);
  outline: 2px dashed color-mix(in srgb, var(--primary-button-background-color), white 90%);
}

.preview-button-disabled {
  background: var(--primary-button-background-color-disabled);
  color: var(--primary-button-text-color-disabled);
  pointer-events: none;
}

.preview-button-secondary {
  background: var(--secondary-button-background-color);
  color: var(--secondary-button-text-color);
}

.preview-button-secondary-hover {
  background: var(--secondary-button-background-color-hover);
  color: var(--secondary-button-text-color-hover);
}

.preview-button-secondary-active {
  background: var(--secondary-button-background-color-active);
  color: var(--secondary-button-text-color-active);
  outline: 2px dashed color-mix(in srgb, var(--primary-button-background-color), white 90%);
}

.preview-button-secondary-disabled {
  background: var(--secondary-button-background-color-disabled);
  color: var(--secondary-button-text-color-disabled);
  pointer-events: none;
}

.setting-group-body-item {
  border: 1px dashed rgba(57, 115, 144, 0.1);
  margin-bottom: 10px;
  padding: 10px;
  border-radius: var(--panel-border-radius);
}

.preview-card-preview-area-item {
  position: relative;
  padding: 10px;
  border-radius: var(--panel-border-radius);
  transition: background-color 0.3s ease;
  border: 1px dashed rgba(57, 115, 144, 0.1);
}
</style>

<!-- not scoped because the targets are created programmatically -->
<style lang="scss">
.connector-svg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 10;
  /* Make sure background is transparent */
  background-color: transparent;
}

/* Styling for the connector path */
.connector-path {
  opacity: 0.7;
  transition: all 0.5s ease-in-out;
  opacity: 1;
  stroke: var(--link-color);
  stroke-width: 0;
  &.highlighted {
    opacity: 1;
    stroke: var(--link-color-hover);
    stroke-width: 2;
  }
}
</style>
