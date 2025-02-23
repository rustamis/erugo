import { createApp } from 'vue'
import './style.scss'
import 'vue-color-kit/dist/vue-color-kit.css'
import App from './App.vue'
import Toast, { POSITION } from 'vue-toastification'
import 'vue-toastification/dist/index.css'

let clickOutsideEvent = null

const showHelpTip = (event, helpTipId) => {
  const helpTip = document.querySelector(helpTipId)
  const cssShiftAmount = -15
  if (!helpTip) return

  helpTip.classList.add('no-transition')
  let top = event.clientY - helpTip.offsetHeight - 10
  let left = event.clientX - helpTip.offsetWidth / 2

  //will this go out of bounds?
  if (top + cssShiftAmount < 0) {
    top = 0
  }
  if (left < 0) {
    left = 0
  }
  if (left + helpTip.offsetWidth > window.innerWidth) {
    left = window.innerWidth - helpTip.offsetWidth
  }
  if (top + cssShiftAmount + helpTip.offsetHeight > window.innerHeight) {
    top = window.innerHeight - helpTip.offsetHeight - cssShiftAmount
  }

  helpTip.style.top = `${top}px`
  helpTip.style.left = `${left}px`

  setTimeout(() => {
    helpTip.classList.remove('no-transition')
    helpTip.classList.add('visible')
    helpTip.classList.add('shifted-up')
  }, 1)

  //register click-outside event
  clickOutsideEvent = document.addEventListener('click', (event) => {
    if (!helpTip.contains(event.target)) {
      hideHelpTip(helpTipId)
    }
  })
}

const hideHelpTip = (helpTipId) => {
  const helpTip = document.querySelector(helpTipId)
  if (!helpTip) return

  helpTip.classList.remove('visible')
  helpTip.classList.remove('shifted-up')

  //unregister click-outside event
  document.removeEventListener('click', clickOutsideEvent)
}

//use a mutation observer to recursively watch the dom for any .slide-text elements that might appear
const observer = new MutationObserver((mutations) => {
  const processNode = (node) => {
    if (node.classList && node.classList.contains('slide-text')) {
      const innerContent = node.querySelector('.content')
      if (innerContent) {
        const outerWidth = node.offsetWidth
        const innerWidth = innerContent.offsetWidth
        if (innerWidth > outerWidth) {
          innerContent.classList.add('needs-slide')
          node.classList.add('has-slide')
          innerContent.style.setProperty('--content-width', `${innerWidth}px`)
          const animationTime = Math.max(innerWidth / 100, 20) // Ensure a minimum animation time of 20 seconds
          innerContent.style.setProperty('--animation-time', `${animationTime}s`)
        }
      }
    }
    node.childNodes.forEach((child) => processNode(child))
  }

  mutations.forEach((mutation) => {
    mutation.addedNodes.forEach((node) => {
      processNode(node)
    })
  })
})

observer.observe(document.body, { childList: true, subtree: true })

createApp(App)
  .use(Toast, {
    position: POSITION.BOTTOM_RIGHT
  })
  .provide('showHelpTip', showHelpTip)
  .provide('hideHelpTip', hideHelpTip)
  .mount('#app')
