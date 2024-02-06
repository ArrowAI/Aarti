

import { CssBaseline, PaletteMode, ThemeProvider, useMediaQuery } from '@mui/material'
import React, { PropsWithChildren, useContext, useEffect, useState } from 'react'
import { Helmet, HelmetProvider } from 'react-helmet-async'
import { configureTheme } from './configureTheme'
import { UiMode } from './theme'

const LOCAL_STORAGE_KEY = 'ui-mode'

const modeFromLocalStorage = (): UiMode => {
  const mode = localStorage.getItem(LOCAL_STORAGE_KEY)
  switch (mode) {
    case 'light':
    case 'dark':
    case undefined:
      return mode
    case '':
    default:
      return undefined
  }
}

interface ModeContext {
  mode: UiMode
  setMode: (mode: UiMode) => void
}

const modeContext = React.createContext<ModeContext>({
  mode: undefined,
  setMode: () => {
  },
})

export const useUiMode = () => {
  const ctx = useContext(modeContext)
  return { mode: ctx.mode, setMode: ctx.setMode }
}

export const ColorModeThemeProvider = ({ children }: PropsWithChildren<any>) => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)')

  const restoredMode = modeFromLocalStorage()
  // console.log('mode: restore', restoredMode || 'system')
  const [storedMode, _setStoredMode] = useState<UiMode>(restoredMode)
  const setStoredMode = (mode: UiMode) => {
    localStorage.setItem(LOCAL_STORAGE_KEY, mode || '')
    _setStoredMode(mode)
  }

  const [appliedMode, setAppliedMode] = useState<PaletteMode>(restoredMode || (prefersDarkMode ? 'dark' : 'light'))
  // console.log('mode: apply', restoredMode || 'system')
  const setMode = (mode: UiMode) => {
    const m = mode || (prefersDarkMode ? 'dark' : 'light')
    setStoredMode(mode)
    setAppliedMode(m)
  }

  useEffect(() => {
    if (storedMode) {
      return
    }
    setMode(undefined)
  }, [prefersDarkMode])

  const theme = React.useMemo(() => configureTheme(appliedMode), [appliedMode])
  return (
    <modeContext.Provider value={{ mode: appliedMode, setMode: setMode }}>
      <HelmetProvider>
        <ThemeProvider theme={theme}>
          <Helmet>
            <meta name='theme-color' content={appliedMode === 'dark' ? '#000000' : '#FFFFFF'} />
          </Helmet>
          <CssBaseline />
          {children}
        </ThemeProvider>
      </HelmetProvider>
    </modeContext.Provider>
  )
}
