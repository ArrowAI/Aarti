 

import { Stack, Tab as MuiTab, Tabs as MuiTabs, Typography } from '@mui/material'
import { styled } from '@mui/material/styles'
import React, { PropsWithChildren } from 'react'
import { usePersistedState } from '../hooks'
import { Void } from '../utils'
import { Code, CodeProps } from './Code'

const Tabs = styled(MuiTabs)(({ theme }) => ({
  minHeight: 0,
  // background: 'black',
  '& .MuiTabs-indicator': {
    display: 'none',
  },
  '& .MuiTabs-flexContainer': {
    justifyContent: 'flex-end',
  },
}))

const Tab = styled(MuiTab)(({ theme }) => ({
  textTransform: 'none',
  // color: 'white',
  padding: `0 ${theme.spacing(2)}`,
  minWidth: 0,
  minHeight: 32,
  [theme.breakpoints.up('sm')]: {
    minWidth: 0,
    minHeight: 32,
  },
}))

export interface MultiLangCodeItemProps extends CodeProps {
  label: string
}

export interface MultiLangCodeContext {
  storageKey: string
  value?: string
  setValue: (v: string) => void
}

const multiLangCodeContext = React.createContext<MultiLangCodeContext>({ storageKey: '', value: '', setValue: Void })

export interface MultiLangCodeProviderProps extends PropsWithChildren<any>, Omit<MultiLangCodeContext, 'setValue'> {}

export const MultiLangCodeProvider = ({ children, storageKey, value: _value }: MultiLangCodeProviderProps) => {
  const [value, setValue] = usePersistedState(_value, 'MultiLangCode-' + storageKey)
  return (
    <multiLangCodeContext.Provider value={{ storageKey, value, setValue }}>
      {children}
    </multiLangCodeContext.Provider>
  )
}

export const useMultiLangCode = () => {
  return React.useContext(multiLangCodeContext)
}

export const MultiLangCodeItem = (props: MultiLangCodeItemProps) => <Code {...props} />

export interface MultiLangCodeProps {
  storageKey: string
  title?: string
  children: React.ReactElement<MultiLangCodeItemProps> | React.ReactElement<MultiLangCodeItemProps>[]
}

export const MultiLangCode = ({ title, children }: MultiLangCodeProps) => {
  const { value, setValue } = useMultiLangCode()
  const handleChange = (_: React.SyntheticEvent, newValue: string) => {
    setValue(newValue)
  }
  const e = Array.isArray(children) ? children.find(c => c.props.label === value) : children
  return (
    <>
      <Stack direction='row' justifyContent='space-between' alignItems='center'>
        <Typography variant='body2'>{title}</Typography>
        <Tabs
          value={value}
          onChange={handleChange}
        >
          {Array.isArray(children) ? children.map((e, i) => <Tab key={i} label={e.props.label}
                                                                 value={e.props.label} />) :
            <Tab label={children.props.label} value={children.props.label} />}
        </Tabs>
      </Stack>
      {e && <Code sx={{ marginTop: 0 }} {...e.props} />}
    </>
  )
}
