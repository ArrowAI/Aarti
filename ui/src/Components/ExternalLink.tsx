 

import { BoxProps } from '@mui/material'
import Box from '@mui/material/Box'
import { LinkProps } from '@mui/material/Link'
import { styled } from '@mui/material/styles'
import { forwardRef } from 'react'

const Link = forwardRef<{}, Omit<BoxProps, 'component' | 'target'> & LinkProps>((props, ref) => (
  <Box
    {...props} target='_blank'
    component='a'
    ref={ref}
  />
))
export const ExternalLink = styled(Link, { shouldForwardProp: prop => prop != 'component' })(({ theme }) => ({
  color: theme.palette.text.secondary,
  textDecoration: 'none',
  '&:hover': {
    color: theme.palette.primary.main,
  },
}))
