<script setup>
  import { ref, onMounted, nextTick } from 'vue'
  import { LogOut, Settings as SettingsIcon } from 'lucide-vue-next'
  import Uploader from './components/uploader.vue'
  import Downloader from './components/downloader.vue'
  import Auth from './components/auth.vue'
  import Settings from './components/settings.vue'
  import { unsplashImages } from './unsplashImages'
  import { getApiUrl } from './utils'
  import { domData } from './domData'
  import { emitter, store } from './store'
  import { logout } from './api'

  const apiUrl = getApiUrl()

  const logoUrl = `${apiUrl}/logo`
  const version = ref()
  const logoWidth = ref(100)

  const auth = ref(null)
  const downloadShareCode = ref('')
  const settingsPanel = ref(null)

  onMounted(() => {
    setMode()
    changeBackground()
    setTimeout(changeBackground, 180000) //change every 3 minutes
    version.value = domData().version
    logoWidth.value = domData().logo_width
    emitter.on('showPasswordResetForm', () => {
      settingsPanel.value.setActiveTab('myProfile')
      nextTick(() => {
        store.setSettingsOpen(true)
        nextTick(() => {
          emitter.emit('profileEditActive')
        })
      })
    })
  })

  const setMode = () => {
    if (window.location.pathname.includes('shares')) {
      store.setMode('download')
      downloadShareCode.value = window.location.pathname.split('/').pop()
      setPageTitle('Download Share')
    } else {
      store.setMode('upload')
      setPageTitle('Create Share')
    }
  }

  const setPageTitle = title => {
    let currentTitle = document.title
    document.title = `${currentTitle} - ${title}`
  }

  const handleLogoutClick = () => {
    logout()
  }

  const changeBackground = async () => {
    let backgrounds = document.querySelectorAll('.backgrounds-item')
    backgrounds.forEach(background => {
      background.classList.remove('active')
    })
    backgrounds[Math.floor(Math.random() * backgrounds.length)].classList.add('active')
  }

  const openSettings = () => {
    store.setSettingsOpen(true)
  }
</script>

<template>
  <div class="backgrounds">
    <div class="backgrounds-item" v-for="image in unsplashImages" :key="image" :style="{ backgroundImage: `url(https://images.unsplash.com/${image.id}?q=80)` }">
      <div class="backgrounds-item-credit" v-html="image.credit"></div>
    </div>
  </div>
  <button class="logout" @click="handleLogoutClick" v-if="store.isLoggedIn()"><LogOut /></button>
  <button class="settings-button" @click="openSettings"><SettingsIcon /></button>
  <div class="wrapper">
    <div class="left-panel">
      <div class="logo-container">
        <img :src="logoUrl" alt="Erugo" id="logo" :style="{ width: `${logoWidth}px` }" />
      </div>

      <div class="ui-container">
        <template v-if="store.mode === 'upload'">
          <Uploader v-if="store.isLoggedIn()" />
          <Auth v-show="!store.isLoggedIn()" ref="auth" />
        </template>
        <Downloader v-if="store.mode === 'download'" :downloadShareCode="downloadShareCode" />
      </div>
    </div>
  </div>
  <div class="version-info">
    <div class="version-info-text">
      Powered by
      <a href="https://github.com/deanward/erugo">erugo</a>
      {{ version }}
    </div>
  </div>
  <Settings ref="settingsPanel" />
</template>
