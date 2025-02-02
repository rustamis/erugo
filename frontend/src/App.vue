<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { GithubIcon, LogOut, Settings as SettingsIcon } from 'lucide-vue-next'
  import Uploader from './components/uploader.vue'
  import Downloader from './components/downloader.vue'
  import Auth from './components/auth.vue'
  import Settings  from './components/settings.vue' 
  import { unsplashImages } from './unsplashImages'
  import { getApiUrl } from './utils'
  import { store } from './store'

  const auth = ref(null)
  const downloadShareCode = ref('')

  onMounted(() => {
    setMode()
    changeBackground()
    setTimeout(changeBackground, 180000) //change every 3 minutes
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

  const setPageTitle = (title) => {
    document.title = `erugo shares - ${title}`
  }

  const logout = () => {
    auth.value.logout()
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
  <button class="logout" @click="logout" v-if="store.isLoggedIn()"><LogOut /></button>
  <button class="settings-button" @click="openSettings" v-if="store.isAdmin()"><SettingsIcon /></button>
  <div class="wrapper">
    <div class="left-panel">
      <div class="logo-container">
        <img src="/erugo.png" alt="Erugo" id="logo" />
      </div>

      <div class="ui-container">
        <template v-if="store.mode === 'upload'">
          <Uploader v-if="store.isLoggedIn()"/>
          <Auth v-show="!store.isLoggedIn()" ref="auth"/>
        </template>
        <Downloader v-if="store.mode === 'download'" :downloadShareCode="downloadShareCode" />
      </div>
    </div>
  </div>
  <div class="version-info">
    <div class="version-info-text">erugo v0.0.1</div>
    <div class="github-link">
      <a href="https://github.com/deanward/erugo">
        <GithubIcon />
        github.com/deanward/erugo
      </a>
    </div>
  </div>
  <Settings />
</template>