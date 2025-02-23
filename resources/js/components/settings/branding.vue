<script setup>
import { ref, onMounted, watch, defineExpose } from 'vue'
import { Pipette, Image, Ruler, Tag, X, Dice5, Images } from 'lucide-vue-next'
import { ColorPicker } from 'vue-color-kit'
import {
  getSettingsByGroup,
  saveSettingsById,
  saveLogo,
  getBackgroundImages,
  saveBackgroundImage,
  deleteBackgroundImage
} from '../../api'
import FileInput from '../fileInput.vue'
import { useToast } from 'vue-toastification'
import { niceFileName, mapSettings } from '../../utils'

const toast = useToast()

const settings = ref({
  logo: null,
  logo_width: '',
  css_primary_color: '',
  css_secondary_color: '',
  css_accent_color: '',
  css_accent_color_light: '',
  use_my_backgrounds: false,
  show_powered_by: true,
  application_name: ''
})

const newBackgroundImage = ref(null)
const backgroundImages = ref([])

const settingsLoaded = ref(false)
const saving = ref(false)

const emit = defineEmits(['navItemClicked'])

onMounted(async () => {
  await loadSettings()
})

const loadSettings = async () => {
  try {
    settings.value = {
      ...mapSettings(await getSettingsByGroup('ui.*'))
    }

    settingsLoaded.value = true
  } catch (error) {
    toast.error('Failed to load settings')
    console.error(error)
  }

  loadBackgroundImages()
}



const loadBackgroundImages = async () => {
  getBackgroundImages().then((data) => {
    backgroundImages.value = data.files
  })
}

const saveSettings = async () => {
  console.log('saving settings')
  saving.value = true
  try {

    if (settings.value.logo instanceof File) {
      saveLogo(settings.value.logo)
    }

    await saveSettingsById(settings.value)
    applySettingsWithoutRefresh()

    saving.value = false
    toast.success('Settings saved successfully')
  } catch (error) {
    saving.value = false
    toast.error('Failed to save settings')
    console.error(error)
  }
}

const applySettingsWithoutRefresh = () => {
  //find the style tag #erugo-css-variables
  const styleTag = document.getElementById('erugo-css-variables')
  if (styleTag) {
    //update the css variables
    styleTag.innerHTML = `
      :root {
        --primary-color: ${settings.value.css_primary_color};
        --secondary-color: ${settings.value.css_secondary_color};
        --accent-color: ${settings.value.css_accent_color};
        --accent-color-light: ${settings.value.css_accent_color_light};
      }
      `
  } else {
    //add the style tag
    const styleTag = document.createElement('style')
    styleTag.id = 'erugo-css-variables'
    styleTag.innerHTML = `
      :root {
        --primary-color: ${settings.value.css_primary_color};
        --secondary-color: ${settings.value.css_secondary_color};
        --accent-color: ${settings.value.css_accent_color};
        --accent-color-light: ${settings.value.css_accent_color_light};
      }
      `
    document.head.appendChild(styleTag)
  }

  //update the logo width
  const logo = document.getElementById('logo')
  if (logo) {
    logo.style.width = `${settings.value.logo_width}`
  }
}

//watch newBackgroundImage and upload it to the server
watch(newBackgroundImage, async () => {
  if (newBackgroundImage.value) {
    saveBackgroundImage(newBackgroundImage.value)
      .then((data) => {
        loadBackgroundImages()
        newBackgroundImage.value = null
        toast.success('Background image uploaded successfully')
      })
      .catch((error) => {
        toast.error('Failed to upload background image')
      })
  }
})

watch(backgroundImages, () => {
  if (backgroundImages.value.length === 0) {
    useMyBackgrounds.value = false
  }
})

const handleDeleteBackgroundImage = (file) => {
  const reallyDelete = confirm('Are you sure you want to delete this background image?')
  if (!reallyDelete) {
    return
  }
  deleteBackgroundImage(file)
    .then((data) => {
      loadBackgroundImages()
      toast.success('Background image deleted successfully')
    })
    .catch((error) => {
      toast.error('Failed to delete background image')
    })
}

const handleNavItemClicked = (item) => {
  emit('navItemClicked', item)
}

