<script setup>
import { ref, onMounted, watch, defineExpose } from 'vue'
import { Settings, Tag, Share2, Send, AtSign } from 'lucide-vue-next'
import { ColorPicker } from 'vue-color-kit'
import {
  getSettingsByGroup,
  saveSettingsById,
  saveLogo,
  getSettingById,
  getBackgroundImages,
  saveBackgroundImage,
  deleteBackgroundImage
} from '../../api'

import { useToast } from 'vue-toastification'
import { mapSettings } from '../../utils'
import { domData } from '../../domData'

const toast = useToast()

const settings = ref({
  application_name: '',
  application_url: '',
  login_message: '',
  max_expiry_time: '',
  max_share_size: '',
  max_share_size_unit: '',
  clean_files_after_days: '',
  emails_share_downloaded_enabled: '',
  smtp_host: '',
  smtp_port: '',
  smtp_username: '',
  smtp_password: '',
  smtp_sender_name: '',
  smtp_sender_address: ''
})

const settingsLoaded = ref(false)
const saving = ref(false)

const emit = defineEmits(['navItemClicked'])

onMounted(async () => {
  await loadSettings()
})

const loadSettings = async () => {
  try {
    settings.value = {
      ...mapSettings(await getSettingsByGroup('system.*')),
      ...mapSettings(await getSettingsByGroup('ui.*'))
    }

    settingsLoaded.value = true
  } catch (error) {
    toast.error('Failed to load settings')
    console.error(error)
  }
}

