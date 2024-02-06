 

import { SvgIcon, SvgIconProps, useTheme } from '@mui/material'

export const VersionIcon = (props: SvgIconProps) => {
  const theme = useTheme()
  return (
    <SvgIcon {...props} viewBox='0 0 21 21'>
      <g fill='none' fillRule='evenodd' stroke={theme.palette.text.primary}
         strokeLinecap='round' strokeLinejoin='round' transform='translate(2 4)'>
        <path d='m.5 8.5 8 4 8.017-4'></path>
        <path d='m.5 4.657 8.008 3.843 8.009-3.843-8.009-4.157z'></path>
      </g>
    </SvgIcon>
  )
}
