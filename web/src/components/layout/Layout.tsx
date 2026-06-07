import { Outlet } from 'react-router-dom'
import { Header } from './Header'
import { Footer } from './Footer'
import { isEmbedMode } from '../../embedHandoff'

export function Layout() {
  const embed = isEmbedMode()

  return (
    <>
      {!embed && <Header />}
      <main className={`main-content${embed ? ' embed-content' : ''}`}>
        <Outlet />
      </main>
      {!embed && <Footer />}
    </>
  )
}
