<script setup>
import { store } from '../store'
import { CircleX, Settings, Users as UsersIcon, UserPlus, Save, Palette, User, Boxes, Plus } from 'lucide-vue-next'
import { ref, onMounted, computed, watch } from 'vue'
import Users from './settings/users.vue'
import BrandingSettings from './settings/branding.vue'
import SystemSettings from './settings/system.vue'
import MyProfile from './settings/myProfile.vue'
import MyShares from './settings/myShares.vue'

//settings panels
const usersPanel = ref(null)
const brandingSettings = ref(null)
const systemSettings = ref(null)
// Create refs for the tab contents
const tabContents = ref({
  branding: ref(null),
  system: ref(null),
  users: ref(null),
  myProfile: ref(null),
  myShares: ref(null)
})

onMounted(() => {
  activeTab.value = getInitialTab()
})

const closeSettings = () => {
  store.setSettingsOpen(false)
}

const clickOutside = (e) => {
  if (e.target === e.currentTarget) {
    closeSettings()
  }
}

const setActiveTab = (tab) => {
  activeTab.value = tab
}

const getInitialTab = () => {
  return 'myShares'
}

// Track active tab
const activeTab = ref(null)

defineExpose({
  setActiveTab
})

const createShare = () => {
  store.setSettingsOpen(false)
}

const handleNavItemClicked = (item) => {
  const scrollableElement = document.querySelector('.tab-content-body')
  const element = document.getElementById(item)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
    scrollableElement.scrollTo({
      top: element.offsetTop - 100,
      behavior: 'smooth'
    })
  }
}

const getSettingsTitle = () => {
  switch (activeTab.value) {
    case 'branding':
      return 'Branding &amp; UI'
    case 'system':
      return 'System Settings'
    case 'users':
      return 'Users'
    case 'myProfile':
      return 'My Profile'
    case 'myShares':
      return 'My Shares'
    default:
      return 'erugo Settings' 
  }
}
</script>

<template>
  <div class="settings-overlay" :class="{ active: store.settingsOpen }" @click="clickOutside">
    <div class="settings-container">
      <div class="settings-header">
        <h1>
          <Settings />
           <span>Manage <span v-html="getSettingsTitle()" /></span>
        </h1>
        <button class="close-settings-button icon-only" @click="closeSettings">
          <CircleX />
        </button>
      </div>
      <div class="settings-tabs-wrapper">
        <div class="settings-tabs-container">
          <div class="settings-tab" :class="{ active: activeTab === 'branding' }" @click="setActiveTab('branding')" v-if="store.isAdmin()">
            <h2>
              <Palette />
              Branding &amp; UI
            </h2>
          </div>
          <div class="settings-tab" :class="{ active: activeTab === 'system' }" @click="setActiveTab('system')" v-if="store.isAdmin()">
            <h2>
              <Settings />
              System Settings
            </h2>
          </div>
          <div class="settings-tab" :class="{ active: activeTab === 'users' }" @click="setActiveTab('users')" v-if="store.isAdmin()">
            <h2>
              <UsersIcon />
              Users
            </h2>
          </div>
          <div class="settings-tab" :class="{ active: activeTab === 'myShares' }" @click="setActiveTab('myShares')">
            <h2>
              <Boxes />
              My Shares
            </h2>
          </div>
          <div class="settings-tab" :class="{ active: activeTab === 'myProfile' }" @click="setActiveTab('myProfile')">
            <h2>
              <User />
              My Profile
            </h2>
          </div>
        </div>
        <div class="settings-tabs-content-container">
          <Transition name="fade">
            <div v-if="activeTab === 'branding'" class="settings-tab-content" ref="tabContents.branding" key="branding">
              <div class="tab-content-header">
                <h2 class="d-none d-md-flex">
                  <Palette />
                  <span>
                    Branding &amp; UI
                    <small>Customise the UI of your erugo instance and add your own branding.</small>
                  </span>
                </h2>
                <div class="user-actions">
                  <button @click="$refs['brandingSettings'].saveSettings()">
                    <Save />
                    Save Branding Settings
                  </button>
                </div>
              </div>
              <div class="tab-content-body">
                <BrandingSettings ref="brandingSettings" v-if="store.settingsOpen" @navItemClicked="handleNavItemClicked" />
              </div>
            </div>
            <div v-else-if="activeTab === 'system'" class="settings-tab-content" ref="tabContents.system" key="system">
              <div class="tab-content-header">
                <h2 class="d-none d-md-flex">
                  <Settings />
                  <span>
                    System Settings
                    <small>Manage the behaviour of your erugo instance.</small>
                  </span>
                </h2>
                <div class="user-actions">
                  <button @click="$refs['systemSettings'].saveSettings()">
                    <Save />
                    Save System Settings
                  </button>
                </div>
              </div>
              <div class="tab-content-body">
                <SystemSettings ref="systemSettings" v-if="store.settingsOpen" @navItemClicked="handleNavItemClicked" />
              </div>
            </div>
            <div v-else-if="activeTab === 'users'" class="settings-tab-content" ref="tabContents.users" key="users">
              <div class="tab-content-header">
                <h2 class="d-none d-md-flex">
                  <UsersIcon />
                  <span>
                    Users
                    <small>Manage your users.</small>
                  </span>
                </h2>
                <div class="user-actions">
                  <button @click="usersPanel.addUser">
                    <UserPlus />
                    Add User
                  </button>
                </div>
              </div>
              <div class="tab-content-body">
                <Users ref="usersPanel" v-if="store.settingsOpen" />
              </div>
            </div>
            <div v-else-if="activeTab === 'myProfile'" class="settings-tab-content" ref="tabContents.myProfile" key="myProfile">
              <div class="tab-content-header">
                <h2 class="d-none d-md-flex">
                  <User />
                  <span>
                    My Profile
                    <small>Manage your profile.</small>
                  </span>
                </h2>
              </div>
              <div class="tab-content-body">
                <MyProfile ref="myProfilePanel" v-if="store.settingsOpen" />
              </div>
            </div>
            <div v-else-if="activeTab === 'myShares'" class="settings-tab-content" ref="tabContents.myShares" key="myShares">
              <div class="tab-content-header">
                <h2 class="d-none d-md-flex">
                  <Boxes />
                  <span>
                    My Shares
                    <small>Manage your shares.</small>
                  </span>
                </h2>
                <div class="user-actions">
                  <button @click="createShare">
                    <Plus />
                    Create Share
                  </button>
                </div>
              </div>
              <div class="tab-content-body">
                <MyShares ref="mySharesPanel" v-if="store.settingsOpen" />
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  background-color: transparent;
  width: 100%;
  height: 100%;
  z-index: 210;
  pointer-events: none;
  transition: all 300ms ease-in-out;
  transition-delay: 300ms;

  .settings-container {

    position: absolute;
    bottom: 0;
    left: 0;
    transform: translateX(calc(50vw - var(--settings-width) / 2)) translateY(100%);

    width: var(--settings-width);
    height: var(--settings-height);

    transition: all 300ms ease-in-out;
    transition-delay: 0s;

    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: flex-start;

  }

  &.active {
    background-color: var(--overlay-background-color);
    pointer-events: auto;
    transition-delay: 0s;
    backdrop-filter: blur(10px);

    .settings-container {
      transform: translateX(calc(50vw - var(--settings-width) / 2)) translateY(0);
      transition-delay: 100ms;
    }
  }
}

