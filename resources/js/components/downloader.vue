<script setup>
import { ref, onMounted } from "vue";
import {
    niceFileSize,
    timeUntilExpiration,
    getApiUrl,
    niceFileType,
    niceFileName,
} from "../utils";
import { FileIcon, HeartCrack, TrendingDown } from "lucide-vue-next";
import { getShare } from "../api";
const apiUrl = getApiUrl();

const share = ref(null);
const showFilesCount = ref(5);
const shareExpired = ref(false);
const downloadLimitReached = ref(false);

//define props
const props = defineProps({
    downloadShareCode: {
        type: String,
        required: true,
    },
});

onMounted(() => {
    fetchShare();
});

const fetchShare = async () => {
    try {
        share.value = await getShare(props.downloadShareCode);
    } catch (error) {
        if (error.message == 'Download limit reached') {
            downloadLimitReached.value = true;
        } else if (error.message == 'Share expired') {
            shareExpired.value = true;
        }
    }
};

const downloadFiles = () => {
    const downloadUrl = `${apiUrl}/api/shares/${props.downloadShareCode}/download`;
    window.open(downloadUrl, "_blank");
};
</script>

<template>
    <div class="download-panel-content">
        <template v-if="share">
            <h1>
                Share contains {{ share.file_count }} file{{
                    share.file_count > 1 ? "s" : ""
                }}
            </h1>
            <div class="total-size">
                Total size: {{ niceFileSize(share.size) }}
            </div>
            <div class="share-expires">
                Expires in {{ timeUntilExpiration(share.expires_at) }}
            </div>
            <div class="file-list">
                <div
                    v-for="file in share.files.slice(0, showFilesCount)"
                    :key="file"
                    class="file-item"
                >
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
                <div
                    v-if="share.files.length > showFilesCount"
                    class="file-item more-files"
                >
                    <div class="file-name more-files">
                        and {{ share.files.length - showFilesCount }} more
                    </div>
                </div>
            </div>
            <div class="download-button-container mt-3">
                <button class="download-button" @click="downloadFiles">
                    Download {{ share.file_count }} file{{
                        share.file_count > 1 ? "s" : ""
                    }}
                </button>
            </div>
        </template>
        <template v-else>
            <template v-if="shareExpired">
                <h1><HeartCrack /> Share Expired</h1>
                <p>
                    The share you are trying to download has expired. Please ask
                    the share creator to create a new share.
                </p>
            </template>
            <template v-else-if="downloadLimitReached">
                <h1><TrendingDown /> Download Limit Reached</h1>
                <p>
                    The share you are trying to download has reached its download limit. Please ask
                    the share creator to increase the download limit.
                </p>
            </template>
            <h1 v-else>Loading...</h1>
        </template>
    </div>
</template>