const saveSettings = async () => {
  console.log('saving settings', settings.value)
  saving.value = true
  try {
    await saveSettingsById({
      ...settings.value
    })

    saving.value = false
    toast.success('Settings saved successfully')
    await loadSettings()
  } catch (error) {
    saving.value = false
    toast.error('Failed to save settings')
    console.error(error)
  }
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
    <div class="row mb-5">
      <div class="col-2 d-none d-md-block">
        <ul class="settings-nav pt-5">
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('general')">
              <Settings />
              {{ $t('settings.system.general') }}
            </a>
          </li>
          <li>
            <a href="" @click.prevent="handleNavItemClicked('shares')">
              <Share2 />
              {{ $t('settings.system.shares') }}
            </a>
          </li>
          <li>
            <a href="" @click.prevent="handleNavItemClicked('emails')">
              <Send />
              {{ $t('settings.system.emails') }}
            </a>
          </li>
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('smtp')">
              <AtSign />
              {{ $t('settings.system.smtp') }}
            </a>
          </li>
        </ul>
      </div>
      <div class="col-12 col-md-8 pt-5">
        <div class="row mb-5">
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="general">
              <div class="setting-group-header">
                <h3>
                  <Settings />
                  {{ $t('settings.system.general') }}
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="application_name">{{ $t('settings.system.application_name') }}</label>
                  <input type="text" id="application_name" v-model="settings.application_name" />
                </div>

                <div class="setting-group-body-item">
                  <label for="application_url">{{ $t('settings.system.application_url') }}</label>
                  <input type="text" id="application_url" v-model="settings.application_url" />
                </div>

                <div class="setting-group-body-item">
                  <label for="login_message">{{ $t('settings.system.login_message') }}</label>
                  <input
                    type="text"
                    id="login_message"
                    v-model="settings.login_message"
                    placeholder="Login to your account to upload files."
                  />
                </div>

                <div class="setting-group-body-item">
                  <label for="default_language">{{ $t('settings.system.default_language') }}</label>
                  <select id="default_language" v-model="settings.default_language">
                    <!-- English-->
                    <option value="en">English</option>
                    <!-- German-->
                    <option value="de">Deutsch</option>
                    <!-- French-->
                    <option value="fr">Français</option>
                    <!-- Italian-->
                    <option value="it">Italiano</option>
                    <!-- Dutch-->
                    <option value="nl">Nederlands</option>
                  </select>
                </div>

                <div class="setting-group-body-item mt-3">
                  <div class="checkbox-container">
                    <input type="checkbox" id="show_language_selector" v-model="settings.show_language_selector" />
                    <label for="show_language_selector">{{ $t('settings.system.show_language_selector') }}</label>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h6>{{ $t('settings.system.application_name') }}</h6>
              <p>{{ $t('settings.system.application_name_description') }}</p>

              <h6>{{ $t('settings.system.application_url') }}</h6>
              <p>{{ $t('settings.system.application_url_description') }}</p>

              <h6>{{ $t('settings.system.login_message') }}</h6>
              <p>{{ $t('settings.system.login_message_description') }}</p>

              <h6>{{ $t('settings.system.default_language') }}</h6>
              <p>{{ $t('settings.system.default_language_description') }}</p>

              <h6>{{ $t('settings.system.show_language_selector') }}</h6>
              <p>{{ $t('settings.system.show_language_selector_description') }}</p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="shares">
              <div class="setting-group-header">
                <h3>
                  <Share2 />
                  {{ $t('settings.system.shares') }}
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="max_expiry_time">
                    {{ $t('settings.system.max_expiry_time') }}
                    <small>({{ $t('settings.system.days') }})</small>
                  </label>
                  <input type="number" id="max_expiry_time" v-model="settings.max_expiry_time" placeholder="∞" />
                </div>
                <div class="setting-group-body-item">
                  <div class="row">
                    <div class="col pe-0">
                      <label for="max_share_size">{{ $t('settings.system.max_share_size') }}</label>
                      <input type="number" id="max_share_size" v-model="settings.max_share_size" />
                    </div>
                    <div class="col-auto ps-1">
                      <label for="max_share_size_unit">&nbsp;</label>
                      <select
                        name="max_share_size_unit"
                        id="max_share_size_unit"
                        v-model="settings.max_share_size_unit"
                      >
                        <option value="MB">MB</option>
                        <option value="GB">GB</option>
                      </select>
                    </div>
                  </div>
                </div>
                <div class="setting-group-body-item">
                  <label for="clean_files_after_days">
                    {{ $t('settings.system.clean_files_after') }}
                    <small>({{ $t('settings.system.days') }})</small>
                  </label>
                  <input
                    type="number"
                    id="clean_files_after_days"
                    v-model="settings.clean_files_after_days"
                    placeholder="30"
                  />
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h6>{{ $t('settings.system.max_expiry_time') }}</h6>
              <p>{{ $t('settings.system.max_expiry_time_description') }}</p>
              <h6>{{ $t('settings.system.max_share_size') }}</h6>
              <p>{{ $t('settings.system.max_share_size_description') }}</p>
              <h6>{{ $t('settings.system.clean_files_after') }}</h6>
              <p>{{ $t('settings.system.clean_files_after_description') }}</p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="emails">
              <div class="setting-group-header">
                <h3>
                  <Send />
                  {{ $t('settings.system.emails') }}
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">

                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_downloaded_enabled"
                      v-model="settings.emails_share_downloaded_enabled"
                    />
                    <label for="emails_share_downloaded_enabled">{{ $t('settings.system.enable_share_downloaded_emails') }}</label>
                  </div>
                </div>

                <div class="setting-group-body-item">

                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_expiry_warning_enabled"
                      v-model="settings.emails_share_expiry_warning_enabled"
                    />
                    <label for="emails_share_expiry_warning_enabled">{{ $t('settings.system.enable_share_expiry_warning_emails') }}</label>
                  </div>
                </div>

                <div class="setting-group-body-item">

                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_expired_warning_enabled"
                      v-model="settings.emails_share_expired_warning_enabled"
                    />
                    <label for="emails_share_expired_warning_enabled">{{ $t('settings.system.enable_share_expired_warning_emails') }}</label>
                  </div>
                </div>

                <div class="setting-group-body-item">

                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_deletion_warning_enabled"
                      v-model="settings.emails_share_deletion_warning_enabled"
                    />
                    <label for="emails_share_deletion_warning_enabled">{{ $t('settings.system.enable_share_deletion_warning_emails') }}</label>
                  </div>
                </div>

                <div class="setting-group-body-item">

                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_deleted_enabled"
                      v-model="settings.emails_share_deleted_enabled"
                    />
                    <label for="emails_share_deleted_enabled">{{ $t('settings.system.enable_share_deleted_emails') }}</label>
                  </div>
                </div>

              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h6>{{ $t('settings.system.enable_share_downloaded_emails') }}</h6>
              <p>{{ $t('settings.system.share_downloaded_emails_description') }}</p>

              <h6>{{ $t('settings.system.enable_share_expiry_warning_emails') }}</h6>
              <p>{{ $t('settings.system.share_expiry_warning_emails_description') }}</p>

              <h6>{{ $t('settings.system.enable_share_expired_warning_emails') }}</h6>
              <p>{{ $t('settings.system.share_expired_warning_emails_description') }}</p>

              <h6>{{ $t('settings.system.enable_share_deletion_warning_emails') }}</h6>
              <p>{{ $t('settings.system.share_deletion_warning_emails_description') }}</p>

              <h6>{{ $t('settings.system.enable_share_deleted_emails') }}</h6>
              <p>{{ $t('settings.system.share_deleted_emails_description') }}</p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="smtp">
              <div class="setting-group-header">
                <h3>
                  <AtSign />
                  {{ $t('settings.system.smtp') }}
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="smtp_host">{{ $t('settings.system.smtp_host') }}</label>
                  <input type="text" id="smtp_host" v-model="settings.smtp_host" />
                </div>
                <div class="setting-group-body-item">
                  <label for="smtp_port">{{ $t('settings.system.smtp_port') }}</label>
                  <input type="number" id="smtp_port" v-model="settings.smtp_port" />
                </div>
                <div class="setting-group-body-item">
                  <label for="smtp_username">{{ $t('settings.system.smtp_username') }}</label>
                  <input type="text" id="smtp_username" v-model="settings.smtp_username" />
                </div>
                <div class="setting-group-body-item">
                  <label for="smtp_password">{{ $t('settings.system.smtp_password') }}</label>
                  <input type="password" id="smtp_password" v-model="settings.smtp_password" />
                </div>
                <div class="setting-group-body-item">
                  <label for="smtp_sender_name">{{ $t('settings.system.smtp_sender_name') }}</label>
                  <input type="text" id="smtp_sender_name" v-model="settings.smtp_sender_name" />
                </div>
                <div class="setting-group-body-item">
                  <label for="smtp_sender_address">{{ $t('settings.system.smtp_sender_address') }}</label>
                  <input type="text" id="smtp_sender_address" v-model="settings.smtp_sender_address" />
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h6>{{ $t('settings.system.smtp_host') }}</h6>
              <p>{{ $t('settings.system.smtp_host_description') }}</p>

              <h6>{{ $t('settings.system.smtp_port') }}</h6>
              <p>{{ $t('settings.system.smtp_port_description') }}</p>

              <h6>{{ $t('settings.system.smtp_username') }}</h6>
              <p>{{ $t('settings.system.smtp_username_description') }}</p>

              <h6>{{ $t('settings.system.smtp_password') }}</h6>
              <p>{{ $t('settings.system.smtp_password_description') }}</p>

              <h6>{{ $t('settings.system.smtp_sender_name') }}</h6>
              <p>{{ $t('settings.system.smtp_sender_name_description') }}</p>

              <h6>{{ $t('settings.system.smtp_sender_address') }}</h6>
              <p>{{ $t('settings.system.smtp_sender_address_description') }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
