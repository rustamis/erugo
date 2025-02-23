<script setup>
  import { ref, computed } from 'vue'
  import { Search } from 'lucide-vue-next'
  const props = defineProps({
    modelValue: {
      type: [File, null, String],
      required: true
    },
    accept: {
      type: String,
      default: '*'
    },
    label: {
      type: String,
      default: 'Select File'
    }
  })

  const fileInput = ref(null)

  //emit
  const emit = defineEmits(['update:modelValue'])

  //methods
  const handleFileChange = event => {
    emit('update:modelValue', event.target.files[0])
  }

  const triggerFileInput = () => {
    fileInput.value.click()
  }

  const buttonMessage = computed(() => {
    let message = ''
    if (props.modelValue) {
      message = sensibleButtonMessage(typeof props.modelValue === 'string' ? props.modelValue : props.modelValue.name)
    } else {
      message = props.label
    }
    return message
  })

  const sensibleButtonMessage = (message) => {
    const maxLength = 25
    if (message.length > maxLength) {
      return message.slice(0, maxLength) + `&hellip;`
    }
    return message
  }
</script>

<template>
  <div class="file-input">
    <input type="file" @change="handleFileChange" ref="fileInput" :accept="accept" />
    <div class="file-input-button" @click="triggerFileInput">
      <div class="file-label" v-html="buttonMessage"></div>
      <button>
        <Search />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .file-input {
    position: relative;
    width: 100%;

  }
  input[type='file'] {
    display: none;
  }
  .file-input-button {
    width: 100%;
    padding: 10px;
    border-radius: 5px;
    border: none;
    background-color: var(--accent-color-light);
    color: var(--secondary-color);
    margin-bottom: 10px;
    border: 1px solid color-mix(in srgb, var(--secondary-color), rgba(255,255,255,0) 80%);
    transition: all 0.3s ease-in-out;
    outline: 2px solid transparent;
    height: 50px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    .file-label {
      flex: 1;
      font-style: italic;
      padding-right: 10px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    button {
      width: 40px;
      height: 40px;
      svg {
        width: 20px;
        height: 20px;
        margin-right: 0;
        margin-top: 0;
      }
    }
  }
</style>
