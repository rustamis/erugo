<script setup>
  import { Pipette, Image, MessageCircleQuestion, Ruler, Tag } from 'lucide-vue-next'
  import { ref, onMounted, defineExpose } from 'vue'
  import { ColorPicker } from 'vue-color-kit'
  import { getSettingsByGroup, saveSettingsById, saveLogo, getSettingById } from '../../api'
  import FileInput from '../fileInput.vue'
  import { useToast } from 'vue-toastification'

  const toast = useToast()

  //These are the variables that will be updated and sent to the backend
  const primaryColor = ref('#000000')
  const secondaryColor = ref('#000000')
  const accentColor = ref('#000000')
  const accentColorLight = ref('#000000')
  const logoFile = ref(null)
  const logoWidth = ref(100)

  const applicationName = ref('')

  //These are the variables that will be passed to the colour pickers
  //we cannot use the variables above for fear of passing back mutated values to the pickers
  const initalPrimaryColor = ref('#000000')
  const initalSecondaryColor = ref('#000000')
  const initalAccentColor = ref('#000000')
  const initalAccentColorLight = ref('#000000')

  const settingsLoaded = ref(false)
  const saving = ref(false)
  const logo = ref(null)
  const errors = ref({
    primaryColor: null,
    secondaryColor: null,
    accentColor: null,
    accentColorLight: null
  })

  function rgbToHex(rgb) {
    //if this is already a hex string, return it
    if (rgb.startsWith('#')) {
      return rgb
    }
    // Extract numbers from rgb string
    const [r, g, b] = rgb.match(/\d+/g).map(Number)
    // Convert to hex
    return (
      '#' +
      [r, g, b]
        .map(x => {
          const hex = x.toString(16)
          return hex.length === 1 ? '0' + hex : hex
        })
        .join('')
    )
  }

  onMounted(async () => {
    getSettingsByGroup('ui.css').then(data => {
      mapSettingsToVariables(data.settings)
    })
    getSettingById('application_name').then(data => {
      applicationName.value = data.value
    })
    getSettingById('logo_width').then(data => {
      logoWidth.value = data.value
    })
  })

  const mapSettingsToVariables = settings => {
    settings.forEach(setting => {
      switch (setting.id) {
        case 'css_primary_color':
          initalPrimaryColor.value = rgbToHex(setting.value)
          primaryColor.value = rgbToHex(setting.value)
          break
        case 'css_secondary_color':
          initalSecondaryColor.value = rgbToHex(setting.value)
          secondaryColor.value = rgbToHex(setting.value)
          break
        case 'css_accent_color':
          initalAccentColor.value = rgbToHex(setting.value)
          accentColor.value = rgbToHex(setting.value)
          break
        case 'css_accent_color_light':
          initalAccentColorLight.value = rgbToHex(setting.value)
          accentColorLight.value = rgbToHex(setting.value)
          break
      }
    })
    settingsLoaded.value = true
  }

  const saveSettings = async () => {
    saving.value = true
    try {
      await saveSettingsById([
        {
          id: 'css_primary_color',
          value: primaryColor.value
        },
        {
          id: 'css_secondary_color',
          value: secondaryColor.value
        },
        {
          id: 'css_accent_color',
          value: accentColor.value
        },
        {
          id: 'css_accent_color_light',
          value: accentColorLight.value
        },
        {
          id: 'application_name',
          value: applicationName.value
        },
        {
          id: 'logo_width',
          value: logoWidth.value + ''
        }
      ])
      applySettingsWithoutRefresh()

      //if the user has changed the logo, save the new logo
      if (logoFile.value) {
        saveLogo(logoFile.value)
      }

      saving.value = false
      toast.success('Settings saved successfully')
    } catch (error) {
      saving.value = false
      toast.error('Failed to save settings')
    }
  }

  const applySettingsWithoutRefresh = () => {
    //find the style tag #erugo-css-variables
    const styleTag = document.getElementById('erugo-css-variables')
    if (styleTag) {
      //update the css variables
      styleTag.innerHTML = `
      :root {
        --primary-color: ${primaryColor.value};
        --secondary-color: ${secondaryColor.value};
        --accent-color: ${accentColor.value};
        --accent-color-light: ${accentColorLight.value};
      }
      `
    } else {
      //add the style tag
      const styleTag = document.createElement('style')
      styleTag.id = 'erugo-css-variables'
      styleTag.innerHTML = `
      :root {
        --primary-color: ${primaryColor.value};
        --secondary-color: ${secondaryColor.value};
        --accent-color: ${accentColor.value};
        --accent-color-light: ${accentColorLight.value};
      }
      `
      document.head.appendChild(styleTag)
    }

    //update the logo width
    const logo = document.getElementById('logo')
    if (logo) {
      logo.style.width = `${logoWidth.value}px`
    }
  }

  //define exposed methods
  defineExpose({
    saveSettings
  })
</script>
<template>
  <div class="setting-groups">
    <div class="setting-group">
      <div class="setting-group-header">
        <h3>
          <Pipette />
          UI Colours
        </h3>
      </div>
      <div class="settings-group-info">
        <MessageCircleQuestion />
        <p>Customize the UI colours to match your brand.</p>
      </div>

      <div class="setting-group-body mt-3" v-if="settingsLoaded">
        <div>
          <h6>Primary Colour</h6>
          <ColorPicker theme="light" :color="initalPrimaryColor" :sucker-hide="false" @changeColor="primaryColor = $event.hex" />
        </div>
        <div>
          <h6>Secondary Colour</h6>
          <ColorPicker theme="light" :color="initalSecondaryColor" :sucker-hide="false" @changeColor="secondaryColor = $event.hex" />
        </div>
        <div>
          <h6>Accent Colour</h6>
          <ColorPicker theme="light" :color="initalAccentColor" :sucker-hide="false" @changeColor="accentColor = $event.hex" />
        </div>
        <div>
          <h6>Accent Colour Light</h6>
          <ColorPicker theme="light" :color="initalAccentColorLight" :sucker-hide="false" @changeColor="accentColorLight = $event.hex" />
        </div>
      </div>
    </div>
    <div class="setting-groups">
      <div class="setting-group">
        <div class="setting-group-header">
          <h3>
            <Image />
            Logo Image
          </h3>
        </div>
        <div class="settings-group-info">
          <MessageCircleQuestion />
          <p>Display your logo throughout erugo.</p>
        </div>
        <div class="setting-group-body mt-3">
          <div class="setting-group-body-item">
            <h6>
              Select Logo
              <small>(Supports png only)</small>
            </h6>
            <FileInput v-model="logoFile" accept="image/png" />
          </div>
        </div>
      </div>
      <div class="setting-group">
        <div class="setting-group-header">
          <h3>
            <Ruler />
            Logo Width
          </h3>
        </div>
        <div class="settings-group-info">
          <MessageCircleQuestion />
          <p>Change the width of your logo.</p>
        </div>
        <div class="setting-group-body mt-3">
          <div class="setting-group-body-item">
            <h6>
              Logo Width
              <small>(px)</small>
            </h6>
            <input type="number" v-model="logoWidth" />
          </div>
        </div>
      </div>
      <div class="setting-group">
        <div class="setting-group-header">
          <h3>
            <Tag />
            Application Name
          </h3>
        </div>
        <div class="settings-group-info">
          <MessageCircleQuestion />
          <p>Change the name of your erugo instance.</p>
        </div>
        <div class="setting-group-body mt-3">
          <div class="setting-group-body-item">
            <h6>Application Name</h6>
            <input type="text" v-model="applicationName" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
