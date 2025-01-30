<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { CircleSlash2, FilePlus, FolderPlus, GithubIcon, Upload, Trash, Copy, X, Loader, LogOut } from 'lucide-vue-next'
  import Uploader from './components/uploader.vue'
  import Downloader from './components/downloader.vue'
  import Auth from './components/auth.vue'
  import { unsplashImages } from './unsplashImages'
  import { getApiUrl } from './utils'

  const auth = ref(null)

  const apiUrl = getApiUrl()

  const mode = ref('upload')
  const downloadShareCode = ref('')
  const loggedIn = ref(false)
  const jwt = ref('')

  onMounted(() => {
    setMode()
    changeBackground()
    setTimeout(changeBackground, 180000) //change every 3 minutes
  })

  const setMode = () => {
    if (window.location.pathname.includes('shares')) {
      mode.value = 'download'
      downloadShareCode.value = window.location.pathname.split('/').pop()
      setPageTitle('Download Share')
    } else {
      mode.value = 'upload'
      setPageTitle('Create Share')
    }
  }

  const authSuccess = (data) => {
    loggedIn.value = true
    jwt.value = data.data.token
    console.log(jwt.value)
    console.log(data)
  }

  const setPageTitle = (title) => {
    document.title = `erugo shares - ${title}`
  }

  const logout = () => {
    auth.value.logout()
  }

  const loggedout = () => {
    loggedIn.value = false
    jwt.value = ''
  }

  const changeBackground = async () => {
    let backgrounds = document.querySelectorAll('.backgrounds-item')
    backgrounds.forEach(background => {
      background.classList.remove('active')
    })
    backgrounds[Math.floor(Math.random() * backgrounds.length)].classList.add('active')
  }
</script>

<template>
  <div class="backgrounds">
    <div class="backgrounds-item" v-for="image in unsplashImages" :key="image" :style="{ backgroundImage: `url(https://images.unsplash.com/${image.id}?q=80)` }">
      <div class="backgrounds-item-credit" v-html="image.credit"></div>
    </div>
  </div>
  <button class="logout" @click="logout"><LogOut /></button>
  <div class="wrapper">
    <div class="left-panel">
      <div class="logo-container">
        <img src="/erugo.png" alt="Erugo" id="logo" />
      </div>

      <div class="ui-container">
        <template v-if="mode === 'upload'">
          <Uploader v-if="loggedIn" :jwt="jwt" />
          <Auth v-show="!loggedIn" @auth-success="authSuccess" @logout="loggedout" ref="auth"/>
        </template>
        <Downloader v-if="mode === 'download'" :downloadShareCode="downloadShareCode" />
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
</template>

<style scoped>
  .logo {
    height: 6em;
    padding: 1.5em;
    will-change: filter;
    transition: filter 300ms;
  }
  .logo:hover {
    filter: drop-shadow(0 0 2em #646cffaa);
  }
  .logo.vue:hover {
    filter: drop-shadow(0 0 2em #42b883aa);
  }
</style>