//define exposed methods
defineExpose({
  saveSettings
})
</script>
<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-2 d-none d-md-block">
        <ul class="settings-nav">
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('background-images')">
              <Images />
              Background Images
            </a>
          </li>
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('logo-settings')">
              <Image />
              Logo
            </a>
          </li>
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('application-name')">
              <Tag />
              Application Name
            </a>
          </li>
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('other-ui-settings')">
              <Dice5 />
              Other UI Settings
            </a>
          </li>
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('ui-colours')">
              <Pipette />
              UI Colours
            </a>
          </li>
        </ul>
      </div>
      <div class="col-12 col-md-8">
        <div class="row mb-5">
          <!-- backgrounds -->
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="background-images">
              <div class="setting-group-header">
                <h3>
                  <Images />
                  Background Images
                </h3>
                <div class="settings-group-info">
                  <p>Manage your background images.</p>
                </div>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <div class="background-images">
                    <div class="background-image" v-for="image in backgroundImages" :key="image">
                      <img :src="`/api/backgrounds/${image}/thumb`" />
                      <div class="name">
                        {{ niceFileName(image) }}
                      </div>
                      <button class="delete" @click="handleDeleteBackgroundImage(image)">
                        <X />
                      </button>
                    </div>
                  </div>

                  <FileInput
                    v-model="newBackgroundImage"
                    accept="image/png, image/jpeg, image/webp"
                    label="Upload Background Image"
                    class="mt-3"
                  />

                  <div class="checkbox-container" :class="{ disabled: backgroundImages.length === 0 }">
                    <input type="checkbox" v-model="settings.use_my_backgrounds" id="useMyBackgrounds" />
                    <label for="useMyBackgrounds">Use my backgrounds</label>
                    <p class="help-text">
                      Use the backgrounds you have uploaded.
                      <br />
                      If not checked, Unsplash backgrounds will be used.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h5>Background images</h5>
              <p>
                Upload custom background images to replace Unsplash defaults. Images rotate randomly every 3 minutes to
                create dynamic branding.
              </p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <!-- logo -->
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="logo-settings">
              <div class="setting-group-header">
                <h3>
                  <Image />
                  Logo
                </h3>
                <div class="settings-group-info">
                  <p>Modify the logo that is displayed throughout erugo.</p>
                </div>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="logoFile">Logo Image</label>
                  <FileInput v-model="settings.logo" accept="image/png" />
                </div>

                <div class="setting-group-body-item">
                  <label for="logoWidth">
                    Logo Width
                    <small>(in pixels)</small>
                  </label>
                  <input type="number" v-model="settings.logo_width" />
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h5>Logo</h5>
              <p>
                Display your company logo in PNG format in the top left corner. Set your desired width in pixels -
                height will adjust automatically to maintain aspect ratio. Use a source image at least as wide as your
                display width to ensure quality.
              </p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <!-- Application name -->
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="application-name">
              <div class="setting-group-header">
                <h3>
                  <Tag />
                  Application Name
                </h3>
                <div class="settings-group-info">
                  <p>Change the name of your erugo instance.</p>
                </div>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <input type="text" v-model="settings.application_name" />
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h5>Application name</h5>
              <p>Customise the displayed name of the application in places like title bars and dialogues.</p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <!-- Other UI settings -->
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="other-ui-settings">
              <div class="setting-group-header">
                <h3>
                  <Dice5 />
                  Other UI Settings
                </h3>
                <div class="settings-group-info">
                  <p>Miscellaneous settings for the UI.</p>
                </div>
              </div>

              <div class="checkbox-container">
                <input type="checkbox" v-model="settings.show_powered_by" id="showPoweredBy" />
                <label for="showPoweredBy">Show Powered By erugo</label>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h5>Show or hide powered by text</h5>
              <p>
                Show or hide the "Powered by erugo" text on the bottom of the page. This text allows erugo to be
                discovered by other users and helps us to grow, however, you can hide it if you wish.
              </p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <!-- UI Colours -->
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="ui-colours">
              <div class="setting-group-header">
                <h3>
                  <Pipette />
                  UI Colours
                </h3>
                <div class="settings-group-info">
                  <p>Customize the UI colours to match your brand.</p>
                </div>
              </div>

              <div class="setting-group-body" v-if="settingsLoaded">
                <div class="row">
                  <div class="col-auto">
                    <h6>Primary Colour</h6>
                    <ColorPicker
                      theme="light"
                      :color="settings.css_primary_color"
                      :sucker-hide="false"
                      @changeColor="settings.css_primary_color = $event.hex"
                    />
                  </div>
                  <div class="col-auto">
                    <h6>Secondary Colour</h6>
                    <ColorPicker
                      theme="light"
                      :color="settings.css_secondary_color"
                      :sucker-hide="false"
                      @changeColor="settings.css_secondary_color = $event.hex"
                    />
                  </div>
                </div>
                <div class="row mt-4">
                  <div class="col-auto">
                    <h6>Accent Colour Light</h6>
                    <ColorPicker
                      theme="light"
                      :color="settings.css_accent_color_light"
                      :sucker-hide="false"
                      @changeColor="settings.css_accent_color_light = $event.hex"
                    />
                  </div>
                  <div class="col-auto">
                    <h6>Accent Colour</h6>
                    <ColorPicker
                      theme="light"
                      :color="settings.css_accent_color"
                      :sucker-hide="false"
                      @changeColor="settings.css_accent_color = $event.hex"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h5>UI colours</h5>
              <p>
                Customize the colours of the UI to match your brand. The primary colour is used for buttons, links, and
                other primary elements. The secondary colour is used for secondary elements and text. The accent colour
                is used for accents and highlights. The easiest way to figure out what each colour does it to experiment
                with them.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
