<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { CircleSlash2, FilePlus, FolderPlus, Upload, Trash, Copy, X, Loader, Check } from 'lucide-vue-next'
  import { niceFileSize, niceFileType, simpleUUID, getApiUrl } from '../utils'
  import { createShare, getHealth } from '../api'

  const apiUrl = getApiUrl()
  const fileInput = ref(null)
  const sharePanelVisible = ref(false)
  const shareUrl = ref('')
  const currentlyUploading = ref(false)
  const uploadBasket = ref([])
  const maxShareSize = ref(0)
  const uploadProgress = ref(0)
  const uploadedBytes = ref(0)
  const totalBytes = ref(0)

  onMounted(async () => {
    const health = await getHealth()
    maxShareSize.value = health.max_share_size
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
    //Simple UUID-like string to track upload progress via SSE
    const uploadId = simpleUUID()
    currentlyUploading.value = true

    if (totalSize.value > maxShareSize.value) {
      alert(`Total size of files is greater than the max share size of ${niceFileSize(maxShareSize.value)}`)
      return
    }

    setTimeout(() => {
      const eventSource = createEventSource(uploadId)
    }, 1)

    try {
      const share = await createShare(uploadBasket.value, 'test', 'test', uploadId)
      console.log(share)
      showSharePanel(createShareURL(share.long_id))
      uploadBasket.value = []
    } catch (error) {
      console.error(error)
    } finally {
      currentlyUploading.value = false
      setTimeout(() => {
        uploadProgress.value = 0
        uploadedBytes.value = 0
        totalBytes.value = 0
      }, 1000)
    }
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
    }, 10)
  }

  const createEventSource = uploadId => {
    const eventSource = new EventSource(`${apiUrl}/api/shares/progress/${uploadId}`)
    console.log(eventSource)
    eventSource.onmessage = event => {
      const progress = JSON.parse(event.data)
      uploadProgress.value = progress.totalProgress
      uploadedBytes.value = progress.totalBytesRead
      totalBytes.value = progress.totalFileSize
    }

    //if we get an error 10 times, close the event source
    let errorCount = 0
    eventSource.onerror = () => {
      errorCount++
      if (errorCount >= 10) {
        eventSource.close()
        currentlyUploading.value = false
      }
    }
    return eventSource
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
    <div>
      <div class="progress-bar-container" :class="{ visible: currentlyUploading }">
        <div class="progress-bar">
          <div class="progress-bar-fill" :style="{ width: `${uploadProgress}%` }"></div>
        </div>
        <div class="progress-bar-text">
          {{ Math.round(uploadProgress) }}%
          <div class="progress-bar-text-sub">{{ niceFileSize(uploadedBytes) }} / {{ niceFileSize(totalBytes) }}</div>
        </div>
      </div>
    </div>
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

<style scoped lang="scss">
  .progress-bar-container {
    margin-top: -20px;
    // width: 300px;
    // height: 30px;
    background-color: var(--accent-color-light);
    border-radius: 5px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    position: absolute;
    opacity: 0;
    transition: all 0.3s ease-in-out;
    left: 0;
    right: 0;
    top: 20px;
    bottom: 0;
    z-index: 1000;
    pointer-events: none;

    &.visible {
      opacity: 1;
    }

    .progress-bar {
      height: 100%;
      width: 100%;
      background: transparent;
      .progress-bar-fill {
        background-color: var(--primary-color);
        border-radius: 5px;
        transition: all 1s linear;
        height: 100%;
      }
    }
    .progress-bar-text {
      font-size: 24px;
      color: var(--secondary-color);
      font-weight: 600;
      position: absolute;
      left: 0;
      top: 0;
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      text-align: center;
      flex-direction: column;
      .progress-bar-text-sub {
        font-size: 10px;
        color: var(--secondary-color);
        font-weight: 400;
      }
    }
  }
</style>
