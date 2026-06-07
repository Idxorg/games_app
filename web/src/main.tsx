import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import { useUserStore } from './stores/userStore'
import { isEmbedMode, initEmbedHandoff } from './embedHandoff'
import './index.css'

// Initialize auth before mounting the app
async function bootstrap() {
  if (isEmbedMode()) {
    try {
      const session = await initEmbedHandoff()
      // loginFromEmbed handles token exchange AND stores the handoff theme
      await useUserStore.getState().loginFromEmbed(session)
    } catch (err) {
      console.error('Embed auth init failed:', err)
    }
  } else {
    // Standalone mode — initialize store (user can browse, auth optional for demo)
    await useUserStore.getState().initialize()
  }

  ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  )
}

bootstrap()
