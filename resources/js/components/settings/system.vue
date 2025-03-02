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
              General
            </a>
          </li>
          <li>
            <a href="" @click.prevent="handleNavItemClicked('shares')">
              <Share2 />
              Shares
            </a>
          </li>
          <li>
            <a href="" @click.prevent="handleNavItemClicked('emails')">
              <Send />
              Emails
            </a>
          </li>
          <li>
            <a href="#" @click.prevent="handleNavItemClicked('smtp')">
              <AtSign />
              SMTP
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
                  General settings
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="application_name">Application Name</label>
                  <input type="text" id="application_name" v-model="settings.application_name" />
                </div>

                <div class="setting-group-body-item mt-3">
                  <label for="application_url">Application URL</label>
                  <input type="text" id="application_url" v-model="settings.application_url" />
                </div>

                <div class="setting-group-body-item mt-3">
                  <label for="login_message">Login Message</label>
                  <input
                    type="text"
                    id="login_message"
                    v-model="settings.login_message"
                    placeholder="Login to your account to upload files."
                  />
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h6>Application name</h6>
              <p>Customise the displayed name of the application in places like title bars and dialogues.</p>

              <h6>Application URL</h6>
              <p>Used to allow emails to link back to the application.</p>

              <h6>Login Message</h6>
              <p>Customise the message displayed on the login screen.</p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="shares">
              <div class="setting-group-header">
                <h3>
                  <Share2 />
                  Share settings
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="max_expiry_time">
                    Max expiry time
                    <small>(days)</small>
                  </label>
                  <input type="number" id="max_expiry_time" v-model="settings.max_expiry_time" placeholder="∞" />
                </div>
                <div class="setting-group-body-item">
                  <div class="row">
                    <div class="col pe-0">
                      <label for="max_share_size">Max share size</label>
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
                <div class="setting-group-body-item mt-3">
                  <label for="clean_files_after_days">
                    Clean files after
                    <small>(days)</small>
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
              <h6>Max expiry time</h6>
              <p>
                Set the maximum nunber of days a user can select for share expiry upon creation. Empty the value to
                allow unlimited expiry time.
              </p>
              <h6>Max share size</h6>
              <p>
                Set the maximum size of a share that a user can create. The value is in megabytes or gigabytes. The
                largest value you can enter here is xx and it limited by PHP settings.
                <a href="https://erugo.app/docs/max-share-size">Learn more.</a>
              </p>
              <h6>Clean files after</h6>
              <p>How long after a share expires before files are deleted from the system.</p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="emails">
              <div class="setting-group-header">
                <h3>
                  <Send />
                  Emails
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="emails_share_downloaded_enabled">Share downloaded emails</label>
                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_downloaded_enabled"
                      v-model="settings.emails_share_downloaded_enabled"
                    />
                    <label for="emails_share_downloaded_enabled">Enable share downloaded emails</label>
                    <p>Sent the first time a share is downloaded.</p>
                  </div>
                </div>

                <div class="setting-group-body-item">
                  <label for="emails_share_expiry_warning_enabled">Share expiry warning emails</label>
                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_expiry_warning_enabled"
                      v-model="settings.emails_share_expiry_warning_enabled"
                    />
                    <label for="emails_share_expiry_warning_enabled">Enable share expiry warning emails</label>
                    <p>Sent when a share is about to expire.</p>
                  </div>
                </div>

                <div class="setting-group-body-item">
                  <label for="emails_share_expired_warning_enabled">Share expired warning emails</label>
                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_expired_warning_enabled"
                      v-model="settings.emails_share_expired_warning_enabled"
                    />
                    <label for="emails_share_expired_warning_enabled">Enable share expired warning emails</label>
                    <p>Sent when a share has expired.</p>
                  </div>
                </div>

                <div class="setting-group-body-item">
                  <label for="emails_share_deletion_warning_enabled">Share deletion warning emails</label>
                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_deletion_warning_enabled"
                      v-model="settings.emails_share_deletion_warning_enabled"
                    />
                    <label for="emails_share_deletion_warning_enabled">Enable share deletion warning emails</label>
                    <p>Sent when a share is about to be deleted.</p>
                  </div>
                </div>

                <div class="setting-group-body-item">
                  <label for="emails_share_deleted_enabled">Share deleted emails</label>
                  <div class="checkbox-container">
                    <input
                      type="checkbox"
                      id="emails_share_deleted_enabled"
                      v-model="settings.emails_share_deleted_enabled"
                    />
                    <label for="emails_share_deleted_enabled">Enable share deleted emails</label>
                    <p>Sent when a share is deleted.</p>
                  </div>
                </div>

              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help">
              <h6>Share notifications</h6>
              <p>Decide whether to send emails to share owners when a share is downloaded, expired or deleted.</p>
            </div>
          </div>
        </div>

        <div class="row mb-5">
          <div class="col-12 col-md-6 pe-0 ps-0 ps-md-3">
            <div class="setting-group" id="smtp">
              <div class="setting-group-header">
                <h3>
                  <AtSign />
                  SMTP settings
                </h3>
              </div>

              <div class="setting-group-body">
                <div class="setting-group-body-item">
                  <label for="smtp_host">Host</label>
                  <input type="text" id="smtp_host" v-model="settings.smtp_host" />
                </div>
                <div class="setting-group-body-item mt-3">
                  <label for="smtp_port">Port</label>
                  <input type="number" id="smtp_port" v-model="settings.smtp_port" />
                </div>
                <div class="setting-group-body-item mt-3">
                  <label for="smtp_username">Username</label>
                  <input type="text" id="smtp_username" v-model="settings.smtp_username" />
                </div>
                <div class="setting-group-body-item mt-3">
                  <label for="smtp_password">Password</label>
                  <input type="password" id="smtp_password" v-model="settings.smtp_password" />
                </div>
                <div class="setting-group-body-item mt-3">
                  <label for="smtp_sender_name">Sender Name</label>
                  <input type="text" id="smtp_sender_name" v-model="settings.smtp_sender_name" />
                </div>
                <div class="setting-group-body-item mt-3">
                  <label for="smtp_sender_address">Sender Address</label>
                  <input type="text" id="smtp_sender_address" v-model="settings.smtp_sender_address" />
                </div>
              </div>
            </div>
          </div>
          <div class="d-none d-md-block col ps-0">
            <div class="section-help"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
