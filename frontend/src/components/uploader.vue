<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { CircleSlash2, FilePlus, FolderPlus, Upload, Trash, Copy, X, Loader, Check } from 'lucide-vue-next'
  import { niceFileSize, niceFileType, getApiUrl } from '../utils'
  import { store } from '../store'



  const apiUrl = getApiUrl()
  
  const fileInput = ref(null)
  const sharePanelVisible = ref(false)
  const shareUrl = ref('')
  const currentlyUploading = ref(false)
  const uploadBasket = ref([])
  const maxShareSize = ref(0)

  onMounted(async () => {
    const response = await fetch(`${apiUrl}/api/health`)
    const data = await response.json()
    maxShareSize.value = data.data.max_share_size
  })


  const showFilePicker = () => {
    fileInput.value.webkitdirectory = false
    fileInput.value.click()
  }

  const showFolderPicker = () => {
    fileInput.value.webkitdirectory = true
    fileInput.value.click()
  }

  const resetFileInput = () => {
    fileInput.value.value = null
  }

  const pushFile = file => {
    if (!uploadBasket.value.some(item => item.name === file.name)) {
      uploadBasket.value.push(file)
    }
  }

  const handleFileSelect = event => {
    if (event.target.files.length > 1) {
      for (let i = 0; i < event.target.files.length; i++) {
        pushFile(event.target.files[i])
      }
    }
    if (event.target.files.length === 1) {
      pushFile(event.target.files[0])
    }
    resetFileInput()
  }

  const removeFile = file => {
    uploadBasket.value = uploadBasket.value.filter(item => item.name !== file.name)
  }

  const totalSize = computed(() => {
    return uploadBasket.value.reduce((acc, file) => acc + file.size, 0)
  })

  const uploadFiles = async () => {
    if (totalSize.value > maxShareSize.value) {
      alert(`Total size of files is greater than the max share size of ${niceFileSize(maxShareSize.value)}`)
      return
    }
    const formData = new FormData()
    uploadBasket.value.forEach(file => {
      formData.append('files', file)
    })

    formData.append('name', 'test')
    formData.append('description', 'test')

    currentlyUploading.value = true
    const response = await fetch(`${apiUrl}/api/shares`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${store.jwt}`
      },
      body: formData
    })

    const data = await response.json()
    if (data.status_code === 200) {
      showSharePanel(createShareURL(data.data.share.long_id))
      uploadBasket.value = []
    }
    currentlyUploading.value = false
  }

  const createShareURL = longId => {
    const currentUrl = window.location.href
    const baseUrl = currentUrl.split('/').slice(0, -1).join('/')
    return `${baseUrl}/shares/${longId}`
  }

  const showSharePanel = url => {
    sharePanelVisible.value = true
    shareUrl.value = url
  }

  const showCopySuccess = ref(false)

  const copyShareUrl = () => {
    navigator.clipboard.writeText(shareUrl.value)
    showCopySuccess.value = true
    setTimeout(() => {
      showCopySuccess.value = false
    }, 2000)
  }
</script>

<template>
  <div class="upload-form">
    <div class="buttons">
      <button class="upload-files" @click="showFilePicker">
        <FilePlus />
        Add Files
      </button>
      <button class="upload-folder" @click="showFolderPicker">
        <FolderPlus />
        Add Folders
      </button>
    </div>
    <div class="max-size-label">{{ niceFileSize(totalSize) }} / {{ niceFileSize(maxShareSize) }}</div>
  </div>

  <div class="upload-basket">
    <div class="upload-basket-item" v-for="file in uploadBasket" :key="file.name" v-if="uploadBasket.length > 0">
      <div class="name">
        {{ file.name }}
      </div>
      <div class="meta">
        <div class="size">
          {{ niceFileSize(file.size) }}
        </div>
        <div class="type">
          {{ niceFileType(file.type) }}
        </div>
      </div>
      <div class="hover-actions">
        <button class="hover-action" @click="removeFile(file)">
          <Trash />
        </button>
      </div>
    </div>

    <div class="upload-basket-empty" v-else>
      <div class="upload-basket-empty-text">
        <CircleSlash2 />
        No files added yet
      </div>
    </div>
  </div>

  <div class="upload-button-container">
    <button class="upload-button" :disabled="uploadBasket.length === 0 || currentlyUploading" @click="uploadFiles" :class="{ uploading: currentlyUploading }">
      <div class="loader" v-if="currentlyUploading">
        <Loader />
      </div>
      <Upload v-else />
      <template v-if="uploadBasket.length > 0">Upload{{ currentlyUploading ? 'ing' : '' }} {{ uploadBasket.length }} file{{ uploadBasket.length > 1 ? 's' : '' }}</template>
      <template v-else>No files added yet</template>
    </button>
  </div>
  <input type="file" @change="handleFileSelect" style="display: none" webkitdirectory directory ref="fileInput" multiple />
  <div class="sharePanel" :class="{ visible: sharePanelVisible }">
    <div class="sharePanel-content">
      <div class="sharePanel-close" @click="sharePanelVisible = false">
        <X />
      </div>
      <div class="sharePanel-title">Share URL</div>
      <div class="sharePanel-url">
        {{ shareUrl }}
        <button class="sharePanel-copy-button" @click="copyShareUrl">
          <Check v-if="showCopySuccess" />
          <Copy v-else />
        </button>
      </div>
    </div>
  </div>
</template>
