 

import { styled } from '@mui/material/styles'
import ReactTypingEffect from 'react-typing-effect'

export const TypingMessage = styled(ReactTypingEffect)(({ theme }) => ({
  width: '100%',

  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',

  fontFamily: 'Courier New,Courier,Lucida Sans Typewriter,Lucida Typewriter,monospace',
  backgroundColor: theme.palette.background.default,
  color: theme.palette.text.secondary,
  fontSize: '3rem',
  marginBottom: theme.spacing(8),
  '& > span': {
    backgroundColor: theme.palette.background.default,
  },
}))
