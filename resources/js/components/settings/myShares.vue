<script setup>
import { ref, onMounted, inject } from 'vue'
import { getMyShares, expireShare, extendShare, setDownloadLimit } from '../../api'
import {
  SquareArrowOutUpRight,
  CalendarPlus,
  CalendarX2,
  HardDriveDownload,
  MessageCircleQuestion,
} from 'lucide-vue-next'
import { useToast } from 'vue-toastification'
import { niceFileSize, niceDate, niceFileName, niceNumber } from '../../utils'
import HelpTip from '../helpTip.vue'

const showHelpTip = inject('showHelpTip')
const hideHelpTip = inject('hideHelpTip')

const toast = useToast()
const maxFilesToShow = 4

const shares = ref([])
onMounted(async () => {
  loadShares()
})

const loadShares = async () => {
  shares.value = await getMyShares()
}

const handleExpireShareClick = async (share) => {
  expireShare(share.id)
    .then(() => {
      toast.success('Share expired')
      loadShares()
    })
    .catch((error) => {
      toast.error('Failed to expire share')
    })
}

const handleExtendShareClick = async (share) => {
  extendShare(share.id)
    .then(() => {
      toast.success('Share extended')
      loadShares()
    })
    .catch((error) => {
      toast.error('Failed to extend share')
    })
}

const handleDownloadLimitChange = async (share) => {
  let newLimit = null
  if (share.download_limit == '' || share.download_limit == null) {
    newLimit = -1
  } else {
    newLimit = parseInt(share.download_limit)
  }


  if (isNaN(newLimit)) {
    toast.error('Invalid download limit')
    return
  }
  setDownloadLimit(share.id, newLimit)
    .then(() => {
      toast.success('Download limit changed')
      loadShares()
    })
    .catch((error) => {
      toast.error('Failed to change download limit')
    })
}

const downloadShare = async (share) => {
  window.location.href = `/api/shares/${share.long_id}/download`
}

const enableExpireShareButton = (share) => {
  return !share.expired && !share.deleted
}

const enableExtendShareButton = (share) => {
  return !share.deleted
}

const enableDownloadButton = (share) => {
  return !share.expired && !share.deleted
}

</script>

<template>
  <div>
    <HelpTip id="download-limit-help-tip" header="Download Limit Help">
      <p>
        You can set the download limit for any share by clicking the right-side number and typing in the new limit.
        Clear the box to remove the limit.
      </p>
      <p>Changes will be saved automatically when you click outside the box.</p>
    </HelpTip>
    <table v-if="shares.length > 0">
      <thead>
        <tr>
          <th>Name</th>
          <th>Files</th>
          <th>
            Downloads
            <MessageCircleQuestion @click.stop="showHelpTip($event, '#download-limit-help-tip')" />
          </th>
          <th>Dates</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="share in shares" :key="share.id">
          <td width="1" style="white-space: nowrap">
            <div class="slide-text">
              <strong class="content">{{ share.name }}</strong>
            </div>
            <a :href="`/shares/${share.long_id}`" target="_blank" class="share_long_id">
              <SquareArrowOutUpRight />
              {{ share.long_id }}
            </a>
          </td>
          <td style="vertical-align: top">
            <h6 class="file-count">
              {{ share.files.length }} file{{ share.files.length > 1 ? 's' : '' }}
              {{ share.files.maxFilesToShow > 1 ? 'including:' : '' }}
            </h6>
            <div class="files-container pt-1">
              <div class="file" v-for="file in share.files.slice(0, maxFilesToShow)" :key="file.id">
                <div class="file-name" :title="file.name">{{ niceFileName(file.name) }}</div>
                <div class="file-size">
                  {{ niceFileSize(file.size) }}
                </div>
              </div>
              <div class="some-more" v-if="share.files.length > maxFilesToShow">
                <span>And {{ share.files.length - maxFilesToShow }} more</span>
              </div>
            </div>
          </td>
          <td width="1" style="white-space: nowrap" class="text-center">
            <div class="download_limit_manager">
              <div class="download_count">
                <label class="count_label">Downloads</label>
                {{ niceNumber(share.download_count) }}
                <span>/</span>
              </div>
              <input
                class="download_limit_input"
                v-model="share.download_limit"
                @change="handleDownloadLimitChange(share)"
                placeholder="∞"
              />
            </div>
          </td>
          <td width="1" style="white-space: nowrap">
            <div class="date-container">
              <div class="date">
                <span>Created:</span>
                {{ niceDate(share.created_at) }}
              </div>
              <div class="date">
                <span>Expires:</span>
                <template v-if="share.expired">
                  <strong class="ps-1 text-danger">Expired</strong>
                </template>
                <template v-else>
                  {{ niceDate(share.expires_at) }}
                </template>
              </div>
              <div class="date">
                <span>Deletes:</span>
                <template v-if="share.deleted">
                  <strong class="ps-1 text-danger">Deleted</strong>
                </template>
                <template v-else>
                  {{ niceDate(share.deletes_at) }}
                </template>
              </div>
            </div>
          </td>
          <td width="1" style="white-space: nowrap">
            <button @click="handleExpireShareClick(share)" class="clear-button" :disabled="!enableExpireShareButton(share)">
              <CalendarX2 />
              Expire Now
            </button>
            <button @click="handleExtendShareClick(share)" class="secondary" :disabled="!enableExtendShareButton(share)">
              <CalendarPlus />
              Extend
            </button>
            <button @click="downloadShare(share)" class="secondary icon-only" title="Download all files" :disabled="!enableDownloadButton(share)">
              <HardDriveDownload style="margin-right: 0" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style lang="scss" scoped>
