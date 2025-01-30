<script setup>
  import { ref, onMounted } from 'vue'
  import { niceFileSize, niceExpirationDate, timeUntilExpiration, getApiUrl } from '../utils'
  import { FileIcon, HeartCrack } from 'lucide-vue-next'

  const apiUrl = getApiUrl()

  const share = ref(null)
  const showFilesCount = ref(5)
  const shareExpired = ref(false)

  //define props
  const props = defineProps({
    downloadShareCode: {
      type: String,
      required: true
    }
  })

  onMounted(() => {
    fetchShare()
  })

  const fetchShare = async () => {
    const response = await fetch(`${apiUrl}/api/shares/${props.downloadShareCode}/`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    const data = await response.json()
    if (data.status_code === 200) {
      share.value = data.data.share
    } else {
      shareExpired.value = true
    }
  }

  const downloadFiles = () => {
    const downloadUrl = `${apiUrl}/api/shares/${props.downloadShareCode}/download/`
    window.open(downloadUrl, '_blank')
  }
</script>

<template>
  <div class="download-panel-content">
    <template v-if="share">
      <h1>Share contains {{ share.num_files }} files</h1>
      <div class="total-size">Total size: {{ niceFileSize(share.total_size) }}</div>
      <div class="share-expires">Expires in {{ timeUntilExpiration(share.expiration_date) }}</div>
      <div class="file-list">
        <div v-for="file in share.files.slice(0, showFilesCount)" :key="file" class="file-item">
          <div class="file-name">
            <FileIcon class="file-icon" />
            {{ file }}
          </div>
        </div>
        <div v-if="share.files.length > showFilesCount" class="file-item more-files">
          <div class="file-name more-files">and {{ share.files.length - showFilesCount }} more</div>
        </div>
      </div>
      <div class="download-button-container mt-3">
        <button class="download-button" @click="downloadFiles">Download</button>
      </div>
    </template>
    <template v-else>
      <template v-if="shareExpired">
        <h1><HeartCrack /> Share Expired</h1>
        <p>The share you are trying to download has expired. Please ask the share creator to create a new share.</p>
      </template>
      <h1 v-else>Loading...</h1>
    </template>
  </div>
</template>
