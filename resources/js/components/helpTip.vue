<script setup>
import { X } from 'lucide-vue-next'
import { inject } from 'vue'

const hideHelpTip = inject('hideHelpTip')

//define props
const props = defineProps({
  id: {
    type: String,
    required: true
  },
  header: {
    type: String,
    required: true
  }
})
</script>

<template>
  <div class="help-tip" :id="id">
      <div class="help-tip-header">
        <h6>{{ header }}</h6>
        <X @click="hideHelpTip(`#${id}`)" class="help-tip-close" />
      </div>
      <div class="help-tip-content">
        <slot></slot>
      </div>
    </div>
</template>