.files-container {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;
  .file {
    display: flex;
    flex-direction: column;
    background-color: color-mix(in srgb, var(--primary-color), rgba(255, 255, 255, 0) 95%);
    border-radius: 5px;
    padding: 5px 10px;
    gap: 1px;
    .file-name {
      font-size: 0.85rem;
      font-weight: bold;
      color: var(--secondary-color);
    }
    .file-size {
      font-size: 0.7rem;
      color: var(--secondary-color);
    }
  }
  .some-more {
    font-size: 0.7rem;
    color: var(--secondary-color);
    margin-left: 10px;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.date-container {
  display: flex;
  flex-direction: column;
  gap: 5px;
  .date {
    background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 95%);
    border-radius: 5px;
    padding: 5px 10px;
    gap: 5px;
    span {
      display: inline-block;
      font-weight: bold;
      background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 95%);
      border-radius: 5px;
      padding: 5px 10px;
      margin-left: -10px;
      margin-bottom: -5px;
      margin-top: -5px;
      height: calc(100% + 10px);
      width: 100px;
      margin-right: 10px;
    }
  }
}

.share_long_id {
  display: block;
  font-size: 1rem;
  color: var(--primary-color);
  text-decoration: none;
  font-weight: bold;

  svg {
    width: 1rem;
    height: 1rem;
    margin-right: 5px;
    margin-top: -2px;
    color: var(--secondary-color);
  }
}

.file-count {
  background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 95%);
  margin-left: -10px;
  margin-top: -10px;
  margin-right: -10px;
  padding: 5px 10px;

  color: var(--secondary-color);
  font-weight: 500;
}

td {
  a {
    color: var(--primary-color);
    text-decoration: none;
    cursor: pointer;
    font-size: 0.75rem;
    margin-top: 10px;
    display: block;
    &:hover {
      text-decoration: underline;
    }
  }
}

.download_limit_manager {
  --height: 40px;
  display: flex;
  flex-direction: row;
  align-items: center;
  background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 95%);
  height: var(--height);
  border-radius: 5px;
  .download_count {
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
    width: 90px;
    padding-bottom: 6px!important;
    span {
      position: absolute;
      left: 86.5px;
      opacity: 0.3;
    }
    .count_label {
      position: absolute;
      left: 0;
      right: 0;
      top: 0;
      bottom: 0;
      display: flex;
      align-items: flex-end;
      justify-content: center;
      opacity: 0.3;
      font-size: 0.5rem;
      font-weight: normal;
      padding-bottom: 1.5px;
    }
    &:after {
      content: 'Limit';
      position: absolute;
      left: 90px;
      width: 90px;
      top: 0;
      bottom: 0;
      display: flex;
      align-items: flex-end;
      justify-content: center;
      opacity: 0.3;
      font-size: 0.5rem;
      font-weight: normal;
      padding-bottom: 1.5px;
    }
  }
  .download_limit_input {
    position: relative!important;
    background-color: color-mix(in srgb, var(--secondary-color), rgba(255, 255, 255, 0) 100%);
    height: var(--height);
    border: none;
    border-radius: 0 3px 3px 0;
    text-align: center;
    margin: 0;
    width: 90px;
    padding-bottom: 16px!important;
    &:focus {
      outline: none;
    }
  }
}
</style>
