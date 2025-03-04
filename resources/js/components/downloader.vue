<script setup>
import { ref, onMounted } from 'vue'
import { niceFileSize, timeUntilExpiration, getApiUrl, niceFileType, niceFileName } from '../utils'
import { FileIcon, HeartCrack, TrendingDown } from 'lucide-vue-next'
import { getShare } from '../api'
const apiUrl = getApiUrl()

const share = ref(null)
const showFilesCount = ref(5)
const shareExpired = ref(false)
const downloadLimitReached = ref(false)

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
  try {
    share.value = await getShare(props.downloadShareCode)
  } catch (error) {
    if (error.message == 'Download limit reached') {
      downloadLimitReached.value = true
    } else if (error.message == 'Share expired') {
      shareExpired.value = true
    }
  }
}

const downloadFiles = () => {
  const downloadUrl = `${apiUrl}/api/shares/${props.downloadShareCode}/download`
  window.open(downloadUrl, '_blank')
}

const splitFullName = (fullName) => {
  if (!fullName) {
    return 'creator'
  }
  const nameParts = fullName.split(' ')
  return nameParts[0]
}
</script>

<template>
  <div class="download-panel-content">
    <template v-if="share">
      <h1>
        {{ share.name }}
      </h1>
      <div class="total-size">{{ $t('total_size') }}: {{ niceFileSize(share.size) }}</div>
      <div class="file-count">
        {{ $t('share.contains.count', 'Contains: {value} files', { value: share.file_count }) }}
      </div>
      <div class="share-expires">
        {{
          $t('share.expires.in', {
            days: timeUntilExpiration(share.expires_at).days,
            hours: timeUntilExpiration(share.expires_at).hours,
            minutes: timeUntilExpiration(share.expires_at).minutes
          })
        }}
      </div>
      <div class="file-list">
        <div v-for="file in share.files.slice(0, showFilesCount)" :key="file" class="file-item">
          <div class="file-name">
            <div class="name">
              {{ niceFileName(file.name) }}
              <div class="size">
                {{ niceFileSize(file.size) }}
              </div>
            </div>
          </div>

          <div class="type">
            {{ niceFileType(file.type) }}
          </div>
        </div>
        <div v-if="share.files.length > showFilesCount" class="file-item more-files">
          <div class="file-name more-files">and {{ share.files.length - showFilesCount }} more</div>
        </div>
      </div>
      <div class="share-message mt-3" v-if="share.description">
        <h6>{{ $t('message.from', { name: splitFullName(share.user.name) }) }}</h6>
        {{ share.description }}
      </div>
      <div class="download-button-container mt-3">
        <button class="download-button" @click="downloadFiles">
          {{ $t('download.files', 'Download {value} files', { value: share.file_count }) }}
        </button>
      </div>
    </template>
    <template v-else>
      <template v-if="shareExpired">
        <h1>
          <HeartCrack />
          {{ $t('share.expired') }}
        </h1>
        <p>{{ $t('share.expired.message') }}</p>
      </template>
      <template v-else-if="downloadLimitReached">
        <h1>
          <TrendingDown />
          {{ $t('share.download_limit_reached') }}
        </h1>
        <p>
          {{ $t('share.download_limit_reached.message') }}
        </p>
      </template>
      <h1 v-else>{{ $t('share.data_loading') }}</h1>
    </template>
  </div>
</template>
