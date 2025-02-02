<script setup>
  import { store } from '../store'
  import { CircleX, Settings, SlidersHorizontal, Users as UsersIcon, UserPlus, Save } from 'lucide-vue-next'
  import { ref, onMounted } from 'vue'
  import Users from './settings/users.vue'

  //settins panels
  const usersPanel = ref(null)


  // Track active tab
  const activeTab = ref('general')

  // Create refs for the tab contents
  const tabContents = ref({
    general: ref(null),
    users: ref(null)
  })

  const closeSettings = () => {
    store.setSettingsOpen(false)
  }

  const clickOutside = e => {
    if (e.target === e.currentTarget) {
      closeSettings()
    }
  }

  const setActiveTab = tab => {
    activeTab.value = tab
  }
</script>

<template>
  <div class="settings-overlay" :class="{ active: store.settingsOpen }" @click="clickOutside">
    <div class="settings-container">
      <div class="settings-header">
        <h1>
          <Settings />
          erugo Settings
        </h1>
        <button class="close-settings-button" @click="closeSettings"><CircleX /></button>
      </div>
      <div class="settings-tabs-wrapper">
        <div class="settings-tabs-container">
          <div class="settings-tab" :class="{ active: activeTab === 'general' }" @click="setActiveTab('general')">
            <h2><SlidersHorizontal /> General</h2>
          </div>
          <div class="settings-tab" :class="{ active: activeTab === 'users' }" @click="setActiveTab('users')">
            <h2><UsersIcon /> Users</h2>
          </div>
        </div>
        <div class="settings-tabs-content-container">
          <Transition name="fade">
            <div v-if="activeTab === 'general'" class="settings-tab-content" ref="tabContents.general" key="general">
              <div class="tab-content-header">
                <h2><SlidersHorizontal /> General</h2>
                <div class="user-actions">
                  <button>
                    <Save />
                    Save General Settings
                  </button>
                </div>
              </div>
            </div>
            <div v-else-if="activeTab === 'users'" class="settings-tab-content" ref="tabContents.users" key="users">
              <div class="tab-content-header">
                <h2><UsersIcon /> Users</h2>
                <div class="user-actions">
                  <button @click="usersPanel.addUser">
                    <UserPlus />
                    Add User
                  </button>
                </div>
              </div>
              <div class="tab-content-body">
                <Users ref="usersPanel" />
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
    background-color: rgba(0, 0, 0, 0);
    width: 100%;
    height: 100%;
    z-index: 9999999999;
    pointer-events: none;
    transition: all 300ms ease-in-out;
    transition-delay: 300ms;

    .settings-container {
      --settings-width: calc(100vw - 100px);
      --settings-height: calc(100vh - 100px);
      position: absolute;
      bottom: 0;
      left: 0;
      transform: translateX(calc(50vw - var(--settings-width) / 2)) translateY(100%);
      background-color: rgba(255, 255, 255, 0.9);
      border-radius: 10px 10px 0 0;
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
      background-color: rgba(0, 0, 0, 0.6);
      pointer-events: auto;
      transition-delay: 0s;
      backdrop-filter: blur(10px);
      

      .settings-container {
        transform: translateX(calc(50vw - var(--settings-width) / 2)) translateY(0);
        transition-delay: 100ms;
        box-shadow: 0 -10px 30px 0 rgba(0, 0, 0, 0.5);
      }
    }
  }

  .settings-header {
    background-color: #eec154;
    border-radius: 5px 5px 0 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 80px;
    width: 100%;
    h1 {
      font-size: 20px;
      font-weight: 600;
      color: #222222;
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
    background-color: #e2e2e2;
    width: 100%;
    .settings-tab {
      background-color: #222222;
      margin-top: 10px;
      padding: 10px;
      border-radius: 5px 5px 0 0;
      box-shadow: inset 0 -5px 2px 0 rgba(0, 0, 0, 0.05);
      cursor: pointer;
      transition: background-color 0.2s ease;

      h2 {
        font-size: 16px;
        font-weight: 600;
        color: #fff;
        margin: 0;
        display: flex;
        align-items: center;
        gap: 10px;
        svg {
          width: 20px;
          height: 20px;
        }
      }

      &.active {
        background-color: #eec154;
        h2 {
          color: #222222;
        }
      }
    }
  }

  .settings-tabs-content-container {
    position: relative;
    flex-grow: 1;
    width: 100%;
    background-color: #eaeaea;
    border-radius: 5px;

    .settings-tab-content {
      position: absolute;
      width: 100%;
      height: 100%;
      padding: 0px;

      display: flex;
      flex-direction: column;
      align-items: flex-start;
      justify-content: flex-start;

      .tab-content-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        background-color: #f4f4f4;
        padding: 20px;
        width: 100%;
        h2 {
          font-size: 1.4rem;
          color: #222222;
          margin: 0;
          display: flex;
          align-items: center;
          gap: 10px;
          svg {
            width: 20px;
            height: 20px;
          }
        }
        p {
          font-size: 1.2rem;
          color: #222222;
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