.settings-header {
  background-color: var(--panel-header-background-color);
  border-radius: 5px 5px 0 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 80px;
  width: 100%;
  h1 {
    font-size: 20px;
    font-weight: 600;
    color: var(--panel-header-text-color);
    padding-left: 20px;
    display: flex;
    align-items: center;
    gap: 10px;
    svg {
      width: 20px;
      height: 20px;
    }
  }
}

.settings-tabs-wrapper {
  width: 100%;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
}

.settings-tabs-container {
  display: flex;
  gap: 5px;
  padding-left: 20px;
  padding-right: 20px;
  background-color: var(--tabs-bar-background-color);
  width: 100%;
  .settings-tab {
    background-color: var(--tabs-tab-background-color);
    margin-top: 10px;
    padding: 10px;
    border-radius: var(--tabs-border-radius);
    box-shadow: inset 0 -5px 2px 0 rgba(0, 0, 0, 0.05);
    cursor: pointer;
    transition: all 300ms ease-in-out;

    h2 {
      font-size: 16px;
      font-weight: 600;
      color: var(--tabs-tab-text-color);
      margin: 0;
      display: flex;
      align-items: center;
      gap: 10px;
      transition: all 100ms ease-in-out;
      svg {
        width: 20px;
        height: 20px;
        display: none;
        @media (min-width: 768px) {
          display: block;
        }
      }
    }

    &.active {
      background-color: var(--tabs-tab-background-color-active);
      h2 {
        color: var(--tabs-tab-text-color-active);
      }
    }

    &:hover {
      background-color: var(--tabs-tab-background-color-hover);
      h2 {
        color: var(--tabs-tab-text-color-hover);
      }
    }
  }
}

.settings-tabs-content-container {
  position: relative;
  flex-grow: 1;
  width: 100%;
  border-radius: 5px;
  background-color: var(--panel-background-color);

  .settings-tab-content {
    position: absolute;
    width: 100%;
    height: 100%;
    padding: 0px;

    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: flex-start;

    background-color: var(--panel-background-color);

    .tab-content-header {
      display: flex;
      justify-content: flex-end;
      @media (min-width: 768px) {
        justify-content: space-between;
      }
      align-items: center;
      background: var(--panel-subheader-background-color);
      padding: 20px;
      width: 100%;
      h2 {
        font-size: 1.4rem;
        color: var(--panel-subheader-text-color);
        margin: 0;
        display: flex;
        align-items: center;
        gap: 10px;
        svg {
          width: 20px;
          height: 20px;
        }
        small {
          display: block;
          font-size: 0.8rem;
          color: var(--panel-subheader-text-color);
          margin: 0;
        }
      }
      p {
        font-size: 1rem;
        color: var(--panel-subheader-text-color);
        margin: 0;
      }
    }

    .tab-content-body {
      display: block;
      padding: 0px;
      overflow-y: auto;
      flex-grow: 1;
      width: 100%;
    }
  }
}

// Cross-fade transition
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active {
  z-index: 1;
}
</style>
